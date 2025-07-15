package token

import (
	"fmt"
	"marketplace-service/internal/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	signingKey     []byte
	expirationTime time.Duration
	logger         logger.Logger
}

func NewService(secretKey string, expirationTime time.Duration, l logger.Logger) *Service {
	return &Service{
		signingKey:     []byte(secretKey),
		expirationTime: expirationTime,
		logger: l,
	}
}

func (s *Service) GenerateToken(userID string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expirationTime)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.signingKey)
}

func (s *Service) ValidateToken(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.signingKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	return claims.Subject, nil
}
