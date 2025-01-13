package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email_address"`
	CreatedAt    int64  `json:"created_at"`
	DeletedAt    int64  `json:"deleted_at"`
	MergedAt     int64  `json:"merged_at"`
	ParentUserID int    `json:"parent_user_id"`
}

func main() {
	app := fiber.New()

	app.Get("/users", getUsers)

	log.Fatal(app.Listen(":3001"))
}

func getUsers(c *fiber.Ctx) error {
	// Query parameters
	firstName := c.Query("first_name")
	lastName := c.Query("last_name")
	email := c.Query("email")

	// Pagination parameters
	page, err := strconv.Atoi(c.Query("page", "1")) // Default to page 1
	if err != nil || page < 1 {
		return c.Status(400).SendString("Invalid page value")
	}

	pageSize, err := strconv.Atoi(c.Query("page_size", "10")) // Default to 10 items per page
	if err != nil || pageSize < 1 {
		return c.Status(400).SendString("Invalid page_size value")
	}

	// Use SCAN to iterate over keys in Redis
	var cursor uint64
	var users []User
	start := (page - 1) * pageSize
	end := start + pageSize

	for {
		// SCAN command to fetch keys
		keys, nextCursor, err := rdb.Scan(ctx, cursor, "*", 100).Result()
		if err != nil {
			log.Printf("Error scanning Redis keys: %v\n", err)
			return c.Status(500).SendString("Error reading from Redis")
		}
		cursor = nextCursor

		// Fetch and filter users
		for _, key := range keys {
			val, err := rdb.Get(ctx, key).Result()
			if err != nil {
				log.Printf("Error fetching key %s from Redis: %v\n", key, err)
				continue
			}

			// Deserialize JSON
			var user User
			if err := json.Unmarshal([]byte(val), &user); err != nil {
				log.Printf("Error unmarshalling Redis value: %v\n", err)
				continue
			}

			// Apply filters (if any)
			if firstName != "" && user.FirstName != firstName {
				continue
			}
			if lastName != "" && user.LastName != lastName {
				continue
			}
			if email != "" && user.Email != email {
				continue
			}

			// Add user to results
			users = append(users, user)

			// Stop if we have enough results for the current page
			if len(users) >= end {
				break
			}
		}

		// Stop if no more keys or we have enough results
		if cursor == 0 || len(users) >= end {
			break
		}
	}

	// Paginate results
	if start > len(users) {
		users = []User{} // No results for the requested page
	} else {
		users = users[start:end]
	}

	// Return paginated results
	return c.JSON(fiber.Map{
		"page":      page,
		"page_size": pageSize,
		"data":      users,
	})
}
