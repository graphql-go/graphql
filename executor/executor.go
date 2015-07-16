package executor

import (
	"reflect"

	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/kinds"
	"github.com/chris-ramon/graphql-go/language/od"
	"github.com/chris-ramon/graphql-go/types"
)

type ExecutionResult struct {
	Data   interface{}
	Errors []errors.GraphQLFormattedError
}

type ExecuteParams struct {
	Schema        types.GraphQLSchema
	Result        types.GraphQLResult
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
	Result    types.GraphQLResult
	Fragments map[string]ast.FragmentDefinition
	Root      map[string]interface{}
	Operation od.OperationDefinition
	Variables map[string]interface{}
	Errors    []error
}

type ExecuteOperationParams struct {
	ExecutionContext ExecutionContext
	Root             map[string]interface{}
	Operation        od.OperationDefinition
}

func executeOperation(p ExecuteOperationParams, r chan types.GraphQLResult) {
	//TODO: mutation operation
	operationType := getOperationRootType(p.ExecutionContext.Schema, p.Operation)
	collectFieldsParams := CollectFieldsParams{
		ExeContext:    p.ExecutionContext,
		OperationType: operationType,
		SelectionSet:  p.Operation.SelectionSet,
	}
	fields := collectFields(collectFieldsParams)
	executeFieldsParams := ExecuteFieldsParams{
		ExecutionContext: p.ExecutionContext,
		ParentType:       operationType,
		Source:           p.Root,
		Fields:           fields,
	}
	executeFields(executeFieldsParams, r)
}

func getOperationRootType(schema types.GraphQLSchema, operation od.OperationDefinition) (r types.GraphQLObjectType) {
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
	ExecutionContext ExecutionContext
	ParentType       types.GraphQLObjectType
	Source           map[string]interface{}
	Fields           map[string][]ast.Field
}

func executeFields(p ExecuteFieldsParams, r chan types.GraphQLResult) {
	var result types.GraphQLResult
	mutable := reflect.ValueOf(p.ExecutionContext.Result.Data).Elem()
	mutable.FieldByName("Name").SetString("R2-D2")
	r <- result
}

type BuildExecutionCtxParams struct {
	Schema        types.GraphQLSchema
	Result        types.GraphQLResult
	Root          map[string]interface{}
	AST           ast.Document
	OperationName string
	Args          map[string]string
	Errors        []error
}

func buildExecutionContext(p BuildExecutionCtxParams) (eCtx ExecutionContext) {
	eCtx.Schema = p.Schema
	eCtx.Result = p.Result
	var operations map[string]ast.Definition
	var fragments map[string]ast.Definition
	for _, statement := range p.AST.Definitions {
		switch statement.GetKind() {
		case kinds.OperationDefinition:
			key := ""
			if statement.GetName().Value != "" {
				key = statement.GetName().Value
			}
			operations[key] = statement
			break
		case kinds.FragmentDefinition:
			fragments[statement.GetName().Value] = statement
			break
		}
	}
	if (p.OperationName == "") && (len(operations) != 1) {
		p.Result.Errors = append(p.Result.Errors, errors.GraphQLFormattedError{
			Message: "Must provide operation name if query contains multiple operations",
		})
		return eCtx
	}
	return eCtx
}
