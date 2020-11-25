package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"log"
)

// NullString to be used in place of sql.NullString
type NullString struct {
	sql.NullString
}

// MarshalJSON from the json.Marshaler interface
func (v NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON from the json.Unmarshaler interface
func (v *NullString) UnmarshalJSON(data []byte) error {
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.String = *x
		v.Valid = true
	} else {
		v.Valid = false
	}
	return nil
}

// NewNullString create a new null string. Empty string evaluates to an
// "invalid" NullString
func NewNullString(value string) *NullString {
	var null NullString
	if value != "" {
		null.String = value
		null.Valid = true
		return &null
	}
	null.Valid = false
	return &null
}

// SerializeNullString serializes `NullString` to a string
func SerializeNullString(value interface{}) interface{} {
	switch value := value.(type) {
	case NullString:
		return value.String
	case *NullString:
		v := *value
		return v.String
	default:
		return nil
	}
}

// ParseNullString parses GraphQL variables from `string` to `CustomID`
func ParseNullString(value interface{}) interface{} {
	switch value := value.(type) {
	case string:
		return NewNullString(value)
	case *string:
		return NewNullString(*value)
	default:
		return nil
	}
}

// ParseLiteralNullString parses GraphQL AST value to `NullString`.
func ParseLiteralNullString(valueAST ast.Value) interface{} {
	switch valueAST := valueAST.(type) {
	case *ast.StringValue:
		return NewNullString(valueAST.Value)
	default:
		return nil
	}
}

// NullableString graphql *Scalar type based of NullString
var NullableString = graphql.NewScalar(graphql.ScalarConfig{
	Name:         "NullableString",
	Description:  "The `NullableString` type repesents a nullable SQL string.",
	Serialize:    SerializeNullString,
	ParseValue:   ParseNullString,
	ParseLiteral: ParseLiteralNullString,
})

/*
CREATE TABLE persons (
	favorite_dog TEXT -- is a nullable field
	);

*/

// Person noqa
type Person struct {
	Name        string      `json:"name"`
	FavoriteDog *NullString `json:"favorite_dog"` // Some people don't like dogs ¯\_(ツ)_/¯
}

// PersonType noqa
var PersonType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Person",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"favorite_dog": &graphql.Field{
			Type: NullableString,
		},
	},
})

func main() {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"people": &graphql.Field{
					Type: graphql.NewList(PersonType),
					Args: graphql.FieldConfigArgument{
						"favorite_dog": &graphql.ArgumentConfig{
							Type: NullableString,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						dog, dogOk := p.Args["favorite_dog"].(*NullString)
						people := []Person{
							Person{Name: "Alice", FavoriteDog: NewNullString("Yorkshire Terrier")},
							// `Bob`'s favorite dog will be saved as null in the database
							Person{Name: "Bob", FavoriteDog: NewNullString("")},
							Person{Name: "Chris", FavoriteDog: NewNullString("French Bulldog")},
						}
						switch {
						case dogOk:
							log.Printf("favorite_dog from arguments: %+v", dog)
							dogPeople := make([]Person, 0)
							for _, p := range people {
								if p.FavoriteDog.Valid {
									if p.FavoriteDog.String == dog.String {
										dogPeople = append(dogPeople, p)
									}
								}
							}
							return dogPeople, nil
						default:
							return people, nil
						}
					},
				},
			},
		}),
	})
	if err != nil {
		log.Fatal(err)
	}
	query := `
query {
  people {
    name
    favorite_dog
    }
}`
	queryWithArgument := `
query {
  people(favorite_dog: "Yorkshire Terrier") {
    name
    favorite_dog
  }
}`
	r1 := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	r2 := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: queryWithArgument,
	})
	if len(r1.Errors) > 0 {
		log.Fatal(r1)
	}
	if len(r2.Errors) > 0 {
		log.Fatal(r1)
	}
	b1, err := json.MarshalIndent(r1, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	b2, err := json.MarshalIndent(r2, "", "  ")
	if err != nil {
		log.Fatal(err)

	}
	fmt.Printf("\nQuery: %+v\n", string(query))
	fmt.Printf("\nResult: %+v\n", string(b1))
	fmt.Printf("\nQuery (with arguments): %+v\n", string(queryWithArgument))
	fmt.Printf("\nResult (with arguments): %+v\n", string(b2))
}

/* Output:
Query:
query {
  people {
    name
    favorite_dog
    }
}

Result: {
  "data": {
    "people": [
      {
        "favorite_dog": "Yorkshire Terrier",
        "name": "Alice"
      },
      {
        "favorite_dog": "",
        "name": "Bob"
      },
      {
        "favorite_dog": "French Bulldog",
        "name": "Chris"
      }
    ]
  }
}

Query (with arguments):
query {
  people(favorite_dog: "Yorkshire Terrier") {
    name
    favorite_dog
  }
}

Result (with arguments): {
  "data": {
    "people": [
      {
        "favorite_dog": "Yorkshire Terrier",
        "name": "Alice"
      }
    ]
  }
}
*/
