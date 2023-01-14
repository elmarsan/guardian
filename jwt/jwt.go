package jwt

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// jwtSecretKey reads secret used for signing jwt tokens.
// If it's missing, program will ends.
func jwtSecretKey() []byte {
	jwtKey, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		log.Fatal("Missing JWT_KEY environment variable")
	}

	return []byte(jwtKey)
}

// jwtExpirationTime return policy expiration time.
// By default 1 hour, it can be overwritten using JWT_EXPIRATION_TIME env var
func jwtExpirationTime() time.Time {
	t := time.Now()
	mins := 60

	expTime, ok := os.LookupEnv("JWT_EXPIRATION_TIME")
	if ok {
		// Ignore parsing errors, in fail case take use time
		parsedMins, _ := strconv.Atoi(expTime)
		if parsedMins > 0 {
			mins = 60
		}
	}

	return t.Add(time.Duration(mins) * time.Minute)
}

// NewToken generates new jwt token valid for the next hour.
func NewToken() (string, error) {
	expirationTime := jwtExpirationTime()
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := jwtSecretKey()
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates given jwt token.
func ValidateToken(token string) error {
	jwtKey := jwtSecretKey()
	jwt, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return err
	}

	if !jwt.Valid {
		return fmt.Errorf("Invalid jwt token (%s)", token)
	}

	return nil
}
