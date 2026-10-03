package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	SecretKey []byte
}

func NewJwtService(secretKey string) JwtService {
	return JwtService{
		SecretKey: []byte(secretKey),
	}
}

func (s *JwtService) GenerateToken(organizationId, userId string) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"userId":         userId,
		"organizationId": organizationId,
		"iat":            now.Unix(),
		"exp":            now.Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.SecretKey)
}
