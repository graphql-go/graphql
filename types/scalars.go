package types

import (
	"fmt"
	"github.com/chris-ramon/graphql-go/language/ast"
	"strconv"
)

type SerializeFn func(value interface{}) interface{}
type ParseValueFn func(value interface{}) interface{}
type ParseLiteralFn func(valueAST ast.Value) interface{}
type GraphQLScalarTypeConfig struct {
	Name         string
	Description  string
	Serialize    SerializeFn
	ParseValue   ParseValueFn
	ParseLiteral ParseLiteralFn
}

// GraphQLScalarType implements GraphQLType, GraphQLInputType, GraphQLNamedType,
// 								GraphQLOutputType, etc
// (TODO: find out what other interfaces GraphQLScalarType implements)
type GraphQLScalarType struct {
	Name        string
	Description string

	scalarConfig GraphQLScalarTypeConfig

	err error
}

func NewGraphQLScalarType(config GraphQLScalarTypeConfig) *GraphQLScalarType {
	st := &GraphQLScalarType{}
	err := invariant(config.Name != "", "Type must be named.")
	if err != nil {
		st.err = err
		return st
	}

	err = assertValidName(config.Name)
	if err != nil {
		st.err = err
		return st
	}

	st.Name = config.Name
	st.Description = config.Description

	err = invariant(
		config.Serialize != nil,
		fmt.Sprintf(`%v must provide "serialize" function. If this custom Scalar is `+
			`also used as an input type, ensure "parseValue" and "parseLiteral" `+
			`functions are also provided.`, st),
	)
	if err != nil {
		st.err = err
		return st
	}

	st.scalarConfig = config
	return st
}

func (st *GraphQLScalarType) Serialize(value interface{}) interface{} {
	if st.scalarConfig.Serialize == nil {
		return value
	}
	return st.scalarConfig.Serialize(value)
}
func (st *GraphQLScalarType) ParseValue(value interface{}) interface{} {
	if st.scalarConfig.ParseValue == nil {
		return value
	}
	return st.scalarConfig.ParseValue(value)
}
func (st *GraphQLScalarType) ParseLiteral(valueAST ast.Value) interface{} {
	if st.scalarConfig.ParseLiteral == nil {
		return valueAST
	}
	return st.scalarConfig.ParseLiteral(valueAST)
}

func (st *GraphQLScalarType) GetName() string {
	return st.Name
}
func (st *GraphQLScalarType) SetName(name string) {
	st.Name = name
}
func (st *GraphQLScalarType) GetDescription() string {
	return st.Description

}
func (st *GraphQLScalarType) String() string {
	return st.Name
}

// TODO: GraphQLScalarType.Coerce() Check if we need this
func (st *GraphQLScalarType) Coerce(value interface{}) (r interface{}) {
	return value

}
func (st *GraphQLScalarType) CoerceLiteral(value interface{}) (r interface{}) {
	return value
}

func coerceInt(value interface{}) interface{} {
	switch value := value.(type) {
	case bool:
		if value == true {
			return int(1)
		}
		return int(0)
	case int:
		return value
	case float32:
		return int(value)
	case float64:
		return int(value)
	case string:
		val, err := strconv.ParseFloat(value, 0)
		if err != nil {
			return int(0)
		}
		return coerceInt(val)
	}
	return int(0)
}

var GraphQLInt *GraphQLScalarType = NewGraphQLScalarType(GraphQLScalarTypeConfig{
	Name:       "Int",
	Serialize:  coerceInt,
	ParseValue: coerceInt,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		// TODO: can move this into each ast.Value.GetValue() implementation
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			if intValue, err := strconv.Atoi(valueAST.Value); err == nil {
				return intValue
			}
		}
		return nil
	},
})

func coerceFloat32(value interface{}) interface{} {
	switch value := value.(type) {
	case bool:
		if value == true {
			return float32(1)
		}
		return float32(0)
	case int:
		return float32(value)
	case float32:
		return value
	case float64:
		return float32(value)
	case string:
		val, err := strconv.ParseFloat(value, 0)
		if err != nil {
			return float32(0)
		}
		return coerceFloat32(val)
	}
	return float32(0)
}

var GraphQLFloat *GraphQLScalarType = NewGraphQLScalarType(GraphQLScalarTypeConfig{
	Name:       "Float",
	Serialize:  coerceFloat32,
	ParseValue: coerceFloat32,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		// TODO: can move this into each ast.Value.GetValue() implementation
		switch valueAST := valueAST.(type) {
		case *ast.FloatValue:
			if floatValue, err := strconv.ParseFloat(valueAST.Value, 32); err == nil {
				return floatValue
			}
		case *ast.IntValue:
			if floatValue, err := strconv.ParseFloat(valueAST.Value, 32); err == nil {
				return floatValue
			}
		}
		return float32(0)
	},
})

func coerceString(value interface{}) interface{} {
	return fmt.Sprintf("%v", value)
}

var GraphQLString *GraphQLScalarType = NewGraphQLScalarType(GraphQLScalarTypeConfig{
	Name:       "String",
	Serialize:  coerceString,
	ParseValue: coerceString,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		// TODO: can move this into each ast.Value.GetValue() implementation
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return valueAST.Value
		}
		return ""
	},
})

func coerceBool(value interface{}) interface{} {
	switch value := value.(type) {
	case bool:
		return value
	case string:
		if value == "true" {
			return true
		}
		return false
	case float64:
		if value != 0 {
			return true
		}
		return false
	case float32:
		if value != 0 {
			return true
		}
		return false
	case int:
		if value != 0 {
			return true
		}
		return false
	}
	return false
}

var GraphQLBoolean *GraphQLScalarType = NewGraphQLScalarType(GraphQLScalarTypeConfig{
	Name:       "Boolean",
	Serialize:  coerceBool,
	ParseValue: coerceBool,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		// TODO: can move this into each ast.Value.GetValue() implementation
		switch valueAST := valueAST.(type) {
		case *ast.BooleanValue:
			return valueAST.Value
		}
		return false
	},
})

var GraphQLID *GraphQLScalarType = NewGraphQLScalarType(GraphQLScalarTypeConfig{
	Name:       "ID",
	Serialize:  coerceString,
	ParseValue: coerceString,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		// TODO: can move this into each ast.Value.GetValue() implementation
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			return valueAST.Value
		case *ast.StringValue:
			return valueAST.Value
		}
		return ""
	},
})
