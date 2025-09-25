package graphql

import (
	"context"
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"
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

// Subscribe performs a subscribe operation on the given query and schema
// To finish a subscription you can simply close the channel from inside the `Subscribe` function
// currently does not support extensions hooks
func Subscribe(p Params) chan *Result {
	return SubscribeWithPool(p, &SimpleResultPool{})
}

func SubscribeWithPool(p Params, resultPool ResultPool) chan *Result {
	source := source.NewSource(&source.Source{
		Body: []byte(p.RequestString),
		Name: "GraphQL request",
	})

	// TODO run extensions hooks

	// parse the source
	AST, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {

		// merge the errors from extensions and the original error from parser
		return sendOneResultAndClose(injectRequest(AST, &Result{
			Errors: gqlerrors.FormatErrors(err),
		}))
	}

	// validate document
	validationResult := ValidateDocument(&p.Schema, AST, nil)

	if !validationResult.IsValid {
		// run validation finish functions for extensions
		return sendOneResultAndClose(injectRequest(AST, &Result{
			Errors: validationResult.Errors,
		}))

	}
	return ExecuteSubscriptionWithPool(ExecuteParams{
		Schema:        p.Schema,
		Root:          p.RootObject,
		AST:           AST,
		OperationName: p.OperationName,
		Args:          p.VariableValues,
		Context:       p.Context,
	}, resultPool)
}

func sendOneResultAndClose(res *Result) chan *Result {
	resultChannel := make(chan *Result, 1)
	resultChannel <- res
	close(resultChannel)
	return resultChannel
}

func injectRequest(request *ast.Document, result *Result) *Result {
	if result != nil {
		result.Request = request
	}
	return result
}

// ExecuteSubscription is similar to graphql.Execute but returns a channel instead of a Result
// currently does not support extensions
func ExecuteSubscription(p ExecuteParams) chan *Result {
	return ExecuteSubscriptionWithPool(p, &SimpleResultPool{})
}

func ExecuteSubscriptionWithPool(p ExecuteParams, resultPool ResultPool) chan *Result {
	if p.Context == nil {
		p.Context = context.Background()
	}

	var mapSourceToResponse = func(payload interface{}) *Result {
		return injectRequest(p.AST, ExecuteWithPool(ExecuteParams{
			Schema:        p.Schema,
			Root:          payload,
			AST:           p.AST,
			OperationName: p.OperationName,
			Args:          p.Args,
			Context:       p.Context,
		}, resultPool))
	}
	var resultChannel = make(chan *Result)
	go func() {
		defer close(resultChannel)
		defer func() {
			if err := recover(); err != nil {
				e, ok := err.(error)
				if !ok {
					return
				}
				result := resultPool.Get()
				result.Errors = gqlerrors.FormatErrors(e)
				resultChannel <- injectRequest(p.AST, result)
			}
			return
		}()

		exeContext, err := buildExecutionContext(buildExecutionCtxParams{
			Schema:        p.Schema,
			Root:          p.Root,
			AST:           p.AST,
			OperationName: p.OperationName,
			Args:          p.Args,
			Context:       p.Context,
		})

		if err != nil {
			result := resultPool.Get()
			result.Errors = gqlerrors.FormatErrors(err)
			resultChannel <- injectRequest(p.AST, result)

			return
		}

		operationType, err := getOperationRootType(p.Schema, exeContext.Operation)
		if err != nil {
			result := resultPool.Get()
			result.Errors = gqlerrors.FormatErrors(err)
			resultChannel <- injectRequest(p.AST, result)

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
			result := resultPool.Get()
			result.Errors = gqlerrors.FormatErrors(fmt.Errorf("the subscription field %q is not defined", fieldName))
			resultChannel <- injectRequest(p.AST, result)

			return
		}

		resolveFn := fieldDef.Subscribe

		if resolveFn == nil {
			result := resultPool.Get()
			result.Errors = gqlerrors.FormatErrors(fmt.Errorf("the subscription function %q is not defined", fieldName))
			resultChannel <- injectRequest(p.AST, result)
			return
		}

		args := getArgumentValues(fieldDef.Args, fieldNode.Arguments, exeContext.VariableValues)
		info := ResolveInfo{
			FieldName:      fieldName,
			FieldASTs:      fieldNodes,
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
			result := resultPool.Get()
			result.Errors = gqlerrors.FormatErrors(err)
			resultChannel <- injectRequest(p.AST, result)

			return
		}

		if fieldResult == nil {
			result := resultPool.Get()
			result.Errors = gqlerrors.FormatErrors(fmt.Errorf("no field result"))
			resultChannel <- injectRequest(p.AST, result)

			return
		}

		switch sub := fieldResult.(type) {
		case chan interface{}:
			for {
				select {
				case <-p.Context.Done():
					return

				case res, more := <-sub:
					if !more {
						return
					}
					resultChannel <- mapSourceToResponse(res)
				}
			}
		case <-chan interface{}:
			for {
				select {
				case <-p.Context.Done():
					return

				case res, more := <-sub:
					if !more {
						return
					}
					resultChannel <- mapSourceToResponse(res)
				}
			}
		default:
			channel := reflect.ValueOf(sub)
			if channel.Kind() != reflect.Chan || (channel.Type().ChanDir()&reflect.RecvDir) == 0 {
				resultChannel <- mapSourceToResponse(fieldResult)
				return
			}
			cases := []reflect.SelectCase{{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(p.Context.Done()),
			}, {
				Dir:  reflect.SelectRecv,
				Chan: channel,
			}}
			for {
				chosen, value, ok := reflect.Select(cases)
				if chosen == 0 || !ok {
					return
				}
				if value.CanInterface() {
					resultChannel <- mapSourceToResponse(value.Interface())
				} else {
					resultChannel <- mapSourceToResponse(nil)
				}
			}
		}
	}()

	// return a result channel
	return resultChannel
}
