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
func Bind(bindTo interface{}) *Field {
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
					inputs[*ctxIn] = reflect.ValueOf(&p.Context)
				} else {
					if p.Context == nil {
						inputs[*ctxIn] = reflect.New(ctxType).Elem()
					} else {
						inputs[*ctxIn] = reflect.ValueOf(p.Context).Elem()
					}
				}
			}
			if inputIn != nil {
				var inputType, baseType reflect.Type
				inputType = tipe.In(*inputIn)
				isPtr := tipe.In(*inputIn).Kind() == reflect.Ptr
				if isPtr {
					baseType = inputType.Elem()
				} else {
					baseType = inputType
				}
				input := reflect.New(baseType).Interface()
				j, err := json.Marshal(p.Args)
				if err == nil {
					err = json.Unmarshal(j, &input)
				}
				if err != nil {
					return nil, err
				}

				if isPtr {
					inputs[*inputIn] = reflect.ValueOf(input)
				} else {
					if input == nil {
						inputs[*inputIn] = reflect.New(baseType).Elem()
					} else {
						inputs[*inputIn] = reflect.ValueOf(input).Elem()
					}
				}
			}

			results := val.Call(inputs)
			if errOut != nil {
				val := results[*errOut].Interface()
				if val != nil {
					err = val.(error)
				}
			}
			if outputOut != nil {
				output = results[*outputOut].Interface()
			}
			return output, err
		}

		var outputType Output
		if outputOut != nil {
			outputType = BindType(tipe.Out(*outputOut))
		}

		field := &Field{
			Type:    outputType,
			Resolve: resolve,
			Args:    queryArgs,
		}

		return field
	} else if tipe.Kind() == reflect.Struct {
		fieldType := BindType(reflect.TypeOf(bindTo))
		field := &Field{
			Type: fieldType,
			Resolve: func(p ResolveParams) (data interface{}, err error) {
				return bindTo, nil
			},
		}
		return field
	} else {
		return &Field{
			Type: getGraphType(tipe),
			Resolve: func(p ResolveParams) (data interface{}, err error) {
				return bindTo, nil
			},
		}
	}
}

func intP(i int) *int {
	return &i
}
