package graphql

import (
	"github.com/graphql-go/graphql"
	"graphql-api/pkg/data/models"
	"graphql-api/pkg/graphql/resolvers"
)

/*
Contact Types
*/
var ContactGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Contact",
	Fields: graphql.Fields{
		"contact_id": &graphql.Field{Type: graphql.Int},
		"name": &graphql.Field{Type: graphql.String},
		"first_name": &graphql.Field{Type: graphql.String},
		"last_name": &graphql.Field{Type: graphql.String},
		"gender_id": &graphql.Field{Type: graphql.Int},
		"dob": &graphql.Field{Type: graphql.DateTime},
		"email": &graphql.Field{Type: graphql.String},
		"phone": &graphql.Field{Type: graphql.String},
		"address": &graphql.Field{Type: graphql.String},
		"photo_path": &graphql.Field{Type: graphql.String},
		"created_at": &graphql.Field{Type: graphql.DateTime},
		"created_by": &graphql.Field{Type: graphql.String},
        // Add field here
	},
})


type ContactQueries struct {
	Gets    func(string) ([]*models.ContactModel, error) `json:"gets"`
}

// Define the ContactQueries type
var ContactQueriesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ContactQueries",
	Fields: graphql.Fields{
		"gets": &graphql.Field{
			Type:    graphql.NewList(ContactGraphQLType),
			Args:    SearhTextQueryArgument,
			Resolve: resolvers.GetContactResolve,
		},
	},
})

/*
Image Type
*/
var ImageGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Image",
	Fields: graphql.Fields{
		"image_url": &graphql.Field{Type: graphql.String},
        // Add field here
	},
})


type ImageQueries struct {
	Get    func(string) (*models.ImageModel, error) `json:"get"`
}

var ImageQueriesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ImageQueries",
	Fields: graphql.Fields{
		"get": &graphql.Field{
			Type:    ImageGraphQLType,
			Resolve: resolvers.GetImageResolve,
		},
	},
})
