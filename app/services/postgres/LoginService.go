package postgres

import (
	models "backendUAS/app/models/postgres"
	repositories "backendUAS/app/repositories/postgres"
	"backendUAS/utils"

	"github.com/gofiber/fiber/v2"
)

func LoginService(c *fiber.Ctx) error {

	var Request models.LoginRequest

	err := c.BodyParser(&Request)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "body tidak valid",
			"error":   err.Error(),
		})
	}

	if Request.Email == "" || Request.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "email dan password wajib diisi",
		})
	}

	User, err := repositories.Authenticate(Request.Email, Request.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "request tidak valid",
			"error":   err.Error(),
		})
	}

	token, err := utils.CreateToken(*User)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak bisa generate token",
			"error":   err.Error(),
		})
	}

	refreshToken, err := utils.RefreshToken(*User)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak bisa generate token refresh",
			"status":  err.Error(),
		})
	}

	loginResponse := models.LoginResponse{
		ID: User.ID.String(),
		Email: User.Email,
		Username: User.Username,
		FullName: User.FullName,
		RoleName: User.RoleName,
	}

	response := models.ApiResponse{
		Status: "success",
		Data: fiber.Map{
			"token":        token,
			"tokenRefresh": refreshToken,
			"user":         loginResponse,
		},
	}

	return c.JSON(response)

}
