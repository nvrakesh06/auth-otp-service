package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nvrakesh06/auth-otp-service/services"
)

func ProtectedRoute(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "You have access to this protected route!"})
}

func LoginWithOTP(c *fiber.Ctx) error {
	type LoginRequest struct {
		Phone string `json:"phone"`
		OTP   string `json:"otp"`
	}

	var request LoginRequest

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

	return c.JSON(fiber.Map{"token": token})
}
