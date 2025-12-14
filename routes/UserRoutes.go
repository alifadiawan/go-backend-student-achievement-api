package routes

import (
	services "backendUAS/app/services/postgres"
	middleware "backendUAS/middlewares"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	users := v1.Group("/users", middleware.AuthRequired())
	users.Get("/", middleware.Permission("user:manage", services.GetAllUserService))
	users.Get("/:user_id", services.GetUsersByIdService)
	users.Post("/", middleware.Permission("user:manage", services.StoreUserService))
	users.Put("/:user_id", middleware.Permission("user:manage", services.UpdateUserService))
	users.Put("/role/:user_id", middleware.Permission("user:manage", services.UpdateUserRoleService))
	users.Delete("/:id", middleware.Permission("user:manage", services.DeleteUserService))

}
