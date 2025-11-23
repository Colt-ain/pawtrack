package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	client := NewTestClient(BaseURL)
	client.SetT(t)

	timestamp := time.Now().UnixNano()
	email := fmt.Sprintf("test_auth_%d@example.com", timestamp)
	t.Run("Register Owner", func(t *testing.T) {
		token, err := client.RegisterAndLogin("Test Owner", email, "password123", "owner")
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("Login", func(t *testing.T) {
		// Login is already covered by RegisterAndLogin which verifies login works.
		// But if we want to test explicit login again:
		loginBody := map[string]string{
			"email":    email,
			"password": "password123",
		}
		var resp map[string]interface{}
		status := client.Post("/auth/login", loginBody, &resp)
		require.Equal(t, 200, status)
		require.NotEmpty(t, resp["token"])
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
