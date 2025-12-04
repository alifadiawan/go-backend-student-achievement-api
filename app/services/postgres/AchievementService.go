package postgres

import (
	models "backendUAS/app/models/postgres"
	repo "backendUAS/app/repositories/postgres"

	"github.com/gofiber/fiber/v2"
)

func GetAllAchievementService(c *fiber.Ctx) error {

	var result []models.AchievementGabung
	var err error

	result, err = repo.GetAllAchievementRepo()
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "gagal mengambil data prestasi mahasiswa",
			})
		}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"count":  len(result),
		"data":   result,
	})
}

func GetAchievementByIDService(c *fiber.Ctx) error {
	studentID := c.Params("student_id")
	if studentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user id tidak valid",
		})
	}

	achievements, err := repo.GetAchievementByIDRepo(studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "tidak bisa mengambil data prestasi",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   achievements,
	})
}
