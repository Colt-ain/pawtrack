package e2e

import (
	"bytes"
	"encoding/json"
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

func NewTestClient(t *testing.T, baseURL string) *TestClient {
	return &TestClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
		t:       t,
	}
}

func (c *TestClient) SetToken(token string) {
	c.Token = token
}

func (c *TestClient) Post(path string, body interface{}, response interface{}) int {
	jsonBody, err := json.Marshal(body)
	require.NoError(c.t, err)

	req, err := http.NewRequest("POST", c.BaseURL+path, bytes.NewBuffer(jsonBody))
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

func (c *TestClient) Get(path string, response interface{}) int {
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	require.NoError(c.t, err)

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

func (c *TestClient) Delete(path string) int {
	req, err := http.NewRequest("DELETE", c.BaseURL+path, nil)
	require.NoError(c.t, err)

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.Client.Do(req)
	require.NoError(c.t, err)
	defer resp.Body.Close()

	return resp.StatusCode
}

func (c *TestClient) RegisterAndLogin(name, email, password string) {
	// Register
	regBody := map[string]string{
		"name":     name,
		"email":    email,
		"password": password,
	}
	c.Post("/auth/register/owner", regBody, nil)

	// Login
	loginBody := map[string]string{
		"email":    email,
		"password": password,
	}
	var resp map[string]interface{}
	c.Post("/auth/login", loginBody, &resp)

	token, ok := resp["token"].(string)
	require.True(c.t, ok, "token should be a string")
	c.SetToken(token)
}
