package forge

import (
	"fmt"
	"os"
	"path/filepath"

	"forge/internal/executor"
	"forge/internal/fileops"
	"forge/internal/template"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init <template-path> [target-directory]",
	Short: "Initialize a new project from a template",
	Long: `Initialize a new project by:
1. Running commands directly in the target directory
2. Copying template files
3. Applying append-only patches

Commands inherit your terminal's stdin/stdout/stderr, so interactive
commands (like npm init, cargo init) work naturally.

Target directory defaults to current working directory if not specified.`,
	Args: cobra.RangeArgs(1, 2),
	Run:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	templatePath := args[0]

	// Determine target directory
	targetDir := "."
	if len(args) == 2 {
		targetDir = args[1]
	}

	// Convert to absolute path
	absTargetDir, err := filepath.Abs(targetDir)
	if err != nil {
		exitWithError("failed to resolve target directory", err)
	}

	// Check if target directory exists and is empty
	if err := validateTargetDirectory(absTargetDir); err != nil {
		exitWithError("target directory validation failed", err)
	}

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

	fmt.Printf("Initializing project from template: %s\n", tmpl.Name)

	// Create target directory if it doesn't exist
	if err := os.MkdirAll(absTargetDir, 0755); err != nil {
		exitWithError("failed to create target directory", err)
	}

	fmt.Printf("Working in target directory: %s\n", absTargetDir)

	// Execute commands directly in target directory
	if len(tmpl.Commands) > 0 {
		fmt.Println("\nExecuting commands:")
		exec := executor.New(absTargetDir, false, false) // false for testMode = forge init mode
		for i, cmdDef := range tmpl.Commands {
			fmt.Printf("  [%d/%d] %s\n", i+1, len(tmpl.Commands), cmdDef.String())
			if err := exec.Run(cmdDef); err != nil {
				exitWithError(fmt.Sprintf("command failed: %s", cmdDef.String()), err)
			}
		}
	}

	// Apply file operations directly in target directory
	if tmpl.HasFileOps() {
		fmt.Println("\nApplying file operations:")
		fops := fileops.New(absTargetDir, resolvedTemplatePath)

		if err := fops.CopyFiles(tmpl.Files.Copy); err != nil {
			exitWithError("failed to copy files", err)
		}

		if err := fops.ApplyAppends(tmpl.Files.Append); err != nil {
			exitWithError("failed to apply patches", err)
		}
	}

	fmt.Printf("\n✓ Project initialized successfully at: %s\n", absTargetDir)
}

func validateTargetDirectory(dir string) error {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		// Directory doesn't exist - we'll create it during commit
		return nil
	}
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("target path exists but is not a directory")
	}

	// Check if directory is empty
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	if len(entries) > 0 {
		return fmt.Errorf("target directory is not empty")
	}

	return nil
}
