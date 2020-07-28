package graphql

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Bind will create a Field around a function formatted a certain way, or any value.
// The function can be formatted:
// - func MyFunction(ctx *context.Context) (output *MyFunctionOutput, err error)
// - func MyFunction(ctx *context.Context, input *MyFunctionInput) (output *MyFunctionOutput, err error)
// The input type, if provided, will be bound to fields, and the output type will be automatically bound as well.
func Bind(bindTo interface{}) *Field {
	val := reflect.ValueOf(bindTo)
	tipe := reflect.TypeOf(bindTo)
	if tipe.Kind() == reflect.Func {
		in := tipe.NumIn()
		out := tipe.NumOut()

		if in != 1 && in != 2 {
			panic(fmt.Sprintf("Mismatch on number of arguments. Expected 1 or 2, got %d.", tipe.NumIn()))
		}

		if out != 1 && out != 2 {
			panic(fmt.Sprintf("Mismatch on number of outputs. Expected 1 or 2, got %d.", tipe.NumOut()))
		}

		var inputType reflect.Type
		if in == 1 {
			inputType = reflect.TypeOf(struct{}{})
		} else {
			inputType = tipe.In(1)
		}

		if inputType.Kind() == reflect.Ptr {
			inputType = inputType.Elem()
		}
		inputFields := BindFields(reflect.New(inputType).Interface())
		queryArgs := FieldConfigArgument{}
		for key, inputField := range inputFields {
			queryArgs[key] = &ArgumentConfig{
				Type: inputField.Type,
			}
		}

		outputType := tipe.Out(0)
		if outputType.Kind() == reflect.Ptr {
			outputType = outputType.Elem()
		}
		output := BindType(outputType)

		resolve := func(p ResolveParams) (data interface{}, err error) {
			input := reflect.New(inputType).Interface()
			j, err := json.Marshal(p.Args)

			if err == nil {
				err = json.Unmarshal(j, &input)
			}

			if err != nil {
				return nil, err
			}

			var results []reflect.Value
			if in == 1 {
				// Simple field
				results = val.Call([]reflect.Value{
					reflect.ValueOf(&p.Context),
				})
			} else {
				// Query with argument
				results = val.Call([]reflect.Value{
					reflect.ValueOf(&p.Context),
					reflect.ValueOf(input),
				})
			}

			if out == 2 && results[1].Interface() != nil {
				// Error
				err = results[1].Interface().(error)
				return nil, err
			} else {
				// Success
				result := results[0].Interface()
				return result, nil
			}
		}

		field := &Field{
			Type:    output,
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
