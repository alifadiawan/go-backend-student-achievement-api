package postgres

import (
	modelsMongo "backendUAS/app/models/mongo"
	models "backendUAS/app/models/postgres"
	repoMongo "backendUAS/app/repositories/mongo"
	repo "backendUAS/app/repositories/postgres"

	"github.com/gofiber/fiber/v2"
)


// @Summary Get all achievements
// @Description Admin dapat melihat semua prestasi, mahasiswa hanya melihat prestasinya sendiri
// @Tags Achievement
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data prestasi"
// @Failure 400 {object} map[string]interface{} "Gagal mengambil data prestasi"
// @Security BearerAuth
// @Router /api/v1/achievements [get]
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


// @Summary Get achievement by ID
// @Description Mengambil detail prestasi berdasarkan achievement ID
// @Tags Achievement
// @Accept json
// @Produce json
// @Param AchievementID path string true "Achievement ID"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data prestasi"
// @Failure 400 {object} map[string]interface{} "Achievement ID tidak valid"
// @Failure 500 {object} map[string]interface{} "Gagal mengambil data prestasi"
// @Security BearerAuth
// @Router /api/v1/achievements/{AchievementID} [get]
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

// @Summary Add new achievement
// @Description Mahasiswa menambahkan prestasi baru (MongoDB + PostgreSQL)
// @Tags Achievement
// @Accept json
// @Produce json
// @Param request body modelsMongo.Achievement true "Achievement Request"
// @Success 200 {object} map[string]interface{} "Berhasil menambahkan prestasi"
// @Failure 400 {object} map[string]interface{} "Body tidak valid"
// @Failure 403 {object} map[string]interface{} "Bukan mahasiswa"
// @Failure 500 {object} map[string]interface{} "Gagal menyimpan data prestasi"
// @Security BearerAuth
// @Router /api/v1/achievements [post]

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


// @Summary Delete achievement
// @Description Menghapus prestasi. Hanya admin atau pemilik prestasi
// @Tags Achievement
// @Accept json
// @Produce json
// @Param achievement_references_id path string true "Achievement Reference ID"
// @Success 200 {object} map[string]interface{} "Berhasil menghapus prestasi"
// @Failure 400 {object} map[string]interface{} "ID tidak valid atau gagal menghapus"
// @Failure 403 {object} map[string]interface{} "Tidak memiliki akses"
// @Failure 404 {object} map[string]interface{} "Achievement tidak ditemukan"
// @Security BearerAuth
// @Router /api/v1/achievements/{achievement_references_id} [delete]
func DeleteAchievementService(c *fiber.Ctx) error {

	role := c.Locals("role")
	CurrectLoggedStudentID := c.Locals("student_id")
	achievement_references_id := c.Params("achievement_references_id")

	if achievement_references_id == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID tidak valid",
		})
	}

	ownerID, err := repo.GetUserIDofAchievementRepo(achievement_references_id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Achievement tidak ditemukan",
		})
	}

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


// @Summary Submit achievement
// @Description Mahasiswa submit prestasi untuk diverifikasi
// @Tags Achievement
// @Accept json
// @Produce json
// @Param achievement_references_id path string true "Achievement Reference ID"
// @Success 200 {object} map[string]interface{} "Berhasil submit prestasi"
// @Failure 400 {object} map[string]interface{} "Gagal submit prestasi"
// @Failure 403 {object} map[string]interface{} "Bukan pemilik prestasi"
// @Failure 404 {object} map[string]interface{} "Achievement tidak ditemukan"
// @Security BearerAuth
// @Router /api/v1/achievements/{achievement_references_id}/submit [post]
func SubmitAchievementService(c *fiber.Ctx) error {
	achievementID := c.Params("achievement_references_id")
	role := c.Locals("role")
	loggedStudentID := c.Locals("student_id").(string)

	ownerID, err := repo.GetUserIDofAchievementRepo(achievementID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "tidak bisa mendapatkan owner", "error": err.Error()})
	}

	if ownerID == "" {
		return c.Status(404).JSON(fiber.Map{"message": "achievement tidak ditemukan"})
	}

	if role == "mahasiswa" && loggedStudentID != ownerID {
		return c.Status(403).JSON(fiber.Map{"message": "maaf ya bukan punya anda"})
	}

	ok, err := repo.SubmitAchievementRepository(achievementID, ownerID)
	if err != nil || !ok {
		return c.Status(400).JSON(fiber.Map{"message": "tidak bisa submit achievement", "error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"status": ok, "message": "berhasil submit achievement"})
}


// @Summary Approve achievement
// @Description Admin atau dosen menyetujui prestasi mahasiswa
// @Tags Achievement
// @Accept json
// @Produce json
// @Param achievement_references_id path string true "Achievement Reference ID"
// @Success 200 {object} map[string]interface{} "Berhasil approve prestasi"
// @Failure 400 {object} map[string]interface{} "Gagal approve prestasi"
// @Failure 403 {object} map[string]interface{} "Akses ditolak"
// @Security BearerAuth
// @Router /api/v1/achievements/{achievement_references_id}/approve [post]
func ApproveAchievmentService(c *fiber.Ctx) error {

	AchievementID := c.Params("achievement_references_id")
	role := c.Locals("role")

	if AchievementID == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "achievement ID tidak valid",
		})
	}

	if role == "mahasiswa" {
		return c.Status(403).JSON(fiber.Map{
			"message": "maap, anda tidak boleh akses route ini ",
		})
	}

	query, err := repo.ApproveAchievmentRepository(AchievementID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak dapat approve achievement",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "berhasil approve achievement",
		"result":  query,
	})

}


// @Summary Verify achievement
// @Description Verifikasi prestasi oleh dosen pembimbing atau admin
// @Tags Achievement
// @Accept json
// @Produce json
// @Param achievement_references_id path string true "Achievement Reference ID"
// @Success 200 {object} map[string]interface{} "Berhasil verify prestasi"
// @Failure 400 {object} map[string]interface{} "Gagal verify prestasi"
// @Failure 403 {object} map[string]interface{} "Akses ditolak"
// @Failure 404 {object} map[string]interface{} "Achievement tidak ditemukan"
// @Security BearerAuth
// @Router /api/v1/achievements/{achievement_references_id}/verify [post]
func VerifyAchievementService(c *fiber.Ctx) error {

	achievementID := c.Params("achievement_references_id")
	if achievementID == "" {
		return c.Status(404).JSON(fiber.Map{
			"message": "ID achievement references tidak valid",
		})
	}

	// RBAC
	role := c.Locals("role").(string)
	if role != "dosen" && role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"message": "akses ditolak",
		})
	}

	// Jika dosen harus cek bahawa achievement milik mahasiswa bimbingannya
	if role == "dosen" {

		studentID, err := repo.GetUserIDofAchievementRepo(achievementID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"message": "achievement tidak ditemukan",
			})
		}

		advisorID, err := repo.GetAdvisorFromStudent(studentID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "gagal mengambil advisor mahasiswa",
				"error":   err.Error(),
			})
		}

		loggedLecturerID := c.Locals("lecturer_id").(string)
		if advisorID != loggedLecturerID {
			return c.Status(403).JSON(fiber.Map{
				"message":          "maaf, anda bukan dosen pembimbing mahasiswa ini",
				"loggedLecturerID": loggedLecturerID,
				"advisorID":        advisorID,
			})
		}
	}

	result, err := repo.VerifyAchievementRepo(achievementID, c.Locals("user_id").(string))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error ketika verify achievement",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Berhasil update achievement ke verified",
		"data":    result,
	})
}


// @Summary Reject achievement
// @Description Menolak prestasi dengan catatan penolakan
// @Tags Achievement
// @Accept json
// @Produce json
// @Param achievement_references_id path string true "Achievement Reference ID"
// @Param request body models.AchievmentRejectRequest true "Rejection Note"
// @Success 200 {object} map[string]interface{} "Berhasil reject prestasi"
// @Failure 400 {object} map[string]interface{} "Gagal reject prestasi"
// @Failure 403 {object} map[string]interface{} "Akses ditolak"
// @Failure 404 {object} map[string]interface{} "Achievement tidak ditemukan"
// @Security BearerAuth
// @Router /api/v1/achievements/{achievement_references_id}/reject [post]
func RejectAchievementService(c *fiber.Ctx) error {

	var rejection_note models.AchievmentRejectRequest
	achievementID := c.Params("achievement_references_id")

	role := c.Locals("role")
	if role == "mahasiswa" {
		return c.Status(403).JSON(fiber.Map{
			"message": "anda tidak boleh, maaf",
		})
	}

	if role == "dosen" {
		studentID, err := repo.GetUserIDofAchievementRepo(achievementID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"message": "achievement tidak ditemukan",
			})
		}

		advisorID, err := repo.GetAdvisorFromStudent(studentID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "gagal mengambil advisor mahasiswa",
				"error":   err.Error(),
			})
		}

		loggedLecturerID := c.Locals("lecturer_id").(string)
		if advisorID != loggedLecturerID {
			return c.Status(403).JSON(fiber.Map{
				"message":          "maaf, anda bukan dosen pembimbing mahasiswa ini",
				"loggedLecturerID": loggedLecturerID,
				"advisorID":        advisorID,
			})
		}
	}

	err := c.BodyParser(&rejection_note)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "rejection note tidak valid",
		})
	}

	query, err := repo.RejectAchievementRepo(achievementID, rejection_note.RejectionNote)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak dapat reject achievement",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "berhasil reject achievement",
		"status":  "success",
		"data":    query,
	})

}


// @Summary Get achievement history
// @Description Mengambil riwayat status prestasi (draft, submitted, rejected)
// @Tags Achievement
// @Accept json
// @Produce json
// @Param achievement_references_id path string true "Achievement Reference ID"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil history prestasi"
// @Failure 404 {object} map[string]interface{} "Achievement tidak ditemukan"
// @Security BearerAuth
// @Router /api/v1/achievements/{achievement_references_id}/history [get]
func HistoryAchievementService(c *fiber.Ctx) error {

	achievement_referencens_id := c.Params("achievement_references_id")

	if achievement_referencens_id == "" {
		return c.Status(404).JSON(fiber.Map{
			"message": "achievement references id tidak valid",
		})
	}

	data, err := repo.GetAchievementByIDRepo(achievement_referencens_id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "data tidak ditemukan",
		})
	}

	history := []map[string]interface{}{}

	history = append(history, map[string]interface{}{
		"status":    "draft",
		"timestamp": data.Achievement.CreatedAt,
	})

	if !data.Achievement.SubmittedAt.IsZero() {
		history = append(history, map[string]interface{}{
			"status":    "submitted",
			"timestamp": data.Achievement.SubmittedAt,
		})
	}

	if data.Achievement.RejectionNote != nil {
		history = append(history, map[string]interface{}{
			"status":    "rejected",
			"timestamp": data.Achievement.UpdatedAt,
			"note":      *data.Achievement.RejectionNote,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"history": history,
	})

}
