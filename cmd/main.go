package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/nvrakesh06/auth-otp-service/config"
	"github.com/nvrakesh06/auth-otp-service/routes"
	"github.com/nvrakesh06/auth-otp-service/database"
)

func main() {
	config.LoadEnv()

	database.ConnectPostgres()
	database.ConnectRedis()
	database.ExpireOldOTPs() // TODO: test this flow and check if there's better way to handle this

	app := fiber.New()

	routes.SetupRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, OTP Service is running!")
	})

	log.Fatal(app.Listen(":8080"))
}
