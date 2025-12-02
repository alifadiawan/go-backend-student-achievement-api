package routes

import (
	services "backendUAS/app/services/postgres"
	middleware "backendUAS/middlewares"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	users := v1.Group("/users", middleware.AuthRequired(), middleware.OnlyAdmin)
	users.Get("/", middleware.OnlyAdmin, services.GetAllUserService)
	users.Get("/:id", middleware.OnlyAdmin, services.GetUsersByIdService)
	users.Post("/", middleware.OnlyAdmin, services.StoreUserService)

}
