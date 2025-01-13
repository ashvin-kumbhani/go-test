package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var (
	db  *sql.DB
	rdb *redis.Client
)

func TestPostgresConnection(t *testing.T) {
	var err error

	// Connect to PostgreSQL
	db, err = sql.Open("postgres", "postgres://postgres:postgres@postgres:5432/userdb?sslmode=disable")
	assert.NoError(t, err, "Failed to connect to PostgreSQL")

	// Verify PostgreSQL connection
	err = db.Ping()
	assert.NoError(t, err, "PostgreSQL connection error")
}

func TestPostgresTableCreation(t *testing.T) {
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

	_, err := db.Exec(createTableQuery)
	assert.NoError(t, err, "Failed to create table")
}

func TestRedisConnection(t *testing.T) {
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	assert.NoError(t, err, "Redis connection error")
}
