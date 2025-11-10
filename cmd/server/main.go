package main

import (
	"log"

	"workshop4-backend/internal/app"
)

func main() {
	// Initialize database
	db := app.InitDatabase()
	defer db.Close()

	// Setup and start server
	server := app.SetupServer(db)
	log.Fatal(server.Listen(":3000"))
}