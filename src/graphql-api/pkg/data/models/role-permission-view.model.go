package models

import "time"

// UserRolePermissions represents the vw_user_role_permissions view
type UserRolePermissionsViewModel struct {
	UserID                    int       `json:"user_id"`
	UserName                  string    `json:"user_name"`
	RoleID                    int       `json:"role_id"`
	RoleName                  string    `json:"role_name"`
	RoleDesc                  *string   `json:"role_desc,omitempty"` // Nullable field
	IsSuperAdmin              bool      `json:"is_super_admin"`
	UserRoleID                int       `json:"user_role_id"`
	UserRoleCreatedAt         time.Time `json:"user_role_created_at"`
	RolePermissionID          *int      `json:"role_permission_id,omitempty"` // Nullable field
	ResourceTypeID            *int      `json:"resource_type_id,omitempty"`  // Nullable field
	ResolveName               *string   `json:"resolve_name,omitempty"` // Nullable field
	CanExecute                *bool     `json:"can_execute,omitempty"` // Nullable field
	CanRead                   *bool     `json:"can_read,omitempty"`   // Nullable field
	CanWrite                  *bool     `json:"can_write,omitempty"`  // Nullable field
	CanDelete                 *bool     `json:"can_delete,omitempty"` // Nullable field
	RolePermissionCreatedAt   *time.Time `json:"role_permission_created_at,omitempty"` // Nullable field
}

// UserRolePermissions represents the vw_user_role_permissions view
type UserRoleResolvePermissionsModel struct {
	UserID                    int       `json:"user_id"`
	IsSuperAdmin              bool      `json:"is_super_admin"`
	ResolveName               *string   `json:"resolve_name,omitempty"` // Nullable field
	CanExecute                *bool     `json:"can_execute,omitempty"` // Nullable field
}
