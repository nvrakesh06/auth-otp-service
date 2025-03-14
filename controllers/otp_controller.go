package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nvrakesh06/auth-otp-service/services"
)

// OTPRequest struct for API request payload
type OTPRequest struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// VerifyOTPRequest struct for OTP verification
type VerifyOTPRequest struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}

// SendOTP generates and stores an OTP
func SendOTP(c *fiber.Ctx) error {
	var request OTPRequest

	// Parse JSON request
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// Generate OTP
	otp := services.GenerateOTP()

	// Save OTP in Redis
	err := services.SaveOTP(request.Phone, otp)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to store OTP"})
	}

	log.Println("Generated OTP:", otp)

	return c.JSON(fiber.Map{"message": "OTP sent successfully!"})
}

func GetStoredOTP(c *fiber.Ctx) error {
	phone := c.Query("phone")

	if phone == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Phone number is required"})
	}

	otp, err := services.GetOTP(phone)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "OTP not found or expired"})
	}

	log.Println("Retrieved OTP for", phone, ":", otp)

	return c.JSON(fiber.Map{"otp": otp})
}

func VerifyOTP(c *fiber.Ctx) error {
	var request VerifyOTPRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request format"})
	}

	isValid, err := services.VerifyOTP(request.Phone, request.OTP)
	if err != nil || !isValid {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid OTP or expired"})
	}

	token, err := services.GenerateJWT(request.Phone)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"message": "OTP Verified!",
		"token":   token,
	})
}
