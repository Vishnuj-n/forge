package executor

import (
	"os"
	"testing"

	"forge/internal/template"
)

func TestExecutorRun(t *testing.T) {
	// Create temp workspace
	wsDir, err := os.MkdirTemp("", "ws-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(wsDir)

	// Initialize git repo in workspace first
	exec := New(wsDir, false)
	initCmd := template.Command{Cmd: []string{"git", "init"}}
	if err := exec.Run(initCmd); err != nil {
		t.Fatalf("git init error = %v", err)
	}

	// Test git config command
	cmd := template.Command{Cmd: []string{"git", "config", "user.name", "Test"}}
	if err := exec.Run(cmd); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestExecutorRunFailure(t *testing.T) {
	// Create temp workspace
	wsDir, err := os.MkdirTemp("", "ws-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(wsDir)

	exec := New(wsDir, false)

	// Test command that should fail
	cmd := template.Command{Cmd: []string{"nonexistent-command-12345"}}
	if err := exec.Run(cmd); err == nil {
		t.Fatal("Run() should fail for non-existent command")
	}
}

func TestExecutorRunEmptyCommand(t *testing.T) {
	wsDir, err := os.MkdirTemp("", "ws-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(wsDir)

	exec := New(wsDir, false)

	// Test empty command
	cmd := template.Command{Cmd: []string{}}
	if err := exec.Run(cmd); err == nil {
		t.Fatal("Run() should fail for empty command")
	}
}
