package handlers

import (
	"fiber-project/config"
	"fiber-project/database"
	"fiber-project/models"
	"log"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	TOKEN_DURATION = time.Hour * 24
)

type LoginInput struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

func isEmail(email string) bool {
	// define regular expression for email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// match email against regular expression
	return emailRegex.MatchString(email)
}

func getUserByEmail(email string) (*models.User, error) {
	// get user by email
	db := database.DB
	var user models.User

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func getUserName(username string) (*models.User, error) {
	// get user by username
	db := database.DB
	var user models.User

	if err := db.Where("user_name = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func CheckPassword(userPassword string, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(inputPassword))
	return err == nil
}

func Login(c *fiber.Ctx) error {

	var userInput LoginInput

	if err := c.BodyParser(&userInput); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to parse input", "error": err.Error()})
	}

	if userInput.ID == "" || userInput.Password == "" {
		return c.Status(500).JSON(fiber.Map{"message": "failed to validate input", "error": "id and password are required"})
	}

	// check if id is email
	// if mail, get user by email
	// else get user by id

	userModel, err := new(models.User), *new(error)

	if isEmail(userInput.ID) {
		// get user by email
		userModel, err = getUserByEmail(userInput.Password)
	} else {
		// get user by id
		userModel, err = getUserName(userInput.ID)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to get user", "error": err.Error()})
	}

	if userModel == nil {
		return c.Status(404).JSON(fiber.Map{"message": "No user found with ID", "data": nil})
	}

	if !CheckPassword(userModel.Password, userInput.Password) {
		return c.Status(500).JSON(fiber.Map{"message": "incorrect password", "data": nil})
	}

	// if user exists
	// generate jwt token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userModel.ID
	claims["exp"] = time.Now().Add(TOKEN_DURATION).Unix()

	t, err := token.SignedString([]byte(config.GetEnvoirmentVariable("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("Token: %v", token)

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": fiber.Map{"token": t}})

}

func VerifyEmail(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DB
	var user models.User

	db.First(&user, id)

	if user.Email == "" && user.UserName == "" {
		return c.Status(404).JSON(fiber.Map{"message": "No user found with ID", "data": nil})
	}

	if user.Verified {
		return c.Status(500).JSON(fiber.Map{"message": "user email already verified", "data": nil})
	}

	user.Verified = true

	db.Save(&user)

	return c.Status(200).JSON(fiber.Map{"message": "user email verified", "data": nil})
}
