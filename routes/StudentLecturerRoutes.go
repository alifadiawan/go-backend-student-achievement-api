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
	student.Get("/", middlewares.Permission("student:view",service.GetStudentsService))
	student.Get("/:id", middlewares.Permission("student:view",service.GetStudentByIDService))
	student.Get("/:id/achievements", middlewares.Permission("student:view",service.GetStudentAchievementByIDService))
	student.Get("/:id/advisor",middlewares.Permission("student:view", service.UpdateStudentAdvisorService))
	

	lecturer := v1.Group("/lecturer", middlewares.AuthRequired())
	lecturer.Get("/", middlewares.Permission("lecturer:view",service.GetLecturerService))
	
}