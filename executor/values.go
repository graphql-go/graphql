package executor

import (
	"fmt"
	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/types"
	"github.com/kr/pretty"
)

//export function getVariableValues(
//schema: GraphQLSchema,
//definitionASTs: Array<VariableDefinition>,
//inputs: { [key: string]: any }
//): { [key: string]: any } {
//return definitionASTs.reduce((values, defAST) => {
//var varName = defAST.variable.name.value;
//values[varName] = getVariableValue(schema, defAST, inputs[varName]);
//return values;
//}, {});
//}

func GetVariableValues(schema types.GraphQLSchema, definitionASTs []*ast.VariableDefinition, inputs map[string]interface{}) (map[string]interface{}, error) {

	pretty.Println("GetVariableValues", schema, definitionASTs, inputs)
	values := map[string]interface{}{}
	for _, defAST := range definitionASTs {
		if defAST == nil {
			continue
		}
		if defAST.Variable == nil {
			continue
		}
		if defAST.Variable.Name == nil {
			continue
		}
		varName := defAST.Variable.Name.Value
		varValue, err := getVariableValue(schema, defAST, inputs[varName])
		if err != nil {
			return values, err
		}
		values[varName] = varValue
	}
	pretty.Println("GetVariableValues", values)
	return values, nil
}

// Given a variable definition, and any value of input, return a value which
// adheres to the variable definition, or throw an error.
func getVariableValue(schema types.GraphQLSchema, definitionAST *ast.VariableDefinition, input interface{}) (interface{}, error) {
	pretty.Println("getVariableValue input", input)

	ttype := typeFromAST(schema, definitionAST.Type)
	pretty.Println("getVariableValue ttype", ttype)
	variable := definitionAST.Variable
	if ttype == nil {
		return "", graphqlerrors.NewGraphQLError(
			fmt.Sprintf(`Variable "$%v" expected value of type
			"%v" which cannot be used as an input type.`, variable.Name.Value, definitionAST.Type),
			[]ast.Node{definitionAST},
			"",
			nil,
			[]int{},
		)
	}
	if isValidInputValue(input, ttype) {
		if isNullish(input) {
			defaultValue := definitionAST.DefaultValue
			if defaultValue != nil {
				//  TODO: Note: enforce that GraphQLType implements GraphQLInputType
				variables := map[string]interface{}{}
				return valueFromAST(defaultValue, ttype, variables)
			}
			return coerceValue(ttype, input)
		}
	}
	if isNullish(input) {
		return "", graphqlerrors.NewGraphQLError(
			fmt.Sprintf(`Variable "$%v" of required type
			"%v" was not provided.`, variable.Name.Value, definitionAST.Type),
			[]ast.Node{definitionAST},
			"",
			nil,
			[]int{},
		)
	}
	return "", graphqlerrors.NewGraphQLError(
		fmt.Sprintf(`Variable "$%v" expected value of type
			"%v" but got: %v.`, variable.Name.Value, definitionAST.Type, input),
		[]ast.Node{definitionAST},
		"",
		nil,
		[]int{},
	)
}

// Given a type and any value, return a runtime value coerced to match the type.
func coerceValue(ttype types.GraphQLInputType, value interface{}) (interface{}, error) {
	return value, nil
}

// graphql-js/src/utilities.js`

func typeFromAST(schema types.GraphQLSchema, inputTypeAST ast.Type) types.GraphQLType {
	switch inputTypeAST := inputTypeAST.(type) {
	case *ast.ListType:
		innerType := typeFromAST(schema, inputTypeAST.Type)
		return types.NewGraphQLList(innerType)
	case *ast.NonNullType:
		innerType := typeFromAST(schema, inputTypeAST.Type)
		return types.NewGraphQLList(innerType)
	case *ast.NamedType:
		nameValue := ""
		if inputTypeAST.Name != nil {
			nameValue = inputTypeAST.Name.Value
		}
		return schema.GetType(nameValue)
	}
	// TODO: do invariant check here
	return nil
}

// isValidInputValue alias isValidJSValue
// Given a JavaScript value and a GraphQL type, determine if the value will be
// accepted for that type. This is primarily useful for validating the
// runtime values of query variables.
func isValidInputValue(input interface{}, ttype types.GraphQLInputType) bool {
	// TODO: isValidInputValue
	return true
}

// Returns true if a value is null, undefined, or NaN.
func isNullish(value interface{}) bool {
	// TODO: rethink isNullish. Do we need it?
	return true
}

/**
 * Produces a value given a GraphQL Value AST.
 *
 * A GraphQL type must be provided, which will be used to interpret different
 * GraphQL Value literals.
 *
 * | GraphQL Value        | JSON Value    |
 * | -------------------- | ------------- |
 * | Input Object         | Object        |
 * | List                 | Array         |
 * | Boolean              | Boolean       |
 * | String / Enum Value  | String        |
 * | Int / Float          | Number        |
 *
 */
func valueFromAST(valueAST ast.Value, ttype types.GraphQLInputType, variables map[string]interface{}) (interface{}, error) {
	//  TODO: Note: enforce that GraphQLType implements GraphQLInputType

	return valueAST, nil
}
