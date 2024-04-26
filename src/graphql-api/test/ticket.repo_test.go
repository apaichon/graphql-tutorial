package test

import (
	"graphql-api/internal/ticket"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTickets(t *testing.T) {
	repo:=  ticket.NewTicketRepo()
	tickets, err := repo.GetTicketsBySearchText("", 10, 0)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	assert.Greater(t,len(tickets),0)
}

func TestGetTicketEvents(t *testing.T) {
	repo:=  ticket.NewTicketRepo()
	tickets, err := repo.GetTicketEventsBySearchText("", 10, 0)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	assert.Greater(t,len(tickets),0)
}