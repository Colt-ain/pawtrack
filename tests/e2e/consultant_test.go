package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConsultantFlow(t *testing.T) {
	client := NewTestClient(BaseURL)
	client.SetT(t)

	// 1. Register Owner and create a Dog
	ownerEmail := fmt.Sprintf("owner_%d@example.com", time.Now().UnixNano())
	ownerToken, err := client.RegisterAndLogin("Owner", ownerEmail, "password", "owner")
	require.NoError(t, err)
	client.SetToken(ownerToken)

	dogID, err := client.CreateDog("Buddy", "Labrador", "2020-01-01T00:00:00Z")
	require.NoError(t, err)

	// 2. Register Consultant
	consultantEmail := fmt.Sprintf("consultant_%d@example.com", time.Now().UnixNano())
	consultantToken, err := client.RegisterAndLogin("Consultant", consultantEmail, "password", "consultant")
	require.NoError(t, err)

	// 3. Consultant updates profile and capture their ID
	client.SetToken(consultantToken)
	profileData := map[string]interface{}{
		"description": "Expert dog trainer with 10 years experience",
		"services":    "Training, Walking, Sitting",
		"breeds":      "Labrador, Poodle, Bulldog",
		"location":    "New York",
		"surname":     "Smith",
	}
	var profileResp map[string]interface{}
	status := client.Put("/consultants/profile", profileData, &profileResp)
	require.Equal(t, http.StatusOK, status)
	require.Equal(t, "Smith", profileResp["surname"])
	
	// Capture consultant ID from profile response
	consultantID := uint(profileResp["user_id"].(float64))

	// 4. Owner searches for consultant (verify search works)
	client.SetToken(ownerToken)
	var searchResp map[string]interface{}
	status = client.Get("/consultants?query=Smith&services=Training", &searchResp)
	require.Equal(t, http.StatusOK, status)
	
	data := searchResp["data"].([]interface{})
	require.NotEmpty(t, data, "Search should find the consultant")

	// 5. Owner invites consultant
	inviteReq := map[string]interface{}{
		"dog_id": dogID,
	}
	var inviteResp map[string]interface{}
	status = client.Post(fmt.Sprintf("/consultants/%d/invite", consultantID), inviteReq, &inviteResp)
	require.Equal(t, http.StatusCreated, status)
	
	inviteToken := inviteResp["token"].(string)
	require.NotEmpty(t, inviteToken)

	// 6. Consultant accepts invite
	client.SetToken(consultantToken)
	var acceptResp map[string]interface{}
	status = client.Post(fmt.Sprintf("/invites/accept?token=%s", inviteToken), nil, &acceptResp)
	require.Equal(t, http.StatusOK, status, "Accept invite failed: %v", acceptResp)

	// 7. Verify Consultant has access to the dog
	var dogResp map[string]interface{}
	status = client.Get(fmt.Sprintf("/dogs/%d", dogID), &dogResp)
	require.Equal(t, http.StatusOK, status)
	require.Equal(t, "Buddy", dogResp["name"])
}
