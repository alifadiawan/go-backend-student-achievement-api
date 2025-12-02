package postgres

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`

	RoleID   uuid.UUID `json:"role_id"`
	RoleName string    `json:role_name`
	IsActive bool      `json:"is_active"`
	Permission []string `json:"permission"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


type UserRequest struct {
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
}
