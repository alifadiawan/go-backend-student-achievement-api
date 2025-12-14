package middlewares

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func CheckPermission(perms map[string]bool, permission string) error {
	if !perms[permission] {
		return errors.New("access denied")
	}
	return nil
}

func Permission(permission string, next func(*fiber.Ctx) error) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		perms, ok := c.Locals("permissions").(map[string]bool)
		if !ok {
			return c.Status(403).JSON(fiber.Map{
				"message": "gagal loading permission",
			})
		}

		if !perms[permission] {
			return c.Status(403).JSON(fiber.Map{
				"message": "maaf, anda gabole akses route ini",
			})
		}

		return next(c)
	}
}
