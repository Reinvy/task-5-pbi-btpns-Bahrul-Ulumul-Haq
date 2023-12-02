package database

import (
	"log"

	"upload-photos/app"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB initializes	 the database connection
func InitDB() {
	// Get the database configuration from the environment variables
	// host := os.Getenv("DB_HOST")
	// port := os.Getenv("DB_PORT")
	// user := os.Getenv("DB_USER")
	// password := os.Getenv("DB_PASSWORD")
	// name := os.Getenv("DB_NAME")

	// Build the connection string
	// dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
	dsn := "host=localhost port=5432 user=postgres password=root dbname=upload-photos sslmode=disable"

	// Connect to the database
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the database schema
	db.AutoMigrate(&app.User{}, &app.Photo{})
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}
