package commit

import (
	"fmt"
	"os"
	"path/filepath"

	"forge/internal/workspace"
)

// Committer handles committing workspace to target directory
type Committer struct{}

// New creates a new committer
func New() *Committer {
	return &Committer{}
}

// Commit moves the workspace contents to the target directory
func (c *Committer) Commit(workspacePath, targetPath string) error {
	// Check if target exists
	targetExists := false
	if info, err := os.Stat(targetPath); err == nil {
		targetExists = true
		if !info.IsDir() {
			return fmt.Errorf("target path exists but is not a directory")
		}
	}
	
	// Check volume compatibility for atomic move
	sameVol := workspace.SameVolume(workspacePath, targetPath)
	
	if sameVol {
		// Atomic move possible
		if targetExists {
			// Target exists and is empty (validated earlier), remove it
			if err := os.Remove(targetPath); err != nil {
				return fmt.Errorf("failed to remove empty target directory: %w", err)
			}
		}
		
		// Atomic rename
		if err := os.Rename(workspacePath, targetPath); err != nil {
			return fmt.Errorf("failed to move workspace to target: %w", err)
		}
		
		fmt.Println("  ✓ Committed atomically")
		return nil
	}
	
	// Cross-volume: best-effort copy
	fmt.Println("  ⚠ Warning: Cross-volume commit detected - using best-effort copy instead of atomic move")
	
	// Ensure target directory exists
	if !targetExists {
		if err := os.MkdirAll(targetPath, 0755); err != nil {
			return fmt.Errorf("failed to create target directory: %w", err)
		}
	}
	
	// Copy all contents
	if err := copyDirContents(workspacePath, targetPath); err != nil {
		return fmt.Errorf("failed to copy workspace contents: %w", err)
	}
	
	fmt.Println("  ✓ Committed (best-effort copy)")
	return nil
}

// copyDirContents copies all contents from src to dst
func copyDirContents(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		
		if entry.IsDir() {
			// Create directory and copy recursively
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := copyDirContents(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// copyFile copies a single file with permissions
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	
	if _, err := dstFile.ReadFrom(srcFile); err != nil {
		return err
	}
	
	// Copy permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	
	return os.Chmod(dst, srcInfo.Mode())
}
