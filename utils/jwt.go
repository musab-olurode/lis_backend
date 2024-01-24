package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GetJWT(id string, time *jwt.NumericDate) (string, error) {
	signingKey := []byte(os.Getenv("JWT_SECRET"))

	claims := &jwt.RegisteredClaims{
		ID:        id,
		ExpiresAt: time,
		Issuer:    "LIS",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func GetAuthToken(id string) (string, error) {
	jwtExpiry := os.Getenv("JWT_AUTH_TOKEN_EXPIRY")

	duration, err := time.ParseDuration(jwtExpiry)
	if err != nil {
		return "", err
	}

	time := jwt.NewNumericDate(time.Now().Add(duration))
	jwtToken, err := GetJWT(id, time)

	return jwtToken, err
}

func GetRefreshToken(id string) (string, error) {
	jwtExpiry := os.Getenv("JWT_REFRESH_TOKEN_EXPIRY")

	duration, err := time.ParseDuration(jwtExpiry)
	if err != nil {
		return "", err
	}

	time := jwt.NewNumericDate(time.Now().Add(duration))
	jwtToken, err := GetJWT(id, time)

	return jwtToken, err
}

func GetJWTClaims(tokenString string) (*jwt.MapClaims, error) {
	signingKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("invalid authorization")
		}
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &claims, nil
	} else {
		log.Printf("error getting token claims")
		return nil, fmt.Errorf("invalid authorization")
	}
}
