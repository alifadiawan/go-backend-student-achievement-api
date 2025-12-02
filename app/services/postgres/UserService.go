package postgres

import (
	models "backendUAS/app/models/postgres"
	repo "backendUAS/app/repositories/postgres"

	"github.com/gofiber/fiber/v2"
)

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

func GetUsersByIdService (c *fiber.Ctx) error {
	UserID := c.Params("user_id")
	
	var User models.User
	
	User, err := repo.GetUsersByIdRepository(UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
			"message": "tidak dapat mengambil user dengan id " + UserID,
		})
	}
	
	return c.Status(200).JSON(fiber.Map{
		"status" : "success",
		"data": User,
	})
	
}


func StoreUserService(c *fiber.Ctx) error {
	var UserRequest models.UserRequest
	
	err := c.BodyParser(&UserRequest)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "body tidak valid",
			"error": err.Error(),
		})
	}	
	
	_, err = repo.StoreUserRepository(UserRequest)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak dapat menambahkan user",
			"error": err.Error(),
		})
	}
	
	return c.Status(200).JSON(fiber.Map{
		"message": "success added user",
		"status": "success",
	})
	
}


func UpdateUserService(c *fiber.Ctx) error {
	userid := c.Locals("user_id").(string) 
	var userRequest models.UserRequest

	err := c.BodyParser(&userRequest)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "body tidak valid",
			"error": err.Error(),
		})
	}

	hasil, err := repo.UpdateUserRepository(userid, userRequest)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "tidak dapat mengupdate user",
			"error" : err.Error(),
		})
	}

	if !hasil {
		return c.Status(404).JSON(fiber.Map{
			"message": "user tidak ditemukan",
			"error": err.Error(),
		})
	}

	return c.Status(400).JSON(fiber.Map{
		"message": "berhasil mengupdate user",
		"data": hasil,
	})

} 
