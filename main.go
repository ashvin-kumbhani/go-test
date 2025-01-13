package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// User struct for representing user data
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

	// Route to serve the SSE stream
	app.Get("/events", sseHandler)

	// Route to handle CSV upload
	app.Post("/upload", uploadCSV)

	log.Fatal(app.Listen(":3000"))
}

// sseHandler sends user updates as Server-Sent Events
func sseHandler(c *fiber.Ctx) error {
	// Set headers for SSE
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	// Continuously send data (simulate user data stream)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Simulating a new user every tick
			user := User{
				ID:           1,
				FirstName:    "John",
				LastName:     "Doe",
				Email:        "john.doe@example.com",
				CreatedAt:    time.Now().Unix(),
				DeletedAt:    0,
				MergedAt:     0,
				ParentUserID: 0,
			}
			userJSON, err := json.Marshal(user)
			if err != nil {
				log.Println("Error marshalling user:", err)
				continue
			}

			// Send user data as SSE
			c.WriteString("data: " + string(userJSON) + "\n\n")
		}
	}
}

// uploadCSV handler to handle CSV file upload (similar to your existing code)
func uploadCSV(c *fiber.Ctx) error {
	// Your existing CSV upload logic
	return c.SendString("File processed and data sent to queue")
}
