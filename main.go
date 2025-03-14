package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/nvrakesh06/auth-otp-service/database"
	"github.com/nvrakesh06/auth-otp-service/controllers"
)

type OTPRequest struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func sendOTP(c *fiber.Ctx) error {
	var request OTPRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	otpCode := controllers.GenerateOTP()

	err := database.SaveOTP(request.Phone, otpCode)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to store OTP"})
	}

	// TODO: Send OTP via SMS & Email (will be implemented next)
	log.Println("Generated OTP:", otpCode)

	return c.JSON(fiber.Map{"message": "OTP sent successfully!"})
}

func getStoredOTP(c *fiber.Ctx) error {
	phone := c.Query("phone")

	log.Println("Fetching OTP for phone:", phone)

	otp, err := database.GetOTP(phone)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "OTP not found or expired"})
	}

	return c.JSON(fiber.Map{"otp": otp})
}


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
	app.Post("/send-otp", sendOTP)
	app.Get("/get-otp", getStoredOTP)

	log.Fatal(app.Listen(":8080"))
}