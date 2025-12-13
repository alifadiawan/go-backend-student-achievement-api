package tests

import (
	mongoService "backendUAS/app/services/mongo"
	postgresService "backendUAS/app/services/postgres"
	"backendUAS/middlewares"
	
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var achievementApp *fiber.App
var studentToken string

func setupAchievementTest() {
	app = fiber.New()
	api := app.Group("/api")
	v1 := api.Group("/v1")

	ach := v1.Group("/achievement", middlewares.AuthRequired())
	ach.Post("/", postgresService.AddAchievementService)
	ach.Delete("/:achievement_references_id", postgresService.DeleteAchievementService)
	ach.Post("/submit/:achievement_references_id", postgresService.SubmitAchievementService)
	ach.Post("/approve/:achievement_references_id", postgresService.ApproveAchievmentService)
	ach.Post("/reject/:achievement_references_id", postgresService.RejectAchievementService)
	ach.Post("/attachments/:achievement_references_id", mongoService.UploadAchievementService)

	// Dummy student token
	studentToken = "Bearer " + "fake-student-jwt"
}

// helper request
func requestAchievement(method, url string, body interface{}, token string, isMultipart bool) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	var req *http.Request

	if isMultipart {
		writer := multipart.NewWriter(&buf)
		if m, ok := body.(map[string]string); ok {
			for k, v := range m {
				writer.WriteField(k, v)
			}
		}
		writer.Close()
		req = httptest.NewRequest(method, url, &buf)
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else if body != nil {
		json.NewEncoder(&buf).Encode(body)
		req = httptest.NewRequest(method, url, &buf)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, url, nil)
	}

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, _ := app.Test(req, -1)
	rec := httptest.NewRecorder()
	rec.Result().StatusCode = resp.StatusCode
	return rec
}

func TestAchievementRoutes(t *testing.T) {
	setupAchievementTest()

	// dummy achievement references ID
	achievementRefID := uuid.New().String()

	t.Run("Add Achievement", func(t *testing.T) {
		payload := map[string]interface{}{
			"student_id":       uuid.New().String(),
			"achievement_type": "academic",
			"title":            "Test Achievement",
			"description":      "Test Description",
			"details": map[string]interface{}{
				"score": 100,
			},
		}
		resp := request("POST", "/api/v1/achievement", payload, studentToken)
		if resp.Code != 201 && resp.Code != 200 {
			t.Fatalf("expected 201 or 200, got %d", resp.Code)
		}
	})

	t.Run("Submit Achievement", func(t *testing.T) {
		resp := request("POST", fmt.Sprintf("/api/v1/achievement/submit/%s", achievementRefID), nil, studentToken)
		if resp.Code != 200 {
			t.Fatalf("expected 200, got %d", resp.Code)
		}
	})

	t.Run("Approve Achievement", func(t *testing.T) {
		resp := request("POST", fmt.Sprintf("/api/v1/achievement/approve/%s", achievementRefID), nil, studentToken)
		if resp.Code != 200 {
			t.Fatalf("expected 200, got %d", resp.Code)
		}
	})

	t.Run("Reject Achievement", func(t *testing.T) {
		payload := map[string]interface{}{
			"rejection_note": "Not valid",
		}
		resp := request("POST", fmt.Sprintf("/api/v1/achievement/reject/%s", achievementRefID), payload, studentToken)
		if resp.Code != 200 {
			t.Fatalf("expected 200, got %d", resp.Code)
		}
	})

	t.Run("Upload Attachment", func(t *testing.T) {
		// simulasikan multipart form tanpa file actual
		payload := map[string]string{
			"student_id": uuid.New().String(),
		}
		resp := request("POST", fmt.Sprintf("/api/v1/achievement/attachments/%s", achievementRefID), payload, studentToken)
		if resp.Code != 200 {
			t.Fatalf("expected 200, got %d", resp.Code)
		}
	})

	t.Run("Delete Achievement", func(t *testing.T) {
		resp := request("DELETE", fmt.Sprintf("/api/v1/achievement/%s", achievementRefID), nil, studentToken)
		if resp.Code != 200 && resp.Code != 204 {
			t.Fatalf("expected 200 or 204, got %d", resp.Code)
		}
	})
}
