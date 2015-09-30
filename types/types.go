package types

import (
	"github.com/chris-ramon/graphql-go/errors"
)

type Schema interface{}

type GraphQLResult struct {
	Data   interface{}                           `json:"data"`
	Errors []graphqlerrors.GraphQLFormattedError `json:"errors"`
}

func (gqR *GraphQLResult) HasErrors() bool {
	return (len(gqR.Errors) > 0)
}
