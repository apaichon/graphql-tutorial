
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
			"biding": &graphql.Field{
				Type: BidingQueriesType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return &BidingQueries{}, nil
				},
			},
			// Add other queries as needed
		},
	},
)

// RootMutation represents the root GraphQL query.
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
			"bidingMutations": &graphql.Field{
				Type: BidingMutationsType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return &BidingMutations{}, nil
				},
			},
			// Add other queries as needed
		},
	},
)

// Root Subscription represents the root GraphQL query.
var RootSubscription = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootSubscription",
		Fields: graphql.Fields{
			"contactSubscriptions": &graphql.Field{
				Type: ContactSubscriptionsType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return &ContactSubscription{}, nil
				},
			},
			// Add other queries as needed
		},
	},
)