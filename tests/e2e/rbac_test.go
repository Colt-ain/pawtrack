package e2e

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRBAC(t *testing.T) {
	clientA := NewTestClient(BaseURL)
	clientA.SetT(t)

	clientA.SetT(t)

	clientB := NewTestClient(BaseURL)
	clientB.SetT(t)

	// Register Owner A
	emailA := fmt.Sprintf("owner_a_%d@example.com", time.Now().UnixNano())
	clientA.RegisterAndLogin("Owner A", emailA, "password123", "owner")

	// Register Owner B
	emailB := fmt.Sprintf("owner_b_%d@example.com", time.Now().UnixNano())
	clientB.RegisterAndLogin("Owner B", emailB, "password123", "owner")

	// Owner A creates Dog
	var dogID float64
	t.Run("Owner A creates Dog", func(t *testing.T) {
		body := map[string]string{
			"name":       "Rex",
			"breed":      "GSD",
			"birth_date": "2020-01-01T00:00:00Z",
		}
		var resp map[string]interface{}
		status := clientA.Post("/dogs", body, &resp)
		require.Equal(t, 201, status)
		dogID = resp["id"].(float64)
	})

	// Owner B tries to get Owner A's Dog
	t.Run("Owner B cannot see Owner A's Dog", func(t *testing.T) {
		status := clientB.Get(fmt.Sprintf("/dogs/%.0f", dogID), nil)
		require.Equal(t, 404, status)
	})

	// Consultant Access
	clientC := NewTestClient(BaseURL)
	clientC.SetT(t)
	emailC := fmt.Sprintf("consultant_%d@example.com", time.Now().UnixNano())

	t.Run("Register Consultant", func(t *testing.T) {
		body := map[string]string{
			"name":     "Dr. Vet",
			"email":    emailC,
			"password": "password123",
		}
		status := clientC.Post("/auth/register/consultant", body, nil)
		require.Equal(t, 201, status)

		// Login
		loginBody := map[string]string{
			"email":    emailC,
			"password": "password123",
		}
		var resp map[string]interface{}
		clientC.Post("/auth/login", loginBody, &resp)
		clientC.SetToken(resp["token"].(string))
	})

	t.Run("Consultant cannot see Dog initially", func(t *testing.T) {
		status := clientC.Get(fmt.Sprintf("/dogs/%.0f", dogID), nil)
		require.Equal(t, 403, status)
	})

	t.Run("Grant Access and Verify", func(t *testing.T) {
		// Grant access via DB
		grantAccess(t, emailC, dogID)

		// Verify access
		var resp map[string]interface{}
		status := clientC.Get(fmt.Sprintf("/dogs/%.0f", dogID), &resp)
		require.Equal(t, 200, status)
	})
}

func grantAccess(t *testing.T, consultantEmail string, dogID float64) {
	dsn := "postgres://pawtrack:pawtrack@localhost:5432/pawtrack?sslmode=disable"
	if envDSN := os.Getenv("E2E_DB_DSN"); envDSN != "" {
		dsn = envDSN
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// Get Consultant ID
	var user struct {
		ID uint
	}
	err = db.Table("users").Where("email = ?", consultantEmail).First(&user).Error
	require.NoError(t, err)

	// Insert Access
	err = db.Exec("INSERT INTO consultant_access (consultant_id, dog_id, granted_at) VALUES (?, ?, NOW())", user.ID, dogID).Error
	require.NoError(t, err)

	// Grant Permissions
	permissions := []string{
		"DOGS_VIEW_ASSIGNED",
		"EVENTS_CREATE_ASSIGNED",
		"EVENTS_VIEW_ASSIGNED",
		"EVENT_COMMENTS_CREATE_ASSIGNED",
		"EVENT_COMMENTS_VIEW_ASSIGNED",
		"CONSULTANT_NOTES_CREATE",
		"CONSULTANT_NOTES_VIEW_OWN",
		"CONSULTANT_NOTES_UPDATE_OWN",
	}

	for _, perm := range permissions {
		var permID int
		err = db.Table("permissions").Select("id").Where("name = ?", perm).Scan(&permID).Error
		require.NoError(t, err)

		err = db.Exec("INSERT INTO user_permissions (user_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING", user.ID, permID).Error
		require.NoError(t, err)
	}
}
