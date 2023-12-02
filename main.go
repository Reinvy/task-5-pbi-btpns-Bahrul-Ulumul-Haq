package main

import (
	"upload-photos/database"
	"upload-photos/router"
)

func main() {
	// Initialize the database connection
	database.InitDB()

	r := router.InitRouter()

	// Run the server
	r.Run()
}
