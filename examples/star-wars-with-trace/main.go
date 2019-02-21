package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/testutil"

	"github.com/graphql-go/graphql"
)

func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		var query string
		if r.Method == http.MethodGet {
			query = r.URL.Query().Get("query")
		} else {
			q := struct {
				Query string `json:"query"`
			}{}
			bs, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			err = json.Unmarshal(bs, &q)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			query = q.Query
		}
		testutil.StarWarsSchema.AddExtensions(&Tracer{})

		result := graphql.Do(graphql.Params{
			Schema:        testutil.StarWarsSchema,
			RequestString: query,
		})
		json.NewEncoder(w).Encode(result)
	})
	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/graphql?query={hero{name}}'")
	http.ListenAndServe(":8080", nil)
}

// Test Tracer

type Tracer struct {
	result    TracingResult
	resolvers map[string]*TracingResolverResult
}

func (t *Tracer) Init(ctx context.Context, _ *graphql.Params) {
	t.result = TracingResult{
		Version: 1,
		Execution: TracingExecutionResult{
			Resolvers: []TracingResolverResult{},
		},
	}
	t.resolvers = make(map[string]*TracingResolverResult)
}

func (t *Tracer) Name() string {
	return "tracing"
}

func (t *Tracer) HasResult() bool {
	return true
}

func (t *Tracer) GetResult(ctx context.Context) interface{} {
	log.Println("getresult was called")
	return t.result
}

func (t *Tracer) ParseDidStart(context.Context) {
	t.result.Parsing.StartOffset = time.Since(t.result.StartTime)
}

func (t *Tracer) ParseEnded(context.Context, error) {
	t.result.Parsing.Duration = time.Since(t.result.StartTime.Add(t.result.Parsing.StartOffset))
}

func (t *Tracer) ValidationDidStart(context.Context) {
	t.result.Validation.StartOffset = time.Since(t.result.StartTime)
}

func (t *Tracer) ValidationEnded(context.Context, []gqlerrors.FormattedError) {
	t.result.Validation.Duration = time.Since(t.result.StartTime.Add(t.result.Validation.StartOffset))
}

func (t *Tracer) ExecutionDidStart(ctx context.Context) {
	t.result.StartTime = time.Now()
	log.Println("Execution did start")
}

func (t *Tracer) ExecutionEnded(ctx context.Context) {
	t.result.EndTime = time.Now()
	t.result.Duration = t.result.EndTime.Sub(t.result.StartTime)
	for _, r := range t.resolvers {
		t.result.Execution.Resolvers = append(t.result.Execution.Resolvers, *r)
	}
	log.Println("Execution ended")
}

func (t *Tracer) ResolveFieldDidStart(ctx context.Context, i *graphql.ResolveInfo) {
	t.resolvers[fmt.Sprint(i.Path.AsArray())] = &TracingResolverResult{
		Path:        i.Path.AsArray(),
		ParentType:  i.ParentType.String(),
		FieldName:   i.FieldName,
		ReturnType:  i.ReturnType.String(),
		StartOffset: time.Since(t.result.StartTime),
	}

	log.Printf("Resolving %+v field started!", i.FieldName)
}

func (t *Tracer) ResolveFieldEnded(ctx context.Context, i *graphql.ResolveInfo) {
	r := t.resolvers[fmt.Sprint(i.Path.AsArray())]
	r.Duration = time.Since(t.result.StartTime.Add(t.resolvers[fmt.Sprint(i.Path.AsArray())].StartOffset))
	log.Printf("Resolving %v field ended!", i.FieldName)
}

type TracingResult struct {
	Version    int                     `json:"version"`
	StartTime  time.Time               `json:"startTime"`
	EndTime    time.Time               `json:"endTime"`
	Duration   time.Duration           `json:"duration"`
	Parsing    TracingParsingResult    `json:"parsing"`
	Validation TracingValidationResult `json:"validation"`
	Execution  TracingExecutionResult  `json:"execution"`
}

type TracingParsingResult struct {
	StartOffset time.Duration `json:"startOffset"`
	Duration    time.Duration `json:"duration"`
}

type TracingValidationResult struct {
	StartOffset time.Duration `json:"startOffset"`
	Duration    time.Duration `json:"duration"`
}

type TracingExecutionResult struct {
	Resolvers []TracingResolverResult `json:"resolvers"`
}

type TracingResolverResult struct {
	Path        []interface{} `json:"path"`
	ParentType  string        `json:"parentType"`
	FieldName   string        `json:"fieldName"`
	ReturnType  string        `json:"returnType"`
	StartOffset time.Duration `json:"startOffset"`
	Duration    time.Duration `json:"duration"`
}
