package postgres

import (
	modelsMongo "backendUAS/app/models/mongo"
	models "backendUAS/app/models/postgres"
	repoMongo "backendUAS/app/repositories/mongo"
	repo "backendUAS/app/repositories/postgres"

	"github.com/gofiber/fiber/v2"
)

func GetAllAchievementService(c *fiber.Ctx) error {

	var result []models.Achievement
	var err error

	role := c.Locals("role")
	if role == "admin" {
		result, err = repo.GetAllAchievementRepo()
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "gagal mengambil data prestasi mahasiswa",
			})
		}
	}

	studentID := c.Locals("student_id").(string)
	if role == "mahasiswa" {
		result, err = repo.GetAllAchievementByStudentIDRepo(studentID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "gagal mengambil data prestasi mahasiswa",
			})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}

func GetAchievementByIDService(c *fiber.Ctx) error {
	AchievementID := c.Params("AchievementID")
	if AchievementID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user id tidak valid",
		})
	}

	achievements, err := repo.GetAchievementByIDRepo(AchievementID)
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

func AddAchievementService(c *fiber.Ctx) error {
	userRole := c.Locals("role")
	if userRole == "" {
		return c.Status(404).JSON(fiber.Map{"message": "tidak ada role"})
	}
	if userRole == "dosen" {
		return c.Status(403).JSON(fiber.Map{"message": "anda bukan mahasiswa"})
	}

	var req modelsMongo.Achievement
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "body tidak valid", "error": err.Error()})
	}

	// Insert ke Mongo
	mongoID, err := repoMongo.AddAchievementRepositoryMongo(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "gagal insert mongo", "error": err.Error()})
	}

	// Insert ke Postgres
	err = repo.InsertAchievementPostgres(req.StudentID, mongoID, "draft")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "gagal insert postgres", "error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "berhasil menambahkan prestasi",
		"data": map[string]interface{}{
			"mongo_id": mongoID,
		},
	})
}

func DeleteAchievementService(c *fiber.Ctx) error {

	role := c.Locals("role")
	CurrectLoggedStudentID := c.Locals("student_id")
	achievement_references_id := c.Params("achievement_references_id")

	if achievement_references_id == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID tidak valid",
		})
	}

	// Ambil owner dari DB
	ownerID, err := repo.GetUserIDofAchievementRepo(achievement_references_id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Achievement tidak ditemukan",
		})
	}

	// Jika BUKAN admin DAN BUKAN pemilik data → blokir
	if role != "admin" && CurrectLoggedStudentID != ownerID {
		return c.Status(403).JSON(fiber.Map{
			"message":                "maaf, kamu tidak berhak menghapus achievement ini",
			"CurrectLoggedStudentID": CurrectLoggedStudentID,
			"ownerID":                ownerID,
		})
	}

	query, err := repo.DeleteAchievementRepo(achievement_references_id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak dapat menghapus user",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  query,
		"message": "Berhasil Hapus Achievement",
	})

}

// func SubmitAchievementService (c *fiber.Ctx) error {

// 	achievement_references_id := c.Params("achievement_references_id")
// 	studentID := c.Locals("student_id")
// 	role := c.Locals("")

// 	ownerID, err := repo.GetUserIDofAchievementRepo(achievement_references_id)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"messsage": "tidak dapat mengambil owner yang punya",
// 			"err": err.Error(),
// 		})
// 	}

// 	if ownerID == "" {
// 		return c.Status(400).JSON(fiber.Map{
// 			"message": "owner ID tidak ",
// 		})
// 	}

// 	query, err := repo.DeleteAchievementRepo(achievement_references_id)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"message": "tidak bisa submit achievement",
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.Status(200).JSON(fiber.Map{
// 		"status": query,
// 		"message": "berhasil delete achievement",
// 	})

// }
