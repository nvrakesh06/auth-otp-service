package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nvrakesh06/auth-otp-service/controllers"
	"github.com/nvrakesh06/auth-otp-service/middleware"
)

func SetupAuthRoutes(app *fiber.App) {
	api := app.Group("/auth")
	api.Post("/login", controllers.LoginWithOTP)

	protected := api.Group("/protected", middleware.JWTMiddleware())
	protected.Get("/", controllers.ProtectedRoute)
}
