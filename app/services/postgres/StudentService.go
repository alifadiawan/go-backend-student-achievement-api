package postgres

import (
	modelsMongo "backendUAS/app/models/postgres"
	repoMongo "backendUAS/app/repositories/mongo"
	repoPG "backendUAS/app/repositories/postgres"

	"github.com/gofiber/fiber/v2"
)

func GetStudentsService(c *fiber.Ctx) error {

	// Ambil role dari JWT
	role, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"message": "invalid token",
		})
	}

	switch role {
	case "admin":
		data, err := repoPG.GetAllStudentsRepo()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "failed to fetch students",
			})
		}

		return c.JSON(fiber.Map{
			"data": data,
		})

	case "dosen":
		lecturerID, ok := c.Locals("lecturer_id").(string)
		if !ok || lecturerID == "" {
			return c.Status(401).JSON(fiber.Map{
				"message": "invalid lecturer context",
			})
		}

		data, err := repoPG.GetStudentsByAdvisorRepo(lecturerID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "failed to fetch advisees",
			})
		}

		return c.JSON(fiber.Map{
			"data": data,
		})

	case "mahasiswa":
		studentID, ok := c.Locals("student_id").(string)
		if !ok || studentID == "" {
			return c.Status(401).JSON(fiber.Map{
				"message": "invalid student context",
			})
		}

		data, err := repoPG.GetStudentByIDRepo(studentID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "failed to fetch student data",
			})
		}

		if data == nil {
			return c.Status(404).JSON(fiber.Map{
				"message": "student not found",
			})
		}

		return c.JSON(fiber.Map{
			"data": data,
		})

	default:
		return c.Status(403).JSON(fiber.Map{
			"message": "forbidden",
		})
	}
}

func GetStudentAchievementByIDService(c *fiber.Ctx) error {
	studentID := c.Params("id")

	if studentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "student id is required",
		})
	}

	studentDetail, err := repoPG.GetStudentByIDRepo(studentID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak bisa ambil student",
			"error": err.Error(),
		})
	}

	achievements, err := repoPG.GetAchievementsByStudentIDRepo(studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get achievements",
			"error":   err.Error(),
		})
	}

	var result []modelsMongo.AchievementMongo

	for _, a := range achievements {
		detail, err := repoMongo.GetAchievementDetailByIdRepo(a.MongoId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to get achievement detail",
				"error":   err.Error(),
			})
		}

		result = append(result, modelsMongo.AchievementMongo{
			Details:     detail,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"student_detail": studentDetail,
		"achievement": achievements,
		"achievement_detail": result,
	})
}


func GetStudentsByAdvisorService(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "dosen" {
		return c.Status(403).JSON(fiber.Map{
			"message": "only lecturer allowed",
		})
	}

	lecturerID := c.Locals("lecturer_id").(string)
	if lecturerID == "" {
		return c.Status(401).JSON(fiber.Map{
			"message": "invalid lecturer context",
		})
	}

	data, err := repoPG.GetStudentsByAdvisorRepo(lecturerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Maaf, tidak bias ambil student",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

func GetStudentByIDService(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role == "mahasiswa" {
		return c.Status(403).JSON(fiber.Map{
			"message": "maaf, tidak boleh akses route ini",
		})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "student id required",
		})
	}

	data, err := repoPG.GetStudentByIDRepo(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to fetch student",
		})
	}

	if data == nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "student not found",
		})
	}

	return c.JSON(fiber.Map{
		"student_detail": data,
	})
}

func UpdateStudentAdvisorService(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"message": "only admin allowed",
		})
	}

	studentID := c.Params("id")

	var req struct {
		AdvisorID string `json:"advisor_id"`
	}

	if err := c.BodyParser(&req); err != nil || req.AdvisorID == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "advisor_id required",
		})
	}

	if err := repoPG.UpdateStudentAdvisorRepo(
		studentID,
		req.AdvisorID,
	); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to update advisor",
		})
	}

	return c.JSON(fiber.Map{
		"message": "advisor updated successfully",
	})
}

func GetAdvisorService(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	var studentID string

	switch role {
	case "admin":
		studentID = c.Params("id")
	case "mahasiswa":
		studentID = c.Locals("student_id").(string)
	default:
		return c.Status(403).JSON(fiber.Map{
			"message": "forbidden",
		})
	}

	data, err := repoPG.GetStudentAdvisorRepo(studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to fetch advisor",
		})
	}

	if data == nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "advisor not assigned",
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}
