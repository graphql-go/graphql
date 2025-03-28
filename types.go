package graphql

import (
	"fmt"

	"github.com/graphql-go/graphql/gqlerrors"
)

// type Schema interface{}

// Result has the response, errors and extensions from the resolved schema
type Result struct {
	Data       interface{}                `json:"data"`
	Errors     []gqlerrors.FormattedError `json:"errors,omitempty"`
	Extensions map[string]interface{}     `json:"extensions,omitempty"`
}

// HasErrors just a simple function to help you decide if the result has errors or not
func (r *Result) HasErrors() bool {
	return len(r.Errors) > 0
}

// ErrorsJoined joins and returns the result errors with `:` as character separator.
func (r *Result) ErrorsJoined() error {
	if r.Errors == nil {
		return nil
	}

	var result error
	for _, err := range r.Errors {
		if result == nil {
			result = fmt.Errorf("%w", err)

			continue
		}

		result = fmt.Errorf("%w: %w", err, result)
	}

	return result
}
