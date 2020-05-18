package graphql

import (
	"context"
	"fmt"

	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

// SubscribeParams parameters for subscribing
type SubscribeParams struct {
	Schema        Schema
	RequestString string
	RootValue     interface{}
	// ContextValue    context.Context
	VariableValues  map[string]interface{}
	OperationName   string
	FieldResolver   FieldResolveFn
	FieldSubscriber FieldResolveFn
}

// SubscriptableSchema implements `graphql-transport-ws` `GraphQLService` interface: https://github.com/graph-gophers/graphql-transport-ws/blob/40c0484322990a129cac2f2d2763c3315230280c/graphqlws/internal/connection/connection.go#L53
type SubscriptableSchema struct {
	Schema     Schema
	RootObject map[string]interface{}
}

func (self *SubscriptableSchema) Subscribe(ctx context.Context, queryString string, operationName string, variables map[string]interface{}) (<-chan *Result, error) {
	c := Subscribe(Params{
		Schema:         self.Schema,
		Context:        ctx,
		OperationName:  operationName,
		RequestString:  queryString,
		RootObject:     self.RootObject,
		VariableValues: variables,
	})
	return c, nil
}

// Subscribe performs a subscribe operation
func Subscribe(p Params) chan *Result {

	source := source.NewSource(&source.Source{
		Body: []byte(p.RequestString),
		Name: "GraphQL request",
	})

	// TODO run extensions hooks

	// parse the source
	AST, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {

		// merge the errors from extensions and the original error from parser
		return sendOneResultAndClose(&Result{
			Errors: gqlerrors.FormatErrors(err),
		})
	}

	// validate document
	validationResult := ValidateDocument(&p.Schema, AST, nil)

	if !validationResult.IsValid {
		// run validation finish functions for extensions
		return sendOneResultAndClose(&Result{
			Errors: validationResult.Errors,
		})

	}
	return ExecuteSubscription(ExecuteParams{
		Schema:        p.Schema,
		Root:          p.RootObject,
		AST:           AST,
		OperationName: p.OperationName,
		Args:          p.VariableValues,
		Context:       p.Context,
	})
}

func sendOneResultAndClose(res *Result) chan *Result {
	resultChannel := make(chan *Result, 1) // TODO unbuffered channel does not pass errors, why?
	resultChannel <- res
	close(resultChannel)
	return resultChannel
}

func ExecuteSubscription(p ExecuteParams) chan *Result {

	if p.Context == nil {
		p.Context = context.Background()
	}

	// TODO run executionDidStart functions from extensions

	var mapSourceToResponse = func(payload interface{}) *Result {
		return Execute(ExecuteParams{
			Schema:        p.Schema,
			Root:          payload,
			AST:           p.AST,
			OperationName: p.OperationName,
			Args:          p.Args,
			Context:       p.Context,
		})
	}
	var resultChannel = make(chan *Result)
	go func() {
		defer close(resultChannel)
		defer func() {
			if err := recover(); err != nil {
				e, ok := err.(error)
				if !ok {
					fmt.Println("strange program path")
					return
				}
				resultChannel <- &Result{
					Errors: gqlerrors.FormatErrors(e),
				}
			}
			return
		}()

		exeContext, err := buildExecutionContext(buildExecutionCtxParams{
			Schema:        p.Schema,
			Root:          p.Root,
			AST:           p.AST,
			OperationName: p.OperationName,
			Args:          p.Args,
			Result:        &Result{}, // TODO what is this?
			Context:       p.Context,
		})

		if err != nil {
			resultChannel <- &Result{
				Errors: gqlerrors.FormatErrors(err),
			}

			return
		}

		operationType, err := getOperationRootType(p.Schema, exeContext.Operation)
		if err != nil {
			resultChannel <- &Result{
				Errors: gqlerrors.FormatErrors(err),
			}

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
			resultChannel <- &Result{
				Errors: gqlerrors.FormatErrors(fmt.Errorf("the subscription field %q is not defined", fieldName)),
			}

			return
		}

		resolveFn := fieldDef.Subscribe

		if resolveFn == nil {
			resultChannel <- &Result{
				Errors: gqlerrors.FormatErrors(fmt.Errorf("the subscription function %q is not defined", fieldName)),
			}

			return
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
			Source:  p.Root,
			Args:    args,
			Info:    info,
			Context: p.Context,
		})
		if err != nil {
			resultChannel <- &Result{
				Errors: gqlerrors.FormatErrors(err),
			}

			return
		}

		if fieldResult == nil {
			resultChannel <- &Result{
				Errors: gqlerrors.FormatErrors(fmt.Errorf("no field result")),
			}

			return
		}

		switch fieldResult.(type) {
		case chan interface{}:
			sub := fieldResult.(chan interface{})
			for {
				select {
				case <-p.Context.Done():
					println("context cancelled")
					// TODO send the context error to the resultchannel?
					return

				case res, more := <-sub:
					if !more {

						return
					}
					resultChannel <- mapSourceToResponse(res)
				}
			}
		default:
			fmt.Println(fieldResult)
			resultChannel <- mapSourceToResponse(fieldResult)
			return
		}
	}()

	// return a result channel
	return resultChannel
}
