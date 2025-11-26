package postgres

import (
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Email    string `json:username`
	Password string `json:password`
}

type LoginResponse struct {
	Email    string `json:email`
	Username string `json:username`
	FullName string `json:full_name`
	RoleID   string `json:role_id`
	RoleName   string `json:role_name`
	IsActive bool   `json:is_active`
	Token    string `json:token`
}

type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
