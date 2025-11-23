package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestClient struct {
	BaseURL string
	Client  *http.Client
	Token   string
	t       *testing.T
}

func NewTestClient(baseURL string) *TestClient {
	return &TestClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
		t:       nil, // t is set per test if needed, or we can pass it in methods. 
		// Actually, better to pass t in NewTestClient if we want to use require.
		// But the existing code uses c.t. Let's fix NewTestClient signature in tests if needed, 
		// or better, let's keep it simple and just use the one from the file.
		// Wait, the previous file had NewTestClient(t, baseURL).
		// But my new test calls NewTestClient(BaseURL).
		// Let's make NewTestClient take just BaseURL and methods take t? 
		// Or better, let's stick to the pattern used in other tests.
		// Checking other tests... they use `client := NewTestClient(BaseURL)` and then `client.t = t` isn't set?
		// Ah, I see in my previous `client_test.go` (before I broke it) it didn't have `t` in struct?
		// No, it did.
		// Let's look at `auth_test.go` to see how it's used.
	}
}

// Helper to set T
func (c *TestClient) SetT(t *testing.T) {
	c.t = t
}

func (c *TestClient) SetToken(token string) {
	c.Token = token
}

func (c *TestClient) doRequest(method, path string, body interface{}, response interface{}) int {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(c.t, err)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	require.NoError(c.t, err)

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.Client.Do(req)
	require.NoError(c.t, err)
	defer resp.Body.Close()

	if response != nil {
		respBody, err := io.ReadAll(resp.Body)
		require.NoError(c.t, err)
		if len(respBody) > 0 {
			err = json.Unmarshal(respBody, response)
			require.NoError(c.t, err, "Failed to unmarshal response: %s", string(respBody))
		}
	}

	return resp.StatusCode
}

func (c *TestClient) Post(path string, body interface{}, result interface{}) int {
	return c.doRequest("POST", path, body, result)
}

func (c *TestClient) Put(path string, body interface{}, result interface{}) int {
	return c.doRequest("PUT", path, body, result)
}

func (c *TestClient) Get(path string, result interface{}) int {
	return c.doRequest("GET", path, nil, result)
}

func (c *TestClient) Delete(path string) int {
	return c.doRequest("DELETE", path, nil, nil)
}

// RegisterAndLogin helper
func (c *TestClient) RegisterAndLogin(name, email, password, role string) (string, error) {
	// Register
	regBody := map[string]string{
		"name":     name,
		"email":    email,
		"password": password,
	}
	
	endpoint := "/auth/register/owner"
	if role == "consultant" {
		endpoint = "/auth/register/consultant"
	}
	
	c.Post(endpoint, regBody, nil)

	// Login
	loginBody := map[string]string{
		"email":    email,
		"password": password,
	}
	var resp map[string]interface{}
	status := c.Post("/auth/login", loginBody, &resp)
	if status != 200 {
		return "", fmt.Errorf("login failed with status %d", status)
	}

	token, ok := resp["token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}
	c.SetToken(token) // Auto-set token
	return token, nil
}

// CreateDog helper
func (c *TestClient) CreateDog(name, breed, birthDate string) (uint, error) {
	req := map[string]string{
		"name":       name,
		"breed":      breed,
		"birth_date": birthDate,
	}
	var resp map[string]interface{}
	status := c.Post("/dogs", req, &resp)
	if status != 201 {
		return 0, fmt.Errorf("create dog failed with status %d", status)
	}
	
	id, ok := resp["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("dog id not found")
	}
	return uint(id), nil
}
