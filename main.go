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

	// databases
	databases.ConnectToPostgres(app)

	log.Fatal(app.Listen(":3000"))

}
