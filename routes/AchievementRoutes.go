package routes

import (
	servicesMongo "backendUAS/app/services/mongo"
	services "backendUAS/app/services/postgres"
	middleware "backendUAS/middlewares"
	"github.com/gofiber/fiber/v2"
)

func AchievementRoutes(app *fiber.App) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	users := v1.Group("/achievements", middleware.AuthRequired())
	users.Get("/", middleware.Permission("achievement:view", services.GetAllAchievementService))
	users.Get("/:AchievementID", middleware.Permission("achievement:view", services.GetAchievementByIDService))
	users.Post("/", middleware.Permission("achievement:create",services.AddAchievementService))
	users.Delete("/:achievement_references_id", middleware.Permission("achievement:delete",services.DeleteAchievementService))

	users.Post("/submit/:achievement_references_id", middleware.Permission("achievement:submit", services.SubmitAchievementService))
	users.Post("/approve/:achievement_references_id",middleware.Permission("achievement:approve",  services.ApproveAchievmentService))
	users.Post("/verify/:achievement_references_id", middleware.Permission("achievement:verify", services.VerifyAchievementService))
	users.Post("/reject/:achievement_references_id", middleware.Permission("achievement:reject", services.RejectAchievementService))
	users.Post("/attachments/:achievement_references_id", middleware.Permission("achievement:attachment", servicesMongo.UploadAchievementService))
	users.Get("/history/:achievement_references_id", middleware.Permission("achievement:history", services.HistoryAchievementService))

}
