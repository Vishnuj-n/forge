package main

// This file uses the "test binary as CLI" pattern.
//
// When `go test -c` is run, Go produces forge.test.exe — a standalone binary
// that contains both the test runner and the CLI itself.
//
// How it works:
//   - TestMain checks for the FORGE_TEST_HELPER env var.
//   - If set, it runs the real CLI (main()) and exits — acting as forge.exe.
//   - If not set, it runs the test suite normally.
//   - Each test re-invokes os.Args[0] (this same binary) with FORGE_TEST_HELPER=1
//     and the desired CLI arguments, capturing output without needing `go run .`.
//
// This means forge.test.exe works on any machine — no Go toolchain required.
// Build it once with:   go test -c -o forge.test.exe
// Run tests with:       ./forge.test.exe -test.v

import (
	"os"
	"os/exec"
	"testing"
)

// TestMain is the entry point for the compiled test binary.
// When FORGE_TEST_HELPER=1 is set, this binary acts as the forge CLI.
func TestMain(m *testing.M) {
	if os.Getenv("FORGE_TEST_HELPER") == "1" {
		main()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// forgeCmd returns an exec.Cmd that runs this same test binary as the forge CLI.
// It sets FORGE_TEST_HELPER=1 so TestMain routes to main() instead of tests.
func forgeCmd(args ...string) *exec.Cmd {
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = append(os.Environ(), "FORGE_TEST_HELPER=1")
	return cmd
}

func TestHelpCommand(t *testing.T) {
	cmd := forgeCmd("--help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("forge --help failed: %v\nOutput: %s", err, string(out))
	}
	if !containsAny(string(out), "Usage:", "Forge - A safety-first project bootstrapper") {
		t.Errorf("Expected help output, got: %s", string(out))
	}
}

func TestVersionFlag(t *testing.T) {
	cmd := forgeCmd("--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("forge --version failed: %v\nOutput: %s", err, string(out))
	}
	if !containsAny(string(out), "forge", "version", "development") {
		t.Errorf("Expected version output, got: %s", string(out))
	}
}

func TestUnknownCommand(t *testing.T) {
	cmd := forgeCmd("notacommand")
	out, err := cmd.CombinedOutput()
	// Cobra exits non-zero for unknown commands — that's expected
	if err == nil {
		t.Errorf("Expected non-zero exit for unknown command, got success.\nOutput: %s", string(out))
	}
	if !containsAny(string(out), "unknown command", "Error") {
		t.Errorf("Expected error message for unknown command, got: %s", string(out))
	}
}

// containsAny returns true if s contains any of the provided substrings.
func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if containsStr(s, sub) {
			return true
		}
	}
	return false
}

func containsStr(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
