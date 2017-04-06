package graphql

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"reflect"
)

const TAG = "json"

func BindFields(obj interface{}) Fields {
	v := reflect.ValueOf(obj)
	fields := make(map[string]*Field)

	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)

		tag := typeField.Tag.Get(TAG)
		if typeField.Type.Kind() == reflect.Struct && tag == "" {
			structFields := BindFields(v.Field(i).Interface())

			fields = appendFields(fields, structFields)
			continue
		}

		if tag == "" {
			continue
		}

		graphType := getGraphType(typeField.Type)
		fields[tag] = &Field{
			Type: graphType,
			Resolve: func(p ResolveParams) (interface{}, error) {
				return extractValue(tag, p.Source), nil
			},
		}
	}
	return fields
}

func getGraphType(tipe reflect.Type) Output {
	kind := tipe.Kind()
	switch kind {
	case reflect.String:
		return String
	case reflect.Int:
		return Int
	case reflect.Float32:
	case reflect.Float64:
		return Float
	case reflect.Bool:
		return Boolean
	case reflect.Slice:
		return getGraphList(tipe)
	}
	return String
}

func getGraphList(tipe reflect.Type) *List {
	switch tipe {
	case reflect.TypeOf([]int{}):
	case reflect.TypeOf([]int8{}):
	case reflect.TypeOf([]int32{}):
	case reflect.TypeOf([]int64{}):
		return NewList(Int)
	case reflect.TypeOf([]bool{}):
		return NewList(Boolean)
	case reflect.TypeOf([]float32{}):
	case reflect.TypeOf([]float64{}):
		return NewList(Float)
	}
	return NewList(String)
}

func appendFields(dest, origin Fields) Fields {
	for key, value := range origin {
		dest[key] = value
	}
	return dest
}

func extractValue(originTag string, obj interface{}) interface{} {
	val := reflect.ValueOf(obj)

	for j := 0; j < val.NumField(); j++ {
		typeField := val.Type().Field(j)
		if typeField.Type.Kind() == reflect.Struct {
			res := extractValue(originTag, val.Field(j).Interface())
			if res != nil {
				return res
			}
		}
		curTag := typeField.Tag
		if originTag == curTag.Get(TAG) {
			return val.Field(j).Interface()
		}
	}
	return nil
}

// lazy way of binding args
func BindArg(obj interface{}, tags ...string) FieldConfigArgument {
	v := reflect.ValueOf(obj)
	var config = make(FieldConfigArgument)
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)

		mytag := typeField.Tag.Get(TAG)
		if inArray(tags, mytag) {
			config[mytag] = &ArgumentConfig{
				Type: getGraphType(typeField.Type),
			}
		}
	}
	return config
}

func UnmarshalArgs(args map[string]interface{}, pointer interface{}) error {
	js, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("ini error 1 %v", err)
	}
	err = json.Unmarshal(js, pointer)
	if err != nil {
		return fmt.Errorf("ini error 2 %v", err)
	}
	return nil
}

func inArray(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("inArray() given a non-slice type")
	}

	for i := 0; i < s.Len(); i++ {
		if reflect.DeepEqual(item, s.Index(i).Interface()) {
			return true
		}
	}
	return false
}
