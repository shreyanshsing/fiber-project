package handlers

import (
	"fiber-project/database"
	gomail_lib "fiber-project/lib/mail"
	"fiber-project/models"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserUpdateInput struct {
	Password string `json:"password"`
	Names    string `json:"names"`
}

type UserDeleteInput struct {
	Password string `json:"password"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validateUser(id string, password string) bool {
	db := database.DB
	var user models.User

	db.First(&user, id)

	if user.UserName == "" && user.Email == "" {
		return false
	}

	return CheckPassword(user.Password, password)
}

func GetAllUsers(c *fiber.Ctx) error {

	db := database.DB
	var users []models.User

	if err := db.Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to get users", "error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"message": "all users", "data": users})
}

func GetUser(c *fiber.Ctx) error {

	id := c.Params("id")
	db := database.DB
	var user models.User

	if err := db.Find(&user, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to get user", "error": err.Error()})
	}

	if user.UserName == "" && user.Email == "" {
		return c.Status(404).JSON(fiber.Map{"message": "No user found with ID", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"message": "user found", "data": user})
}

func CreateUser(c *fiber.Ctx) error {

	user := new(models.User)

	// parse-input
	if err := c.BodyParser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to parse input", "error": err.Error()})
	}

	// validate input
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to validate input", "error": err.Error()})
	}

	// hash password
	hashedPassword, err := hashPassword(user.Password)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to hash password", "error": err.Error()})
	}

	user.Password = hashedPassword
	db := database.DB

	if err := db.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to create user", "error": err.Error()})
	}

	log.Println("User created successfully!")

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered from panic in SendMail:", r)
			}
		}()
		gomail_lib.SendMail(user)
	}()

	return c.Status(200).JSON(fiber.Map{"message": "user created", "data": user})
}

func UpdateUser(c *fiber.Ctx) error {
	var input UserUpdateInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to parse input", "error": err.Error()})
	}

	id := c.Params("id")
	db := database.DB
	var user models.User

	db.First(&user, id)

	log.Println(input.Names != "", input.Password)

	if input.Names != "" {
		user.Names = input.Names
	}

	if input.Password != "" {
		// hash password
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "failed to hash password", "error": err.Error()})
		}
		user.Password = hashedPassword
	}

	db.Save(&user)

	return c.Status(200).JSON(fiber.Map{"message": "user updated", "data": user})
}

func DeleteUser(c *fiber.Ctx) error {
	var input UserDeleteInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to parse input", "error": err.Error()})
	}
	id := c.Params("id")

	if !validateUser(id, input.Password) {
		return c.Status(404).JSON(fiber.Map{"message": "No user found with ID", "data": nil})
	}

	db := database.DB
	var user models.User

	db.First(&user, id)
	db.Unscoped().Delete(&user)

	return c.Status(200).JSON(fiber.Map{"message": "user deleted", "data": user})
}
