package forge

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install forge to user bin directory",
	Long: `Install forge to your user bin directory and add to PATH.

This command copies forge.exe to %USERPROFILE%\bin and adds it to your User PATH.
No Administrator privileges required.

For system-wide installation to Program Files, use: forge install --system`,
	Run: runInstall,
}

var systemInstall bool

func init() {
	installCmd.Flags().BoolVar(&systemInstall, "system", false, "Install to Program Files (requires Administrator)")
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) {
	var installDir, exePath string
	
	if systemInstall {
		// System-wide installation
		installDir = filepath.Join("C:", string(filepath.Separator), "Program Files", "Forge")
		exePath = filepath.Join(installDir, "forge.exe")
	} else {
		// User-based installation (default)
		userProfile := os.Getenv("USERPROFILE")
		if userProfile == "" {
			exitWithError("USERPROFILE environment variable not set", nil)
		}
		installDir = filepath.Join(userProfile, "bin")
		exePath = filepath.Join(installDir, "forge.exe")
	}
	
	// Get current executable path
	currentExe, err := os.Executable()
	if err != nil {
		exitWithError("failed to get current executable path", err)
	}
	
	fmt.Println("Installing Forge...")
	fmt.Printf("Source: %s\n", currentExe)
	fmt.Printf("Target: %s\n", exePath)
	if systemInstall {
		fmt.Println("Mode:   System-wide (Program Files)")
	} else {
		fmt.Println("Mode:   User installation")
	}
	fmt.Println()
	
	// Check if already installed
	if _, err := os.Stat(exePath); err == nil {
		fmt.Println("⚠ Forge is already installed at:", exePath)
		fmt.Print("Do you want to reinstall/update? (yes/no): ")
		
		var response string
		fmt.Scanln(&response)
		if response != "yes" && response != "y" {
			fmt.Println("Installation cancelled.")
			return
		}
		fmt.Println()
	}
	
	// Create installation directory
	if err := os.MkdirAll(installDir, 0755); err != nil {
		if systemInstall {
			fmt.Println("\n⚠ NOTE: Installation to Program Files requires Administrator privileges.")
			fmt.Println("\nPlease run:")
			fmt.Println("  1. Open PowerShell as Administrator")
			fmt.Println("  2. Run: forge install --system")
			fmt.Println("\nOr use default user installation without --system flag")
		} else {
			exitWithError("failed to create installation directory", err)
		}
		os.Exit(1)
	}
	
	// Copy executable
	if err := copyFile(currentExe, exePath); err != nil {
		exitWithError("failed to copy forge.exe", err)
	}
	
	if systemInstall {
		fmt.Println("✓ Copied forge.exe to Program Files")
	} else {
		fmt.Println("✓ Copied forge.exe to user bin directory")
	}
	
	// Add to PATH
	if err := addToPath(installDir); err != nil {
		fmt.Printf("⚠ Could not automatically add to PATH: %v\n", err)
		fmt.Println("\nManual PATH setup:")
		fmt.Println("1. Open: Settings → System → About → Advanced system settings")
		fmt.Println("2. Click: Environment Variables")
		fmt.Println("3. Under User variables, edit Path")
		fmt.Printf("4. Add: %s\n", installDir)
		fmt.Println("5. Click OK and restart your terminal")
	} else {
		fmt.Println("✓ Added to User PATH")
	}
	
	fmt.Println("\n✓ Installation complete!")
	
	// Setup global templates directory
	setupGlobalTemplates()
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	
	return os.WriteFile(dst, input, 0755)
}

func setupGlobalTemplates() error {
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		return fmt.Errorf("USERPROFILE environment variable not set")
	}

	globalTemplatesDir := filepath.Join(userProfile, ".forge", "templates")

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Global Templates Directory Setup")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\nForge can use a global templates directory for sharing templates")
	fmt.Println("across all your projects.")
	fmt.Printf("\nProposed location: %s\n", globalTemplatesDir)

	fmt.Print("\nWould you like to set up a global templates directory? (yes/no): ")
	var response string
	fmt.Scanln(&response)

	if response != "yes" && response != "y" {
		fmt.Println("\n✓ Skipped global templates setup")
		fmt.Println("\nUsage:")
		fmt.Println("  forge init TEMPLATE           Initialize a project")
		fmt.Println("  forge test TEMPLATE           Test a template")
		fmt.Println("  forge new TEMPLATE_NAME       Create a new template")
		fmt.Println("  forge list                    List templates")
		fmt.Println("\nTo enable global templates later, run:")
		fmt.Printf("  mkdir %s\n", globalTemplatesDir)
		fmt.Printf("  set FORGE_TEMPLATES=%s\n", globalTemplatesDir)
		fmt.Println("\nImportant: Close and reopen your terminal to use the global 'forge' command")
		return nil
	}

	// Create global templates directory
	if err := os.MkdirAll(globalTemplatesDir, 0755); err != nil {
		fmt.Printf("⚠ Warning: Could not create global templates directory: %v\n", err)
		fmt.Println("You can create it manually later:")
		fmt.Printf("  mkdir %s\n", globalTemplatesDir)
		return nil
	}

	fmt.Printf("✓ Created global templates directory: %s\n", globalTemplatesDir)

	// Set environment variable
	fmt.Println("\nSetting FORGE_TEMPLATES environment variable...")
	psCmd := fmt.Sprintf(`
		[Environment]::SetEnvironmentVariable("FORGE_TEMPLATES", "%s", "User")
		Write-Host "✓ Set FORGE_TEMPLATES environment variable"
	`, globalTemplatesDir)

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCmd)
	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠ Warning: Could not automatically set environment variable: %v\n", err)
		fmt.Println("\nYou can set it manually:")
		fmt.Println("1. Open: Settings → System → About → Advanced system settings")
		fmt.Println("2. Click: Environment Variables")
		fmt.Println("3. Under User variables, click New")
		fmt.Println("4. Variable name: FORGE_TEMPLATES")
		fmt.Printf("5. Variable value: %s\n", globalTemplatesDir)
		fmt.Println("6. Click OK")
		return nil
	}

	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println("Global Templates Setup Complete!")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("\nLocation: %s\n", globalTemplatesDir)
	fmt.Println("\nYou can now:")
	fmt.Println("  1. Create templates globally: forge new my-template")
	fmt.Println("  2. Use them across all projects: forge init my-template")
	fmt.Println("  3. Share templates by copying from this directory")
	fmt.Println("\nUsage:")
	fmt.Println("  forge init TEMPLATE           Initialize a project (uses global templates)")
	fmt.Println("  forge test TEMPLATE           Test a template")
	fmt.Println("  forge new TEMPLATE_NAME       Create a new template (stored globally)")
	fmt.Println("  forge list                    List templates")
	fmt.Println("\nImportant: Close and reopen your terminal to use the global 'forge' command")

	return nil
}

func addToPath(dir string) error {
	// Use PowerShell to add to User PATH
	psCmd := fmt.Sprintf(`
		$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
		if ($userPath -notlike "*%s*") {
			$newPath = "$userPath;%s"
			[Environment]::SetEnvironmentVariable("Path", $newPath, "User")
			Write-Host "Added to PATH"
		} else {
			Write-Host "Already in PATH"
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
