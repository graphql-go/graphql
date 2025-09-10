package graphql

import (
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"
)

// type Schema interface{}

// Result has the response, errors and extensions from the resolved schema
type Result struct {
	Request    *ast.Document              `json:"-" yaml:"-" msgpack:"-"`
	Data       interface{}                `json:"data" yaml:"data" msgpack:"data"`
	Errors     []gqlerrors.FormattedError `json:"errors,omitempty" yaml:"errors,omitempty" msgpack:"errors,omitempty"`
	Extensions map[string]interface{}     `json:"extensions,omitempty" yaml:"extensions,omitempty" msgpack:"extensions,omitempty"`
}

// HasErrors just a simple function to help you decide if the result has errors or not
func (r *Result) HasErrors() bool {
	return len(r.Errors) > 0
}
