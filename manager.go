package graphql

import (
	"fmt"
	"sync"
	"time"

	"github.com/GannettDigital/graphql/gqlerrors"
	"github.com/GannettDigital/graphql/language/ast"
)

const (
	requestQueueBuffer = 10
	workerMaxRequests  = 5000
	workerStartDelay   = 500 * time.Microsecond
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

// resolveManager runs resolve functions and completeValue requests in a set of worker go routines.
// Having a set number of workers limits the churn of go routines while still providing parallel resolving of results
// which is key to performance when some resolving requires a network call.
//
// The nature of the GraphQL resolving code is that a single resolve call could end up calling other resolve calls
// as part of it. This means to avoid a full deadlock a hard limit on the number of workers can't be set, instead
// a slight buffer and a slight delay is added to give preference to reusing workers over creating new.
//
// A spike in traffic can lead to an increase in the number of workers many of which are subsequently idle. To
// allow slowly dropping back down to a reasonable pool workers will exit after processing a set number of requests.
// The original set of workers will always recreate a new worker on exit to avoid the number dropping too much.
type resolveManager struct {
	completeRequests chan completeRequest
	once             sync.Once
	quit             chan bool
	resolveRequests  chan resolveRequest
}

func newResolveManager(startingWorkers int) *resolveManager {
	manager := &resolveManager{
		completeRequests: make(chan completeRequest, requestQueueBuffer),
		once:             sync.Once{},
		quit:             make(chan bool),
		resolveRequests:  make(chan resolveRequest, requestQueueBuffer),
	}

	for i := 0; i < startingWorkers; i++ {
		go manager.newWorker(true)
	}

	return manager
}

func (manager *resolveManager) completeRequest(req completeRequest) {
	t := time.NewTimer(workerStartDelay)
	defer func() {
		if !t.Stop() {
			<-t.C
		}
	}()
	for {
		select {
		case manager.completeRequests <- req:
			return
		case <-t.C:
			go manager.newWorker(false)
			t.Reset(workerStartDelay)
		}
	}
}

func (manager *resolveManager) resolveRequest(name string, response chan<- resolverResponse, fn FieldResolveFn, params ResolveParams) {
	req := resolveRequest{
		fn:       fn,
		name:     name,
		params:   params,
		response: response,
	}

	t := time.NewTimer(workerStartDelay)
	defer func() {
		if !t.Stop() {
			<-t.C
		}
	}()
	for {
		select {
		case manager.resolveRequests <- req:
			return
		case <-t.C:
			go manager.newWorker(false)
			t.Reset(workerStartDelay)
		}
	}
}

func (manager *resolveManager) newWorker(revive bool) {
	// This is a simplistic means of dropping the number of workers back down eventually after a spike but other
	// methods such as added a timer caused a performance impact.
	for i := 0; i < workerMaxRequests; i++ {
		select {
		case <-manager.quit:
			return
		case req := <-manager.completeRequests:
			result := completeValueCatchingError(req.eCtx, req.returnType, req.fieldASTs, req.info, req.value)
			req.response <- completeResponse{index: req.index, result: result}
		case req := <-manager.resolveRequests:
			manager.resolve(req)
		}
	}
	if revive {
		go manager.newWorker(revive)
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

func (manager *resolveManager) stop() {
	manager.once.Do(func() { close(manager.quit) })
}
