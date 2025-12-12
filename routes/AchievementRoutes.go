package routes

import (
	services "backendUAS/app/services/postgres"
	servicesMongo "backendUAS/app/services/mongo"
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
	users.Delete("/:achievement_references_id", services.DeleteAchievementService)


	users.Post("/submit/:achievement_references_id", services.SubmitAchievementService)
	users.Post("/approve/:achievement_references_id", services.ApproveAchievmentService)
	users.Post("/verify/:achievement_references_id", services.VerifyAchievementService)
	users.Post("/reject/:achievement_references_id", services.RejectAchievementService)
	users.Post("/attachments/:achievement_references_id", servicesMongo.UploadAchievementService)	

}
