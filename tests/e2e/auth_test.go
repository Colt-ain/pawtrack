package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	client := NewTestClient(t, BaseURL)

	timestamp := time.Now().UnixNano()
	email := fmt.Sprintf("test_auth_%d@example.com", timestamp)
	password := "password123"

	t.Run("Register Owner", func(t *testing.T) {
		body := map[string]string{
			"name":     "Test Owner",
			"email":    email,
			"password": password,
		}
		status := client.Post("/auth/register/owner", body, nil)
		require.Equal(t, 201, status)
	})

	t.Run("Login", func(t *testing.T) {
		body := map[string]string{
			"email":    email,
			"password": password,
		}
		var resp map[string]interface{}
		status := client.Post("/auth/login", body, &resp)
		require.Equal(t, 200, status)

		token, ok := resp["token"].(string)
		require.True(t, ok, "token should be a string")
		require.NotEmpty(t, token, "token should not be empty")
		
		client.SetToken(token)
	})

	t.Run("Login Invalid Password", func(t *testing.T) {
		body := map[string]string{
			"email":    email,
			"password": "wrongpassword",
		}
		status := client.Post("/auth/login", body, nil)
		require.Equal(t, 401, status)
	})
}
