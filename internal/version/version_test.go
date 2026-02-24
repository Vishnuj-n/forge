package version

import (
	"strings"
	"testing"
)

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		current    string
		latest     string
		wantNewer  bool
		wantErrMsg string
	}{
		// Basic version comparisons
		{"0.1.5", "0.2.0", true, ""},
		{"0.1.5", "0.1.6", true, ""},
		{"0.1.5", "1.0.0", true, ""},
		{"0.2.0", "0.1.5", false, ""},
		{"0.1.5", "0.1.5", false, ""},

		// With "v" prefix
		{"v0.1.5", "v0.2.0", true, ""},
		{"v0.1.5", "0.2.0", true, ""},
		{"0.1.5", "v0.2.0", true, ""},

		// Pre-release handling
		{"0.1.5", "0.2.0-alpha", true, ""},
		{"0.1.5", "0.1.6-beta", true, ""},

		// Error cases
		{"invalid", "0.2.0", false, "invalid current version format"},
		{"0.1.5", "invalid", false, "invalid latest version format"},
		{"0.1", "0.2.0", false, "expected semantic version format"},
	}

	for _, tt := range tests {
		t.Run(tt.current+"_vs_"+tt.latest, func(t *testing.T) {
			got, err := IsNewerVersion(tt.current, tt.latest)
			if tt.wantErrMsg != "" {
				if err == nil {
					t.Errorf("IsNewerVersion(%q, %q) = got nil error, want error containing %q", tt.current, tt.latest, tt.wantErrMsg)
					return
				}
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("IsNewerVersion(%q, %q) error = %v, want error containing %q", tt.current, tt.latest, err, tt.wantErrMsg)
				}
			} else {
				if err != nil {
					t.Errorf("IsNewerVersion(%q, %q) unexpected error: %v", tt.current, tt.latest, err)
				}
				if got != tt.wantNewer {
					t.Errorf("IsNewerVersion(%q, %q) = %v, want %v", tt.current, tt.latest, got, tt.wantNewer)
				}
			}
		})
	}
}
