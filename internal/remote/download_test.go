package remote

import (
	"testing"
)

func TestDownloadRepoZip_InvalidURL(t *testing.T) {
	_, err := DownloadRepoZip(":invalid-url:")
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}
