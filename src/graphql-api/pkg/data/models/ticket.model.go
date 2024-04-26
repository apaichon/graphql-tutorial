package models

import (
	"time"
)

// Ticket represents a Ticket record in the database
type TicketModel struct {
	TicketId int     `json:"ticket_id"`
	Type     string  `json:"type"`
	Price    float64 `json:"price"`
	EventId  int     `json:"event_id"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string `json:"created_by"`
}

type TicketEventModel struct {
	Ticket TicketModel `json:"ticket"`
	Event  EventModel  `json:"event"`
}

// Function to copy a slice of pointers
func CopyEventSlice(src []*EventModel) []EventModel {
	// Create a new slice with the same length and capacity as the source
	dst := make([]EventModel, len(src))

	// Copy pointers from source to destination
	for i, event := range src {
		dst[i] = *event
	}

	return dst
}

// Function to create a map from the Event slice by EventId
func CreateEventMap(events []EventModel) map[int]EventModel {
	eventMap := make(map[int]EventModel)
	for _, event := range events {
		eventMap[event.EventId] = event
	}
	return eventMap
}

// Function to map Tickets with Events by EventId
func MapTicketsWithEvents(tickets []TicketModel, eventMap map[int]EventModel) []TicketEventModel {
	var ticketEvents []TicketEventModel
	for _, ticket := range tickets {
		if event, found := eventMap[ticket.EventId]; found {
			ticketEvents = append(ticketEvents, TicketEventModel{
				Ticket: ticket,
				Event:  event,
			})
		}
	}
	return ticketEvents
}