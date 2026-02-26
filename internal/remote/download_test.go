package remote

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestDownloadRepoZip_InvalidURL(t *testing.T) {
	_, err := DownloadRepoZip(":invalid-url:")
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestDownloadRepoZip_Baseline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	start := time.Now()
	_, _ = DownloadRepoZip(server.URL)
	duration := time.Since(start)

	if duration < 1*time.Second {
		t.Errorf("expected at least 1s duration, got %v", duration)
	}
	t.Logf("Baseline duration: %v", duration)
}

func TestDownloadRepoZip_Timeout(t *testing.T) {
	// Save original client and restore after test
	originalClient := httpClient
	defer func() { httpClient = originalClient }()

	httpClient = &http.Client{
		Timeout: 100 * time.Millisecond,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	_, err := DownloadRepoZip(server.URL)
	if err == nil {
		t.Error("expected timeout error, got nil")
	} else if !strings.Contains(err.Error(), "timeout") && !strings.Contains(err.Error(), "Client.Timeout exceeded") {
		t.Errorf("expected timeout error, got: %v", err)
	}
	t.Logf("Got expected error: %v", err)
}
