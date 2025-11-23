package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEvents(t *testing.T) {
	client := NewTestClient(t, BaseURL)

	timestamp := time.Now().UnixNano()
	email := fmt.Sprintf("test_event_%d@example.com", timestamp)
	client.RegisterAndLogin("Event Owner", email, "password123")

	var dogID float64

	t.Run("Create Dog for Events", func(t *testing.T) {
		body := map[string]string{
			"name":       "Buddy",
			"breed":      "Golden Retriever",
			"birth_date": "2021-01-01T00:00:00Z",
		}
		var resp map[string]interface{}
		client.Post("/dogs", body, &resp)
		dogID = resp["id"].(float64)
	})

	t.Run("Create Events", func(t *testing.T) {
		events := []map[string]interface{}{
			{"dog_id": dogID, "type": "walk", "note": "Morning walk", "at": "2025-01-01T08:00:00Z"},
			{"dog_id": dogID, "type": "feed", "note": "Dinner", "at": "2025-01-01T18:00:00Z"},
		}
		for _, e := range events {
			status := client.Post("/events", e, nil)
			require.Equal(t, 201, status)
		}
	})

	t.Run("List Events", func(t *testing.T) {
		var resp map[string]interface{}
		status := client.Get("/events", &resp)
		require.Equal(t, 200, status)

		events := resp["events"].([]interface{})
		require.Len(t, events, 2)
		require.Equal(t, float64(2), resp["total_count"])
	})

	t.Run("Filter Events by Type", func(t *testing.T) {
		var resp map[string]interface{}
		status := client.Get("/events?types=walk", &resp)
		require.Equal(t, 200, status)

		events := resp["events"].([]interface{})
		require.Len(t, events, 1)
		require.Equal(t, "walk", events[0].(map[string]interface{})["type"])
	})

	t.Run("Search Events", func(t *testing.T) {
		var resp map[string]interface{}
		status := client.Get("/events?search=Dinner", &resp)
		require.Equal(t, 200, status)

		events := resp["events"].([]interface{})
		require.Len(t, events, 1)
		require.Equal(t, "feed", events[0].(map[string]interface{})["type"])
	})
}
