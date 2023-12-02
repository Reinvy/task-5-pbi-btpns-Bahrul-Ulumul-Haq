package middlewares

import (
	"net/http"

	"upload-photos/app"
	"upload-photos/database"
	"upload-photos/helpers"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware that checks the JWT token and sets the current user in the context
func AuthMiddleware(c *gin.Context) {
	// Get the bearer token from the header
	token := helpers.GetBearerToken(c.Request)

	// Parse the token and get the user ID
	userID, err := helpers.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get the user from the database by ID
	db := database.GetDB()
	var user app.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Set the current user in the context
	c.Set("currentUser", user)

	// Proceed to the next handler
	c.Next()
}
