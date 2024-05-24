package graphql

import (
	"github.com/graphql-go/graphql"
)

var SearhTextQueryArgument = graphql.FieldConfigArgument{
	"searchText": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"limit": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
	"offset": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
}

var SearhTextPaginationQueryArgument = graphql.FieldConfigArgument{
	"searchText": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
	"page": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
	"pageSize": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
}

var IdArgument = graphql.FieldConfigArgument{
	"id": &graphql.ArgumentConfig{
		Type: graphql.Int,
	},
}

var CreateContactInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CreateContactInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name":       {Type: graphql.NewNonNull(graphql.String)},
		"first_name": {Type: graphql.NewNonNull(graphql.String)},
		"last_name":  {Type: graphql.NewNonNull(graphql.String)},
		"gender_id":  {Type: graphql.NewNonNull(graphql.Int)},
		"dob":        {Type: graphql.NewNonNull(graphql.DateTime)},
		"email":      {Type: graphql.NewNonNull(graphql.String)},
		"address":    {Type: graphql.String},
		"phone":      {Type: graphql.String},
		"photo_path": {Type: graphql.String},
	},
})

var CreateContactArgument = graphql.FieldConfigArgument{
	"input": &graphql.ArgumentConfig{
		Type: CreateContactInput,
	},
}

var CreateContactsArgument = graphql.FieldConfigArgument{
	"contacts": &graphql.ArgumentConfig{
		Type: graphql.NewNonNull(graphql.NewList(CreateContactInput)),
	},
}

var UpdateContactArgument = graphql.FieldConfigArgument{
	"input": &graphql.ArgumentConfig{
		Type: graphql.NewInputObject(graphql.InputObjectConfig{
			Name: "UpdateContactInput",
			Fields: graphql.InputObjectConfigFieldMap{
				"contact_id": {Type: graphql.NewNonNull(graphql.Int)},
				"name":       {Type: graphql.NewNonNull(graphql.String)},
				"first_name": {Type: graphql.NewNonNull(graphql.String)},
				"last_name":  {Type: graphql.NewNonNull(graphql.String)},
				"gender_id":  {Type: graphql.NewNonNull(graphql.Int)},
				"dob":        {Type: graphql.NewNonNull(graphql.DateTime)},
				"email":      {Type: graphql.NewNonNull(graphql.String)},
				"address":    {Type: graphql.String},
				"phone":      {Type: graphql.String},
				"photo_path": {Type: graphql.String},
			},
		}),
	},
}

var CreateBidingInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CreateBidingInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"room_id":   {Type: graphql.NewNonNull(graphql.Int)},
		"bidder":    {Type: graphql.NewNonNull(graphql.String)},
		"bid_price": {Type: graphql.NewNonNull(graphql.Float)},
		"bid_time":  {Type: graphql.NewNonNull(graphql.DateTime)},
	},
})

var CreateBidingArgument = graphql.FieldConfigArgument{
	"input": &graphql.ArgumentConfig{
		Type: CreateBidingInput,
	},
}
