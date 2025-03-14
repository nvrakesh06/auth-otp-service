package database

import (
	"log"
	"time"
)

func ExpireOldOTPs() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			query := `UPDATE otp_requests 
			          SET status='expired', last_updated=NOW() 
			          WHERE status='alive' AND expires_at < NOW()`
			_, err := DB.Exec(query)
			if err != nil {
				log.Println("⚠️ Failed to update expired OTPs:", err)
			} else {
				log.Println("✅ Expired OTPs updated successfully!")
			}
		}
	}()
}
