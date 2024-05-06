package models

import (
	"time"
)
// RolePermission represents the permissions granted to a role
type RolePermissionsModel struct {
	RolePermissionID   int       `json:"role_permission_id"`
	RolePermissionDesc string    `json:"role_permission_desc"`
	ResourceTypeID     int       `json:"resource_type_id"`
	ResolveName        string    `json:"resolve_name"`
	CanExecute         bool      `json:"can_execute"`
	CanRead            bool      `json:"can_read"`
	CanWrite           bool      `json:"can_write"`
	CanDelete          bool      `json:"can_delete"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          string    `json:"created_by"`
}