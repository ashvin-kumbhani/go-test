package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	rdb *redis.Client
	ctx = context.Background()
)

func init() {
	var err error

	// Connect to PostgreSQL
	db, err = sql.Open("postgres", "postgres://postgres:postgres@postgres:5432/userdb?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Verify PostgreSQL connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("PostgreSQL connection error: %v", err)
	}
	log.Println("Connected to PostgreSQL successfully.")

	// Create table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(100),
		last_name VARCHAR(100),
		email VARCHAR(150) UNIQUE,
		created_at BIGINT,
		deleted_at BIGINT,
		merged_at BIGINT,
		parent_user_id INT
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	log.Println("Users table ensured in PostgreSQL.")

	// Connect to Redis
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	// Verify Redis connection
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis connection error: %v", err)
	}
	log.Println("Connected to Redis successfully.")
}
