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
	Root          map[string]interface{}
	AST           ast.Document
	OperationName string
	Args          map[string]string
}

func Execute(p ExecuteParams, r chan types.GraphQLResult) {
	var errors []error
	params := BuildExecutionCtxParams{
		Schema:        p.Schema,
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

func executeOperation(p ExecuteOperationParams, r chan types.GraphQLResult) {
	var result types.GraphQLResult
	r <- result
}

type ExecutionContext struct {
	Schema    types.GraphQLSchema
	Fragments map[string]ast.FragmentDefinition
	Root      map[string]interface{}
	Operation ast.OperationDefinition
	Variables map[string]interface{}
	Errors    []error
}

type BuildExecutionCtxParams struct {
	Schema        types.GraphQLSchema
	Root          map[string]interface{}
	AST           ast.Document
	OperationName string
	Args          map[string]string
	Errors        []error
}

type ExecuteOperationParams struct {
	ExecutionContext ExecutionContext
	Root             map[string]interface{}
	Operation        ast.OperationDefinition
}

func buildExecutionContext(p BuildExecutionCtxParams) (eCtx ExecutionContext) {
	return eCtx
}
