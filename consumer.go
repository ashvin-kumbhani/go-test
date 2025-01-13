package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
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
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"csv_data_queue", // Name
		true,             // Durable
		false,            // Delete when unused
		false,            // Exclusive
		false,            // No-wait
		nil,              // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare RabbitMQ queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer name
		true,   // Auto-acknowledge
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	// Process messages
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var user User
			err := json.Unmarshal(d.Body, &user)
			if err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}

			// Store in PostgreSQL
			_, err = db.Exec(`
				INSERT INTO users (id, first_name, last_name, email, created_at, deleted_at, merged_at, parent_user_id)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
				user.ID, user.FirstName, user.LastName, user.Email, user.CreatedAt, user.DeletedAt, user.MergedAt, user.ParentUserID,
			)
			if err != nil {
				log.Printf("Error inserting into PostgreSQL: %v", err)
				continue
			}

			// Store in Redis
			err = rdb.Set(ctx, user.Email, d.Body, 0).Err()
			if err != nil {
				log.Printf("Error storing in Redis: %v", err)
			}
		}
	}()

	log.Println("Waiting for messages...")
	<-forever
}
