package mongo

import (
	mongoModels "backendUAS/app/models/mongo"
	mongoRepo "backendUAS/app/repositories/mongo"

	"github.com/gofiber/fiber/v2"
)

func UploadAchievementService(c *fiber.Ctx) error {

	achievement_references_id := c.Params("achievement_references_id")
	if achievement_references_id == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "achievement references id tidak valid",
		})
	}

	file, err := c.FormFile("attachment")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "file attachment tidak ditemukan",
			"error":   err.Error(),
		})
	}

	req := mongoModels.AchievementAttachementRequest{
		AchievementReferencesID: achievement_references_id,
		Attachment:              file,
	}

	result, err := mongoRepo.UploadAchievementRepo(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "tidak dapat upload attachement",
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "berhasil upload attachment",
		"status":  "success",
		"data": result,
	})

}
