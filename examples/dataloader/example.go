package dataloaderexample

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"

	"github.com/bigdrum/godataloader"
	"github.com/graphql-go/graphql"
)

var postDB = map[string]*post{
	"1": &post{
		ID:       "1",
		Content:  "Hello 1",
		AuthorID: "1",
	},
	"2": &post{
		ID:       "2",
		Content:  "Hello 2",
		AuthorID: "1",
	},
	"3": &post{
		ID:       "3",
		Content:  "Hello 3",
		AuthorID: "2",
	},
	"4": &post{
		ID:       "4",
		Content:  "Hello 4",
		AuthorID: "2",
	},
}

var userDB = map[string]*user{
	"1": &user{
		ID:   "1",
		Name: "Mike",
	},
	"2": &user{
		ID:   "2",
		Name: "John",
	},
	"3": &user{
		ID:   "3",
		Name: "Kate",
	},
}

var loaderKey = struct{}{}

type loader struct {
	postLoader *dataloader.DataLoader
	userLoader *dataloader.DataLoader
}

func newLoader(sch *dataloader.Scheduler) *loader {
	return &loader{
		postLoader: dataloader.New(sch, dataloader.Parallel(func(key interface{}) dataloader.Value {
			// In practice, we will make remote request (e.g. SQL) to fetch post.
			// Here we just fake it.
			log.Print("Load post ", key)
			time.Sleep(time.Second)
			id := key.(string)
			return dataloader.NewValue(postDB[id], nil)
		})),
		userLoader: dataloader.New(sch, func(keys []interface{}) []dataloader.Value {
			// In practice, we will make remote request (e.g. SQL) to fetch multiple users.
			// Here we just fake it.
			log.Print("Batch load users ", keys)
			time.Sleep(time.Second)
			var ret []dataloader.Value
			for _, key := range keys {
				id := key.(string)
				ret = append(ret, dataloader.NewValue(userDB[id], nil))
			}
			return ret
		}),
	}
}

type post struct {
	ID       string `json:"id"`
	Content  string `json:"content"`
	AuthorID string `json:"author_id"`
}

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func attachNewDataLoader(parent context.Context, sch *dataloader.Scheduler) context.Context {
	dl := newLoader(sch)
	return context.WithValue(parent, loaderKey, dl)
}

func getDataLoader(ctx context.Context) *loader {
	return ctx.Value(loaderKey).(*loader)
}

func CreateSchema() graphql.Schema {
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	postType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Post",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"content": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					post := p.Source.(*post)
					id := post.AuthorID
					dl := getDataLoader(p.Context)
					return dl.userLoader.Load(id).Unbox()
				},
			},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"post": &graphql.Field{
				Type: postType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(string)
					if !ok {
						return nil, nil
					}
					dl := getDataLoader(p.Context)
					return dl.postLoader.Load(id).Unbox()
				},
			},
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(string)
					if !ok {
						return nil, nil
					}
					dl := getDataLoader(p.Context)
					return dl.userLoader.Load(id).Unbox()
				},
			},
		}})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
	if err != nil {
		panic(err)
	}
	return schema
}

type dataloaderExecutor struct {
	sch *dataloader.Scheduler
}

func (e *dataloaderExecutor) RunMany(fs []func()) {
	if len(fs) == 1 {
		fs[0]()
		return
	}
	if len(fs) == 0 {
		return
	}

	wg := dataloader.NewWaitGroup(e.sch)
	for i := range fs {
		f := fs[i]
		wg.Add(1)
		e.sch.Spawn(func() {
			defer wg.Done()
			f()
		})
	}
	wg.Wait()
}

func RunQuery(query string, schema graphql.Schema) *graphql.Result {
	var result *graphql.Result
	dataloader.RunWithScheduler(func(sch *dataloader.Scheduler) {
		executor := dataloaderExecutor{sch}
		ctx := attachNewDataLoader(context.Background(), sch)
		result = graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
			Context:       ctx,
			Executor:      &executor,
		})
		if len(result.Errors) > 0 {
			fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
		}
	})

	return result
}
