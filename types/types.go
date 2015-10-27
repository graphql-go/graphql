package types

import (
	"github.com/chris-ramon/graphql/errors"
)

// type Schema interface{}

type Result struct {
	Data   interface{}                           `json:"data"`
	Errors []graphqlerrors.FormattedError `json:"errors,omitempty"`
}

func (gqR *Result) HasErrors() bool {
	return (len(gqR.Errors) > 0)
}
