package tests

import (
	services "backendUAS/app/services/postgres"
	middleware "backendUAS/middlewares"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var app *fiber.App
var adminToken string

func setup() {
	app = fiber.New()
	api := app.Group("/api")
	v1 := api.Group("/v1")

	users := v1.Group("/users", middleware.AuthRequired())

	users.Get("/", middleware.Permission("user:manage", services.GetAllUserService))
	users.Get("/:user_id", services.GetUsersByIdService)
	users.Post("/", middleware.Permission("user:manage", services.StoreUserService))
	users.Put("/:user_id", middleware.Permission("user:manage", services.UpdateUserService))
	users.Put("/role/:user_id", middleware.Permission("user:manage", services.UpdateUserRoleService))
	users.Delete("/:id", middleware.Permission("user:manage", services.DeleteUserService))

	// Dummy admin token
	adminToken = "Bearer " + "fake-admin-jwt"
}

func request(method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}

	req := httptest.NewRequest(method, url, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, _ := app.Test(req, -1)
	rec := httptest.NewRecorder()
	rec.Result().StatusCode = resp.StatusCode
	return rec
}

func TestUsersRoutes(t *testing.T) {
	setup()

	t.Run("Get all users", func(t *testing.T) {
		resp := request("GET", "/api/v1/users", nil, adminToken)
		if resp.Code != 200 {
			t.Fatalf("expected 200, got %d", resp.Code)
		}
	})

	t.Run("Get user by ID", func(t *testing.T) {
		id := uuid.New()
		resp := request("GET", fmt.Sprintf("/api/v1/users/%s", id.String()), nil, adminToken)
		if resp.Code != 200 {
			t.Fatalf("expected 200, got %d", resp.Code)
		}
	})

	t.Run("Store user", func(t *testing.T) {
		payload := map[string]interface{}{
			"username":  "testuser",
			"email":     "testuser@example.com",
			"full_name": "Test User",
			"role_id":   uuid.New(),
			"password":  "password123",
		}
		resp := request("POST", "/api/v1/users", payload, adminToken)
		if resp.Code != 200 {
			t.Fatalf("expected 201, got %d", resp.Code)
		}
	})

	t.Run("Update user", func(t *testing.T) {
		id := uuid.New()
		payload := map[string]interface{}{
			"full_name": "Updated Name",
		}
		resp := request("PUT", fmt.Sprintf("/api/v1/users/%s", id.String()), payload, adminToken)
		if resp.Code != 200 {
			t.Fatalf("expected 200, got %d", resp.Code)
		}
	})

	t.Run("Update user role", func(t *testing.T) {
		id := uuid.New()
		payload := map[string]interface{}{
			"role_name": "lecturer",
		}
		resp := request("PUT", fmt.Sprintf("/api/v1/users/role/%s", id.String()), payload, adminToken)
		if resp.Code != 200 {
			t.Fatalf("expected 200, got %d", resp.Code)
		}
	})

	t.Run("Delete user", func(t *testing.T) {
		id := uuid.New()
		resp := request("DELETE", fmt.Sprintf("/api/v1/users/%s", id.String()), nil, adminToken)
		if resp.Code != 200 {
			t.Fatalf("expected 200, got %d", resp.Code)
		}
	})
}
