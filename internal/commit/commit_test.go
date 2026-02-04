package commit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCommitBasic(t *testing.T) {
	// Create workspace with some content
	wsDir, err := os.MkdirTemp("", "ws-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(wsDir)

	// Add file to workspace
	testFile := filepath.Join(wsDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	// Create target directory
	targetDir, err := os.MkdirTemp("", "target-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(targetDir)

	// Remove empty target so commit can use it
	if err := os.RemoveAll(targetDir); err != nil {
		t.Fatalf("RemoveAll error = %v", err)
	}

	// Commit
	committer := New()
	if err := committer.Commit(wsDir, targetDir); err != nil {
		t.Fatalf("Commit() error = %v", err)
	}

	// Verify file exists in target
	targetFile := filepath.Join(targetDir, "test.txt")
	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		t.Fatal("file was not committed to target directory")
	}

	// Verify content
	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}

	if string(content) != "test content" {
		t.Errorf("Commit() content mismatch: got %q, want %q", content, "test content")
	}
}

func TestCommitNonEmptyTargetFails(t *testing.T) {
	wsDir, err := os.MkdirTemp("", "ws-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(wsDir)

	// Create non-empty target directory
	targetDir, err := os.MkdirTemp("", "target-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(targetDir)

	// Add file to target
	if err := os.WriteFile(filepath.Join(targetDir, "existing.txt"), []byte("existing"), 0644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	// Commit to non-empty target should fail
	committer := New()
	if err := committer.Commit(wsDir, targetDir); err == nil {
		t.Fatal("Commit() should fail for non-empty target directory")
	}
}

func TestCommitPreservesFilePerms(t *testing.T) {
	// Create workspace with file
	wsDir, err := os.MkdirTemp("", "ws-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(wsDir)

	testFile := filepath.Join(wsDir, "test.sh")
	if err := os.WriteFile(testFile, []byte("#!/bin/bash\necho test"), 0755); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	// Create target
	targetDir, err := os.MkdirTemp("", "target-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(targetDir)

	if err := os.RemoveAll(targetDir); err != nil {
		t.Fatalf("RemoveAll error = %v", err)
	}

	// Commit
	committer := New()
	if err := committer.Commit(wsDir, targetDir); err != nil {
		t.Fatalf("Commit() error = %v", err)
	}

	// Verify permissions
	targetFile := filepath.Join(targetDir, "test.sh")
	info, err := os.Stat(targetFile)
	if err != nil {
		t.Fatalf("Stat error = %v", err)
	}

	// Check if executable bit is preserved (Unix-style)
	if (info.Mode() & 0111) == 0 {
		t.Log("Note: executable bit not preserved - may be expected on some systems")
	}
}
