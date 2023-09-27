const express = require("express");
const axios = require("axios");
const bodyparser = require("body-parser");
const cors = require("cors");
const { buildSchema } = require("graphql");
const graphqlHTTP = require("express-graphql");
let [query = ""] = process.argv.slice(2);
if (!query.trim()) {
  query = `
  type Query {
    hello: String
  }
`;
}
const schema = buildSchema(query);
const PLAYGROUND_PORT = process.env.PLAYGROUND_PORT || 4000;
const GRAPHQL_PORT = process.env.GRAPHQL_PORT || 8080;

const app = express();

app.use(cors());
app.get(
  "/graphql",
  graphqlHTTP({
    schema: schema,
    rootValue: {},
    graphiql: true,
  })
);
app.post("/*", bodyparser.json(), (req, res) => {
  const options = {
    url: req.path,
    baseURL: `http://localhost:${GRAPHQL_PORT}/`,
    method: "get",
    params: req.body,
  };
  axios(options).then(({ data }) => res.send(data));
});
app.listen(PLAYGROUND_PORT, () => {
  console.log(
    `🚀 GraphQL Express playground server is running on: http://localhost:${PLAYGROUND_PORT}/graphql`
  );
});
