package executor

import (
	"fmt"
	"log"

	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/language/kinds"
	"github.com/chris-ramon/graphql-go/types"
)

type ExecuteParams struct {
	Schema        types.GraphQLSchema
	Root          map[string]interface{}
	AST           ast.Document
	OperationName string
	Args          map[string]string
}

func Execute(p ExecuteParams, resultChan chan *types.GraphQLResult) {
	var errors []error
	var result types.GraphQLResult
	params := BuildExecutionCtxParams{
		Schema:        p.Schema,
		Root:          p.Root,
		AST:           p.AST,
		OperationName: p.OperationName,
		Args:          p.Args,
		Errors:        errors,
		Result:        &result,
		ResultChan:    resultChan,
	}
	exeContext := buildExecutionContext(params)
	if result.HasErrors() {
		return
	}
	eOperationParams := ExecuteOperationParams{
		ExecutionContext: exeContext,
		Root:             p.Root,
		Operation:        exeContext.Operation,
	}
	executeOperation(eOperationParams, resultChan)
}

type ExecuteOperationParams struct {
	ExecutionContext ExecutionContext
	Root             map[string]interface{}
	Operation        ast.Definition
}

func executeOperation(p ExecuteOperationParams, r chan *types.GraphQLResult) {
	//TODO: mutation operation
	operationType := getOperationRootType(p.ExecutionContext.Schema, p.Operation, r)
	collectFieldsParams := CollectFieldsParams{
		ExeContext:    p.ExecutionContext,
		OperationType: operationType,
		SelectionSet:  p.Operation.GetSelectionSet(),
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

func getOperationRootType(schema types.GraphQLSchema, operation ast.Definition, r chan *types.GraphQLResult) (objType types.GraphQLObjectType) {
	switch operation.GetOperation() {
	case "query":
		return schema.GetQueryType()
	case "mutation":
		mutationType := schema.GetMutationType()
		if mutationType.Name != "" {
			var result types.GraphQLResult
			err := graphqlerrors.NewGraphQLFormattedError("Schema is not configured for mutations")
			result.Errors = append(result.Errors, err)
			r <- &result
			return objType
		}
		return mutationType
	default:
		var result types.GraphQLResult
		err := graphqlerrors.NewGraphQLFormattedError("Can only execute queries and mutations")
		result.Errors = append(result.Errors, err)
		r <- &result
		return objType
	}
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

func executeFields(p ExecuteFieldsParams, resultChan chan *types.GraphQLResult) {
	var result types.GraphQLResult
	//mutable := reflect.ValueOf(p.ExecutionContext.Result.Data).Elem()
	//mutable.FieldByName("Name").SetString("R2-D2")
	resultChan <- &result
}

type BuildExecutionCtxParams struct {
	Schema        types.GraphQLSchema
	Root          map[string]interface{}
	AST           ast.Document
	OperationName string
	Args          map[string]string
	Errors        []error
	Result        *types.GraphQLResult
	ResultChan    chan *types.GraphQLResult
}

type ExecutionContext struct {
	Schema    types.GraphQLSchema
	Fragments map[string]ast.Definition
	Root      map[string]interface{}
	Operation ast.Definition
	Variables map[string]interface{}
	Errors    []error
}

func buildExecutionContext(p BuildExecutionCtxParams) (eCtx ExecutionContext) {
	operations := make(map[string]ast.Definition)
	fragments := make(map[string]ast.Definition)
	for _, statement := range p.AST.Definitions {
		switch statement.GetKind() {
		case kinds.OperationDefinition:
			log.Println("kinds.OperationDefinition")
			key := ""
			if statement.GetName().Value != "" {
				key = statement.GetName().Value
			}
			operations[key] = statement
			break
		case kinds.FragmentDefinition:
			log.Println("kinds.FragmentDefinition")
			fragments[statement.GetName().Value] = statement
			break
		default:
			log.Println("default")
		}
	}
	log.Printf("debug - operations: %v", operations)
	if (p.OperationName == "") && (len(operations) != 1) {
		err := graphqlerrors.NewGraphQLFormattedError("Must provide operation name if query contains multiple operations")
		p.Result.Errors = append(p.Result.Errors, err)
		p.ResultChan <- p.Result
		return eCtx
	}
	var opName string
	if p.OperationName == "" {
		for k, _ := range operations {
			opName = k
			break
		}
	}
	operation, found := operations[opName]
	if !found {
		var result types.GraphQLResult
		err := graphqlerrors.NewGraphQLFormattedError(fmt.Sprintf("Unknown operation name: %s", opName))
		result.Errors = append(result.Errors, err)
		return eCtx
	}
	variables := GetVariableValues(p.Schema, operation.GetVariableDefinitions(), p.Args)
	eCtx.Schema = p.Schema
	eCtx.Fragments = fragments
	eCtx.Root = p.Root
	eCtx.Operation = operation
	eCtx.Variables = variables
	eCtx.Errors = p.Errors
	return eCtx
}
