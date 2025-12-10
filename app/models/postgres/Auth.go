package postgres

import (
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID          string   `json:"id"`
	Email       string   `json:"email"`
	StudentID   *string   `json:"-"`
	NIM         *string   `json:"-"`
	LecturerID  *string  `json:"lecturer_id,omitempty"`
	Username    string   `json:"username"`
	FullName    string   `json:"full_name"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type ApiResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type JWTClaims struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	StudentID string `json:"student_id"`
	NIM       string `json:"nim"`
	LecturerID string `json:"lecturer_id,omitempty"`
	jwt.RegisteredClaims
}
