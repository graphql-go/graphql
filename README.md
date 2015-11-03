graphql-go [![Build Status](https://travis-ci.org/chris-ramon/graphql-go.svg)](https://travis-ci.org/chris-ramon/graphql-go) [![GoDoc](https://godoc.org/graphql.co/graphql?status.svg)](https://godoc.org/github.com/chris-ramon/graphql-go) [![Coverage Status](https://coveralls.io/repos/chris-ramon/graphql-go/badge.svg?branch=master&service=github)](https://coveralls.io/github/chris-ramon/graphql-go?branch=master) [![Join the chat at https://gitter.im/chris-ramon/graphql-go](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/chris-ramon/graphql-go?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)


A *work-in-progress* implementation of GraphQL for Go.

### Origin and Current Direction

This project was originally a port of [v0.4.3](https://github.com/graphql/graphql-js/releases/tag/v0.4.3) of [graphql-js](https://github.com/graphql/graphql-js) (excluding the Validator), which was based on the July 2015 GraphQL specification. `graphql-go` is currently several versions behind `graphql-js`, however future efforts will be guided directly by the [latest formal GraphQL specification](https://github.com/facebook/graphql/releases) (currently: [October 2015](https://github.com/facebook/graphql/releases/tag/October2015)).

### Roadmap
- [x] Lexer
- [x] Parser
- [x] Schema Parser
- [x] Printer
- [x] Schema Printer
- [x] Visitor
- [x] Executor
- [ ] Validator
- [ ] Examples
  - [ ] Basic Usage (see: [PR-#21](https://github.com/chris-ramon/graphql-go/pull/21)) 
  - [ ] React/Relay
- [ ] Alpha Release (v0.1)

The `Validator` is optional, per official GraphQL specification, but it would be a useful addition.
