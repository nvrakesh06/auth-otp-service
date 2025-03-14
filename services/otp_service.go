package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/nvrakesh06/auth-otp-service/database"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SaveOTP(phone string, otp string) error {
	ctx := context.Background()

	_, err := database.RedisClient.Get(ctx, phone).Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("redis error: %v", err)
	}

	if err == nil {
		delErr := database.RedisClient.Del(ctx, phone).Err()
		if delErr != nil {
			fmt.Printf("⚠️ Failed to delete existing OTP from Redis: %v\n", delErr)
		}

		updateQuery := `UPDATE otp_requests SET status='duplicate', last_updated=NOW() WHERE phone=$1 AND status='alive'`
		_, dbErr := database.DB.Exec(updateQuery, phone)
		if dbErr != nil {
			return fmt.Errorf("failed to update existing OTP status in PostgreSQL: %v", dbErr)
		}
	}

	expiration := 5 * time.Minute
	err = database.RedisClient.Set(ctx, phone, otp, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to store OTP in Redis: %v", err)
	}

	query := `INSERT INTO otp_requests (phone, otp, status, created_at, expires_at, last_updated) VALUES ($1, $2, 'alive', $3, $4, NOW())`
	_, err = database.DB.Exec(query, phone, otp, time.Now(), time.Now().Add(expiration))
	if err != nil {
		return fmt.Errorf("failed to store OTP in PostgreSQL: %v", err)
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
		return false, fmt.Errorf("invalid OTP")
	}

	err = database.RedisClient.Del(ctx, phone).Err()
	if err != nil {
		return false, fmt.Errorf("failed to delete OTP after verification")
	}

	updateQuery := `UPDATE otp_requests SET status='consumed', last_updated=NOW() WHERE phone=$1 AND status='alive'`
	_, err = database.DB.Exec(updateQuery, phone)
	if err != nil {
		return false, fmt.Errorf("failed to update OTP status")
	}

	return true, nil
}
