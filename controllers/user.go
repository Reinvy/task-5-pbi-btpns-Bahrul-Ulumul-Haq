package controllers

import (
	"net/http"

	"upload-photos/app"
	"upload-photos/database"
	"upload-photos/helpers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterUser registers a new user
func RegisterUser(c *gin.Context) {
	// Bind the request body to a user struct
	var user app.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the user fields
	if err := user.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the user password
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save the user to the database
	db := database.GetDB()
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate a JWT token for the user
	token, err := helpers.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the user and the token
	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

// LoginUser logs in a user
func LoginUser(c *gin.Context) {
	// Bind the request body to a user struct
	var user app.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var rawPassword = user.Password

	// Get the user from the database by email
	db := database.GetDB()
	if err := db.Where("email = ?", user.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	println(rawPassword)
	println(user.Password)
	// Check the user password
	if err := user.CheckPassword(user.Password, rawPassword); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate a JWT token for the user
	token, err := helpers.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the user and the token
	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

// UpdateUser updates a user
func UpdateUser(c *gin.Context) {
	// Get the user ID from the URL parameter
	userID := c.Param("userId")

	// Get the user from the database by ID
	db := database.GetDB()
	var user app.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Get the current user from the context
	currentUser := c.MustGet("currentUser").(app.User)

	// Check if the current user is the same as the user to be updated
	if currentUser.ID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this user"})
		return
	}

	// Bind the request body to a user struct
	var updatedUser app.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the updated user fields
	if err := updatedUser.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the updated user password if changed
	if updatedUser.Password != "" {
		if err := updatedUser.HashPassword(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Update the user in the database
	if err := db.Model(&user).Updates(updatedUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated user
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// DeleteUser deletes a user
func DeleteUser(c *gin.Context) {
	// Get the user ID from the URL parameter
	userID := c.Param("userId")

	// Get the user from the database by ID
	db := database.GetDB()
	var user app.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Get the current user from the context
	currentUser := c.MustGet("currentUser").(app.User)

	// Check if the current user is the same as the user to be deleted
	if currentUser.ID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this user"})
		return
	}

	// Delete the user from the database
	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
