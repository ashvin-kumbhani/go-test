package main

import (
	"encoding/json"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	CreatedAt    int64  `json:"created_at"`
	DeletedAt    int64  `json:"deleted_at"`
	MergedAt     int64  `json:"merged_at"`
	ParentUserID int    `json:"parent_user_id"`
}

func TestConsumer(t *testing.T) {
	// Mock RabbitMQ connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost/")
	assert.NoError(t, err)
	defer conn.Close()

	ch, err := conn.Channel()
	assert.NoError(t, err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"csv_data_queue", // Name
		true,             // Durable
		false,            // Delete when unused
		false,            // Exclusive
		false,            // No-wait
		nil,              // Arguments
	)
	assert.NoError(t, err)

	// Mock message
	user := User{
		ID:           1,
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "john.doe@example.com",
		CreatedAt:    1627847281,
		DeletedAt:    0,
		MergedAt:     0,
		ParentUserID: 0,
	}
	body, err := json.Marshal(user)
	assert.NoError(t, err)

	err = ch.Publish(
		"",     // Exchange
		q.Name, // Routing key
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	assert.NoError(t, err)

	// Consume message
	msgs, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer name
		true,   // Auto-acknowledge
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Arguments
	)
	assert.NoError(t, err)

	// Process message
	for d := range msgs {
		var receivedUser User
		err := json.Unmarshal(d.Body, &receivedUser)
		assert.NoError(t, err)
		assert.Equal(t, user, receivedUser)
		break
	}
}
