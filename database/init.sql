-- Ensure the database is created
CREATE DATABASE otp_service;

\c otp_service;

-- Create ENUM type for OTP status
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'otp_status_enum') THEN
        CREATE TYPE otp_status_enum AS ENUM ('alive', 'consumed', 'expired', 'duplicate');
    END IF;
END $$;

-- Create OTP Requests table if it doesn't exist
CREATE TABLE IF NOT EXISTS otp_requests (
    id SERIAL PRIMARY KEY,
    phone VARCHAR(20) NOT NULL,
    otp VARCHAR(6) NOT NULL,
    status otp_status_enum DEFAULT 'alive',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP + INTERVAL '5 minutes'),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- ✅ Ensure the index exists on `expires_at`
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_otp_expires_at') THEN
        CREATE INDEX idx_otp_expires_at ON otp_requests (expires_at);
    END IF;
END $$;