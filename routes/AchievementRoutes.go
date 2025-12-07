package routes

import (
	services "backendUAS/app/services/postgres"
	middleware "backendUAS/middlewares"
	"github.com/gofiber/fiber/v2"
)

func AchievementRoutes(app *fiber.App) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	users := v1.Group("/achievement", middleware.AuthRequired())
	users.Get("/", services.GetAllAchievementService)
	users.Get("/:AchievementID", services.GetAchievementByIDService)
	users.Post("/", services.AddAchievementService)
	users.Delete("/:AchievementID", services.DeleteAchievementService)

}
