package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"another-graph/models"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)


func main() {
	// Define Name GraphQL Object
	nameType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Name",
		Fields: graphql.Fields{
			"title": &graphql.Field{Type: graphql.String},
			"first": &graphql.Field{Type: graphql.String},
			"last":  &graphql.Field{Type: graphql.String},
		},
	})

	// Define Street GraphQL Object
	streetType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Street",
		Fields: graphql.Fields{
			"number": &graphql.Field{Type: graphql.Int},
			"name":   &graphql.Field{Type: graphql.String},
		},
	})

	// Define Coordinates GraphQL Object
	coordinatesType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Coordinates",
		Fields: graphql.Fields{
			"latitude":  &graphql.Field{Type: graphql.String},
			"longitude": &graphql.Field{Type: graphql.String},
		},
	})

	// Define Timezone GraphQL Object
	timezoneType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Timezone",
		Fields: graphql.Fields{
			"offset":      &graphql.Field{Type: graphql.String},
			"description": &graphql.Field{Type: graphql.String},
		},
	})

	// Define Location GraphQL Object
	locationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Location",
		Fields: graphql.Fields{
			"street":      &graphql.Field{Type: streetType},
			"city":        &graphql.Field{Type: graphql.String},
			"state":       &graphql.Field{Type: graphql.String},
			"country":     &graphql.Field{Type: graphql.String},
			"postcode":    &graphql.Field{Type: graphql.Int},
			"coordinates": &graphql.Field{Type: coordinatesType},
			"timezone":    &graphql.Field{Type: timezoneType},
		},
	})

	// Define Login GraphQL Object
	loginType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Login",
		Fields: graphql.Fields{
			"uuid":     &graphql.Field{Type: graphql.String},
			"username": &graphql.Field{Type: graphql.String},
			"password": &graphql.Field{Type: graphql.String},
			"salt":     &graphql.Field{Type: graphql.String},
			"md5":      &graphql.Field{Type: graphql.String},
			"sha1":     &graphql.Field{Type: graphql.String},
			"sha256":   &graphql.Field{Type: graphql.String},
		},
	})

	// Define DOB GraphQL Object
	dobType := graphql.NewObject(graphql.ObjectConfig{
		Name: "DOB",
		Fields: graphql.Fields{
			"date": &graphql.Field{Type: graphql.String},
			"age":  &graphql.Field{Type: graphql.Int},
		},
	})

	// Define Registered GraphQL Object
	registeredType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Registered",
		Fields: graphql.Fields{
			"date": &graphql.Field{Type: graphql.String},
			"age":  &graphql.Field{Type: graphql.Int},
		},
	})

	// Define ID GraphQL Object
	idType := graphql.NewObject(graphql.ObjectConfig{
		Name: "ID",
		Fields: graphql.Fields{
			"name":  &graphql.Field{Type: graphql.String},
			"value": &graphql.Field{Type: graphql.String},
		},
	})

	// Define Picture GraphQL Object
	pictureType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Picture",
		Fields: graphql.Fields{
			"large":     &graphql.Field{Type: graphql.String},
			"medium":    &graphql.Field{Type: graphql.String},
			"thumbnail": &graphql.Field{Type: graphql.String},
		},
	})

	// Define Result GraphQL Object
	resultType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Result",
		Fields: graphql.Fields{
			"gender":     &graphql.Field{Type: graphql.String},
			"name":       &graphql.Field{Type: nameType},
			"location":   &graphql.Field{Type: locationType},
			"email":      &graphql.Field{Type: graphql.String},
			"login":      &graphql.Field{Type: loginType},
			"dob":        &graphql.Field{Type: dobType},
			"registered": &graphql.Field{Type: registeredType},
			"phone":      &graphql.Field{Type: graphql.String},
			"cell":       &graphql.Field{Type: graphql.String},
			"id":         &graphql.Field{Type: idType},
			"picture":    &graphql.Field{Type: pictureType},
			"nat":        &graphql.Field{Type: graphql.String},
		},
	})

	// Define Info GraphQL Object
	infoType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Info",
		Fields: graphql.Fields{
			"seed":    &graphql.Field{Type: graphql.String},
			"results": &graphql.Field{Type: graphql.Int},
			"page":    &graphql.Field{Type: graphql.Int},
			"version": &graphql.Field{Type: graphql.String},
		},
	})

	// Define Data GraphQL Object
	dataType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Data",
		Fields: graphql.Fields{
			"results": &graphql.Field{
				Type: graphql.NewList(resultType),
			},
			"info": &graphql.Field{Type: infoType},
		},
	})

	// Define the Query Object
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"data": &graphql.Field{
				Type: dataType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Here you should fetch and return the actual data

					data, err := fetchUserData()
					if err != nil {
						return models.Data{}, nil
					}
					// fmt.Printf("result-random:%v", data)
					return data, nil
				},
			},
		},
	})

	// Create the Schema
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// _ = schema // This schema object would be used by your GraphQL server
	// Create a GraphQL handler for HTTP requests
	graphqlHandler := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false, // Disable GraphiQL for subscriptions endpoint
		Playground: true,
	})


	// Serve GraphQL API at /graphql endpoint
	http.Handle("/graphql", graphqlHandler)
	server := &http.Server{
		Addr: fmt.Sprintf(":%v", 4001),
	}
	fmt.Printf("Server is running at http://localhost:%v/graphql\n", 4001)
	server.ListenAndServe()
}

// fetchUserData fetches data from the Random User API
func fetchUserData() (models.Data, error) {
	response, err := http.Get("https://randomuser.me/api/")
	if err != nil {
		return models.Data{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return models.Data{}, err
	}

	var user models.Data
	err = json.Unmarshal(body, &user)
	if err != nil {
		return models.Data{}, err
	}

	// fmt.Printf("resutl:%v", user)

	return user, nil
}
