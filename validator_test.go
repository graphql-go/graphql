package graphql_test

import (
	"testing"

	"github.com/sprucehealth/graphql"
	"github.com/sprucehealth/graphql/language/parser"
	"github.com/sprucehealth/graphql/language/source"
	"github.com/sprucehealth/graphql/testutil"
)

func TestConcurrentValidateDocument(t *testing.T) {
	validate := func() {
		query := `
		query HeroNameAndFriendsQuery {
			hero {
				id
				name
				friends {
					name
				}
			}
		}
	`
		ast, err := parser.Parse(parser.ParseParams{Source: source.New("", query)})
		if err != nil {
			t.Fatal(err)
		}
		r := graphql.ValidateDocument(&testutil.StarWarsSchema, ast, nil)
		if !r.IsValid {
			t.Fatal("Not valid")
		}
	}
	go validate()
	validate()
}

func BenchmarkValidateDocument(b *testing.B) {
	query := `
		query HeroNameAndFriendsQuery {
			hero {
				id
				name
				friends {
					name
				}
			}
		}
	`
	ast, err := parser.Parse(parser.ParseParams{Source: source.New("", query)})
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := graphql.ValidateDocument(&testutil.StarWarsSchema, ast, nil)
		if !r.IsValid {
			b.Fatal("Not valid")
		}
	}
}
