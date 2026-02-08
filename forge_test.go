package main

import (
	"os/exec"
	"testing"
)

func TestHelpCommand(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "--help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("forge --help failed: %v\nOutput: %s", err, string(out))
	}
	if len(out) == 0 || !containsHelp(string(out)) {
		t.Errorf("Expected help output, got: %s", string(out))
	}
}

func containsHelp(out string) bool {
	return (len(out) > 0 && (contains(out, "Usage:") || contains(out, "Forge - A safety-first project bootstrapper")))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr))))
}
