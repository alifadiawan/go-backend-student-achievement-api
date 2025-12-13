package main

import (
	"backendUAS/databases"
	"backendUAS/routes"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/swaggo/fiber-swagger"
	_ "backendUAS/docs"

	"github.com/joho/godotenv"


	"github.com/gofiber/fiber/v2/middleware/cors"
)

// @title Backend UAS - Sistem Kelola Prestasi Mahasiswa
// @version 1.0
// @description Sebuah sistem backend berbasis go yang mengelola prestasi mahasiswa

// @host localhost:3000

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by your JWT token.
func main() {
	app := fiber.New()
	app.Use(cors.New())

	// swagger
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// routes
	routes.AuthRoutes(app)
	routes.UserRoutes(app)
	routes.AchievementRoutes(app)
	routes.StudentLecturerRoute(app)

	// databases
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env not loaded")
	}

	_, err := databases.ConnectToPostgres()
	if err != nil {
		log.Fatal("Postgres tidak Connect: ", err)
	}

	_, _, _, err = databases.ConnectToMongo()
	if err != nil {
		log.Fatal("MongoDB Tidak Connect:", err)
	}

	log.Fatal(app.Listen(":3000"))

}
