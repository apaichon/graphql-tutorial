package resolvers

import (
	"graphql-api/internal/event"
	"github.com/graphql-go/graphql"
)

func GetEventsResolve(params graphql.ResolveParams) (interface{}, error) {
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
	eventRepo :=event.NewEventRepo()

	// Fetch contacts from the database
	events, err := eventRepo.GetEventsBySearchText(searchText, limit, offset)
	if err != nil {
		return nil, err
	}
	return events, nil
}


func GetEventByIdResolve(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(int)
	eventRepo := event.NewEventRepo()

	// Fetch tickets from the database
	event, err := eventRepo.GetEventByID(id)
	if err != nil {
		return nil, err
	}
	return event, nil
}