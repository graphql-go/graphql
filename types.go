package graphql

import (
	"github.com/graphql-go/graphql/gqlerrors"
)

type Result struct {
	Data       interface{}                `json:"data"`
	Errors     []gqlerrors.FormattedError `json:"errors,omitempty"`
	Extensions map[string]interface{}     `json:"extensions,omitempty"`
}

func (r *Result) HasErrors() bool {
	return len(r.Errors) > 0
}
