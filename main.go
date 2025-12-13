package main

import (
	"backendUAS/databases"
	"backendUAS/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// routes
	routes.AuthRoutes(app)
	routes.UserRoutes(app)
	routes.AchievementRoutes(app)
	routes.StudentLecturerRoute(app)

	// databases
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
