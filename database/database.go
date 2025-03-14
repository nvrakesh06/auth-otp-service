package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/go-redis/redis/v8"
)

var (
	DB          *sqlx.DB        
	RedisClient *redis.Client   
	ctx         = context.Background()
)

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	var err error
	for i := 0; i < 5; i++ {
		DB, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			fmt.Println("✅ Connected to PostgreSQL successfully!")
			return
		}
		fmt.Println("⏳ Waiting for PostgreSQL to be ready...")
		time.Sleep(20 * time.Second)
	}

	if err != nil {
		log.Fatal("❌ Failed to connect to PostgreSQL: ", err)
	}
}

func ConnectRedis() { 
	RedisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("❌ Failed to connect to Redis:", err)
	}

	fmt.Println("✅ Connected to Redis successfully!")
}

func SaveOTP(phone string, otp string) error {
	ctx := context.Background()
	expiration := 5 * time.Minute

	log.Println("Saving OTP for phone:", phone)

	err := RedisClient.Set(ctx, phone, otp, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to store OTP: %v", err)
	}

	return nil
}


func GetOTP(phone string) (string, error) {
	ctx := context.Background()

	keys, err := RedisClient.Keys(ctx, "*").Result()
	if err != nil {
		log.Println("Error fetching keys:", err)
	}

	log.Println("Redis stored keys:", keys)

	otp, err := RedisClient.Get(ctx, phone).Result()
	if err != nil {
		log.Println("Error retrieving OTP for", phone, ":", err)
		return "", fmt.Errorf("OTP not found or expired")
	}

	log.Println("OTP retrieved:", otp)
	return otp, nil
}

func DeleteOTP(phone string) error {
	ctx := context.Background()

	err := RedisClient.Del(ctx, phone).Err()
	if err != nil {
		log.Println("Error deleting OTP for phone:", phone, err)
		return fmt.Errorf("failed to delete OTP: %v", err)
	}

	log.Println("OTP deleted for phone:", phone)
	return nil
}