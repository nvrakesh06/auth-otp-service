package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"log"

	"github.com/nvrakesh06/auth-otp-service/database"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SaveOTP(phone string, otp string) error {
	ctx := context.Background()
	expiration := 5 * time.Minute // OTP expires in 5 minutes

	err := database.RedisClient.Set(ctx, phone, otp, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to store OTP: %v", err)
	}

	return nil
}

func GetOTP(phone string) (string, error) {
	ctx := context.Background()

	keys, err := database.RedisClient.Keys(ctx, "*").Result()
	if err != nil {
		log.Println("Error fetching keys:", err)
	}

	log.Println("Redis stored keys:", keys)

	otp, err := database.RedisClient.Get(ctx, phone).Result()
	if err != nil {
		log.Println("Error retrieving OTP for", phone, ":", err)
		return "", fmt.Errorf("OTP not found or expired")
	}

	log.Println("OTP retrieved:", otp)
	return otp, nil
}

func VerifyOTP(phone string, otp string) (bool, error) {
	ctx := context.Background()

	storedOTP, err := database.RedisClient.Get(ctx, phone).Result()
	if err != nil {
		return false, fmt.Errorf("OTP not found or expired")
	}

	if storedOTP != otp {
		return false, fmt.Errorf("Invalid OTP")
	}

	err = database.RedisClient.Del(ctx, phone).Err()
	if err != nil {
		return false, fmt.Errorf("Failed to delete OTP after verification")
	}

	return true, nil
}
