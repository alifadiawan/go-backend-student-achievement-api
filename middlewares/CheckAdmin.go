package middlewares

import (

	"github.com/gofiber/fiber/v2"
)

func OnlyAdmin(c *fiber.Ctx) error {

	role := c.Locals("role")

	if role == nil || role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"message" : "endpoint forbidden",
		})
	}

	return c.Next()

}
