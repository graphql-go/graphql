package executor_test

import (
	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/language/location"
	"github.com/chris-ramon/graphql-go/testutil"
	"github.com/chris-ramon/graphql-go/types"
	"golang.org/x/net/context"
	"reflect"
	"testing"
)

// testNumberHolder maps to numberHolderType
type testNumberHolder struct {
	TheNumber int `json:"theNumber"` // map field to `theNumber` so it can be resolve by the default ResolveFn
}
type testRoot struct {
	NumberHolder *testNumberHolder
}

func newTestRoot(originalNumber int) *testRoot {
	return &testRoot{
		NumberHolder: &testNumberHolder{originalNumber},
	}
}
func (r *testRoot) ImmediatelyChangeTheNumber(newNumber int) *testNumberHolder {
	r.NumberHolder.TheNumber = newNumber
	return r.NumberHolder
}
func (r *testRoot) PromiseToChangeTheNumber(newNumber int) *testNumberHolder {
	return r.ImmediatelyChangeTheNumber(newNumber)
}
func (r *testRoot) FailToChangeTheNumber(newNumber int) *testNumberHolder {
	panic("Cannot change the number")
}
func (r *testRoot) PromiseAndFailToChangeTheNumber(newNumber int) *testNumberHolder {
	panic("Cannot change the number")
}

// numberHolderType creates a mapping to testNumberHolder
var numberHolderType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "NumberHolder",
	Fields: types.GraphQLFieldConfigMap{
		"theNumber": &types.GraphQLFieldConfig{
			Type: types.GraphQLInt,
		},
	},
})

var mutationsTestSchema, _ = types.NewGraphQLSchema(types.GraphQLSchemaConfig{
	Query: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Query",
		Fields: types.GraphQLFieldConfigMap{
			"numberHolder": &types.GraphQLFieldConfig{
				Type: numberHolderType,
			},
		},
	}),
	Mutation: types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Mutation",
		Fields: types.GraphQLFieldConfigMap{
			"immediatelyChangeTheNumber": &types.GraphQLFieldConfig{
				Type: numberHolderType,
				Args: types.GraphQLFieldConfigArgumentMap{
					"newNumber": &types.GraphQLArgumentConfig{
						Type: types.GraphQLInt,
					},
				},
				Resolve: func(ctx context.Context, p types.GQLFRParams) interface{} {
					newNumber := 0
					obj, _ := p.Source.(*testRoot)
					newNumber, _ = p.Args["newNumber"].(int)
					return obj.ImmediatelyChangeTheNumber(newNumber)
				},
			},
			"promiseToChangeTheNumber": &types.GraphQLFieldConfig{
				Type: numberHolderType,
				Args: types.GraphQLFieldConfigArgumentMap{
					"newNumber": &types.GraphQLArgumentConfig{
						Type: types.GraphQLInt,
					},
				},
				Resolve: func(ctx context.Context, p types.GQLFRParams) interface{} {
					newNumber := 0
					obj, _ := p.Source.(*testRoot)
					newNumber, _ = p.Args["newNumber"].(int)
					return obj.PromiseToChangeTheNumber(newNumber)
				},
			},
			"failToChangeTheNumber": &types.GraphQLFieldConfig{
				Type: numberHolderType,
				Args: types.GraphQLFieldConfigArgumentMap{
					"newNumber": &types.GraphQLArgumentConfig{
						Type: types.GraphQLInt,
					},
				},
				Resolve: func(ctx context.Context, p types.GQLFRParams) interface{} {
					newNumber := 0
					obj, _ := p.Source.(*testRoot)
					newNumber, _ = p.Args["newNumber"].(int)
					return obj.FailToChangeTheNumber(newNumber)
				},
			},
			"promiseAndFailToChangeTheNumber": &types.GraphQLFieldConfig{
				Type: numberHolderType,
				Args: types.GraphQLFieldConfigArgumentMap{
					"newNumber": &types.GraphQLArgumentConfig{
						Type: types.GraphQLInt,
					},
				},
				Resolve: func(ctx context.Context, p types.GQLFRParams) interface{} {
					newNumber := 0
					obj, _ := p.Source.(*testRoot)
					newNumber, _ = p.Args["newNumber"].(int)
					return obj.PromiseAndFailToChangeTheNumber(newNumber)
				},
			},
		},
	}),
})

func TestMutations_ExecutionOrdering_EvaluatesMutationsSerially(t *testing.T) {

	root := newTestRoot(6)
	doc := `mutation M {
      first: immediatelyChangeTheNumber(newNumber: 1) {
        theNumber
      },
      second: promiseToChangeTheNumber(newNumber: 2) {
        theNumber
      },
      third: immediatelyChangeTheNumber(newNumber: 3) {
        theNumber
      }
      fourth: promiseToChangeTheNumber(newNumber: 4) {
        theNumber
      },
      fifth: immediatelyChangeTheNumber(newNumber: 5) {
        theNumber
      }
    }`

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"first": map[string]interface{}{
				"theNumber": 1,
			},
			"second": map[string]interface{}{
				"theNumber": 2,
			},
			"third": map[string]interface{}{
				"theNumber": 3,
			},
			"fourth": map[string]interface{}{
				"theNumber": 4,
			},
			"fifth": map[string]interface{}{
				"theNumber": 5,
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: mutationsTestSchema,
		AST:    ast,
		Root:   root,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestMutations_EvaluatesMutationsCorrectlyInThePresenceOfAFailedMutation(t *testing.T) {

	root := newTestRoot(6)
	doc := `mutation M {
      first: immediatelyChangeTheNumber(newNumber: 1) {
        theNumber
      },
      second: promiseToChangeTheNumber(newNumber: 2) {
        theNumber
      },
      third: failToChangeTheNumber(newNumber: 3) {
        theNumber
      }
      fourth: promiseToChangeTheNumber(newNumber: 4) {
        theNumber
      },
      fifth: immediatelyChangeTheNumber(newNumber: 5) {
        theNumber
      }
      sixth: promiseAndFailToChangeTheNumber(newNumber: 6) {
        theNumber
      }
    }`

	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"first": map[string]interface{}{
				"theNumber": 1,
			},
			"second": map[string]interface{}{
				"theNumber": 2,
			},
			"third": nil,
			"fourth": map[string]interface{}{
				"theNumber": 4,
			},
			"fifth": map[string]interface{}{
				"theNumber": 5,
			},
			"sixth": nil,
		},
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Cannot change the number`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 8, Column: 7},
				},
			},
			graphqlerrors.GraphQLFormattedError{
				Message: `Cannot change the number`,
				Locations: []location.SourceLocation{
					location.SourceLocation{Line: 17, Column: 7},
				},
			},
		},
	}
	// parse query
	ast := testutil.Parse(t, doc)

	// execute
	ep := executor.ExecuteParams{
		Schema: mutationsTestSchema,
		AST:    ast,
		Root:   root,
	}
	result := testutil.Execute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	t.Skipf("Testing equality for slice of errors in results")
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
