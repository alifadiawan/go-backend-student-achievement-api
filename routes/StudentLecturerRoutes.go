package routes

import (
	service "backendUAS/app/services/postgres"
	"backendUAS/middlewares"

	"github.com/gofiber/fiber/v2"
)


func StudentLecturerRoute(app *fiber.App) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	student := v1.Group("/student", middlewares.AuthRequired())
	student.Get("/", service.GetStudentsService)
	student.Get("/:id", service.GetStudentByIDService)
	student.Get("/:id/achievements", service.GetStudentAchievementByIDService)
	student.Get("/:id/advisor", service.UpdateStudentAdvisorService)
	

	// lecturer := v1.Group("/lecturer", middlewares.AuthRequired())
	
}