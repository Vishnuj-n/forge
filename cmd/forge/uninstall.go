package forge

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	
	fmt.Println("Forge is installed at:", exePath)
	fmt.Print("\nDo you want to uninstall? (yes/no): ")
	
	var response string
	fmt.Scanln(&response)
	if response != "yes" && response != "y" {
		fmt.Println("Uninstall cancelled.")
		return
	}
	
	fmt.Println("\nUninstalling Forge...")
	
	// Remove global templates directory
	userProfile := os.Getenv("USERPROFILE")
	if userProfile != "" {
		globalTemplatesDir := filepath.Join(userProfile, ".forge", "templates")
		if _, err := os.Stat(globalTemplatesDir); err == nil {
			fmt.Printf("\nRemoving global templates directory: %s\n", globalTemplatesDir)
			if err := os.RemoveAll(globalTemplatesDir); err != nil {
				fmt.Printf("⚠ Warning: Could not remove templates directory: %v\n", err)
			} else {
				fmt.Println("✓ Removed templates directory")
			}
		}
	}
	
	// Remove from PATH
	if err := removeFromPath(installDir); err != nil {
		fmt.Printf("⚠ Warning: Could not automatically remove from PATH: %v\n", err)
		fmt.Println("\nYou may need to manually remove from PATH:")
		fmt.Println("1. Open: Settings → System → About → Advanced system settings")
		fmt.Println("2. Click: Environment Variables")
		fmt.Println("3. Under User variables, edit Path")
		fmt.Printf("4. Remove: %s\n", installDir)
		fmt.Println("5. Click OK and restart your terminal")
	} else {
		fmt.Println("✓ Removed from User PATH")
	}
	
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Forge has been uninstalled.")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("\nPlease delete the executable manually:\n  %s\n", exePath)
	fmt.Println("\nYou can delete it using:")
	fmt.Printf("  del \"%s\"\n", exePath)
	fmt.Println("\nOr from File Explorer:")
	fmt.Println("  1. Open File Explorer")
	fmt.Printf("  2. Navigate to: %s\n", installDir)
	fmt.Println("  3. Right-click forge.exe and delete")
	fmt.Println("\nTo reinstall Forge later, run: forge install")
	fmt.Println(strings.Repeat("=", 60))
}

func removeFromPath(dir string) error {
	// Use PowerShell to remove from User PATH
	psCmd := fmt.Sprintf(`
		$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
		if ($userPath -like "*%s*") {
			$paths = $userPath -split ";"
			$newPaths = $paths | Where-Object { $_ -ne "%s" }
			$newPath = $newPaths -join ";"
			[Environment]::SetEnvironmentVariable("Path", $newPath, "User")
		}
	`, dir, dir)
	
	// Execute PowerShell command
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCmd)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return fmt.Errorf("failed to modify PATH: %w (output: %s)", err, string(output))
	}
	
	return nil
}
