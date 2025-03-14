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
			ensureTableExists()
			return
		}

		log.Println("⏳ Waiting for PostgreSQL to be ready...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("❌ Failed to connect to PostgreSQL:", err)
	}
}

func ensureTableExists() {
	query := `
	DO $$ 
	BEGIN
	    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'otp_status_enum') THEN
	        CREATE TYPE otp_status_enum AS ENUM ('alive', 'consumed', 'expired', 'duplicate');
	    END IF;
	END $$;

	CREATE TABLE IF NOT EXISTS otp_requests (
	    id SERIAL PRIMARY KEY,
	    phone VARCHAR(20) NOT NULL,
	    otp VARCHAR(6) NOT NULL,
	    status otp_status_enum DEFAULT 'alive',
	    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	    expires_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP + INTERVAL '5 minutes'),
	    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	DO $$
	BEGIN
	    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_otp_expires_at') THEN
	        CREATE INDEX idx_otp_expires_at ON otp_requests (expires_at);
	    END IF;
	END $$;
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("❌ Error creating table:", err)
	} else {
		log.Println("✅ Table `otp_requests` ensured in PostgreSQL!")
	}
}