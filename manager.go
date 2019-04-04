package graphql

import (
	"fmt"

	"github.com/GannettDigital/graphql/gqlerrors"
	"github.com/GannettDigital/graphql/language/ast"
)

const (
	requestQueueBuffer = 50 // this also defines the number of permanent workers which is double this number
)

// completeRequest contains the information needed to complete a field.
// A completeRequest is passed to a resolveManager worker which processes the request.
// This is only used when completing multiple items concurrently such as needed for a list
type completeRequest struct {
	index    int
	response chan<- completeResponse

	eCtx       *executionContext
	returnType Type
	fieldASTs  []*ast.Field
	info       ResolveInfo
	value      interface{}
}

// completeResponse is the type containing the completion response which is sent over a channel as workers finish
// processing.
type completeResponse struct {
	index  int
	result interface{}
}

// resolveRequest contains the information needed to resolve a field.
// A resolveRequest is passed to a resolveManager worker which processes the request.
type resolveRequest struct {
	fn       FieldResolveFn
	name     string
	params   ResolveParams
	response chan<- resolverResponse
}

// resolverResponse is the type containing the resolver response which is sent over a channel as workers finish
// processing.
type resolverResponse struct {
	err    error
	name   string
	result interface{}
}

// resolveManager runs resolve functions and completeValue requests with a set of worker go routines.
// Having a set of workers limits the churn of go routines while still providing parallel resolving of results
// which is key to performance when some resolving requires a network call.
//
// The nature of the GraphQL resolving code is that a single resolve call could end up calling other resolve calls
// as part of it. This means to avoid a full deadlock a hard limit on the number of workers can't be set, instead
// a slight buffer and a slight delay is added to give preference to reusing workers over creating new.
//
// A small set of workers are long lived to allow some processing to happen at all times and enable a buffered channel
// Most other workers are not long lived but will only exit when the request channels are empty.
type resolveManager struct {
	completeRequests chan completeRequest
	resolveRequests  chan resolveRequest
}

func newResolveManager() *resolveManager {
	manager := &resolveManager{
		completeRequests: make(chan completeRequest, requestQueueBuffer),
		resolveRequests:  make(chan resolveRequest, requestQueueBuffer),
	}

	for i := 0; i < 2*requestQueueBuffer; i++ {
		go manager.infiniteWorker()
	}
	return manager
}

func (manager *resolveManager) completeRequest(req completeRequest) {
	select {
	case manager.completeRequests <- req:
		return
	default:
		go manager.newWorker()
		manager.completeRequests <- req
	}
}

func (manager *resolveManager) infiniteWorker() {
	for {
		select {
		case req := <-manager.completeRequests:
			result := completeValueCatchingError(req.eCtx, req.returnType, req.fieldASTs, req.info, req.value)
			req.response <- completeResponse{index: req.index, result: result}
		case req := <-manager.resolveRequests:
			manager.resolve(req)
		}
	}
}

func (manager *resolveManager) newWorker() {
	for {
		select {
		case req := <-manager.completeRequests:
			result := completeValueCatchingError(req.eCtx, req.returnType, req.fieldASTs, req.info, req.value)
			req.response <- completeResponse{index: req.index, result: result}
		case req := <-manager.resolveRequests:
			manager.resolve(req)
		default:
			return
		}
	}
}

func (manager *resolveManager) resolve(req resolveRequest) {
	defer func() {
		if r := recover(); r != nil {
			var err error
			if r, ok := r.(string); ok {
				err = NewLocatedError(
					fmt.Sprintf("%v", r),
					FieldASTsToNodeASTs(req.params.Info.FieldASTs),
				)
			}
			if r, ok := r.(error); ok {
				err = gqlerrors.FormatError(r)
			}
			req.response <- resolverResponse{name: req.name, err: err}
		}
	}()

	result, err := req.fn(req.params)
	req.response <- resolverResponse{name: req.name, result: result, err: err}
}

func (manager *resolveManager) resolveRequest(name string, response chan<- resolverResponse, fn FieldResolveFn, params ResolveParams) {
	req := resolveRequest{
		fn:       fn,
		name:     name,
		params:   params,
		response: response,
	}

	select {
	case manager.resolveRequests <- req:
		return
	default:
		go manager.newWorker()
		manager.resolveRequests <- req
	}
}
