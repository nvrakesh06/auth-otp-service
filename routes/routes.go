package routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
	SetupOTPRoutes(app)
	SetupAuthRoutes(app)
}
