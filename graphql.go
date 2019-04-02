package graphql

import (
	"context"
	"fmt"
	"sync"

	"github.com/GannettDigital/graphql/gqlerrors"
	"github.com/GannettDigital/graphql/language/parser"
	"github.com/GannettDigital/graphql/language/source"
)

// TODO I should build a defaultManager and allow for one passed in as part of the params
var (
	manager     *resolveManager
	managerInit sync.Once
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

	// The maximum complexity cost of a query, any query exceeding this will error
	MaxCost int

	// Context may be provided to pass application-specific per-request
	// information to resolve functions.
	Context context.Context
}

func Do(p Params) *Result {
	if manager == nil {
		managerInit.Do(func() {
			manager = newResolveManager()
		})
	}

	source := source.NewSource(&source.Source{
		Body: []byte(p.RequestString),
		Name: "GraphQL request",
	})
	AST, err := parser.Parse(parser.ParseParams{Source: source})
	if err != nil {
		return &Result{
			Errors: gqlerrors.FormatErrors(err),
		}
	}
	validationResult := ValidateDocument(&p.Schema, AST, nil)

	if !validationResult.IsValid {
		return &Result{
			Errors: validationResult.Errors,
		}
	}

	ep := ExecuteParams{
		Schema:        p.Schema,
		Root:          p.RootObject,
		AST:           AST,
		OperationName: p.OperationName,
		Args:          p.VariableValues,
		Context:       p.Context,
		manager:       manager,
	}

	var cost int
	if p.MaxCost > 0 {
		cost, err = QueryComplexity(ep)
		if err != nil {
			return &Result{Errors: []gqlerrors.FormattedError{gqlerrors.FormatError(err)}}
		}
		if cost > p.MaxCost {
			return &Result{
				Errors: []gqlerrors.FormattedError{
					gqlerrors.NewFormattedError(fmt.Sprintf("maximum complexity cost %d exceeded, query cost %d", p.MaxCost, cost)),
				},
			}
		}
	}

	result := Execute(ep)
	result.QueryComplexity = cost

	return result
}
