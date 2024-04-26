package test

import (
	"graphql-api/internal/contact"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContacts(t *testing.T) {
	repo:=  contact.NewContactRepo()
	contacts, err := repo.GetContactsBySearchText("", 10, 0)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	assert.Greater(t,len(contacts),0)
}