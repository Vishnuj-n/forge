package commit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

		// Ensure target directory is empty
		empty, err := isDirEmpty(targetPath)
		if err != nil {
			return fmt.Errorf("failed to check if target directory is empty: %w", err)
		}
		if !empty {
			return fmt.Errorf("target directory is not empty")
		}
	}

	// Check if current working directory is inside the target path
	// If so, temporarily change to parent to avoid Windows directory lock
	originalWd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	needsChdirWorkaround := false
	if targetExists {
		// Check if we're inside the target directory
		absOriginalWd, _ := filepath.Abs(originalWd)
		absTargetPath, _ := filepath.Abs(targetPath)

		if absOriginalWd == absTargetPath || isSubPath(absOriginalWd, absTargetPath) {
			needsChdirWorkaround = true
			// Change to parent directory to release the lock
			parentDir := filepath.Dir(absTargetPath)
			if err := os.Chdir(parentDir); err != nil {
				return fmt.Errorf("failed to change to parent directory: %w", err)
			}
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

		// If we changed directories, change into the new project directory
		if needsChdirWorkaround {
			if err := os.Chdir(targetPath); err != nil {
				// Non-fatal: just warn
				fmt.Printf("  ⚠ Warning: could not change back to project directory: %v\n", err)
			}
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

// isSubPath checks if child path is inside parent path
func isSubPath(child, parent string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	// If relative path starts with "..", it's outside parent
	return !filepath.IsAbs(rel) && !startsWithDotDot(rel)
}

// startsWithDotDot checks if path starts with ".."
func startsWithDotDot(path string) bool {
	if path == ".." {
		return true
	}
	return strings.HasPrefix(path, ".."+string(filepath.Separator))
}

// isDirEmpty checks if the directory at path is empty
func isDirEmpty(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}
