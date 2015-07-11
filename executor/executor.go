package executor

import (
	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/types"
)

type ExecutionResult struct {
	Data   interface{}
	Errors []errors.GraphQLFormattedError
}

type ExecuteParams struct {
	Schema        types.GraphQLSchema
	Result        interface{}
	Root          map[string]interface{}
	AST           ast.Document
	OperationName string
	Args          map[string]string
}

func Execute(p ExecuteParams, r chan types.GraphQLResult) {
	var errors []error
	params := BuildExecutionCtxParams{
		Schema:        p.Schema,
		Result:        p.Result,
		Root:          p.Root,
		AST:           p.AST,
		OperationName: p.OperationName,
		Args:          p.Args,
		Errors:        errors,
	}
	exeContext := buildExecutionContext(params)
	eOperationParams := ExecuteOperationParams{
		ExecutionContext: exeContext,
		Root:             p.Root,
		Operation:        exeContext.Operation,
	}
	executeOperation(eOperationParams, r)
}

type ExecutionContext struct {
	Schema    types.GraphQLSchema
	Result    interface{}
	Fragments map[string]ast.FragmentDefinition
	Root      map[string]interface{}
	Operation ast.OperationDefinition
	Variables map[string]interface{}
	Errors    []error
}

type ExecuteOperationParams struct {
	ExecutionContext ExecutionContext
	Root             map[string]interface{}
	Operation        ast.OperationDefinition
}

func executeOperation(p ExecuteOperationParams, r chan types.GraphQLResult) {
	//TODO: mutation operation
	var result types.GraphQLResult
	//mutable := reflect.ValueOf(p.ExecutionContext.Result).Elem()
	//mutable.FieldByName("Name").SetString("R2-D2")
	operationType := getOperationRootType(p.ExecutionContext.Schema, p.Operation)
	collectFieldsParams := CollectFieldsParams{
		ExeContext:    p.ExecutionContext,
		OperationType: operationType,
		SelectionSet:  p.Operation.SelectionSet,
	}
	fields := collectFields(collectFieldsParams)
	executeFieldsParams := ExecuteFieldsParams{
		ExeContext: p.ExecutionContext,
		ParentType: operationType,
		Source:     p.Root,
		Fields:     fields,
	}
	executeFields(executeFieldsParams)
	r <- result
}

func getOperationRootType(schema types.GraphQLSchema, operation ast.OperationDefinition) (r types.GraphQLObjectType) {
	return r
}

type CollectFieldsParams struct {
	ExeContext           ExecutionContext
	OperationType        types.GraphQLObjectType
	SelectionSet         ast.SelectionSet
	Fields               map[string][]ast.Field
	VisitedFragmentNames map[string]bool
}

func collectFields(p CollectFieldsParams) (r map[string][]ast.Field) {
	return r
}

type ExecuteFieldsParams struct {
	ExeContext ExecutionContext
	ParentType types.GraphQLObjectType
	Source     map[string]interface{}
	Fields     map[string][]ast.Field
}

func executeFields(p ExecuteFieldsParams) (r map[string]interface{}) {
	return r
}

type BuildExecutionCtxParams struct {
	Schema        types.GraphQLSchema
	Result        interface{}
	Root          map[string]interface{}
	AST           ast.Document
	OperationName string
	Args          map[string]string
	Errors        []error
}

func buildExecutionContext(p BuildExecutionCtxParams) (eCtx ExecutionContext) {
	eCtx.Schema = p.Schema
	eCtx.Result = p.Result
	return eCtx
}
