package postgres

import (
	models "backendUAS/app/models/postgres"
	repo "backendUAS/app/repositories/postgres"

	"github.com/gofiber/fiber/v2"
)

func GetAllUserService(c *fiber.Ctx) error {

	var Users []models.User

	Users, err := repo.GetAllUserRepository()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
		"data":    Users,
	})

}
