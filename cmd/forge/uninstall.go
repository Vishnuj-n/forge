package forge

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall forge",
	Long: `Remove forge from your bin directory and PATH.

By default, removes from user bin directory.
For system-wide uninstall, use: forge uninstall --system`,
	Run: runUninstall,
}

var systemUninstall bool

func init() {
	uninstallCmd.Flags().BoolVar(&systemUninstall, "system", false, "Uninstall from Program Files")
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall(cmd *cobra.Command, args []string) {
	var installDir, exePath string
	
	if systemUninstall {
		// System-wide uninstallation
		installDir = filepath.Join("C:", string(filepath.Separator), "Program Files", "Forge")
		exePath = filepath.Join(installDir, "forge.exe")
	} else {
		// User-based uninstallation (default)
		userProfile := os.Getenv("USERPROFILE")
		if userProfile == "" {
			exitWithError("USERPROFILE environment variable not set", nil)
		}
		installDir = filepath.Join(userProfile, "bin")
		exePath = filepath.Join(installDir, "forge.exe")
	}
	
	// Check if installed
	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		if systemUninstall {
			fmt.Println("Forge is not installed in Program Files")
		} else {
			fmt.Println("Forge is not installed in user bin directory")
		}
		return
	}
	
	fmt.Println("Forge is installed at:", installDir)
	fmt.Print("\nDo you want to uninstall? (yes/no): ")
	
	var response string
	fmt.Scanln(&response)
	if response != "yes" && response != "y" {
		fmt.Println("Uninstall cancelled.")
		return
	}
	
	fmt.Println("\nUninstalling Forge...")
	
	// Remove installation directory or file
	var removeErr error
	if systemUninstall {
		removeErr = os.RemoveAll(installDir)
	} else {
		// For user install, just remove the exe
		removeErr = os.Remove(exePath)
	}
	
	if removeErr != nil {
		if systemUninstall {
			fmt.Println("\n⚠ NOTE: Uninstallation from Program Files requires Administrator privileges.")
			fmt.Println("\nTo uninstall:")
			fmt.Println("1. Run PowerShell as Administrator")
			fmt.Println("2. Run: forge uninstall --system")
		} else {
			exitWithError("failed to remove forge.exe", removeErr)
		}
		os.Exit(1)
	}
	
	if systemUninstall {
		fmt.Println("✓ Removed from Program Files")
	} else {
		fmt.Println("✓ Removed from user bin directory")
	}
	
	// Inform about PATH
	fmt.Println("\n⚠ You may need to manually remove from PATH:")
	fmt.Println("1. Open: Settings → System → About → Advanced system settings")
	fmt.Println("2. Click: Environment Variables")
	fmt.Println("3. Under User variables, edit Path")
	fmt.Printf("4. Remove: %s\n", installDir)
	fmt.Println("5. Click OK and restart your terminal")
	
	fmt.Println("\n✓ Uninstall complete!")
}
