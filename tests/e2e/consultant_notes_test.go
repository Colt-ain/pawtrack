package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConsultantNotes(t *testing.T) {
	client := NewTestClient(BaseURL)
	client.SetT(t)

	// Setup: Create consultant with dog access
	consultantEmail := fmt.Sprintf("consultant_notes_%d@example.com", time.Now().UnixNano())
	consultantToken, err := client.RegisterAndLogin("Consultant Notes", consultantEmail, "password", "consultant")
	require.NoError(t, err)
	client.SetToken(consultantToken)

	// Update profile to get consultant ID
	profileData := map[string]interface{}{
		"description": "Tester",
		"services":    "Testing",
		"breeds":      "All",
		"location":    "Test City",
		"surname":     "Tester",
	}
	var profileResp map[string]interface{}
	status := client.Put("/consultants/profile", profileData, &profileResp)
	require.Equal(t, http.StatusOK, status)
	consultantID := uint(profileResp["user_id"].(float64))

	// Create owner and dog
	ownerEmail := fmt.Sprintf("owner_notes_%d@example.com", time.Now().UnixNano())
	ownerToken, err := client.RegisterAndLogin("Owner Notes", ownerEmail, "password", "owner")
	require.NoError(t, err)
	client.SetToken(ownerToken)

	dogID, err := client.CreateDog("TestDog", "Labrador", "2020-01-01T00:00:00Z")
	require.NoError(t, err)

	// Owner invites consultant
	inviteReq := map[string]interface{}{
		"dog_id": dogID,
	}
	var inviteResp map[string]interface{}
	status = client.Post(fmt.Sprintf("/consultants/%d/invite", consultantID), inviteReq, &inviteResp)
	require.Equal(t, http.StatusCreated, status)

	inviteToken := inviteResp["token"].(string)

	// Consultant accepts invite
	client.SetToken(consultantToken)
	var acceptResp map[string]interface{}
	status = client.Post(fmt.Sprintf("/invites/accept?token=%s", inviteToken), nil, &acceptResp)
	require.Equal(t, http.StatusOK, status)

	var noteID float64

	t.Run("Create Note", func(t *testing.T) {
		noteReq := map[string]interface{}{
			"dog_id":  dogID,
			"title":   "First Session",
			"content": "# Session Notes\n\nDog responded well to training.",
		}
		var resp map[string]interface{}
		status := client.Post("/consultant-notes", noteReq, &resp)
		require.Equal(t, http.StatusCreated, status)
		require.Equal(t, "First Session", resp["title"])
		noteID = resp["id"].(float64)
	})

	t.Run("Get Note", func(t *testing.T) {
		var resp map[string]interface{}
		status := client.Get(fmt.Sprintf("/consultant-notes/%.0f", noteID), &resp)
		require.Equal(t, http.StatusOK, status)
		require.Equal(t, "First Session", resp["title"])
		require.Contains(t, resp["content"], "Session Notes")
	})

	t.Run("Update Note", func(t *testing.T) {
		updateReq := map[string]interface{}{
			"title":   "Updated Session",
			"content": "# Updated Notes\n\nAdded more observations.",
		}
		var resp map[string]interface{}
		status := client.Put(fmt.Sprintf("/consultant-notes/%.0f", noteID), updateReq, &resp)
		require.Equal(t, http.StatusOK, status)
		require.Equal(t, "Updated Session", resp["title"])
	})

	t.Run("List Notes", func(t *testing.T) {
		var resp map[string]interface{}
		status := client.Get("/consultant-notes", &resp)
		require.Equal(t, http.StatusOK, status)
		
		notes := resp["notes"].([]interface{})
		require.NotEmpty(t, notes)
		require.GreaterOrEqual(t, int(resp["total_count"].(float64)), 1)
	})

	t.Run("Filter by Dog", func(t *testing.T) {
		var resp map[string]interface{}
		status := client.Get(fmt.Sprintf("/consultant-notes?dog_id=%.0f", float64(dogID)), &resp)
		require.Equal(t, http.StatusOK, status)
		
		notes := resp["notes"].([]interface{})
		require.NotEmpty(t, notes)
	})

	t.Run("Search in Content", func(t *testing.T) {
		var resp map[string]interface{}
		status := client.Get("/consultant-notes?search=observations", &resp)
		require.Equal(t, http.StatusOK, status)
		
		notes := resp["notes"].([]interface{})
		require.NotEmpty(t, notes)
	})

	t.Run("Sort by Updated Date", func(t *testing.T) {
		var resp map[string]interface{}
		status := client.Get("/consultant-notes?sort_by=updated_at&order=desc", &resp)
		require.Equal(t, http.StatusOK, status)
		
		notes := resp["notes"].([]interface{})
		require.NotEmpty(t, notes)
	})

	t.Run("Delete Note", func(t *testing.T) {
		status := client.Delete(fmt.Sprintf("/consultant-notes/%.0f", noteID))
		require.Equal(t, http.StatusNoContent, status)
		
		// Verify deleted
		var resp map[string]interface{}
		status = client.Get(fmt.Sprintf("/consultant-notes/%.0f", noteID), &resp)
		require.Equal(t, http.StatusNotFound, status)
	})

	t.Run("Cannot Create Note Without Access", func(t *testing.T) {
		// Create another dog that consultant doesn't have access to
		client.SetToken(ownerToken)
		dogID2, err := client.CreateDog("AnotherDog", "Poodle", "2021-01-01T00:00:00Z")
		require.NoError(t, err)

		client.SetToken(consultantToken)
		noteReq := map[string]interface{}{
			"dog_id":  dogID2,
			"title":   "Unauthorized",
			"content": "Should fail",
		}
		var resp map[string]interface{}
		status := client.Post("/consultant-notes", noteReq, &resp)
		require.Equal(t, http.StatusForbidden, status)
	})
}
