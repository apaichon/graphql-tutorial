package models

import (
	"time"
)

// Contact represents a Contact record in the database
type ContactModel struct {
	ContactId int64 `json:"contact_id"`
	Name string `json:"name"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	GenderId int `json:"gender_id"`
	Dob time.Time `json:"dob"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Address string `json:"address"`
	PhotoPath string `json:"photo_path"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string `json:"created_by"`
}

type ContactPaginationModel struct {
	Contacts []*ContactModel `json:"contacts"`
	Pagination *PaginationModel `json:"pagination"`
}

