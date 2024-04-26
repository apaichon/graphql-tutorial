package resolvers

import (
	"graphql-api/internal/ticket"
	"github.com/graphql-go/graphql"
)


func GetTicketsResolve(params graphql.ResolveParams) (interface{}, error) {
	// Update limit and offset if provided
	limit, ok := params.Args["limit"].(int)
	if !ok {
		limit = 10
	}
	offset, ok := params.Args["offset"].(int)

	if !ok {
		offset = 0
	}

	searchText, ok := params.Args["searchText"].(string)

	if !ok {
		searchText = ""
	}

	ticketRepo := ticket.NewTicketRepo()

	// Fetch tickets from the database
	tickets, err := ticketRepo.GetTicketsBySearchText(searchText, limit, offset)
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func GetTicketByIdResolve(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(int)
	ticketRepo := ticket.NewTicketRepo()

	// Fetch tickets from the database
	ticket, err := ticketRepo.GetTicketByID(id)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

func GetTicketEventsResolve(params graphql.ResolveParams) (interface{}, error) {
	// Update limit and offset if provided
	limit, ok := params.Args["limit"].(int)
	if !ok {
		limit = 10
	}
	offset, ok := params.Args["offset"].(int)

	if !ok {
		offset = 0
	}

	searchText, ok := params.Args["searchText"].(string)

	if !ok {
		searchText = ""
	}

	ticketRepo := ticket.NewTicketRepo()

	// Fetch tickets from the database
	tickets, err := ticketRepo.GetTicketEventsBySearchText(searchText, limit, offset)
	if err != nil {
		return nil, err
	}
	return tickets, nil
}