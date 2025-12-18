package middlewares

import (
	repoPG "backendUAS/app/repositories/postgres"
	"backendUAS/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token akses diperlukan",
			})
		}

		var token string
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			token = authHeader
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token tidak valid",
			})
		}

		// Simpan data user ke context
		c.Locals("user_id", claims.UserID)
		c.Locals("student_id", claims.StudentID)
		c.Locals("lecturer_id", claims.LecturerID)
		c.Locals("nim", claims.NIM)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		userID := claims.UserID
		if userID == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "ID pengguna tidak valid",
			})
		}

		permissions, err := repoPG.LoadPermissions(userID)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Gagal memuat izin pengguna",
			})
		}

		c.Locals("permissions", permissions)

		return c.Next()
	}
}
