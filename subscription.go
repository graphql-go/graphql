package graphql

import (
	"context"
	"fmt"

	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"
)

// SubscribeParams parameters for subscribing
type SubscribeParams struct {
	Schema          Schema
	Document        *ast.Document
	RootValue       interface{}
	ContextValue    context.Context
	VariableValues  map[string]interface{}
	OperationName   string
	FieldResolver   FieldResolveFn
	FieldSubscriber FieldResolveFn
}

// Subscribe performs a subscribe operation
func Subscribe(ctx context.Context, p SubscribeParams) chan *Result {
	resultChannel := make(chan *Result)

	var mapSourceToResponse = func(payload interface{}) *Result {
		return Execute(ExecuteParams{
			Schema:        p.Schema,
			Root:          payload,
			AST:           p.Document,
			OperationName: p.OperationName,
			Args:          p.VariableValues,
			Context:       p.ContextValue,
		})
	}

	go func() {
		result := &Result{}
		defer func() {
			if err := recover(); err != nil {
				result.Errors = append(result.Errors, gqlerrors.FormatError(err.(error)))
				resultChannel <- result
			}
			close(resultChannel)
		}()

		exeContext, err := buildExecutionContext(buildExecutionCtxParams{
			Schema:        p.Schema,
			Root:          p.RootValue,
			AST:           p.Document,
			OperationName: p.OperationName,
			Args:          p.VariableValues,
			Result:        result,
			Context:       p.ContextValue,
		})

		if err != nil {
			result.Errors = append(result.Errors, gqlerrors.FormatError(err.(error)))
			resultChannel <- result
			return
		}

		operationType, err := getOperationRootType(p.Schema, exeContext.Operation)
		if err != nil {
			result.Errors = append(result.Errors, gqlerrors.FormatError(err.(error)))
			resultChannel <- result
			return
		}

		fields := collectFields(collectFieldsParams{
			ExeContext:   exeContext,
			RuntimeType:  operationType,
			SelectionSet: exeContext.Operation.GetSelectionSet(),
		})

		responseNames := []string{}
		for name := range fields {
			responseNames = append(responseNames, name)
		}
		responseName := responseNames[0]
		fieldNodes := fields[responseName]
		fieldNode := fieldNodes[0]
		fieldName := fieldNode.Name.Value
		fieldDef := getFieldDef(p.Schema, operationType, fieldName)

		if fieldDef == nil {
			err := fmt.Errorf("the subscription field %q is not defined", fieldName)
			result.Errors = append(result.Errors, gqlerrors.FormatError(err.(error)))
			resultChannel <- result
			return
		}

		resolveFn := p.FieldSubscriber
		if resolveFn == nil {
			resolveFn = DefaultResolveFn
		}
		if fieldDef.Subscribe != nil {
			resolveFn = fieldDef.Subscribe
		}
		fieldPath := &ResponsePath{
			Key: responseName,
		}

		args := getArgumentValues(fieldDef.Args, fieldNode.Arguments, exeContext.VariableValues)
		info := ResolveInfo{
			FieldName:      fieldName,
			FieldASTs:      fieldNodes,
			Path:           fieldPath,
			ReturnType:     fieldDef.Type,
			ParentType:     operationType,
			Schema:         p.Schema,
			Fragments:      exeContext.Fragments,
			RootValue:      exeContext.Root,
			Operation:      exeContext.Operation,
			VariableValues: exeContext.VariableValues,
		}

		fieldResult, err := resolveFn(ResolveParams{
			Source:  p.RootValue,
			Args:    args,
			Info:    info,
			Context: p.ContextValue,
		})
		if err != nil {
			result.Errors = append(result.Errors, gqlerrors.FormatError(err.(error)))
			resultChannel <- result
			return
		}

		if fieldResult == nil {
			err := fmt.Errorf("no field result")
			result.Errors = append(result.Errors, gqlerrors.FormatError(err.(error)))
			resultChannel <- result
			return
		}

		switch fieldResult.(type) {
		case chan interface{}:
			sub := fieldResult.(chan interface{})
			for {
				select {
				case <-ctx.Done():
					return

				case res := <-sub:
					resultChannel <- mapSourceToResponse(res)
				}
			}
		default:
			resultChannel <- mapSourceToResponse(fieldResult)
			return
		}
	}()

	// return a result channel
	return resultChannel
}
