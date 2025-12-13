package tests

import (
	service "backendUAS/app/services/postgres"
	"backendUAS/databases"
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(".env not loaded")
	}

	databases.ConnectToMongo()
	databases.ConnectToPostgres()

	os.Exit(m.Run())
}

func TestLoginSuccess(t *testing.T) {
	app := fiber.New()
	app.Post("/api/v1login", service.LoginService)

	body := []byte(`{
		"email": "admin@gmail.com",
		"password": "password"
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1login",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestLoginEmptyBody(t *testing.T) {
	app := fiber.New()
	app.Post("/api/v1login", service.LoginService)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1login",
		nil,
	)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}


func TestLoginMissingField(t *testing.T) {
	app := fiber.New()
	app.Post("/api/v1login", service.LoginService)

	body := []byte(`{
		"email": "",
		"password": ""
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1login",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

