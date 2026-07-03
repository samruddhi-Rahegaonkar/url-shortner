package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitSchema() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		original_url TEXT NOT NULL,
		short_url VARCHAR(50) UNIQUE NOT NULL,
		creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		user_id INTEGER REFERENCES users(id)
	);
	`

	_, err := DB.Exec(context.Background(), query)
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Database schema initialized")
}
func ConnectDB() {
	var err error
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("Database URL is not set")
	}
	DB, err = pgxpool.New(
		context.Background(),
		dbURL,
	)
	if err != nil {
		panic(err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Database Connected successfully")
	InitSchema()
	// Print current database and schema
	var dbName, schema string
	err = DB.QueryRow(
		context.Background(),
		"SELECT current_database(), current_schema()",
	).Scan(&dbName, &schema)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database :", dbName)
	fmt.Println("Schema   :", schema)

	// Check if users table exists
	var exists bool
	err = DB.QueryRow(
		context.Background(),
		`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema='public'
			AND table_name='users'
		)
		`,
	).Scan(&exists)

	if err != nil {
		panic(err)
	}

	fmt.Println("Users table exists:", exists)

	// Count users
	var count int
	err = DB.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM users",
	).Scan(&count)

	if err != nil {
		fmt.Println("❌ Error while counting users:", err)
	} else {
		fmt.Println("Number of users:", count)
	}
}
