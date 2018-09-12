package graphql

import (
	"github.com/GannettDigital/graphql/gqlerrors"
)

// type Schema interface{}

type Result struct {
	Data            interface{}                `json:"data"`
	Errors          []gqlerrors.FormattedError `json:"errors,omitempty"`
	QueryComplexity int                        `json:"query_complexity,omitempty"`
}

func (r *Result) HasErrors() bool {
	return len(r.Errors) > 0
}
