package models

import (
	"time"
)

// UserRole represents the relationship between a user and a role
type UserRolesModel struct {
	UserRoleID int       `json:"user_role_id"`
	RoleID     int       `json:"role_id"`
	UserID     string    `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
}


