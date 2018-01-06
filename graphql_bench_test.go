package graphql_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/benchutil"
)

type B struct {
	Query  string
	Schema graphql.Schema
}

func benchGraphql(bench B, p graphql.Params, t testing.TB) {
	result := graphql.Do(p)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
}

// Benchmark a reasonably large list of small items.
func BenchmarkListQuery_1(b *testing.B) {
	nItemsListQueryBenchmark(1)(b)
}

func BenchmarkListQuery_100(b *testing.B) {
	nItemsListQueryBenchmark(100)(b)
}

func BenchmarkListQuery_1K(b *testing.B) {
	nItemsListQueryBenchmark(1000)(b)
}

func BenchmarkListQuery_10K(b *testing.B) {
	nItemsListQueryBenchmark(10 * 1000)(b)
}

func BenchmarkListQuery_100K(b *testing.B) {
	nItemsListQueryBenchmark(100 * 1000)(b)
}

func nItemsListQueryBenchmark(x int) func(b *testing.B) {
	return func(b *testing.B) {
		schema := benchutil.ListSchemaWithXItems(x)

		bench := B{
			Query: `
				query {
					colors {
						hex
						r
						g
						b
					}
				}
			`,
			Schema: schema,
		}

		for i := 0; i < b.N; i++ {

			params := graphql.Params{
				Schema:        schema,
				RequestString: bench.Query,
			}
			benchGraphql(bench, params, b)
		}
	}
}

func BenchmarkWideQuery_1_1(b *testing.B) {
	nFieldsyItemsQueryBenchmark(1, 1)(b)
}

func BenchmarkWideQuery_10_1(b *testing.B) {
	nFieldsyItemsQueryBenchmark(10, 1)(b)
}

func BenchmarkWideQuery_100_1(b *testing.B) {
	nFieldsyItemsQueryBenchmark(100, 1)(b)
}

func BenchmarkWideQuery_1K_1(b *testing.B) {
	nFieldsyItemsQueryBenchmark(1000, 1)(b)
}

func BenchmarkWideQuery_1_10(b *testing.B) {
	nFieldsyItemsQueryBenchmark(1, 10)(b)
}

func BenchmarkWideQuery_10_10(b *testing.B) {
	nFieldsyItemsQueryBenchmark(10, 10)(b)
}

func BenchmarkWideQuery_100_10(b *testing.B) {
	nFieldsyItemsQueryBenchmark(100, 10)(b)
}

func BenchmarkWideQuery_1K_10(b *testing.B) {
	nFieldsyItemsQueryBenchmark(1000, 10)(b)
}

func nFieldsyItemsQueryBenchmark(x int, y int) func(b *testing.B) {
	return func(b *testing.B) {
		schema := benchutil.WideSchemaWithXFieldsAndYItems(x, y)
		query := benchutil.WideSchemaQuery(x)

		bench := B{
			Query:  query,
			Schema: schema,
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			params := graphql.Params{
				Schema:        schema,
				RequestString: bench.Query,
			}
			benchGraphql(bench, params, b)
		}
	}
}
