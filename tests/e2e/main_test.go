package e2e

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

var BaseURL = "http://localhost:8080/api/v1"

func init() {
	if url := os.Getenv("E2E_BASE_URL"); url != "" {
		BaseURL = url
	}
}

func TestMain(m *testing.M) {
	// Check if server is running
	u, err := url.Parse(BaseURL)
	if err != nil {
		fmt.Printf("Invalid BaseURL: %v\n", err)
		os.Exit(1)
	}
	healthURL := fmt.Sprintf("%s://%s/health", u.Scheme, u.Host)

	if err := waitForServer(healthURL); err != nil {
		fmt.Printf("Server not ready: %v\n", err)
		fmt.Println("Please ensure the server is running (docker compose up)")
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func waitForServer(url string) error {
	for i := 0; i < 30; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("server did not respond after 30 seconds")
}
