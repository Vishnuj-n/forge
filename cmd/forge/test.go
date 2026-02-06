package forge

import (
	"fmt"

	"forge/internal/executor"
	"forge/internal/fileops"
	"forge/internal/template"
	"forge/internal/workspace"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test <template-path>",
	Short: "Test a template without committing",
	Long: `Test a template by running the full workflow in a temporary workspace:
1. Creating an isolated temporary workspace
2. Running commands declared in the template
3. Copying template files
4. Applying append-only patches
5. Displaying the workspace path for inspection

The workspace is NOT committed to any target directory.
The workspace path is displayed so you can inspect the result.`,
	Args: cobra.ExactArgs(1),
	Run:  runTest,
}

func init() {
	testCmd.Flags().BoolVarP(&interactiveFlag, "interactive", "i", false, "Enable interactive mode for commands that require user input")
	rootCmd.AddCommand(testCmd)
}

func runTest(cmd *cobra.Command, args []string) {
	templatePath := args[0]

	// Resolve template path (handles both full paths and template names)
	resolvedTemplatePath, err := template.ResolveTemplatePath(templatePath)
	if err != nil {
		exitWithError("failed to resolve template", err)
	}

	// Load template
	tmpl, err := template.Load(templatePath)
	if err != nil {
		exitWithError("failed to load template", err)
	}

	fmt.Printf("Testing template: %s\n", tmpl.Name)

	// Create workspace
	ws, err := workspace.New()
	if err != nil {
		exitWithError("failed to create workspace", err)
	}
	// Note: We don't defer cleanup here - we want to keep it for inspection

	fmt.Printf("Working in temporary workspace: %s\n", ws.Path())

	// Execute commands
	if len(tmpl.Commands) > 0 {
		fmt.Println("\nExecuting commands:")
		exec := executor.New(ws.Path(), interactiveFlag, true)
		for i, cmdDef := range tmpl.Commands {
			fmt.Printf("  [%d/%d] %s\n", i+1, len(tmpl.Commands), cmdDef.String())
			if err := exec.Run(cmdDef); err != nil {
				exitWithError(fmt.Sprintf("command failed: %s", cmdDef.String()), err)
			}
		}
	}

	// Apply file operations
	if tmpl.HasFileOps() {
		fmt.Println("\nApplying file operations:")
		fops := fileops.New(ws.Path(), resolvedTemplatePath)

		if err := fops.CopyFiles(tmpl.Files.Copy); err != nil {
			exitWithError("failed to copy files", err)
		}

		if err := fops.ApplyAppends(tmpl.Files.Append); err != nil {
			exitWithError("failed to apply patches", err)
		}
	}

	fmt.Println("\n✓ Template test completed successfully")
	fmt.Printf("\nWorkspace location: %s\n", ws.Path())
	fmt.Println("(Workspace will persist for inspection - delete manually when done)")
}
