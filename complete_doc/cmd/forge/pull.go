package forge

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"forge/internal/remote"
)

var pullCmd = &cobra.Command{
	Use:   "pull [template-name]",
	Short: "Download templates from the official repository",
	Long: `Download templates from the official Forge templates repository.

Examples:
  forge pull git              # Download a single template
  forge pull --all            # Download all available templates

Templates are stored in: %USERPROFILE%\.forge\templates`,
	Run: runPull,
}

var pullAll bool

func init() {
	pullCmd.Flags().BoolVar(&pullAll, "all", false, "Pull all templates from repository")
	rootCmd.AddCommand(pullCmd)
}

func runPull(cmd *cobra.Command, args []string) {
	// Validate arguments
	if !pullAll && len(args) == 0 {
		fmt.Println("Error: template name required")
		fmt.Println("Usage: forge pull <template-name>")
		fmt.Println("   or: forge pull --all")
		os.Exit(1)
	}

	if pullAll && len(args) > 0 {
		fmt.Println("Error: cannot specify template name with --all flag")
		os.Exit(1)
	}

	// Get global templates directory
	globalDir, err := getGlobalTemplatesDir()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Ensure global templates directory exists
	if err := os.MkdirAll(globalDir, 0755); err != nil {
		fmt.Printf("Error creating templates directory: %v\n", err)
		os.Exit(1)
	}

	if pullAll {
		// Pull all templates
		if err := pullAllTemplates(globalDir); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Pull single template
		templateName := args[0]
		if err := pullSingleTemplate(templateName, globalDir); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func getGlobalTemplatesDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(home, ".forge", "templates"), nil
}

func pullSingleTemplate(templateName, globalDir string) error {
	fmt.Println("Downloading templates...")

	zipURL := "https://github.com/Vishnuj-n/forge-templates/archive/refs/heads/main.zip"
	zipPath, err := remote.DownloadRepoZip(zipURL)
	if err != nil {
		return err
	}
	defer os.Remove(zipPath)

	fmt.Printf("Installing template '%s'...\n", templateName)
	if err := remote.InstallSingleTemplate(zipPath, templateName, globalDir); err != nil {
		return err
	}

	fmt.Printf("Template '%s' installed successfully.\n", templateName)
	return nil
}

func pullAllTemplates(globalDir string) error {
	fmt.Println("Downloading templates...")

	zipURL := "https://github.com/Vishnuj-n/forge-templates/archive/refs/heads/main.zip"
	zipPath, err := remote.DownloadRepoZip(zipURL)
	if err != nil {
		return err
	}
	defer os.Remove(zipPath)

	fmt.Println("Installing all templates...")
	installed, err := remote.InstallAllTemplates(zipPath, globalDir)
	if err != nil {
		return err
	}

	for _, name := range installed {
		fmt.Printf("✓ %s\n", name)
	}
	fmt.Printf("Completed. %d templates installed or updated.\n", len(installed))
	return nil
}
