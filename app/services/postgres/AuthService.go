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

	refreshToken, err := utils.RefreshToken(models.User{})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak bisa generate token refresh",
			"status":  err.Error(),
		})
	}

	loginResponse := models.LoginResponse{
		ID:          User.ID,
		Email:       User.Email,
		StudentID:   User.StudentID,
		NIM:         User.NIM,
		Username:    User.Username,
		FullName:    User.FullName,
		Role:        User.Role,
		Permissions: User.Permissions,
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

func Profile(c *fiber.Ctx) error {
	UserIDJWT := c.Locals("user_id")
	UserID := UserIDJWT.(string)

	// fmt.Println(UserID)

	if UserID == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "id tidak ditemukan",
		})
	}

	var UserProfile *models.User

	UserProfile, err := repositories.GetProfile(UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": err,
		})
	}

	return c.JSON(UserProfile)

}

func LogoutService(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "logout successful",
	})
}
