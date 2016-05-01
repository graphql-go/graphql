package graphql_test

import (
	"testing"

	"github.com/sprucehealth/graphql"
	"github.com/sprucehealth/graphql/language/parser"
	"github.com/sprucehealth/graphql/language/source"
	"github.com/sprucehealth/graphql/testutil"
)

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
