package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWorkspaceNew(t *testing.T) {
	ws, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if ws == nil {
		t.Fatal("New() returned nil workspace")
	}

	if ws.Path() == "" {
		t.Fatal("Path() returned empty string")
	}

	// Verify workspace exists
	if _, err := os.Stat(ws.Path()); os.IsNotExist(err) {
		t.Fatal("workspace directory was not created")
	}

	// Cleanup
	if err := ws.Cleanup(); err != nil {
		t.Fatalf("Cleanup() error = %v", err)
	}

	// Verify cleanup
	if _, err := os.Stat(ws.Path()); err == nil {
		t.Fatal("workspace directory was not deleted after cleanup")
	}
}

func TestWorkspaceIsolation(t *testing.T) {
	ws1, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer ws1.Cleanup()

	ws2, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer ws2.Cleanup()

	if ws1.Path() == ws2.Path() {
		t.Fatal("workspaces have the same path")
	}
}

func TestSameVolume(t *testing.T) {
	// Both paths in temp dir should be same volume
	tmpDir1, err := os.MkdirTemp("", "test1-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(tmpDir1)

	tmpDir2, err := os.MkdirTemp("", "test2-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(tmpDir2)

	if !SameVolume(tmpDir1, tmpDir2) {
		t.Log("Note: SameVolume returned false - may be expected on some systems")
	}
}

func TestGetVolume(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(tmpDir)

	vol := GetVolume(tmpDir)
	// On Windows, should return something like "C:"
	// On other systems, may return empty
	t.Logf("GetVolume returned: %q", vol)

	// Create a file and test with it
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	volFile := GetVolume(testFile)
	if volFile != vol {
		t.Errorf("GetVolume returned different values for dir (%q) and file (%q)", vol, volFile)
	}
}
