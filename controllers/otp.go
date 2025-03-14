package controllers

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano()) // Ensure randomness
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
