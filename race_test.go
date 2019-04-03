package graphql_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestRace(t *testing.T) {
	tempfile, err := ioutil.TempFile("", "examplerace.*.go")
	if err != nil {
		t.Fatal(err)
	}
	filename := tempfile.Name()
	t.Log(filename)

	defer os.Remove(filename)

	_, err = tempfile.Write([]byte(`
		package main

		import (
			"runtime"
			"sync"

			"github.com/graphql-go/graphql"
		)

		func main() {
			var wg sync.WaitGroup
			wg.Add(2)
			for i := 0; i < 2; i++ {
				go func() {
					defer wg.Done()
					schema, _ := graphql.NewSchema(graphql.SchemaConfig{
						Query: graphql.NewObject(graphql.ObjectConfig{
							Name: "RootQuery",
							Fields: graphql.Fields{
								"hello": &graphql.Field{
									Type: graphql.String,
									Resolve: func(p graphql.ResolveParams) (interface{}, error) {
										return "world", nil
									},
								},
							},
						}),
					})
					runtime.KeepAlive(schema)
				}()
			}

			wg.Wait()
		} 
	`))
	if err != nil {
		t.Fatal(err)
	}

	if err := tempfile.Close(); err != nil {
		t.Fatal(err)
	}

	result, err := exec.Command("go", "run", "-race", filename).CombinedOutput()
	if err != nil || len(result) != 0 {
		t.Log(string(result))
		t.Fatal(err)
	}
}
