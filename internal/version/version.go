package version

import (
	"fmt"
	"strconv"
	"strings"
)

// IsNewerVersion compares two semantic versions (e.g., "0.1.5" vs "0.2.0")
// Returns true if latest is newer than current.
// Example: IsNewerVersion("0.1.5", "0.2.0") returns true
func IsNewerVersion(current, latest string) (bool, error) {
	currParts, err := parseVersion(current)
	if err != nil {
		return false, fmt.Errorf("invalid current version format: %w", err)
	}

	latestParts, err := parseVersion(latest)
	if err != nil {
		return false, fmt.Errorf("invalid latest version format: %w", err)
	}

	// Compare major.minor.patch
	for i := 0; i < 3; i++ {
		if latestParts[i] > currParts[i] {
			return true, nil
		}
		if latestParts[i] < currParts[i] {
			return false, nil
		}
	}

	// Versions are equal
	return false, nil
}

// parseVersion converts a version string like "0.1.5" into [0, 1, 5]
// Handles versions with leading "v" prefix (e.g., "v0.2.0")
func parseVersion(v string) ([3]int, error) {
	var result [3]int

	// Remove common prefixes
	v = strings.TrimPrefix(v, "v")
	v = strings.TrimSpace(v)

	// Split by dots
	parts := strings.Split(v, ".")
	if len(parts) < 3 {
		return result, fmt.Errorf("expected semantic version format (major.minor.patch), got: %s", v)
	}

	// Parse major, minor, patch
	for i := 0; i < 3; i++ {
		// Handle pre-release versions (e.g., "0.2.0-alpha") - just take the numeric part
		numStr := strings.FieldsFunc(parts[i], func(r rune) bool {
			return r == '-' || r == '+' || r == '.'
		})[0]

		num, err := strconv.Atoi(numStr)
		if err != nil {
			return result, fmt.Errorf("invalid version number at position %d: %s", i, parts[i])
		}
		result[i] = num
	}

	return result, nil
}
