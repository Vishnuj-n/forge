package workspace

import (
	"fmt"
	"os"
	"path/filepath"
)

// Workspace represents an isolated temporary workspace
type Workspace struct {
	path string
}

// New creates a new temporary workspace
func New() (*Workspace, error) {
	// Create temp directory with a recognizable prefix
	tempDir, err := os.MkdirTemp("", "forge-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp workspace: %w", err)
	}

	return &Workspace{path: tempDir}, nil
}

// Path returns the absolute path to the workspace
func (w *Workspace) Path() string {
	return w.path
}

// Cleanup removes the workspace directory and all its contents
func (w *Workspace) Cleanup() error {
	if w.path == "" {
		return nil
	}
	return os.RemoveAll(w.path)
}

// GetVolume returns the volume/drive letter of the given path (Windows-specific)
// Returns empty string on non-Windows or if volume cannot be determined
func GetVolume(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return ""
	}

	vol := filepath.VolumeName(absPath)
	return vol
}

// SameVolume checks if two paths are on the same volume
func SameVolume(path1, path2 string) bool {
	vol1 := GetVolume(path1)
	vol2 := GetVolume(path2)

	if vol1 == "" || vol2 == "" {
		return false
	}

	return vol1 == vol2
}
