package controllers

import (
	"net/http"

	"upload-photos/app"
	"upload-photos/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreatePhoto creates a new photo
func CreatePhoto(c *gin.Context) {
	// Bind the request body to a photo struct

	var photoInput app.PhotoInput
	var photo app.Photo
	var user app.User

	if err := c.ShouldBindJSON(&photoInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	if err := db.Where("id = ?", photoInput.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	photo.Title = photoInput.Title
	photo.Caption = photoInput.Caption
	photo.PhotoUrl = photoInput.PhotoUrl
	photo.User = user

	// Validate the photo fields
	if err := photo.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the current user from the context
	currentUser := c.MustGet("currentUser").(app.User)

	// Set the user ID of the photo to the current user ID
	photo.UserID = currentUser.ID

	// Save the photo to the database
	if err := db.Create(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the photo
	c.JSON(http.StatusOK, gin.H{"photo": photo})
}

// GetPhotos gets all photos
func GetPhotos(c *gin.Context) {
	// Get all photos from the database
	db := database.GetDB()
	var photos []app.Photo
	if err := db.Preload("User").Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the photos
	c.JSON(http.StatusOK, gin.H{"photos": photos})
}

// UpdatePhoto updates a photo
func UpdatePhoto(c *gin.Context) {
	// Get the photo ID from the URL parameter
	photoID := c.Param("photoId")

	// Get the photo from the database by ID
	var photoInput app.PhotoInput
	var photoUpdate app.Photo
	var photo app.Photo
	var user app.User

	if err := c.ShouldBindJSON(&photoInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	if err := db.Preload("User").Where("id = ?", photoID).First(&photo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if photo.UserID != photoInput.UserID {
		c.JSON(http.StatusNotFound, gin.H{"error": "User ID Invalid"})
		return
	}

	if err := db.Where("id = ?", photoInput.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Get the current user from the context
	currentUser := c.MustGet("currentUser").(app.User)

	// Check if the current user is the same as the user who created the photo
	if currentUser.ID != photo.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this photo"})
		return
	}

	photoUpdate.Title = photoInput.Title
	photoUpdate.Caption = photoInput.Caption
	photoUpdate.PhotoUrl = photoInput.PhotoUrl
	photoUpdate.User = user

	// Validate the updated photo fields
	if err := photoUpdate.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the photo in the database
	if err := db.Model(&photo).Updates(photoUpdate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated photo
	c.JSON(http.StatusOK, gin.H{"photo": photo})
}

// DeletePhoto deletes a photo
func DeletePhoto(c *gin.Context) {
	// Get the photo ID from the URL parameter
	photoID := c.Param("photoId")

	// Get the photo from the database by ID
	db := database.GetDB()
	var photo app.Photo
	if err := db.Where("id = ?", photoID).First(&photo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Get the current user from the context
	currentUser := c.MustGet("currentUser").(app.User)

	// Check if the current user is the same as the user who created the photo
	if currentUser.ID != photo.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this photo"})
		return
	}

	// Delete the photo from the database
	if err := db.Delete(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}
