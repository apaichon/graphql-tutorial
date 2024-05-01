
package graphql

import (
	"github.com/graphql-go/graphql"
)

// RootQuery represents the root GraphQL query.
var RootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"contacts": &graphql.Field{
				Type: ContactQueriesType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return &ContactQueries{}, nil
				},
			},
			"images": &graphql.Field{
				Type: ImageQueriesType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return &ImageQueries{}, nil
				},
			},
			"events": &graphql.Field{
				Type: EventQueriesType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return &EventQueries{}, nil
				},
			},
			"tickets": &graphql.Field{
				Type: TicketQueriesType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return &TicketQueries{}, nil
				},
			},
			
			"ticketEvents": &graphql.Field{
				Type: TicketEventsQueriesType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return &TicketEventQueries{}, nil
				},
			},
			// Add other queries as needed
		},
	},
)

// RootQuery represents the root GraphQL query.
var RootMutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"contactMutations": &graphql.Field{
				Type: ContactMutationsType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return &ContactMutations{}, nil
				},
			},
			// Add other queries as needed
		},
	},
)