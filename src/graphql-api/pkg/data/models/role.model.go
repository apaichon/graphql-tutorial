package models

import "time"

// Role represents a role in the system
type RoleModel struct {
	RoleID     int       `json:"role_id"`
	RoleName   string    `json:"role_name"`
	RoleDesc   string    `json:"role_desc,omitempty"`
	IsSuperAdmin bool    `json:"is_super_admin"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	StatusID   int       `json:"status_id"`
}

