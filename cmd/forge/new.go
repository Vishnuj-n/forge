package forge

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"forge/internal/scaffold"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <template-name>",
	Short: "Create a new template",
	Long: `Create a new template with a complete directory structure.

Generates:
- template.yaml (with examples and comments)
- README.md (template documentation)
- files/ directory (for template files)
- patches/ directory (for append patches)

Example:
  forge new my-awesome-template`,
	Args: cobra.ExactArgs(1),
	Run:  runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func runNew(cmd *cobra.Command, args []string) {
	templateName := args[0]

	// Determine template directory (prefer global over local)
	templatesDir := determineTemplatesDirForNew()

	// Create template generator
	gen := scaffold.New(templatesDir)

	// Validate name
	if err := scaffold.ValidateName(templateName); err != nil {
		exitWithError("invalid template name", err)
	}

	fmt.Printf("Creating new template: %s\n", templateName)
	fmt.Printf("Location: %s\n", templatesDir)
	fmt.Println()

	// Generate template
	templateDir, err := gen.Generate(templateName)
	if err != nil {
		exitWithError("failed to create template", err)
	}

	// Print success message
	fmt.Println("")
	fmt.Println(scaffold.GetNextSteps(templateDir, templateName))

	// Prompt to open the template folder
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to open the template folder? [y/N]: ")
	resp, _ := reader.ReadString('\n')
	resp = strings.TrimSpace(strings.ToLower(resp))
	if resp == "y" || resp == "yes" {
		openErr := openPath(templateDir)
		if openErr != nil {
			fmt.Printf("Failed to open template folder: %v\n", openErr)
		}
	}
}

func openPath(path string) error {
	// Only use explorer for Windows
	return exec.Command("explorer", path).Start()
}

func determineTemplatesDirForNew() string {
	// Check global templates directory first
	if forgeTemplates := os.Getenv("FORGE_TEMPLATES"); forgeTemplates != "" {
		if _, err := os.Stat(forgeTemplates); err == nil {
			fmt.Println("📁 Using global templates directory")
			return forgeTemplates
		}
	}

	// Check home/.forge/templates
	if home, err := os.UserHomeDir(); err == nil {
		globalDir := filepath.Join(home, ".forge", "templates")
		if _, err := os.Stat(globalDir); err == nil {
			fmt.Println("📁 Using global templates directory")
			return globalDir
		}
	}

	// Fall back to local ./templates
	fmt.Println("📁 Using local templates directory")
	return "templates"
}
