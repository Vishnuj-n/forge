package forge

import (
	"fmt"

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

	// Create template generator
	gen := scaffold.New("templates")

	// Validate name
	if err := scaffold.ValidateName(templateName); err != nil {
		exitWithError("invalid template name", err)
	}

	fmt.Printf("Creating new template: %s\n", templateName)

	// Generate template
	templateDir, err := gen.Generate(templateName)
	if err != nil {
		exitWithError("failed to create template", err)
	}

	// Print success message
	fmt.Println("")
	fmt.Println(scaffold.GetNextSteps(templateDir, templateName))
}
