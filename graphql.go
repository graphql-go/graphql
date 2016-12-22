package graphql

import (
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"golang.org/x/net/context"
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

	// Context may be provided to pass application-specific per-request
	// information to resolve functions.
	Context context.Context

	// Executor allows to control the behavior of how to perform resolving function that
	// can be run concurrently. If not given, they will be executed serially.
	Executor Executor

	// If true, introspection queries are blocked.
	BlockIntrospection bool
}

// Parse parses reuqest string to an AST. It does not validate the AST.
func Parse(requestString string) (*ast.Document, error) {
	source := source.NewSource(&source.Source{
		Body: []byte(requestString),
		Name: "GraphQL request",
	})
	return parser.Parse(parser.ParseParams{Source: source})
}

// Validate validates the AST. If it's not valid, an error result
// is returned, otherwise, it returns nil.
func Validate(AST *ast.Document, schema *Schema) *Result {
	validationResult := ValidateDocument(schema, AST, nil)
	if !validationResult.IsValid {
		return &Result{
			Errors: validationResult.Errors,
		}
	}
	return nil
}

// DoWithAST execute GraphQL request with a parsed and valid AST.
func DoWithAST(p Params, AST *ast.Document) *Result {
	return Execute(ExecuteParams{
		Schema:        p.Schema,
		Root:          p.RootObject,
		AST:           AST,
		OperationName: p.OperationName,
		Args:          p.VariableValues,
		Context:       p.Context,
		Executor:      p.Executor,
		BlockMeta:     p.BlockIntrospection,
	})
}

func Do(p Params) *Result {
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

	return Execute(ExecuteParams{
		Schema:        p.Schema,
		Root:          p.RootObject,
		AST:           AST,
		OperationName: p.OperationName,
		Args:          p.VariableValues,
		Context:       p.Context,
		Executor:      p.Executor,
		BlockMeta:     p.BlockIntrospection,
	})
}
