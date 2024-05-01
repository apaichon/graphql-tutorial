package graphql

import (
	"graphql-api/pkg/data/models"
	"graphql-api/pkg/graphql/resolvers"

	"github.com/graphql-go/graphql"
)

/*
Contact Types
*/
var ContactGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Contact",
	Fields: graphql.Fields{
		"contact_id": &graphql.Field{Type: graphql.Int},
		"name":       &graphql.Field{Type: graphql.String},
		"first_name": &graphql.Field{Type: graphql.String},
		"last_name":  &graphql.Field{Type: graphql.String},
		"gender_id":  &graphql.Field{Type: graphql.Int},
		"dob":        &graphql.Field{Type: graphql.DateTime},
		"email":      &graphql.Field{Type: graphql.String},
		"phone":      &graphql.Field{Type: graphql.String},
		"address":    &graphql.Field{Type: graphql.String},
		"photo_path": &graphql.Field{Type: graphql.String},
		"created_at": &graphql.Field{Type: graphql.DateTime},
		"created_by": &graphql.Field{Type: graphql.String},
		// Add field here
	},
})

var ContactPaginationGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ContactPagination",
	Fields: graphql.Fields{
		"contacts": &graphql.Field{Type: graphql.NewList(ContactGraphQLType)},
		"pagination":       &graphql.Field{Type: PaginationGraphQLType  },	
		// Add field here
	},
})



type ContactQueries struct {
	Gets func(string) ([]*models.ContactModel, error) `json:"gets"`
	GetPagination func(string) (*models.ContactPaginationModel, error) `json:"getPagination"`
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
		"getPagination": &graphql.Field{
			Type:    ContactPaginationGraphQLType,
			Args:    SearhTextPaginationQueryArgument,
			Resolve: resolvers.GetContactsPaginationResolve,
		},
		"getById": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    IdArgument,
			Resolve: resolvers.GetContactByIdResolve,
		},

	},
})

type ContactMutations struct {
	CreateContact func(map[string]interface{}) (*models.ContactModel, error) `json:"createContact"`
	CreateContacts func(map[string]interface{}) ([]*models.ContactModel, error) `json:"createContacts"`
	UpdateContact func(map[string]interface{}) (*models.ContactModel, error) `json:"updateContact"`	
	DeleteContact func(int) (int, error) `json:"deleteContact"`	
}

var StatusGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Status",
	Fields: graphql.Fields{
		"status_id": &graphql.Field{Type: graphql.Int},
		"status": &graphql.Field{Type: graphql.String},
		"message": &graphql.Field{Type: graphql.String},
		// Add field here
	},
})


// Define the ContactMutations type
var ContactMutationsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ContactMutations",
	Fields: graphql.Fields{
		"createContact": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    CreateContactArgument,
			Resolve: resolvers.CretateContactResolve,
		},
		"createContacts": &graphql.Field{
			Type:    StatusGraphQLType,
			Args:    CreateContactsArgument,
			Resolve: resolvers.CretateContactsResolve,
		},
		"updateContact": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    UpdateContactArgument,
			Resolve: resolvers.UpdateContactResolve,
		},
		"deleteContact": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    IdArgument,
			Resolve: resolvers.DeleteContactResolve,
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
	Get func(string) (*models.ImageModel, error) `json:"get"`
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

/*
Event Type
*/
var EventGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Event",
	Fields: graphql.Fields{
		"event_id":        &graphql.Field{Type: graphql.Int},
		"parent_event_id": &graphql.Field{Type: graphql.Int},
		"name":            &graphql.Field{Type: graphql.String},
		"description":     &graphql.Field{Type: graphql.String},
		"start_date":      &graphql.Field{Type: graphql.DateTime},
		"end_date":        &graphql.Field{Type: graphql.DateTime},
		"location_id":     &graphql.Field{Type: graphql.Int},
		"created_at": &graphql.Field{Type: graphql.DateTime},
		"created_by": &graphql.Field{Type: graphql.String},
		// Add field here
	},
})

type EventQueries struct {
	Gets    func(string) ([]*models.EventModel, error) `json:"gets"`
	GetById func(string) (*models.EventModel, error)   `json:"getById"`
}

// Define the EventQueries type
var EventQueriesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "EventQueries",
	Fields: graphql.Fields{
		"gets": &graphql.Field{
			Type:    graphql.NewList(EventGraphQLType),
			Args:    SearhTextQueryArgument,
			Resolve: resolvers.GetEventsResolve,
		},
		"getById": &graphql.Field{
			Type:    EventGraphQLType,
			Args:    IdArgument,
			Resolve: resolvers.GetEventByIdResolve,
		},
	},
})

/*
Ticket Type
*/
var TicketGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Ticket",
	Fields: graphql.Fields{
		"ticket_id": &graphql.Field{Type: graphql.Int},
		"type":      &graphql.Field{Type: graphql.String},
		"price":     &graphql.Field{Type: graphql.Float},
		"event_id":    &graphql.Field{Type: graphql.Int},
		"created_at": &graphql.Field{Type: graphql.DateTime},
		"created_by": &graphql.Field{Type: graphql.String},
		// Add field here
	},
})

type TicketQueries struct {
	Gets    func(string) ([]*models.TicketModel, error) `json:"gets"`
	GetById func(string) (*models.TicketModel, error)   `json:"getById"`
}

// Define the TicketQueries type
var TicketQueriesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TicketQueries",
	Fields: graphql.Fields{
		"gets": &graphql.Field{
			Type:    graphql.NewList(TicketGraphQLType),
			Args:    SearhTextQueryArgument,
			Resolve: resolvers.GetTicketsResolve,
		},
		"getById": &graphql.Field{
			Type:    TicketGraphQLType,
			Args:    IdArgument,
			Resolve: resolvers.GetTicketByIdResolve,
		},
	},
})
/*
Ticket Event Type
*/
var TicketEventGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TicketEvent",
	Fields: graphql.Fields{
		"ticket": &graphql.Field{Type: TicketGraphQLType},
		"event":  &graphql.Field{Type: EventGraphQLType},
		// Add field here
	},
})

type TicketEventQueries struct {
	Gets    func(string) ([]*models.TicketEventModel, error) `json:"gets"`
	GetById func(string) (*models.TicketEventModel, error)   `json:"getById"`
}

// Define the TicketQueries type
var TicketEventsQueriesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TicketEventsQueries",
	Fields: graphql.Fields{
		"gets": &graphql.Field{
			Type:    graphql.NewList(TicketEventGraphQLType),
			Args:    SearhTextQueryArgument,
			Resolve: resolvers.GetTicketEventsResolve,
		},
	},
})

/*
Pagination Type
*/
var PaginationGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Pagination",
	Fields: graphql.Fields{
		"page": &graphql.Field{Type: graphql.Int},
		"pageSize":  &graphql.Field{Type: graphql.Int},
		"totalPages":  &graphql.Field{Type: graphql.Int},
		"totalItems":  &graphql.Field{Type: graphql.Int},
		"hasNext":  &graphql.Field{Type: graphql.Boolean},
		"hasPrevious":  &graphql.Field{Type: graphql.Boolean},
		// Add field here
	},
})

