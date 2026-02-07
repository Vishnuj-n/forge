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
	Name        string    `yaml:"name"`
	Description string    `yaml:"description,omitempty"`
	Version     string    `yaml:"version,omitempty"`
	Commands    []Command `yaml:"commands"`
	Files       FileOps   `yaml:"files"`
}

// Command represents a single command to execute
type Command struct {
	Cmd         []string `yaml:"cmd"`
	Interactive bool     `yaml:"interactive"`
	TestCmd     []string `yaml:"test_cmd"`
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
// It accepts both full paths and template names
// For template names, it searches in:
// 1. ./templates/<name>
// 2. $FORGE_TEMPLATES/<name>
// 3. $HOME/.forge/templates/<name>
func Load(templatePath string) (*Template, error) {
	// Resolve the template path
	resolvedPath, err := ResolveTemplatePath(templatePath)
	if err != nil {
		return nil, fmt.Errorf("template path not found: %w", err)
	}

	return loadFromPath(resolvedPath)
}

// ResolveTemplatePath resolves a template path by:
// 1. Trying the literal path
// 2. If not found, treating it as a name and searching standard locations
func ResolveTemplatePath(templatePath string) (string, error) {
	// First, try the literal path
	if absPath, err := filepath.Abs(templatePath); err == nil {
		if _, err := os.Stat(absPath); err == nil {
			return absPath, nil
		}
	}

	// If not found as literal path, treat it as a template name
	// and search in standard locations
	searchPaths := getSearchPaths(templatePath)
	for _, searchPath := range searchPaths {
		if _, err := os.Stat(searchPath); err == nil {
			return searchPath, nil
		}
	}

	// Not found anywhere
	return "", fmt.Errorf("template '%s' not found in any search location", templatePath)
}

// getSearchPaths returns the list of paths to search for a template name
func getSearchPaths(templateName string) []string {
	var paths []string

	// 1. Check ./templates/<templateName>
	paths = append(paths, filepath.Join("templates", templateName))

	// 2. Check $FORGE_TEMPLATES/<templateName>
	if envPath := os.Getenv("FORGE_TEMPLATES"); envPath != "" {
		paths = append(paths, filepath.Join(envPath, templateName))
	}

	// 3. Check $HOME/.forge/templates/<templateName>
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".forge", "templates", templateName))
	}

	return paths
}

// loadFromPath loads a template from a resolved path
func loadFromPath(resolvedPath string) (*Template, error) {
	// Check if it's a directory or file
	info, err := os.Stat(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("template path not found: %w", err)
	}

	var yamlPath string
	if info.IsDir() {
		// If directory, look for template.yaml inside
		yamlPath = filepath.Join(resolvedPath, "template.yaml")
	} else {
		// If file, use it directly
		yamlPath = resolvedPath
	}

	// Read YAML file
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// Parse using shared Parse function
	tmpl, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template YAML: %w", err)
	}

	return tmpl, nil
}

// Parse parses template YAML data into a Template and validates it.
// This is a public helper to allow callers to parse YAML without duplicating logic.
func Parse(data []byte) (*Template, error) {
	var tmpl Template
	if err := yaml.Unmarshal(data, &tmpl); err != nil {
		return nil, fmt.Errorf("failed to unmarshal template YAML: %w", err)
	}

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
