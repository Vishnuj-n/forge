package update

import (
	"fmt"
	"os"
	"path/filepath"

	"forge/internal/remote"
	"forge/internal/version"
)

const (
	// GitHub repository coordinates for Forge binary
	GitHubOwner = "Vishnuj-n"
	GitHubRepo  = "forge"
)

// CheckUpdate checks if a newer version is available on GitHub
// Returns the new version string, whether an update is available, and any error
func CheckUpdate(currentVersion string) (newVersion string, available bool, err error) {
	tagName, _, err := remote.FetchLatestReleaseVersion(GitHubOwner, GitHubRepo)
	if err != nil {
		return "", false, fmt.Errorf("failed to check for updates: %w", err)
	}

	isNewer, err := version.IsNewerVersion(currentVersion, tagName)
	if err != nil {
		return "", false, fmt.Errorf("failed to compare versions: %w", err)
	}

	return tagName, isNewer, nil
}

// PerformUpdate downloads the latest binary and replaces the current executable
// binaryPath should be the path to the current forge executable
func PerformUpdate(newVersion, binaryPath string) error {
	// Fetch release info and download URL
	_, downloadURL, err := remote.FetchLatestReleaseVersion(GitHubOwner, GitHubRepo)
	if err != nil {
		return fmt.Errorf("failed to fetch release info: %w", err)
	}

	// Create temp file for download
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "forge-update.exe")

	defer os.Remove(tempFile) // Clean up temp file after update

	// Download the new binary
	fmt.Printf("Downloading Forge %s...\n", newVersion)
	if err := remote.DownloadReleaseBinary(downloadURL, tempFile); err != nil {
		return fmt.Errorf("failed to download new binary: %w", err)
	}

	// Replace the old binary with the new one
	// We need to back up the old binary in case something goes wrong
	backupPath := binaryPath + ".backup"

	// On Windows, we can't directly replace a running executable
	// So we rename the old one, move the new one in place, and clean up
	if err := os.Rename(binaryPath, backupPath); err != nil {
		return fmt.Errorf("failed to back up current binary: %w", err)
	}

	if err := os.Rename(tempFile, binaryPath); err != nil {
		// Restore backup on failure
		os.Rename(backupPath, binaryPath)
		return fmt.Errorf("failed to install new binary: %w", err)
	}

	// Clean up backup
	os.Remove(backupPath)

	return nil
}
