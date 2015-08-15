package executor

import (
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/types"
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

func GetVariableValues(schema types.GraphQLSchema, definitionASTs []ast.VariableDefinition, inputs map[string]string) (r map[string]interface{}) {
	//TODO: use reduce
	return r
}
