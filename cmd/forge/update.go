package forge

import (
	"fmt"
	"os"

	"forge/internal/update"

	"github.com/spf13/cobra"
)

var checkOnlyUpdate bool

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update forge to the latest version",
	Long: `Check for a newer version of forge on GitHub and install it.

By default this command downloads and replaces the current binary in-place.
Use --check to only report whether an update is available without installing.`,
	Run: runUpdate,
}

func init() {
	updateCmd.Flags().BoolVar(&checkOnlyUpdate, "check", false, "Only check for updates, do not install")
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) {
	current := Version

	fmt.Printf("Current version: %s\n", current)
	fmt.Print("Checking for updates...\n")

	newVersion, available, err := update.CheckUpdate(current)
	if err != nil {
		exitWithError("failed to check for updates", err)
	}

	if !available {
		fmt.Println("forge is already up to date.")
		return
	}

	if current == "development" {
		fmt.Printf("Development build detected — installing latest release: %s\n", newVersion)
	} else {
		fmt.Printf("New version available: %s\n", newVersion)
	}

	if checkOnlyUpdate {
		fmt.Println("Run 'forge update' without --check to install.")
		return
	}

	// Resolve current executable path
	exePath, err := os.Executable()
	if err != nil {
		exitWithError("could not determine current executable path", err)
	}

	if err := update.PerformUpdate(newVersion, exePath); err != nil {
		exitWithError("update failed", err)
	}

	fmt.Printf("forge updated to %s successfully.\n", newVersion)
}
