package postgres

import (
	repoPG "backendUAS/app/repositories/postgres"
	repoMongo "backendUAS/app/repositories/mongo"

	"github.com/gofiber/fiber/v2"
)

func GetReportStatisticsService(c *fiber.Ctx) error {

	totalByPeriod, err := repoPG.GetTotalAchievementByPeriodRepo()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "gagal mengambil total prestasi per periode",
		})
	}

	topStudents, err := repoPG.GetTopStudentsRepo(10)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "gagal mengambil top mahasiswa berprestasi",
		})
	}

	totalByType, err := repoMongo.GetTotalAchievementByTypeRepo()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "gagal mengambil total prestasi per tipe",
		})
	}

	competitionLevel, err := repoMongo.GetCompetitionLevelDistributionRepo()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "gagal mengambil distribusi tingkat kompetisi",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": fiber.Map{
			"total_by_type":                  totalByType,
			"total_by_period":                totalByPeriod,
			"top_students":                   topStudents,
			"competition_level_distribution": competitionLevel,
		},
	})
}




func GetStudentReportService(c *fiber.Ctx) error {

	role := c.Locals("role")
	paramStudentID := c.Params("id")


	if role == "mahasiswa" {
		studentID := c.Locals("student_id").(string)
		if studentID != paramStudentID {
			return c.Status(403).JSON(fiber.Map{
				"message": "akses ditolak",
			})
		}
	}

	totalVerified, err := repoPG.GetTotalVerifiedByStudentRepo(paramStudentID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "gagal mengambil total prestasi mahasiswa",
		})
	}

	byType, err := repoMongo.GetStudentAchievementByTypeRepo(paramStudentID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "gagal mengambil prestasi per tipe mahasiswa",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"student_id":     paramStudentID,
		"total_verified": totalVerified,
		"by_type":        byType,
	})
}


