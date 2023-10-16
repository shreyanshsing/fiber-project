package router

import (
	"fiber-project/handlers"
	"fiber-project/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRouter(app *fiber.App) {

	// only routes which starts with /api
	// will be processed by this router
	api := app.Group("/api", logger.New())

	// auth routes
	auth := api.Group("/auth")
	auth.Post("/login", handlers.Login)
	auth.Get("/verify/:id", handlers.VerifyEmail)

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Fiber-Project")
	})

	// user routes
	user := api.Group("/user")
	user.Get("/", handlers.GetAllUsers)
	user.Get("/:id", middlewares.ProtectedRoute(), handlers.GetUser)
	user.Post("/", handlers.CreateUser)
	user.Patch("/:id", middlewares.ProtectedRoute(), handlers.UpdateUser)
	user.Delete("/:id", middlewares.ProtectedRoute(), handlers.DeleteUser)
}
