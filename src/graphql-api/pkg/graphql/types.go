package graphql

import (
	// "strconv"
	"graphql-api/internal/auth"
	"graphql-api/internal/cache"
	"graphql-api/internal/subscription"
	"graphql-api/pkg/data/models"
	"graphql-api/pkg/graphql/resolvers"
	"graphql-api/pkg/graphql/scalar"

	"github.com/graphql-go/graphql"
)

/*
Contact Types
*/
var ContactGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Contact",
	Fields: graphql.Fields{
		"contact_id": &graphql.Field{Type: scalar.Int64Type},
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
		"contacts":   &graphql.Field{Type: graphql.NewList(ContactGraphQLType)},
		"pagination": &graphql.Field{Type: PaginationGraphQLType},
		// Add field here
	},
})

type ContactQueries struct {
	Gets          func(string) ([]*models.ContactModel, error)         `json:"gets"`
	GetPagination func(string) (*models.ContactPaginationModel, error) `json:"getPagination"`
}

// Define the ContactQueries type
var ContactQueriesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ContactQueries",
	Fields: graphql.Fields{
		"gets": &graphql.Field{
			Type: graphql.NewList(ContactGraphQLType),
			Args: SearhTextQueryArgument,
			Resolve:// auth.AuthorizeResolverClean("contacts.gets", monitoring.TraceResolver( cache.GetCacheResolver(resolvers.GetContactResolve))),
			resolvers.GetContactResolve,
		},
		"getPagination": &graphql.Field{
			Type:    ContactPaginationGraphQLType,
			Args:    SearhTextPaginationQueryArgument,
			Resolve: auth.AuthorizeResolverClean("contacts.getPagination", cache.GetCacheResolver(resolvers.GetContactsPaginationResolve)),
		},
		"getById": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    IdArgument,
			Resolve: auth.AuthorizeResolverClean("contacts.getById", cache.GetCacheResolver(resolvers.GetContactByIdResolve)),
		},
	},
})

type ContactMutations struct {
	CreateContact  func(map[string]interface{}) (*models.ContactModel, error)   `json:"createContact"`
	CreateContacts func(map[string]interface{}) ([]*models.ContactModel, error) `json:"createContacts"`
	UpdateContact  func(map[string]interface{}) (*models.Status, error)         `json:"updateContact"`
	DeleteContact  func(int) (*models.Status, error)                            `json:"deleteContact"`
}

var StatusGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Status",
	Fields: graphql.Fields{
		"status_id": &graphql.Field{Type: graphql.Int},
		"status":    &graphql.Field{Type: graphql.String},
		"message":   &graphql.Field{Type: graphql.String},
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
			Resolve: auth.AuthorizeResolverClean("contactMutations.createContact", resolvers.CreateContactResolve),
		},
		"createContacts": &graphql.Field{
			Type:    StatusGraphQLType,
			Args:    CreateContactsArgument,
			Resolve: resolvers.CreateContactsResolve,
		},
		"updateContact": &graphql.Field{
			Type:    ContactGraphQLType,
			Args:    UpdateContactArgument,
			Resolve: auth.AuthorizeResolverClean("contactMutations.updateContact", resolvers.UpdateContactResolve),
		},
		"deleteContact": &graphql.Field{
			Type:    StatusGraphQLType,
			Args:    IdArgument,
			Resolve: auth.AuthorizeResolverClean("contactMutations.deleteContact", resolvers.DeleteContactResolve),
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
		"created_at":      &graphql.Field{Type: graphql.DateTime},
		"created_by":      &graphql.Field{Type: graphql.String},
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
		"ticket_id":  &graphql.Field{Type: graphql.Int},
		"type":       &graphql.Field{Type: graphql.String},
		"price":      &graphql.Field{Type: graphql.Float},
		"event_id":   &graphql.Field{Type: graphql.Int},
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
		"page":        &graphql.Field{Type: graphql.Int},
		"pageSize":    &graphql.Field{Type: graphql.Int},
		"totalPages":  &graphql.Field{Type: graphql.Int},
		"totalItems":  &graphql.Field{Type: graphql.Int},
		"hasNext":     &graphql.Field{Type: graphql.Boolean},
		"hasPrevious": &graphql.Field{Type: graphql.Boolean},
		// Add field here
	},
})

var SubscribeMessageGraphQLType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SubscribeMessage",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"content":   &graphql.Field{Type: graphql.String},
		"timestamp": &graphql.Field{Type: graphql.DateTime},
		// Add field here
	},
})

// Define the ContactQueries type
var ContactSubscriptionsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ContactSubscriptions",
	Fields: graphql.Fields{
		"contactCreated": &graphql.Field{
			Type:    SubscribeMessageGraphQLType,
			Resolve: subscription.MessageSubscribe,
		},
	},
})

type ContactSubscription struct {
	ContactCreated func(map[string]interface{}) (*models.ContactModel, error) `json:"contactCreated"`
}

var BidingRoomGraphQLType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "BidingRoom",
		Fields: graphql.Fields{
			"room_id": &graphql.Field{
				Type: graphql.Int,
			},
			"start_date": &graphql.Field{
				Type: graphql.DateTime,
			},
			"end_date": &graphql.Field{
				Type: graphql.DateTime,
			},
			"product_name": &graphql.Field{
				Type: graphql.String,
			},
			"floor_price": &graphql.Field{
				Type: graphql.Float,
			},
			"product_image": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var BidingGraphQLType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Biding",
		Fields: graphql.Fields{
			"bid_id": &graphql.Field{
				Type: graphql.Int,
			},
			"room_id": &graphql.Field{
				Type: graphql.Int,
			},
			"bidder": &graphql.Field{
				Type: graphql.String,
			},
			"bid_price": &graphql.Field{
				Type: graphql.Float,
			},
			"bid_time": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

type BidingQueries struct {
	GetRoomById   func(int) (*models.BidingRoomModel, error) `json:"getRoomById"`
	GetTop5Biding func(int) ([]*models.BidingModel, error)   `json:"getTop5Biding"`
}

// Define the TicketQueries type
var BidingQueriesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "BidingQueries",
	Fields: graphql.Fields{
		"getRoomById": &graphql.Field{
			Type:    BidingRoomGraphQLType,
			Args:    IdArgument,
			Resolve: resolvers.GetBidingRoomByIdResolve,
		},
		"getTop5Biding": &graphql.Field{
			Type:    graphql.NewList(BidingGraphQLType),
			Args:    IdArgument,
			Resolve: resolvers.GetBidingTop5Resolve,
		},
	},
})

type BidingMutations struct {
	CreateBiding  func(map[string]interface{}) (*models.BidingModel, error)   `json:"createBiding"`
}


// Define the ContactMutations type
var BidingMutationsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "BidingMutations",
	Fields: graphql.Fields{
		"createBiding": &graphql.Field{
			Type:    BidingGraphQLType,
			Args:    CreateBidingArgument,
			Resolve:  resolvers.CreateBidingResolve,
		},

	},
})


