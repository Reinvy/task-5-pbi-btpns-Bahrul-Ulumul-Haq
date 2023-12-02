package helpers

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// GenerateToken generates a JWT token for a user ID
func GenerateToken(userID uint) (string, error) {
	// Get the secret key from the environment variable
	secret := os.Getenv("JWT_SECRET")

	// Create a token with the user ID as the claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	// Return the signed token
	return signedToken, nil
}

// ParseToken parses a JWT token and returns the user ID
func ParseToken(tokenString string) (uint, error) {
	// Get the secret key from the environment variable
	secret := os.Getenv("JWT_SECRET")

	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	// Get the user ID from the token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}

	// Return an error if the token is invalid
	return 0, jwt.ErrInvalidKey
}
