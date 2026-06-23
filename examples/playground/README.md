# Express GraphQL playground

This just to run an Express GraphQL client which uses proxy to call the go graphql server.
Just an easy hack to show beautiful user interface.

## Pre-Requisites

Nodejs and NPM

## How it works

```bash
cd examples/playground
npm install

go run main.go
```

## Complex Sample

modify PLAYGROUND_PORT or GRAPHQL_PORT if you want:

```go
// main.go
cmd.Env = append(os.Environ(),
  fmt.Sprintf("GRAPHQL_PORT=%d", GRAPHQL_PORT),       // GRAPHQL_PORT
  fmt.Sprintf("PLAYGROUND_PORT=%d", PLAYGROUND_PORT), // this value is used
)
```

You can pass query schema:

```go
query := `
  type Query {
    hello: String
  }
`
cmd := exec.Command("node", "index.js", query)
```

Once playground server will run you will see output like

🚀 GraphQL Express playground server is running on: <http://localhost:8081/graphql>

Open in browser:  [playground](http://localhost:8081/graphql)


