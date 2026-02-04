package fileops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"forge/internal/template"
)

// FileOps handles file operations (copy and append)
type FileOps struct {
	workspaceDir string
	templateDir  string
}

// New creates a new file operations handler
func New(workspaceDir, templatePath string) *FileOps {
	// Determine template directory
	templateDir := templatePath
	if info, err := os.Stat(templatePath); err == nil && !info.IsDir() {
		templateDir = filepath.Dir(templatePath)
	}
	
	return &FileOps{
		workspaceDir: workspaceDir,
		templateDir:  templateDir,
	}
}

// CopyFiles copies files/directories from template to workspace
func (f *FileOps) CopyFiles(copyPaths []string) error {
	for _, srcPath := range copyPaths {
		// Resolve source path relative to template directory
		absSrc := filepath.Join(f.templateDir, srcPath)
		
		info, err := os.Stat(absSrc)
		if err != nil {
			return fmt.Errorf("failed to stat %s: %w", srcPath, err)
		}
		
		if info.IsDir() {
			// Copy entire directory
			if err := f.copyDir(absSrc, f.workspaceDir); err != nil {
				return fmt.Errorf("failed to copy directory %s: %w", srcPath, err)
			}
			fmt.Printf("  ✓ Copied directory: %s\n", srcPath)
		} else {
			// Copy single file to workspace root
			dstPath := filepath.Join(f.workspaceDir, filepath.Base(srcPath))
			if err := f.copyFile(absSrc, dstPath); err != nil {
				return fmt.Errorf("failed to copy file %s: %w", srcPath, err)
			}
			fmt.Printf("  ✓ Copied file: %s\n", srcPath)
		}
	}
	
	return nil
}

// ApplyAppends applies append-only patches
func (f *FileOps) ApplyAppends(patches []template.AppendPatch) error {
	for _, patch := range patches {
		// Resolve source path relative to template directory
		srcPath := filepath.Join(f.templateDir, patch.Source)
		
		// Resolve target path in workspace
		dstPath := filepath.Join(f.workspaceDir, patch.Target)
		
		// Read source content
		content, err := os.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("failed to read patch source %s: %w", patch.Source, err)
		}
		
		// Check if target exists
		if _, err := os.Stat(dstPath); os.IsNotExist(err) {
			return fmt.Errorf("append target %s does not exist (patches can only append to existing files)", patch.Target)
		}
		
		// Append content to target
		file, err := os.OpenFile(dstPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open target %s: %w", patch.Target, err)
		}
		defer file.Close()
		
		if _, err := file.Write(content); err != nil {
			return fmt.Errorf("failed to append to %s: %w", patch.Target, err)
		}
		
		fmt.Printf("  ✓ Appended to: %s\n", patch.Target)
	}
	
	return nil
}

// copyFile copies a single file
func (f *FileOps) copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	
	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	
	// Copy permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// copyDir recursively copies a directory
func (f *FileOps) copyDir(src, dstBase string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Calculate relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		
		// Calculate destination path
		dstPath := filepath.Join(dstBase, relPath)
		
		if info.IsDir() {
			// Create directory
			return os.MkdirAll(dstPath, info.Mode())
		}
		
		// Copy file
		return f.copyFile(path, dstPath)
	})
}
