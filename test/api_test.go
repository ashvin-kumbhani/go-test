package main

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func getUsers(c *fiber.Ctx) error {
	// Mock response for testing
	users := []string{"John Doe", "Jane Doe"}
	return c.JSON(users)
}

func TestGetUsers(t *testing.T) {
	app := fiber.New()
	app.Get("/users", getUsers)

	tests := []struct {
		description  string
		queryParams  string
		expectedCode int
	}{
		{
			description:  "Valid request with no filters",
			queryParams:  "",
			expectedCode: 200,
		},
		{
			description:  "Valid request with first_name filter",
			queryParams:  "?first_name=John",
			expectedCode: 200,
		},
		{
			description:  "Valid request with last_name filter",
			queryParams:  "?last_name=Doe",
			expectedCode: 200,
		},
		{
			description:  "Valid request with email filter",
			queryParams:  "?email=john.doe@example.com",
			expectedCode: 200,
		},
		{
			description:  "Invalid page value",
			queryParams:  "?page=0",
			expectedCode: 400,
		},
		{
			description:  "Invalid page_size value",
			queryParams:  "?page_size=0",
			expectedCode: 400,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/users"+test.queryParams, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err, test.description)
		assert.Equal(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
