package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEventComments(t *testing.T) {
	client := NewTestClient(BaseURL)
	client.SetT(t)

	// Setup: Create owner with dog and event
	ownerEmail := fmt.Sprintf("owner_comments_%d@example.com", time.Now().UnixNano())
	ownerToken, err := client.RegisterAndLogin("Owner Comments", ownerEmail, "password", "owner")
	require.NoError(t, err)
	client.SetToken(ownerToken)

	dogID, err := client.CreateDog("CommentDog", "Retriever", "2020-01-01T00:00:00Z")
	require.NoError(t, err)

	// Create an event
	eventReq := map[string]interface{}{
		"dog_id": dogID,
		"type":   "walk",
		"note":   "Test walk for comments",
		"at":     "2025-11-23T10:00:00Z",
	}
	var eventResp map[string]interface{}
	status := client.Post("/events", eventReq, &eventResp)
	require.Equal(t, http.StatusCreated, status)
	eventID := uint(eventResp["id"].(float64))

	// Setup: Create consultant with access
	consultantEmail := fmt.Sprintf("consultant_comments_%d@example.com", time.Now().UnixNano())
	consultantToken, err := client.RegisterAndLogin("Consultant Comments", consultantEmail, "password", "consultant")
	require.NoError(t, err)

	// Update consultant profile
	client.SetToken(consultantToken)
	profileData := map[string]interface{}{
		"description": "Tester",
		"services":    "Testing",
		"breeds":      "All",
		"location":    "Test City",
		"surname":     "Tester",
	}
	var profileResp map[string]interface{}
	status = client.Put("/consultants/profile", profileData, &profileResp)
	require.Equal(t, http.StatusOK, status)
	consultantID := uint(profileResp["user_id"].(float64))

	// Owner invites consultant
	client.SetToken(ownerToken)
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

	var ownerCommentID, consultantCommentID float64

	t.Run("Owner creates comment", func(t *testing.T) {
		client.SetToken(ownerToken)
		commentReq := map[string]interface{}{
			"event_id": eventID,
			"content":  "# Great walk!\n\nDog was very energetic today.",
		}
		var resp map[string]interface{}
		status := client.Post(fmt.Sprintf("/events/%d/comments", eventID), commentReq, &resp)
		require.Equal(t, http.StatusCreated, status)
		require.Equal(t, "# Great walk!\n\nDog was very energetic today.", resp["content"])
		ownerCommentID = resp["id"].(float64)
	})

	t.Run("Consultant creates comment", func(t *testing.T) {
		client.SetToken(consultantToken)
		commentReq := map[string]interface{}{
			"event_id": eventID,
			"content":  "## Professional observation\n\n- Good behavior\n- Needs more training",
		}
		var resp map[string]interface{}
		status := client.Post(fmt.Sprintf("/events/%d/comments", eventID), commentReq, &resp)
		require.Equal(t, http.StatusCreated, status)
		consultantCommentID = resp["id"].(float64)
	})

	t.Run("List comments for event", func(t *testing.T) {
		client.SetToken(ownerToken)
		var resp map[string]interface{}
		status := client.Get(fmt.Sprintf("/events/%d/comments", eventID), &resp)
		require.Equal(t, http.StatusOK, status)

		comments := resp["comments"].([]interface{})
		require.Len(t, comments, 2)
		require.Equal(t, float64(2), resp["count"].(float64))
	})

	t.Run("Get comment by ID", func(t *testing.T) {
		client.SetToken(ownerToken)
		var resp map[string]interface{}
		status := client.Get(fmt.Sprintf("/event-comments/%.0f", ownerCommentID), &resp)
		require.Equal(t, http.StatusOK, status)
		require.Contains(t, resp["content"], "Great walk")
	})

	t.Run("Update own comment", func(t *testing.T) {
		client.SetToken(ownerToken)
		updateReq := map[string]interface{}{
			"content": "# Updated: Great walk!\n\nDog was very energetic.",
		}
		var resp map[string]interface{}
		status := client.Put(fmt.Sprintf("/event-comments/%.0f", ownerCommentID), updateReq, &resp)
		require.Equal(t, http.StatusOK, status)
		require.Contains(t, resp["content"], "Updated")
	})

	t.Run("Cannot update others' comment", func(t *testing.T) {
		client.SetToken(ownerToken)
		updateReq := map[string]interface{}{
			"content": "Trying to update consultant's comment",
		}
		var resp map[string]interface{}
		status := client.Put(fmt.Sprintf("/event-comments/%.0f", consultantCommentID), updateReq, &resp)
		require.Equal(t, http.StatusForbidden, status)
	})

	t.Run("Delete own comment", func(t *testing.T) {
		client.SetToken(ownerToken)
		status := client.Delete(fmt.Sprintf("/event-comments/%.0f", ownerCommentID))
		require.Equal(t, http.StatusNoContent, status)

		// Verify deleted
		var resp map[string]interface{}
		status = client.Get(fmt.Sprintf("/event-comments/%.0f", ownerCommentID), &resp)
		require.Equal(t, http.StatusNotFound, status)
	})

	t.Run("Consultant without access cannot comment", func(t *testing.T) {
		// Create another consultant without access
		anotherConsultantEmail := fmt.Sprintf("consultant_no_access_%d@example.com", time.Now().UnixNano())
		anotherToken, err := client.RegisterAndLogin("No Access Consultant", anotherConsultantEmail, "password", "consultant")
		require.NoError(t, err)

		client.SetToken(anotherToken)
		commentReq := map[string]interface{}{
			"event_id": eventID,
			"content":  "Should not work",
		}
		var resp map[string]interface{}
		status := client.Post(fmt.Sprintf("/events/%d/comments", eventID), commentReq, &resp)
		require.Equal(t, http.StatusForbidden, status)
	})
}
