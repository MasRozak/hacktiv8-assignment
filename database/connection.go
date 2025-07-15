package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(dbURL string) {
	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database connected successfully")
}

func CreateTables() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		bio TEXT
	);

	CREATE TABLE IF NOT EXISTS posts (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS likes (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL,
		post_id VARCHAR(36) NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		UNIQUE(user_id, post_id)
	);

	CREATE TABLE IF NOT EXISTS comments (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL,
		post_id VARCHAR(36) NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS follows (
		id VARCHAR(36) PRIMARY KEY,
		follower_id VARCHAR(36) NOT NULL,
		following_id VARCHAR(36) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(follower_id, following_id),
		CHECK (follower_id != following_id)
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}
	log.Println("Tables created successfully")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
