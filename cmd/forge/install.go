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

For system-wide installation to Program Files, use: forge install --system

Flags:
  --force      Re-run full setup and prompts
  --bin-only   Only replace binary, skip all setup
`,
	Run: runInstall,
}

var systemInstall bool
var forceInstall bool
var binOnlyInstall bool

func init() {
	installCmd.Flags().BoolVar(&systemInstall, "system", false, "Install to Program Files (requires Administrator)")
	installCmd.Flags().BoolVar(&forceInstall, "force", false, "Re-run full setup and prompts")
	installCmd.Flags().BoolVar(&binOnlyInstall, "bin-only", false, "Only replace binary, skip all setup")
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) {
	var installDir, exePath string
	userProfile := os.Getenv("USERPROFILE")
	if systemInstall {
		installDir = filepath.Join("C:", string(filepath.Separator), "Program Files", "Forge")
		exePath = filepath.Join(installDir, "forge.exe")
	} else {
		if userProfile == "" {
			exitWithError("USERPROFILE environment variable not set", nil)
		}
		installDir = filepath.Join(userProfile, "bin")
		exePath = filepath.Join(installDir, "forge.exe")
	}

	// Config file path
	configPath := ""
	if userProfile != "" {
		configPath = filepath.Join(userProfile, ".forge", "config.yaml")
	}

	// Get current executable path
	currentExe, err := os.Executable()
	if err != nil {
		exitWithError("failed to get current executable path", err)
	}

	// Check config file
	configExists := false
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			configExists = true
		}
	}

	// Install/replace binary
	fmt.Println("Installing Forge...")
	fmt.Printf("Source: %s\n", currentExe)
	fmt.Printf("Target: %s\n", exePath)
	if systemInstall {
		fmt.Println("Mode:   System-wide (Program Files)")
	} else {
		fmt.Println("Mode:   User installation")
	}
	fmt.Println()

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

	// Install logic
	if binOnlyInstall {
		fmt.Println("\n✓ Binary replaced (bin-only mode). No setup or prompts.")
		return
	}

	if configExists && !forceInstall {
		fmt.Println("\nForge already installed. Skipping setup and prompts.")
		fmt.Println("✓ Done")
		return
	}

	// First install or --force: run full setup
	if err := setupGlobalTemplates(); err != nil {
		fmt.Printf("⚠ Warning: Global templates setup failed: %v\n", err)
	}

	// Write config file
	if configPath != "" {
		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err == nil {
			_ = os.WriteFile(configPath, []byte("templates_initialized: true\n"), 0644)
		}
	}

	fmt.Println("\n✓ Forge installed and configured.")
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

	// Check if global templates directory already exists
	if _, err := os.Stat(globalTemplatesDir); err == nil {
		fmt.Println("\n⚠ Global templates directory already exists!")
		fmt.Print("Do you want to preserve existing templates? (yes/no): ")
		var preserveResponse string
		fmt.Scanln(&preserveResponse)

		if preserveResponse == "yes" || preserveResponse == "y" {
			fmt.Printf("\n✓ Preserving existing templates in: %s\n", globalTemplatesDir)
			// Set environment variable if not already set
			return setForgeTemplatesEnvVar(globalTemplatesDir)
		}

		fmt.Println("⚠ Warning: Existing templates will be removed!")
		fmt.Print("Are you sure you want to continue? (yes/no): ")
		var confirmResponse string
		fmt.Scanln(&confirmResponse)

		if confirmResponse != "yes" && confirmResponse != "y" {
			fmt.Println("\n✓ Installation cancelled. Existing templates preserved.")
			return nil
		}

		// Remove existing directory
		if err := os.RemoveAll(globalTemplatesDir); err != nil {
			fmt.Printf("⚠ Warning: Could not remove existing templates directory: %v\n", err)
			return nil
		}
		fmt.Println("✓ Removed existing templates directory")
	}

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
	return setForgeTemplatesEnvVar(globalTemplatesDir)
}

func setForgeTemplatesEnvVar(globalTemplatesDir string) error {
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
