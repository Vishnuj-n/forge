package template

import (
	"os"
	"testing"
)

func TestTemplateLoad(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
		wantName string
	}{
		{
			name: "valid template",
			yaml: `name: test-template
commands:
  - cmd: ["git", "init"]
files:
  copy:
    - files/
  append:
    - target: ".gitignore"
      source: "patches/gitignore.append"`,
			wantErr: false,
			wantName: "test-template",
		},
		{
			name: "missing name",
			yaml: `commands:
  - cmd: ["git", "init"]`,
			wantErr: true,
		},
		{
			name: "empty command",
			yaml: `name: bad-template
commands:
  - cmd: []`,
			wantErr: true,
		},
		{
			name: "valid minimal template",
			yaml: `name: minimal
commands:
  - cmd: ["echo", "hello"]`,
			wantErr: false,
			wantName: "minimal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpFile, err := os.CreateTemp("", "template-*.yaml")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tt.yaml); err != nil {
				t.Fatalf("failed to write temp file: %v", err)
			}
			tmpFile.Close()

			// Test Load
			tmpl, err := Load(tmpFile.Name())
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tmpl.Name != tt.wantName {
				t.Errorf("Load() got name %q, want %q", tmpl.Name, tt.wantName)
			}
		})
	}
}

func TestCommandString(t *testing.T) {
	cmd := Command{Cmd: []string{"git", "init", "--bare"}}
	expected := "git init --bare"
	if got := cmd.String(); got != expected {
		t.Errorf("Command.String() = %q, want %q", got, expected)
	}
}

func TestHasFileOps(t *testing.T) {
	tests := []struct {
		name string
		tmpl *Template
		want bool
	}{
		{
			name: "no file ops",
			tmpl: &Template{Name: "test", Files: FileOps{}},
			want: false,
		},
		{
			name: "has copy",
			tmpl: &Template{Name: "test", Files: FileOps{Copy: []string{"files/"}}},
			want: true,
		},
		{
			name: "has append",
			tmpl: &Template{Name: "test", Files: FileOps{Append: []AppendPatch{{Target: ".gitignore", Source: "patches/gitignore.append"}}}},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tmpl.HasFileOps(); got != tt.want {
				t.Errorf("HasFileOps() = %v, want %v", got, tt.want)
			}
		})
	}
}
