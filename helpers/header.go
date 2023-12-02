package helpers

import (
	"net/http"
	"strings"
)

// GetBearerToken gets the bearer token from the authorization header
func GetBearerToken(r *http.Request) string {
	// Get the authorization header value
	authHeader := r.Header.Get("Authorization")

	// Split the header value by space
	parts := strings.Split(authHeader, " ")

	// Check if the header value has two parts and the first part is "Bearer"
	if len(parts) == 2 && parts[0] == "Bearer" {
		// Return the second part as the token
		return parts[1]
	}

	// Return an empty string if the header value is not valid
	return ""
}
