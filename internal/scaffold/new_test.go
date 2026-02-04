package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid simple name",
			input:   "my-template",
			wantErr: false,
		},
		{
			name:    "valid with underscore",
			input:   "my_template",
			wantErr: false,
		},
		{
			name:    "valid with numbers",
			input:   "template123",
			wantErr: false,
		},
		{
			name:    "empty name",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid special chars",
			input:   "my@template",
			wantErr: true,
		},
		{
			name:    "invalid spaces",
			input:   "my template",
			wantErr: true,
		},
		{
			name:    "too long",
			input:   string(make([]byte, 51)),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	// Create temp base directory
	baseDir, err := os.MkdirTemp("", "scaffold-test-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(baseDir)

	gen := New(baseDir)

	// Test successful generation
	templateName := "test-template"
	templateDir, err := gen.Generate(templateName)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Verify directory exists
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		t.Fatal("template directory was not created")
	}

	// Verify template.yaml exists
	yamlPath := filepath.Join(templateDir, "template.yaml")
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		t.Fatal("template.yaml was not created")
	}

	// Verify README.md exists
	readmePath := filepath.Join(templateDir, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		t.Fatal("README.md was not created")
	}

	// Verify files directory exists
	filesDir := filepath.Join(templateDir, "files")
	if _, err := os.Stat(filesDir); os.IsNotExist(err) {
		t.Fatal("files directory was not created")
	}

	// Verify patches directory exists
	patchesDir := filepath.Join(templateDir, "patches")
	if _, err := os.Stat(patchesDir); os.IsNotExist(err) {
		t.Fatal("patches directory was not created")
	}

	// Verify template.yaml content
	yamlContent, err := os.ReadFile(yamlPath)
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}

	if !contains(string(yamlContent), "name: test-template") {
		t.Fatal("template.yaml does not contain template name")
	}

	if !contains(string(yamlContent), "commands:") {
		t.Fatal("template.yaml does not contain commands section")
	}
}

func TestGenerateDuplicate(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "scaffold-test-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(baseDir)

	gen := New(baseDir)

	// Generate first template
	_, err = gen.Generate("duplicate-test")
	if err != nil {
		t.Fatalf("First Generate() error = %v", err)
	}

	// Try to generate duplicate
	_, err = gen.Generate("duplicate-test")
	if err == nil {
		t.Fatal("Generate() should fail for duplicate template")
	}
}

func TestGenerateInvalidName(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "scaffold-test-")
	if err != nil {
		t.Fatalf("MkdirTemp error = %v", err)
	}
	defer os.RemoveAll(baseDir)

	gen := New(baseDir)

	_, err = gen.Generate("invalid@name")
	if err == nil {
		t.Fatal("Generate() should fail for invalid name")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr)
}
