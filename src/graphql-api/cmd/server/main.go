package main

import (
    "fmt"
    "net/http"

    "github.com/graphql-go/graphql"
    "github.com/graphql-go/handler"
)

// Define GraphQL schema
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
    Query:    rootQuery,
    Mutation: rootMutation,
})

// Define root query
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
    Name: "RootQuery",
    Fields: graphql.Fields{
        "hello": &graphql.Field{
            Type: graphql.String,
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                return "Hello, GraphQL!", nil
            },
        },
    },
})

// Define root mutation
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
    Name: "RootMutation",
    Fields: graphql.Fields{
        "createEvent": &graphql.Field{
            Type: graphql.String,
            Args: graphql.FieldConfigArgument{
                "name": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.String),
                },
            },
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                // Here you would typically handle the mutation logic
                name, _ := p.Args["name"].(string)
                return fmt.Sprintf("Event created with name: %s", name), nil
            },
        },
    },
})

func main() {
    // Create a GraphQL handler for HTTP requests
    graphqlHandler := handler.New(&handler.Config{
        Schema: &schema,
        Pretty: true,
    })

    // Serve GraphQL API at /graphql endpoint
    http.Handle("/graphql", graphqlHandler)

    // Start the HTTP server
    fmt.Println("Server is running at http://localhost:4000/graphql")
    http.ListenAndServe(":4000", nil)
}
