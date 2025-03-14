package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectPostgres() {
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
			log.Println("✅ Connected to PostgreSQL successfully!")
			return
		}

		log.Println("⏳ Waiting for PostgreSQL to be ready...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("❌ Failed to connect to PostgreSQL:", err)
	}
}
