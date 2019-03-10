package graphql

import (
	"sync"

	"github.com/graphql-go/graphql/gqlerrors"
)

// type Schema interface{}

type Result struct {
	Data   interface{}                `json:"data"`
	Errors []gqlerrors.FormattedError `json:"errors,omitempty"`

	errorsLock sync.Mutex
}

func (r *Result) HasErrors() bool {
	return len(r.Errors) > 0
}

// AppendErrors is the thread-safe way to append error(s) to Result.Errors.
func (r *Result) AppendErrors(errs ...gqlerrors.FormattedError) {
	r.errorsLock.Lock()
	defer r.errorsLock.Unlock()

	r.Errors = append(r.Errors, errs...)
}
