package main

import (
	"fiber-project/database"
	"fiber-project/router"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	database.ConnectToDatabase()
	router.SetupRouter(app)

	err := app.Listen(":8080")

	if err != nil {
		os.Exit(2)
	}
}
