package graphql

import (
	"fmt"
	"strconv"

	"github.com/graphql-go/graphql/language/ast"
)

var (
	MaxInt = 9007199254740991
	MinInt = -9007199254740991
)

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
		return intOrNil(int(value))
	case float64:
		return intOrNil(int(value))
	case string:
		val, err := strconv.ParseFloat(value, 0)
		if err != nil {
			return nil
		}
		return coerceInt(val)
	}
	return int(0)
}

// Integers are only safe when between -(2^53 - 1) and 2^53 - 1 due to being
// encoded in JavaScript and represented in JSON as double-precision floating
// point numbers, as specified by IEEE 754.
func intOrNil(value int) interface{} {
	if value <= MaxInt && value >= MinInt {
		return value
	}
	return nil
}

var Int *Scalar = NewScalar(ScalarConfig{
	Name:       "Int",
	Serialize:  coerceInt,
	ParseValue: coerceInt,
	ParseLiteral: func(valueAST ast.Value) interface{} {
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
			return nil
		}
		return coerceFloat32(val)
	}
	return float32(0)
}

var Float *Scalar = NewScalar(ScalarConfig{
	Name:       "Float",
	Serialize:  coerceFloat32,
	ParseValue: coerceFloat32,
	ParseLiteral: func(valueAST ast.Value) interface{} {
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
		return nil
	},
})

func coerceString(value interface{}) interface{} {
	return fmt.Sprintf("%v", value)
}

var String *Scalar = NewScalar(ScalarConfig{
	Name:       "String",
	Serialize:  coerceString,
	ParseValue: coerceString,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return valueAST.Value
		}
		return nil
	},
})

func coerceBool(value interface{}) interface{} {
	switch value := value.(type) {
	case bool:
		return value
	case string:
		switch value {
		case "", "false":
			return false
		}
		return true
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

var Boolean *Scalar = NewScalar(ScalarConfig{
	Name:       "Boolean",
	Serialize:  coerceBool,
	ParseValue: coerceBool,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.BooleanValue:
			return valueAST.Value
		}
		return nil
	},
})

var ID *Scalar = NewScalar(ScalarConfig{
	Name:       "ID",
	Serialize:  coerceString,
	ParseValue: coerceString,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			return valueAST.Value
		case *ast.StringValue:
			return valueAST.Value
		}
		return nil
	},
})
