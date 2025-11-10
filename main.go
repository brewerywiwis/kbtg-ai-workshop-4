package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"

	"workshop4-backend/adapter"
	"workshop4-backend/domain"
	"workshop4-backend/handler"
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
	createUsersTable := `
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

	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	// Create transfers table
	createTransfersTable := `
	CREATE TABLE IF NOT EXISTS transfers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		from_user_id INTEGER NOT NULL,
		to_user_id INTEGER NOT NULL,
		amount INTEGER NOT NULL CHECK (amount > 0),
		status TEXT NOT NULL CHECK (status IN ('pending','processing','completed','failed','cancelled','reversed')),
		note TEXT,
		idempotency_key TEXT NOT NULL UNIQUE,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		completed_at TEXT,
		fail_reason TEXT,
		FOREIGN KEY (from_user_id) REFERENCES users(id),
		FOREIGN KEY (to_user_id) REFERENCES users(id)
	);`

	_, err = db.Exec(createTransfersTable)
	if err != nil {
		log.Fatal("Failed to create transfers table:", err)
	}

	// Create transfer indexes
	transferIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_transfers_from ON transfers(from_user_id);",
		"CREATE INDEX IF NOT EXISTS idx_transfers_to ON transfers(to_user_id);",
		"CREATE INDEX IF NOT EXISTS idx_transfers_created ON transfers(created_at);",
	}

	for _, idx := range transferIndexes {
		_, err = db.Exec(idx)
		if err != nil {
			log.Fatal("Failed to create transfer index:", err)
		}
	}

	// Create point_ledger table
	createLedgerTable := `
	CREATE TABLE IF NOT EXISTS point_ledger (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		change INTEGER NOT NULL,
		balance_after INTEGER NOT NULL,
		event_type TEXT NOT NULL CHECK (event_type IN ('transfer_out','transfer_in','adjust','earn','redeem')),
		transfer_id INTEGER,
		reference TEXT,
		metadata TEXT,
		created_at TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (transfer_id) REFERENCES transfers(id)
	);`

	_, err = db.Exec(createLedgerTable)
	if err != nil {
		log.Fatal("Failed to create point_ledger table:", err)
	}

	// Create ledger indexes
	ledgerIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_ledger_user ON point_ledger(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_ledger_transfer ON point_ledger(transfer_id);",
		"CREATE INDEX IF NOT EXISTS idx_ledger_created ON point_ledger(created_at);",
	}

	for _, idx := range ledgerIndexes {
		_, err = db.Exec(idx)
		if err != nil {
			log.Fatal("Failed to create ledger index:", err)
		}
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
	sampleUsers := []domain.User{
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

	// Initialize repositories
	userRepo := adapter.NewSqliteUserRepository(db)
	transferRepo := adapter.NewSqliteTransferRepository(db)
	ledgerRepo := adapter.NewSqlitePointLedgerRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	transferService := service.NewTransferService(transferRepo, ledgerRepo, userRepo, db)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	transferHandler := handler.NewTransferHandler(transferService)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})

	// Register routes
	userHandler.RegisterRoutes(app)
	transferHandler.RegisterRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
