package remote

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
)

// GitHubRelease represents a GitHub release API response (simplified)
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// FetchLatestReleaseVersion queries GitHub Releases API to get the latest version
// and returns the version tag and Windows binary download URL
func FetchLatestReleaseVersion(owner, repo string) (version string, downloadURL string, err error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch releases from GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", "", fmt.Errorf("failed to fetch releases from GitHub: http %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return "", "", fmt.Errorf("failed to parse GitHub release response: %w", err)
	}

	if release.TagName == "" {
		return "", "", fmt.Errorf("no releases found")
	}

	// Find Windows binary asset
	expectedAsset := "forge.exe"
	if runtime.GOOS == "darwin" {
		expectedAsset = "forge-darwin"
	} else if runtime.GOOS == "linux" {
		expectedAsset = "forge-linux"
	}

	for _, asset := range release.Assets {
		if asset.Name == expectedAsset {
			return release.TagName, asset.DownloadURL, nil
		}
	}

	return "", "", fmt.Errorf("no suitable binary found for %s in release %s", runtime.GOOS, release.TagName)
}

// DownloadReleaseBinary downloads a binary from the given URL to a temporary file
// Returns the path to the downloaded file. Caller must remove the file.
func DownloadReleaseBinary(downloadURL, tempPath string) error {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to download binary: http %d", resp.StatusCode)
	}

	// Create temp file
	tmp, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file for binary: %w", err)
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to write binary to temp file: %w", err)
	}

	return nil
}
