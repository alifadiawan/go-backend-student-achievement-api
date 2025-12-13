package postgres

import (
	repopg "backendUAS/app/repositories/postgres"

	"github.com/gofiber/fiber/v2"
)


// @Summary Get lecturer data
// @Description Mengambil data dosen. Mahasiswa hanya bisa melihat dosen pembimbingnya, admin dapat melihat semua dosen.
// @Tags Lecturer
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data dosen"
// @Failure 400 {object} map[string]interface{} "Tidak dapat mengambil data dosen"
// @Failure 403 {object} map[string]interface{} "Role tidak memiliki akses"
// @Failure 404 {object} map[string]interface{} "Role tidak valid"
// @Security BearerAuth
// @Router /api/v1/lecturers [get]
func GetLecturerService(c *fiber.Ctx) error {

	role := c.Locals("role")
	if role == "" {
		return c.Status(404).JSON(fiber.Map{
			"message": "role tidak valid",
		})
	}

	if role == "dosen" {
		return c.Status(403).JSON(fiber.Map{
			"message": "maaf, role anda tidak bole",
		})
	}

	if role == "mahasiswa" {
		student_id := c.Locals("student_id").(string)

		lecturer_detail, err := repopg.GetLecturerByStudentIDRepo(student_id)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "tidak dapat mengambil data lecturer",
				"error":   err.Error(),
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"data":   lecturer_detail,
		})

	}

	lecturer_detail, err := repopg.GetAllLecturerRepo()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak dapat mengambil data lecturer",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   lecturer_detail,
	})

}
