package graphql

import (
	"context"
	"sync"

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

type listPool struct {
	bin8  sync.Pool
	bin16 sync.Pool
	bin32 sync.Pool
}

func NewListPool() *listPool {
	return &listPool{
		bin8: sync.Pool{
			New: func() interface{} {
				return make([]interface{}, 0, 8) // Capacity up to 8
			},
		},
		bin16: sync.Pool{
			New: func() interface{} {
				return make([]interface{}, 0, 16) // Capacity up to 16
			},
		},
		bin32: sync.Pool{
			New: func() interface{} {
				return make([]interface{}, 0, 32) // Capacity up to 32
			},
		},
	}
}

func (pool *listPool) Get(capacity int) []interface{} {
	var list []interface{}
	switch {
	case capacity <= 8:
		list = pool.bin8.Get().([]interface{})
	case capacity <= 16:
		list = pool.bin16.Get().([]interface{})
	case capacity <= 32:
		list = pool.bin32.Get().([]interface{})
	default:
		// if no suitable pool, create a new slice
		list = make([]interface{}, 0, capacity)
	}
	return list
}

func (pool *listPool) Put(list []interface{}) {
	list = list[:0] // reset length to 0
	switch {
	case cap(list) <= 8:
		pool.bin8.Put(list)
	case cap(list) <= 16:
		pool.bin16.Put(list)
	case cap(list) <= 32:
		pool.bin32.Put(list)
	default:
		// for very large slices, we choose not to pool them
	}
}

type objectPool struct {
	bin8  sync.Pool
	bin16 sync.Pool
	bin32 sync.Pool
}

func NewObjectPool() *objectPool {
	return &objectPool{
		bin8: sync.Pool{
			New: func() interface{} {
				return make(map[string]interface{}, 8) // Capacity up to 8
			},
		},
		bin16: sync.Pool{
			New: func() interface{} {
				return make(map[string]interface{}, 16) // Capacity up to 16
			},
		},
		bin32: sync.Pool{
			New: func() interface{} {
				return make(map[string]interface{}, 32) // Capacity up to 32
			},
		},
	}
}

func (pool *objectPool) Get(capacity int) map[string]interface{} {
	var object map[string]interface{}
	switch {
	case capacity <= 8:
		object = pool.bin8.Get().(map[string]interface{})
	case capacity <= 16:
		object = pool.bin16.Get().(map[string]interface{})
	case capacity <= 32:
		object = pool.bin32.Get().(map[string]interface{})
	default:
		// if no suitable pool, create a new slice
		object = make(map[string]interface{}, capacity)
	}
	return object
}

func (pool *objectPool) Put(object map[string]interface{}) {
	switch {
	case len(object) <= 8:
		pool.bin8.Put(object)
	case len(object) <= 16:
		pool.bin16.Put(object)
	case len(object) <= 32:
		pool.bin32.Put(object)
	default:
		// for very large slices, we choose not to pool them
		return
	}
	// reset length to 0
	for key := range object {
		delete(object, key)
	}
}

type ResultPool interface {
	Get() *Result
	Put(result *Result)

	getListFor(result *Result, capacity int) []interface{}
	getObjectFor(result *Result, capacity int) map[string]interface{}
}

type SimpleResultPool struct{}

func (pool *SimpleResultPool) Get() *Result {
	return &Result{}
}

func (pool *SimpleResultPool) Put(*Result) {
}

func (pool *SimpleResultPool) getListFor(result *Result, capacity int) []interface{} {
	return make([]interface{}, 0, capacity)
}

func (pool *SimpleResultPool) getObjectFor(result *Result, capacity int) map[string]interface{} {
	return make(map[string]interface{}, capacity)
}

type resultPool struct {
	pool       sync.Pool
	listPool   *listPool
	objectPool *objectPool
}

func NewResultPool() ResultPool {
	return &resultPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &Result{}
			},
		},
		listPool:   NewListPool(),
		objectPool: NewObjectPool(),
	}
}

func (pool *resultPool) Get() *Result {
	return pool.pool.Get().(*Result)
}

func (pool *resultPool) Put(result *Result) {
	for _, list := range result.dataLists {
		pool.listPool.Put(list)
	}
	for _, object := range result.dataObjects {
		pool.objectPool.Put(object)
	}
	*result = Result{} // reset
	pool.pool.Put(result)
}

func (pool *resultPool) getListFor(result *Result, capacity int) []interface{} {
	list := pool.listPool.Get(capacity)
	result.dataLists = append(result.dataLists, list)
	return list
}

func (pool *resultPool) getObjectFor(result *Result, capacity int) map[string]interface{} {
	object := pool.objectPool.Get(capacity)
	result.dataObjects = append(result.dataObjects, object)
	return object
}

func Do(p Params) *Result {
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
		return &Result{
			Errors: extErrs,
		}
	}

	extErrs, parseFinishFn := handleExtensionsParseDidStart(&p)
	if len(extErrs) != 0 {
		return &Result{
			Errors: extErrs,
		}
	}

	// parse the source
	AST, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		// run parseFinishFuncs for extensions
		extErrs = parseFinishFn(err)

		// merge the errors from extensions and the original error from parser
		extErrs = append(extErrs, gqlerrors.FormatErrors(err)...)
		return &Result{
			Errors: extErrs,
		}
	}

	// run parseFinish functions for extensions
	extErrs = parseFinishFn(err)
	if len(extErrs) != 0 {
		return &Result{
			Request: AST,
			Errors:  extErrs,
		}
	}

	// notify extensions about the start of the validation
	extErrs, validationFinishFn := handleExtensionsValidationDidStart(&p)
	if len(extErrs) != 0 {
		return &Result{
			Request: AST,
			Errors:  extErrs,
		}
	}

	// validate document
	validationResult := ValidateDocument(&p.Schema, AST, nil)

	if !validationResult.IsValid {
		// run validation finish functions for extensions
		extErrs = validationFinishFn(validationResult.Errors)

		// merge the errors from extensions and the original error from parser
		extErrs = append(extErrs, validationResult.Errors...)
		return &Result{
			Request: AST,
			Errors:  extErrs,
		}
	}

	// run the validationFinishFuncs for extensions
	extErrs = validationFinishFn(validationResult.Errors)
	if len(extErrs) != 0 {
		return &Result{
			Request: AST,
			Errors:  extErrs,
		}
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
