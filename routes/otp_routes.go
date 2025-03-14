package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nvrakesh06/auth-otp-service/controllers"
)

func SetupOTPRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/send-otp", controllers.SendOTP)
	api.Get("/get-otp", controllers.GetStoredOTP)
	api.Post("/verify-otp", controllers.VerifyOTP)
}