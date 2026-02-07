package forge

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"forge/internal/template"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [templates-directory]",
	Short: "List available templates",
	Long: `List all available templates from a directory.

If no directory is specified, looks for templates in:
1. ./templates (current directory)
2. $FORGE_TEMPLATES environment variable
3. $HOME/.forge/templates

Templates can be:
- Directories containing template.yaml
- Direct .yaml files`,
	Args: cobra.MaximumNArgs(1),
	Run:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	// Determine templates directory
	var templatesDir string

	if len(args) == 1 {
		templatesDir = args[0]
	} else {
		// Try default locations
		templatesDir = findTemplatesDir()
	}

	if templatesDir == "" {
		fmt.Println("No templates directory found.")
		fmt.Println("\nSearched locations:")
		fmt.Println("  - ./templates")
		if envPath := os.Getenv("FORGE_TEMPLATES"); envPath != "" {
			fmt.Printf("  - %s (from FORGE_TEMPLATES)\n", envPath)
		}
		if home, err := os.UserHomeDir(); err == nil {
			fmt.Printf("  - %s\n", filepath.Join(home, ".forge", "templates"))
		}
		fmt.Println("\nCreate a templates directory or specify one: forge list <path>")
		return
	}

	// Check if directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		fmt.Printf("Templates directory not found: %s\n", templatesDir)
		return
	}

	fmt.Printf("Templates in: %s\n\n", templatesDir)

	// Find and list templates
	templates, err := discoverTemplates(templatesDir)
	if err != nil {
		exitWithError("failed to discover templates", err)
	}

	if len(templates) == 0 {
		fmt.Println("No templates found.")
		return
	}

	// Display templates in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tCOMMANDS\tFILE OPS\tPATH")
	fmt.Fprintln(w, "----\t--------\t--------\t----")

	for _, tmplInfo := range templates {
		fmt.Fprintf(w, "%s\t%d\t%d\t%s\n",
			tmplInfo.Name,
			len(tmplInfo.Template.Commands),
			len(tmplInfo.Template.Files.Copy)+len(tmplInfo.Template.Files.Append),
			tmplInfo.RelPath)
	}

	w.Flush()
}

type templateInfo struct {
	Name     string
	RelPath  string
	Template *template.Template
}

func discoverTemplates(baseDir string) ([]templateInfo, error) {
	var templates []templateInfo

	// Walk the directory
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Look for template.yaml or *.yaml files
		if !info.IsDir() && (info.Name() == "template.yaml" || filepath.Ext(info.Name()) == ".yaml") {
			// Read file bytes and parse directly to avoid duplicated I/O
			data, err := os.ReadFile(path)
			if err != nil {
				return nil // skip unreadable
			}

			tmpl, err := template.Parse(data)
			if err != nil {
				return nil // skip invalid templates
			}

			// Calculate relative path
			relPath, err := filepath.Rel(baseDir, filepath.Dir(path))
			if err != nil {
				relPath = filepath.Dir(path)
			}
			if relPath == "." {
				relPath = filepath.Base(filepath.Dir(path))
			}

			templates = append(templates, templateInfo{
				Name:     tmpl.Name,
				RelPath:  relPath,
				Template: tmpl,
			})
		}

		return nil
	})

	return templates, err
}

func findTemplatesDir() string {
	// 1. Try ./templates
	if _, err := os.Stat("templates"); err == nil {
		return "templates"
	}

	// 2. Try FORGE_TEMPLATES env var
	if envPath := os.Getenv("FORGE_TEMPLATES"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath
		}
	}

	// 3. Try $HOME/.forge/templates
	if home, err := os.UserHomeDir(); err == nil {
		homeTemplates := filepath.Join(home, ".forge", "templates")
		if _, err := os.Stat(homeTemplates); err == nil {
			return homeTemplates
		}
	}

	return ""
}
