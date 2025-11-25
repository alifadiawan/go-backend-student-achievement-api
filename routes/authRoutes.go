package routes

import (
	services "backendUAS/app/services/postgres"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {

    api := app.Group("/api")

    v1 := api.Group("/v1")

    auth := v1.Group("/auth")
	auth.Post("/login", services.LoginService)
	

}
