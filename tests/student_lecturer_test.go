package tests

import (
	service "backendUAS/app/services/postgres"
	"backendUAS/middlewares"

	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var studentLecturerApp *fiber.App
var userToken string

func setupStudentLecturerTest() {
	studentLecturerApp = fiber.New()
	api := studentLecturerApp.Group("/api")
	v1 := api.Group("/v1")

	// student routes
	student := v1.Group("/student", middlewares.AuthRequired())
	student.Get("/", service.GetStudentsService)
	student.Get("/:id", service.GetStudentByIDService)
	student.Get("/:id/achievements", service.GetStudentAchievementByIDService)
	student.Get("/:id/advisor", service.UpdateStudentAdvisorService)

	// lecturer routes
	lecturer := v1.Group("/lecturer", middlewares.AuthRequired())
	lecturer.Get("/", service.GetLecturerService)

	// Dummy token
	userToken = "Bearer " + "fake-user-jwt"
}

// helper request
func requestSL(method, url string, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, nil)
	req.Header.Set("Authorization", token)
	resp, _ := studentLecturerApp.Test(req, -1)
	rec := httptest.NewRecorder()
	rec.Result().StatusCode = resp.StatusCode
	return rec
}

func TestStudentLecturerRoutes(t *testing.T) {
	setupStudentLecturerTest()

	t.Run("Get all students", func(t *testing.T) {
		resp := requestSL("GET", "/api/v1/student", userToken)
		if resp.Code != 200 {
			t.Fatalf("expected 200, got %d", resp.Code)
		}
	})

	t.Run("Get student by ID", func(t *testing.T) {
		id := uuid.New()
		resp := requestSL("GET", fmt.Sprintf("/api/v1/student/%s", id.String()), userToken)
		if resp.Code != 200 {
			t.Fatalf("expected 200, got %d", resp.Code)
		}
	})

	t.Run("Get student achievements", func(t *testing.T) {
		id := uuid.New()
		resp := requestSL("GET", fmt.Sprintf("/api/v1/student/%s/achievements", id.String()), userToken)
		if resp.Code != 200 {
			t.Fatalf("expected 200, got %d", resp.Code)
		}
	})

	t.Run("Get student advisor", func(t *testing.T) {
		id := uuid.New()
		resp := requestSL("GET", fmt.Sprintf("/api/v1/student/%s/advisor", id.String()), userToken)
		if resp.Code != 200 {
			t.Fatalf("expected 200, got %d", resp.Code)
		}
	})

	t.Run("Get all lecturers", func(t *testing.T) {
		resp := requestSL("GET", "/api/v1/lecturer", userToken)
		if resp.Code != 200 {
			t.Fatalf("expected 200, got %d", resp.Code)
		}
	})
}
