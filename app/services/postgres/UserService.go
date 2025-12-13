package postgres

import (
	models "backendUAS/app/models/postgres"
	repo "backendUAS/app/repositories/postgres"

	"github.com/gofiber/fiber/v2"
)


// @Summary Get all users
// @Description Mengambil semua data user
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Berhasil mengambil semua user"
// @Failure 400 {object} map[string]interface{} "Gagal mengambil data user"
// @Security BearerAuth
// @Router /api/v1/users [get]
func GetAllUserService(c *fiber.Ctx) error {

	var Users []models.User

	Users, err := repo.GetAllUserRepository()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
			"status":  "error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
		"data":    Users,
	})

}



// @Summary Get user by ID
// @Description Mengambil data user berdasarkan user_id. Hanya admin atau user itu sendiri yang dapat mengakses.
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil user"
// @Failure 400 {object} map[string]interface{} "Tidak dapat mengambil user dengan ID tersebut"
// @Failure 403 {object} map[string]interface{} "Hanya admin atau user itu sendiri yang boleh mengakses"
// @Security BearerAuth
// @Router /api/v1/users/{user_id} [get]
func GetUsersByIdService(c *fiber.Ctx) error {
	UserID := c.Params("user_id")

	var User models.User

	UserIdViaJWT := c.Locals("user_id")
	role := c.Locals("role")
	if role != "admin" && UserID != UserIdViaJWT {
		return c.Status(403).JSON(fiber.Map{
			"message": "maaf, hanya admin yang boleh yaw",
		})
	}

	User, err := repo.GetUsersByIdRepository(UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "tidak dapat mengambil user dengan id " + UserID,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   User,
	})

}



// @Summary Store new user
// @Description Menambahkan user baru
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.UpdateUser true "User Request"
// @Success 200 {object} map[string]interface{} "Berhasil menambahkan user"
// @Failure 400 {object} map[string]interface{} "Body tidak valid atau gagal menambahkan user"
// @Security BearerAuth
// @Router /api/v1/users [post]
func StoreUserService(c *fiber.Ctx) error {
	var UserRequest models.UpdateUser

	err := c.BodyParser(&UserRequest)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "body tidak valid",
			"error":   err.Error(),
		})
	}

	_, err = repo.StoreUserRepository(UserRequest)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak dapat menambahkan user",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success added user",
		"status":  "success",
	})

}



// @Summary Update user
// @Description Mengupdate data user berdasarkan user_id
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body models.UpdateUser true "User Request"
// @Success 200 {object} map[string]interface{} "Berhasil mengupdate user"
// @Failure 400 {object} map[string]interface{} "Body tidak valid atau gagal mengupdate user"
// @Failure 404 {object} map[string]interface{} "User tidak ditemukan"
// @Security BearerAuth
// @Router /api/v1/users/{user_id} [put]
func UpdateUserService(c *fiber.Ctx) error {
	userid := c.Params("user_id")
	var userRequest models.UpdateUser

	err := c.BodyParser(&userRequest)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "body tidak valid",
			"error":   err.Error(),
		})
	}

	hasil, err := repo.UpdateUserRepository(userid, userRequest)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak dapat mengupdate user",
			"error":   err.Error(),
		})
	}

	if !hasil {
		return c.Status(404).JSON(fiber.Map{
			"message": "user tidak ditemukan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "berhasil mengupdate user",
		"data":    hasil,
	})

}


// @Summary Update user role
// @Description Mengupdate role user berdasarkan user_id
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body models.UpdateUserRole true "Update User Role Request"
// @Success 200 {object} map[string]interface{} "Berhasil mengupdate role user"
// @Failure 400 {object} map[string]interface{} "Body tidak valid atau gagal mengupdate role user"
// @Failure 404 {object} map[string]interface{} "User tidak ditemukan"
// @Security BearerAuth
// @Router /api/v1/users/{user_id}/role [put]
func UpdateUserRoleService(c *fiber.Ctx) error {
	userid := c.Params("user_id")
	var userRequest models.UpdateUserRole

	err := c.BodyParser(&userRequest)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "body tidak valid",
			"error":   err.Error(),
		})
	}

	hasil, err := repo.UpdateUserRoleRepository(userid, userRequest.RoleName)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak dapat mengupdate role user",
			"error":   err.Error(),
		})
	}

	if !hasil {
		return c.Status(404).JSON(fiber.Map{
			"message": "user tidak ditemukan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "berhasil mengupdate role user",
		"data":    hasil,
	})

}



// @Summary Delete user
// @Description Menghapus user berdasarkan user_id
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "Berhasil menghapus user"
// @Failure 400 {object} map[string]interface{} "Gagal menghapus user"
// @Security BearerAuth
// @Router /api/v1/users/{id} [delete]
func DeleteUserService(c *fiber.Ctx) error {
	userid := c.Params("id")

	hasil, err := repo.DeleteUserRepository(userid)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "gagal delete user",
			"error":   err.Error(),
		})
	}

	if !hasil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak dapat delete user",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Berhasil Delete User",
	})

}
