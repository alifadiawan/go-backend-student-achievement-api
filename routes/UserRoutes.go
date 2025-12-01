package routes

import (
	services "backendUAS/app/services/postgres"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	users := v1.Group("/users")
	users.Get("/", services.GetAllUserService)

}
