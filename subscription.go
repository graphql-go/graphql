package graphql

import (
	"context"
	"fmt"
	"sync"

	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"
)

// ResultIteratorParams parameters passed to the result iterator handler
type ResultIteratorParams struct {
	ResultCount int64   // number of results this iterator has processed
	Result      *Result // the current result
	Done        func()  // Removes the current handler
	Cancel      func()  // Cancels the iterator, same as iterator.Cancel()
}

// ResultIteratorFn a result iterator handler
type ResultIteratorFn func(p ResultIteratorParams)

// holds subscription handler data
type subscriptionHanlderConfig struct {
	handler  ResultIteratorFn
	doneFunc func()
}

// ResultIterator handles processing results from a chan *Result
type ResultIterator struct {
	currentHandlerID int64
	count            int64
	wg               sync.WaitGroup
	ctx              context.Context
	ch               chan *Result
	cancelFunc       context.CancelFunc
	cancelled        bool
	handlers         map[int64]*subscriptionHanlderConfig
}

// NewResultIterator creates a new iterator and starts handling message on the result channel
func NewResultIterator(ctx context.Context, ch chan *Result) *ResultIterator {
	if ctx == nil {
		ctx = context.Background()
	}

	cctx, cancelFunc := context.WithCancel(ctx)
	iterator := &ResultIterator{
		currentHandlerID: 0,
		count:            0,
		ctx:              cctx,
		ch:               ch,
		cancelFunc:       cancelFunc,
		cancelled:        false,
		handlers:         map[int64]*subscriptionHanlderConfig{},
	}

	go func() {
		for {
			select {
			case <-iterator.ctx.Done():
				return
			case res := <-iterator.ch:
				if iterator.cancelled {
					return
				}
				iterator.wg.Wait()
				iterator.wg.Add(1)
				iterator.count++
				iterator.wg.Done()
				for _, h := range iterator.handlers {
					iterator.wg.Wait()
					h.handler(ResultIteratorParams{
						ResultCount: iterator.count,
						Result:      res,
						Done:        h.doneFunc,
						Cancel:      iterator.Cancel,
					})
				}
			}
		}
	}()

	return iterator
}

// adds a new handler
func (c *ResultIterator) addHandler(handler ResultIteratorFn) {
	c.wg.Add(1)
	handlerID := c.currentHandlerID + 1
	c.currentHandlerID = handlerID
	c.handlers[handlerID] = &subscriptionHanlderConfig{
		handler: handler,
		doneFunc: func() {
			c.removeHandler(handlerID)
		},
	}
	c.wg.Done()
}

// removes a handler and cancels if no more handlers exist
func (c *ResultIterator) removeHandler(handlerID int64) {
	c.wg.Add(1)
	delete(c.handlers, handlerID)
	if len(c.handlers) == 0 {
		c.Cancel()
	}
	c.wg.Done()
}

// ForEach adds a handler and handles each message as they come
func (c *ResultIterator) ForEach(handler ResultIteratorFn) {
	c.addHandler(handler)
}

// Cancel cancels the iterator
func (c *ResultIterator) Cancel() {
	c.cancelled = true
	c.cancelFunc()
}

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
func Subscribe(p SubscribeParams) *ResultIterator {
	resultChannel := make(chan *Result)
	// Use background context if no context was provided
	ctx := p.ContextValue
	if ctx == nil {
		ctx = context.Background()
	}

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
			}
			resultChannel <- result
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
			Context: exeContext.Context,
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
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("done context called")
					return
				case res := <-fieldResult.(chan interface{}):

					resultChannel <- mapSourceToResponse(res)
				}
			}
		default:
			resultChannel <- mapSourceToResponse(fieldResult)
			return
		}
	}()

	// return a result iterator
	return NewResultIterator(p.ContextValue, resultChannel)
}
