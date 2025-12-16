package routes

import (
	service "backendUAS/app/services/postgres"
	"backendUAS/middlewares"

	"github.com/gofiber/fiber/v2"
)


func AnalyticsRoute(app *fiber.App) {

    api := app.Group("/api")

    v1 := api.Group("/v1")

    analytics := v1.Group("/reports", middlewares.AuthRequired())
	analytics.Get("/statistics", service.GetReportStatisticsService)
	analytics.Get("/student/:id",service.GetStudentReportService)

}