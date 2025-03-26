// models/user_response.go
package models

type PermissionEntry struct {
	ServiceID  string `json:"service_id"`
	Service    string `json:"service"`
	Permission string `json:"permission"`
}

type UserResponse struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Email       string            `json:"email"`
	Role        string            `json:"role"`
	Permissions []PermissionEntry `json:"permissions"`
}
