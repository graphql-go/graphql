package graphql_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
)

func TestIntrospection_SchemaWithMutationAndSubscription(t *testing.T) {
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "MutationRoot",
		Fields: graphql.Fields{
			"doSomething": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	subscriptionType := graphql.NewObject(graphql.ObjectConfig{
		Name: "SubscriptionRoot",
		Fields: graphql.Fields{
			"onEvent": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "QueryRoot",
			Fields: graphql.Fields{
				"f1": &graphql.Field{Type: graphql.String},
			},
		}),
		Mutation:     mutationType,
		Subscription: subscriptionType,
	})
	if err != nil {
		t.Fatalf("unexpected error creating schema: %v", err)
	}

	query := `
		{
			__schema {
				queryType { name }
				mutationType { name }
				subscriptionType { name }
			}
		}
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"__schema": map[string]interface{}{
				"queryType": map[string]interface{}{
					"name": "QueryRoot",
				},
				"mutationType": map[string]interface{}{
					"name": "MutationRoot",
				},
				"subscriptionType": map[string]interface{}{
					"name": "SubscriptionRoot",
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestIntrospection_SchemaWithCustomDirective(t *testing.T) {
	customDirective := graphql.NewDirective(graphql.DirectiveConfig{
		Name: "customDirective",
		Locations: []string{
			graphql.DirectiveLocationField,
			graphql.DirectiveLocationSubscription,
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "QueryRoot",
			Fields: graphql.Fields{
				"f1": &graphql.Field{Type: graphql.String},
			},
		}),
		Directives: append(graphql.SpecifiedDirectives, customDirective),
	})
	if err != nil {
		t.Fatalf("unexpected error creating schema: %v", err)
	}

	query := `
		{
			__schema {
				directives {
					name
					onOperation
					onFragment
					onField
				}
			}
		}
	`
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if result.HasErrors() {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	data, ok := result.Data.(map[string]interface{})
	if !ok {
		t.Fatal("expected data map")
	}
	schemaData, ok := data["__schema"].(map[string]interface{})
	if !ok {
		t.Fatal("expected __schema in data")
	}
	directives, ok := schemaData["directives"].([]interface{})
	if !ok {
		t.Fatal("expected directives list")
	}
	found := false
	for _, d := range directives {
		dir, ok := d.(map[string]interface{})
		if !ok {
			continue
		}
		if dir["name"] == "customDirective" {
			found = true
			onOp, ok := dir["onOperation"].(bool)
			if !ok || !onOp {
				t.Fatal("expected onOperation=true for directive with SUBSCRIPTION location")
			}
			break
		}
	}
	if !found {
		t.Fatal("expected customDirective in schema directives")
	}
}

func TestIntrospection_EnumValuesWithDeprecated(t *testing.T) {
	testEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "TestEnum",
		Values: graphql.EnumValueConfigMap{
			"NONDEPRECATED": &graphql.EnumValueConfig{
				Value: 0,
			},
			"DEPRECATED": &graphql.EnumValueConfig{
				Value:             1,
				DeprecationReason: "Removed in 2.0",
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "QueryRoot",
			Fields: graphql.Fields{
				"testEnum": &graphql.Field{
					Type: testEnum,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("unexpected error creating schema: %v", err)
	}

	query := `
		{
			__type(name: "TestEnum") {
				enumValues(includeDeprecated: true) {
					name
					isDeprecated
					deprecationReason
				}
			}
		}
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"__type": map[string]interface{}{
				"enumValues": []interface{}{
					map[string]interface{}{
						"name":              "NONDEPRECATED",
						"isDeprecated":      false,
						"deprecationReason": nil,
					},
					map[string]interface{}{
						"name":              "DEPRECATED",
						"isDeprecated":      true,
						"deprecationReason": "Removed in 2.0",
					},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if !testutil.ContainSubset(result.Data.(map[string]interface{}), expected.Data.(map[string]interface{})) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
