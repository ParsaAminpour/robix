package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateJWTSecret() {
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		fmt.Printf("Error generating random bytes: %v\n", err)
		return
	}

	// Encode as base64
	secret := base64.StdEncoding.EncodeToString(bytes)
	fmt.Printf("Generated JWT Secret: %s\n", secret)
	fmt.Printf("\nAdd this to your .env files in both user-auth and matchmaking services:\n")
	fmt.Printf("JWT_SECRET=%s\n", secret)
}
