package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/nvrakesh06/auth-otp-service/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	database.ConnectDatabase()
	database.ConnectRedis()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OTP Service is Running!")
	})

	log.Fatal(app.Listen(":8080"))
}