package executor_test

import (
	"reflect"
	"testing"

	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/executor"
	"github.com/chris-ramon/graphql-go/gqltypes"
	"github.com/chris-ramon/graphql-go/language/location"
	"github.com/chris-ramon/graphql-go/testutil"
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
var numberHolderType = gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
	Name: "NumberHolder",
	Fields: gqltypes.GraphQLFieldConfigMap{
		"theNumber": &gqltypes.GraphQLFieldConfig{
			Type: gqltypes.GraphQLInt,
		},
	},
})

var mutationsTestSchema, _ = gqltypes.NewGraphQLSchema(gqltypes.GraphQLSchemaConfig{
	Query: gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Query",
		Fields: gqltypes.GraphQLFieldConfigMap{
			"numberHolder": &gqltypes.GraphQLFieldConfig{
				Type: numberHolderType,
			},
		},
	}),
	Mutation: gqltypes.NewGraphQLObjectType(gqltypes.GraphQLObjectTypeConfig{
		Name: "Mutation",
		Fields: gqltypes.GraphQLFieldConfigMap{
			"immediatelyChangeTheNumber": &gqltypes.GraphQLFieldConfig{
				Type: numberHolderType,
				Args: gqltypes.GraphQLFieldConfigArgumentMap{
					"newNumber": &gqltypes.GraphQLArgumentConfig{
						Type: gqltypes.GraphQLInt,
					},
				},
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					newNumber := 0
					obj, _ := p.Source.(*testRoot)
					newNumber, _ = p.Args["newNumber"].(int)
					return obj.ImmediatelyChangeTheNumber(newNumber)
				},
			},
			"promiseToChangeTheNumber": &gqltypes.GraphQLFieldConfig{
				Type: numberHolderType,
				Args: gqltypes.GraphQLFieldConfigArgumentMap{
					"newNumber": &gqltypes.GraphQLArgumentConfig{
						Type: gqltypes.GraphQLInt,
					},
				},
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					newNumber := 0
					obj, _ := p.Source.(*testRoot)
					newNumber, _ = p.Args["newNumber"].(int)
					return obj.PromiseToChangeTheNumber(newNumber)
				},
			},
			"failToChangeTheNumber": &gqltypes.GraphQLFieldConfig{
				Type: numberHolderType,
				Args: gqltypes.GraphQLFieldConfigArgumentMap{
					"newNumber": &gqltypes.GraphQLArgumentConfig{
						Type: gqltypes.GraphQLInt,
					},
				},
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
					newNumber := 0
					obj, _ := p.Source.(*testRoot)
					newNumber, _ = p.Args["newNumber"].(int)
					return obj.FailToChangeTheNumber(newNumber)
				},
			},
			"promiseAndFailToChangeTheNumber": &gqltypes.GraphQLFieldConfig{
				Type: numberHolderType,
				Args: gqltypes.GraphQLFieldConfigArgumentMap{
					"newNumber": &gqltypes.GraphQLArgumentConfig{
						Type: gqltypes.GraphQLInt,
					},
				},
				Resolve: func(p gqltypes.GQLFRParams) interface{} {
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

	expected := &gqltypes.GraphQLResult{
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

	expected := &gqltypes.GraphQLResult{
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
