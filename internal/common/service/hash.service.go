package CommonService

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashString(input string) (string, error) {
	start := time.Now()

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(input),
		bcrypt.DefaultCost,
	)

	log.Printf("password hashing: %v", time.Since(start))
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CompareHash(hash, input string) error {
	start := time.Now()

	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(input),
	)

	log.Printf("password hash comparison: %v", time.Since(start))
	if err != nil {
		return err
	}

	return nil
}

func GenerateRandomToken(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	token := make([]byte, length)

	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("generate random token: %w", err)
	}

	for i := range token {
		token[i] = charset[int(randomBytes[i])%len(charset)]
	}

	return string(token), nil
}
