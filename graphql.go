package graphql

import (
	"context"

	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

type Params struct {
	// The GraphQL type system to use when validating and executing a query.
	Schema Schema

	// A GraphQL language formatted string representing the requested operation.
	RequestString string

	// The value provided as the first argument to resolver functions on the top
	// level type (e.g. the query object type).
	RootObject map[string]interface{}

	// A mapping of variable name to runtime value to use for all variables
	// defined in the requestString.
	VariableValues map[string]interface{}

	// The name of the operation to use if requestString contains multiple
	// possible operations. Can be omitted if requestString contains only
	// one operation.
	OperationName string

	// Context may be provided to pass application-specific per-request
	// information to resolve functions.
	Context context.Context
}

type ResultPool interface {
	Get() *Result
	Put(result *Result)

	GetListFor(result *Result, capacity int) []interface{}
	GetObjectFor(result *Result, capacity int) map[string]interface{}
}

// SimpleResultPool does nothing but let the GC take over
type SimpleResultPool struct{}

func (pool *SimpleResultPool) Get() *Result {
	return &Result{}
}

func (pool *SimpleResultPool) Put(*Result) {
	// does nothing, let the GC release it
}

func (pool *SimpleResultPool) GetListFor(result *Result, capacity int) []interface{} {
	return make([]interface{}, 0, capacity)
}

func (pool *SimpleResultPool) GetObjectFor(result *Result, capacity int) map[string]interface{} {
	return make(map[string]interface{}, capacity)
}

func Do(p Params) *Result {
	// by using SimpleResultPool here preserves the original interface and behavior
	// uses do not need to call Put on the returned result
	return DoWithPool(p, &SimpleResultPool{})
}

func DoWithPool(p Params, resultPool ResultPool) *Result {
	source := source.NewSource(&source.Source{
		Body: []byte(p.RequestString),
		Name: "GraphQL request",
	})

	// run init on the extensions
	extErrs := handleExtensionsInits(&p)
	if len(extErrs) != 0 {
		result := resultPool.Get()
		result.Errors = extErrs
		return result
	}

	extErrs, parseFinishFn := handleExtensionsParseDidStart(&p)
	if len(extErrs) != 0 {
		result := resultPool.Get()
		result.Errors = extErrs
		return result
	}

	// parse the source
	AST, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		// run parseFinishFuncs for extensions
		extErrs = parseFinishFn(err)

		// merge the errors from extensions and the original error from parser
		extErrs = append(extErrs, gqlerrors.FormatErrors(err)...)
		result := resultPool.Get()
		result.Errors = extErrs
		return result
	}

	// run parseFinish functions for extensions
	extErrs = parseFinishFn(err)
	if len(extErrs) != 0 {
		result := resultPool.Get()
		result.Request = AST
		result.Errors = extErrs
		return result
	}

	// notify extensions about the start of the validation
	extErrs, validationFinishFn := handleExtensionsValidationDidStart(&p)
	if len(extErrs) != 0 {
		result := resultPool.Get()
		result.Request = AST
		result.Errors = extErrs
		return result
	}

	// validate document
	validationResult := ValidateDocument(&p.Schema, AST, nil)

	if !validationResult.IsValid {
		// run validation finish functions for extensions
		extErrs = validationFinishFn(validationResult.Errors)

		// merge the errors from extensions and the original error from parser
		extErrs = append(extErrs, validationResult.Errors...)
		result := resultPool.Get()
		result.Request = AST
		result.Errors = extErrs
		return result
	}

	// run the validationFinishFuncs for extensions
	extErrs = validationFinishFn(validationResult.Errors)
	if len(extErrs) != 0 {
		result := resultPool.Get()
		result.Request = AST
		result.Errors = extErrs
		return result
	}

	result := ExecuteWithPool(ExecuteParams{
		Schema:        p.Schema,
		Root:          p.RootObject,
		AST:           AST,
		OperationName: p.OperationName,
		Args:          p.VariableValues,
		Context:       p.Context,
	}, resultPool)
	result.Request = AST
	return result
}
