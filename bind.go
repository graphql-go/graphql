package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

var ctxType = reflect.TypeOf((*context.Context)(nil)).Elem()
var errType = reflect.TypeOf((*error)(nil)).Elem()

/*
	Bind will create a Field around a function formatted a certain way, or any value.

	The input parameters can be, in any order,
	- context.Context, or *context.Context (optional)
	- An input struct, or pointer (optional)

	The output parameters can be, in any order,
	- A primitive, an output struct, or pointer (required for use in schema)
	- error (optional)

	Input or output types provided will be automatically bound using BindType.
*/
func Bind(bindTo interface{}, additionalFields ...Fields) *Field {
	combinedAdditionalFields := MergeFields(additionalFields...)
	val := reflect.ValueOf(bindTo)
	tipe := reflect.TypeOf(bindTo)
	if tipe.Kind() == reflect.Func {
		in := tipe.NumIn()
		out := tipe.NumOut()

		var ctxIn *int
		var inputIn *int

		var errOut *int
		var outputOut *int

		queryArgs := FieldConfigArgument{}

		if in > 2 {
			panic(fmt.Sprintf("Mismatch on number of inputs. Expected 0, 1, or 2. got %d.", tipe.NumIn()))
		}

		if out > 2 {
			panic(fmt.Sprintf("Mismatch on number of outputs. Expected 0, 1, or 2, got %d.", tipe.NumOut()))
		}

		// inTypes := make([]reflect.Type, in)
		// outTypes := make([]reflect.Type, out)

		for i := 0; i < in; i++ {
			t := tipe.In(i)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			switch t {
			case ctxType:
				if ctxIn != nil {
					panic(fmt.Sprintf("Unexpected multiple *context.Context inputs."))
				}
				ctxIn = intP(i)
			default:
				if inputIn != nil {
					panic(fmt.Sprintf("Unexpected multiple inputs."))
				}
				inputType := tipe.In(i)
				if inputType.Kind() == reflect.Ptr {
					inputType = inputType.Elem()
				}
				inputFields := BindFields(reflect.New(inputType).Interface())
				for key, inputField := range inputFields {
					queryArgs[key] = &ArgumentConfig{
						Type: inputField.Type,
					}
				}

				inputIn = intP(i)
			}
		}

		for i := 0; i < out; i++ {
			t := tipe.Out(i)
			switch t.String() {
			case errType.String():
				if errOut != nil {
					panic(fmt.Sprintf("Unexpected multiple error outputs"))
				}
				errOut = intP(i)
			default:
				if outputOut != nil {
					panic(fmt.Sprintf("Unexpected multiple outputs"))
				}
				outputOut = intP(i)
			}
		}

		resolve := func(p ResolveParams) (output interface{}, err error) {
			inputs := make([]reflect.Value, in)
			if ctxIn != nil {
				isPtr := tipe.In(*ctxIn).Kind() == reflect.Ptr
				if isPtr {
					if p.Context == nil {
						inputs[*ctxIn] = reflect.New(ctxType)
					} else {
						inputs[*ctxIn] = reflect.ValueOf(&p.Context)
					}
				} else {
					if p.Context == nil {
						inputs[*ctxIn] = reflect.New(ctxType).Elem()
					} else {
						inputs[*ctxIn] = reflect.ValueOf(p.Context).Convert(ctxType).Elem()
					}
				}
			}
			if inputIn != nil {
				var inputType, inputBaseType, sourceType, sourceBaseType reflect.Type
				sourceVal := reflect.ValueOf(p.Source)
				sourceExists := !sourceVal.IsZero()
				if sourceExists {
					sourceType = sourceVal.Type()
					if sourceType.Kind() == reflect.Ptr {
						sourceBaseType = sourceType.Elem()
					} else {
						sourceBaseType = sourceType
					}
				}
				inputType = tipe.In(*inputIn)
				isPtr := tipe.In(*inputIn).Kind() == reflect.Ptr
				if isPtr {
					inputBaseType = inputType.Elem()
				} else {
					inputBaseType = inputType
				}
				var input interface{}
				if sourceExists && sourceBaseType.AssignableTo(inputBaseType) {
					input = sourceVal.Interface()
				} else {
					input = reflect.New(inputBaseType).Interface()
					j, err := json.Marshal(p.Args)
					if err == nil {
						err = json.Unmarshal(j, &input)
					}
					if err != nil {
						return nil, err
					}
				}

				inputs[*inputIn], err = convertValue(reflect.ValueOf(input), inputType)
				if err != nil {
					return nil, err
				}
			}
			results := val.Call(inputs)
			if errOut != nil {
				val := results[*errOut].Interface()
				if val != nil {
					err = val.(error)
				}
				if err != nil {
					return output, err
				}
			}
			if outputOut != nil {
				var val reflect.Value
				val, err = convertValue(results[*outputOut], tipe.Out(*outputOut))
				if err != nil {
					return nil, err
				}
				if !val.IsZero() {
					output = val.Interface()
				}
			}
			return output, err
		}

		var outputType Output
		if outputOut != nil {
			outputType = BindType(tipe.Out(*outputOut))
			extendType(outputType, combinedAdditionalFields)
		}

		field := &Field{
			Type:    outputType,
			Resolve: resolve,
			Args:    queryArgs,
		}

		return field
	} else if tipe.Kind() == reflect.Struct {
		fieldType := BindType(reflect.TypeOf(bindTo))
		extendType(fieldType, combinedAdditionalFields)
		field := &Field{
			Type: fieldType,
			Resolve: func(p ResolveParams) (data interface{}, err error) {
				return bindTo, nil
			},
		}
		return field
	} else {
		if len(additionalFields) > 0 {
			panic("Cannot add field resolvers to a scalar type.")
		}
		return &Field{
			Type: getGraphType(tipe),
			Resolve: func(p ResolveParams) (data interface{}, err error) {
				return bindTo, nil
			},
		}
	}
}

func extendType(t Type, fields Fields) {
	switch t.(type) {
	case *Object:
		object := t.(*Object)
		for fieldName, fieldConfig := range fields {
			object.AddFieldConfig(fieldName, fieldConfig)
		}
		return
	case *List:
		list := t.(*List)
		extendType(list.OfType, fields)
		return
	}
}

func convertValue(value reflect.Value, targetType reflect.Type) (ret reflect.Value, err error) {
	if !value.IsValid() || value.IsZero() {
		return reflect.Zero(targetType), nil
	}
	if value.Type().Kind() == reflect.Ptr {
		if targetType.Kind() == reflect.Ptr {
			return value, nil
		} else {
			return value.Elem(), nil
		}
	} else {
		if targetType.Kind() == reflect.Ptr {
			// Will throw an informative error
			return value.Convert(targetType), nil
		} else {
			return value, nil
		}
	}
}

func intP(i int) *int {
	return &i
}
