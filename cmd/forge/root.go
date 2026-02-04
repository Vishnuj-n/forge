package forge

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const Version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "forge",
	Short: "Forge - A safety-first project bootstrapper",
	Long: `Forge is a Windows-only CLI tool that bootstraps new projects by:
- Running ecosystem-native initialization commands
- Applying template file overlays
- Applying safe, append-only file patches
- Committing the final result atomically to a user directory`,
	Version: Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func exitWithError(msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s: %v\n", msg, err)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	}
	os.Exit(1)
}
