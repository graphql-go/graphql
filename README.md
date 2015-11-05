# graphql [![Build Status](https://travis-ci.org/graphql-go/graphql.svg)](https://travis-ci.org/graphql-go/graphql) [![GoDoc](https://godoc.org/graphql.co/graphql?status.svg)](https://godoc.org/github.com/graphql-go/graphql) [![Coverage Status](https://coveralls.io/repos/graphql-go/graphql/badge.svg?branch=master&service=github)](https://coveralls.io/github/graphql-go/graphql?branch=master) [![Join the chat at https://gitter.im/chris-ramon/graphql](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/graphql-go/graphql?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)


A *work-in-progress* implementation of GraphQL for Go.

### Origin and Current Direction

This project was originally a port of [v0.4.3](https://github.com/graphql/graphql-js/releases/tag/v0.4.3) of [graphql-js](https://github.com/graphql/graphql-js) (excluding the Validator), which was based on the July 2015 GraphQL specification. `graphql` is currently several versions behind `graphql-js`, however future efforts will be guided directly by the [latest formal GraphQL specification](https://github.com/facebook/graphql/releases) (currently: [October 2015](https://github.com/facebook/graphql/releases/tag/October2015)).

### Third Party Libraries
| Name          | Author        | Description  |
|:-------------:|:-------------:|:------------:|
| [graphql-go-handler](https://github.com/graphql-go/graphql-go-handler) | [Hafiz Ismail](https://github.com/sogko) | Middleware to handle GraphQL queries through HTTP requests. |
| [graphql-relay-go](https://github.com/graphql-go/graphql-relay-go) | [Hafiz Ismail](https://github.com/sogko) | Lib to construct a graphql-go server supporting react-relay. |
| [golang-relay-starter-kit](https://github.com/graphql-go/golang-relay-starter-kit) | [Hafiz Ismail](https://github.com/sogko) | Barebones starting point for a Relay application with Golang GraphQL server. |

### Blog Posts
- [Golang + GraphQL + Relay](http://wehavefaces.net/)

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
  - [ ] Basic Usage (see: [PR-#21](https://github.com/graphql-go/graphql/pull/21)) 
  - [ ] React/Relay
- [ ] Alpha Release (v0.1)

The `Validator` is optional, per official GraphQL specification, but it would be a useful addition.
