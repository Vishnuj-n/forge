package fileops

import (
	"os"
	"path/filepath"
	"testing"

	"forge/internal/template"
)

func TestCopyFile(t *testing.T) {
	// Create source file
	srcDir, err := os.MkdirTemp("", "src-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(srcDir)

	srcFile := filepath.Join(srcDir, "test.txt")
	content := []byte("test content")
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	// Create destination directory
	dstDir, err := os.MkdirTemp("", "dst-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(dstDir)

	// Create FileOps
	fops := &FileOps{workspaceDir: dstDir, templateDir: srcDir}

	// Test copyFile
	dstFile := filepath.Join(dstDir, "test.txt")
	if err := fops.copyFile(srcFile, dstFile); err != nil {
		t.Fatalf("copyFile error = %v", err)
	}

	// Verify file was copied
	dstContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}

	if string(dstContent) != string(content) {
		t.Errorf("copyFile content mismatch: got %q, want %q", dstContent, content)
	}
}

func TestCopyDir(t *testing.T) {
	// Create source directory structure
	srcDir, err := os.MkdirTemp("", "src-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(srcDir)

	// Create subdirectory and file
	subDir := filepath.Join(srcDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}

	if err := os.WriteFile(filepath.Join(subDir, "test.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	// Create destination
	dstDir, err := os.MkdirTemp("", "dst-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(dstDir)

	fops := &FileOps{workspaceDir: dstDir, templateDir: srcDir}

	// Test copyDir
	dstSubDir := filepath.Join(dstDir, "subdir")
	if err := fops.copyDir(srcDir, dstDir); err != nil {
		t.Fatalf("copyDir error = %v", err)
	}

	// Verify directory structure was copied
	if _, err := os.Stat(dstSubDir); os.IsNotExist(err) {
		t.Fatal("subdirectory was not copied")
	}

	if _, err := os.Stat(filepath.Join(dstSubDir, "test.txt")); os.IsNotExist(err) {
		t.Fatal("file in subdirectory was not copied")
	}
}

func TestApplyAppends(t *testing.T) {
	// Create workspace
	wsDir, err := os.MkdirTemp("", "ws-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(wsDir)

	// Create template directory
	tmplDir, err := os.MkdirTemp("", "tmpl-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(tmplDir)

	// Create patches directory
	patchDir := filepath.Join(tmplDir, "patches")
	if err := os.MkdirAll(patchDir, 0755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}

	// Create patch file
	patchFile := filepath.Join(patchDir, "gitignore.append")
	patchContent := []byte("\n# Added by template\n.env\n")
	if err := os.WriteFile(patchFile, patchContent, 0644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	// Create target file in workspace
	targetFile := filepath.Join(wsDir, ".gitignore")
	baseContent := []byte("*.log\n")
	if err := os.WriteFile(targetFile, baseContent, 0644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	// Apply patches
	fops := New(wsDir, tmplDir)
	patches := []template.AppendPatch{
		{Target: ".gitignore", Source: "patches/gitignore.append"},
	}

	if err := fops.ApplyAppends(patches); err != nil {
		t.Fatalf("ApplyAppends error = %v", err)
	}

	// Verify content was appended
	finalContent, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}

	expected := string(baseContent) + string(patchContent)
	if string(finalContent) != expected {
		t.Errorf("ApplyAppends result mismatch:\ngot:\n%q\n\nwant:\n%q", finalContent, expected)
	}
}

func TestApplyAppendsNonExistentTarget(t *testing.T) {
	// Create workspace
	wsDir, err := os.MkdirTemp("", "ws-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(wsDir)

	// Create template directory
	tmplDir, err := os.MkdirTemp("", "tmpl-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(tmplDir)

	// Create patches directory
	patchDir := filepath.Join(tmplDir, "patches")
	if err := os.MkdirAll(patchDir, 0755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}

	// Create patch file
	patchFile := filepath.Join(patchDir, "gitignore.append")
	if err := os.WriteFile(patchFile, []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	// Apply patches to non-existent target
	fops := New(wsDir, tmplDir)
	patches := []template.AppendPatch{
		{Target: ".gitignore", Source: "patches/gitignore.append"},
	}

	if err := fops.ApplyAppends(patches); err == nil {
		t.Fatal("ApplyAppends should fail for non-existent target")
	}
}
