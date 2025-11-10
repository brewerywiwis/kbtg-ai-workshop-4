package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"

	"workshop4-backend/adapter"
	"workshop4-backend/handler"
	"workshop4-backend/models"
	"workshop4-backend/service"
)

// Database connection
var db *sql.DB

// initDatabase initializes the SQLite database and creates the users table
func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "users.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create users table
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		phone TEXT,
		email TEXT,
		member_since TEXT,
		membership_level TEXT,
		member_id TEXT,
		points INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// Insert sample data if table is empty
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatal("Failed to check table:", err)
	}

	if count == 0 {
		insertSampleData()
	}
}

// insertSampleData adds initial sample users to the database
func insertSampleData() {
	sampleUsers := []models.User{
		{
			Name:            "สมชาย ใจดี",
			Phone:           "081-234-5678",
			Email:           "somchai@example.com",
			MemberSince:     "15/6/2566",
			MembershipLevel: "Gold",
			MemberID:        "LBK001234",
			Points:          15420,
		},
		{
			Name:            "สมหญิง ดีใจ",
			Phone:           "081-567-8901",
			Email:           "somying@example.com",
			MemberSince:     "20/7/2566",
			MembershipLevel: "Silver",
			MemberID:        "LBK001235",
			Points:          8500,
		},
	}

	for _, user := range sampleUsers {
		_, err := db.Exec(`
			INSERT INTO users (name, phone, email, member_since, membership_level, member_id, points)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			user.Name, user.Phone, user.Email, user.MemberSince, user.MembershipLevel, user.MemberID, user.Points)
		if err != nil {
			log.Printf("Failed to insert sample user: %v", err)
		}
	}
}

func main() {
	initDatabase()
	defer db.Close()

	dbRepo := adapter.NewSqliteUserRepository(db)
	userService := service.NewUserService(dbRepo)
	userHandler := handler.NewUserHandler(userService)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})

	userHandler.RegisterRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
