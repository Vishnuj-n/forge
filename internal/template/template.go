package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Template represents a project template configuration
type Template struct {
	Name     string       `yaml:"name"`
	Commands []Command    `yaml:"commands"`
	Files    FileOps      `yaml:"files"`
}

// Command represents a single command to execute
type Command struct {
	Cmd []string `yaml:"cmd"`
}

// FileOps represents file operations (copy and append)
type FileOps struct {
	Copy   []string      `yaml:"copy"`
	Append []AppendPatch `yaml:"append"`
}

// AppendPatch represents an append-only patch operation
type AppendPatch struct {
	Target string `yaml:"target"`
	Source string `yaml:"source"`
}

// String returns a human-readable representation of the command
func (c Command) String() string {
	return strings.Join(c.Cmd, " ")
}

// Load loads and validates a template from the given path
func Load(templatePath string) (*Template, error) {
	// Resolve template path
	absPath, err := filepath.Abs(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve template path: %w", err)
	}
	
	// Check if it's a directory or file
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("template path not found: %w", err)
	}
	
	var yamlPath string
	if info.IsDir() {
		// If directory, look for template.yaml inside
		yamlPath = filepath.Join(absPath, "template.yaml")
	} else {
		// If file, use it directly
		yamlPath = absPath
	}
	
	// Read YAML file
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}
	
	// Parse YAML
	var tmpl Template
	if err := yaml.Unmarshal(data, &tmpl); err != nil {
		return nil, fmt.Errorf("failed to parse template YAML: %w", err)
	}
	
	// Validate
	if err := tmpl.validate(); err != nil {
		return nil, fmt.Errorf("template validation failed: %w", err)
	}
	
	return &tmpl, nil
}

// validate checks if the template is valid
func (t *Template) validate() error {
	if t.Name == "" {
		return fmt.Errorf("template name is required")
	}
	
	// Validate commands
	for i, cmd := range t.Commands {
		if len(cmd.Cmd) == 0 {
			return fmt.Errorf("command %d: cmd array is empty", i)
		}
		if cmd.Cmd[0] == "" {
			return fmt.Errorf("command %d: first element (executable) cannot be empty", i)
		}
	}
	
	// Validate append patches
	for i, patch := range t.Files.Append {
		if patch.Target == "" {
			return fmt.Errorf("append patch %d: target is required", i)
		}
		if patch.Source == "" {
			return fmt.Errorf("append patch %d: source is required", i)
		}
	}
	
	return nil
}

// HasFileOps returns true if the template has any file operations
func (t *Template) HasFileOps() bool {
	return len(t.Files.Copy) > 0 || len(t.Files.Append) > 0
}
