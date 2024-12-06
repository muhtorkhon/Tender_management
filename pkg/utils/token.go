package utils

import (
	"tender_management/config"
	"tender_management/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


func GenerateToken(email, role string) (string, error) {
	claims := &models.Claims{
		Email: email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var key config.Config
	return token.SignedString(key.SecretKey)
}

func ValidateToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}
	key := &config.Config{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key.SecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}