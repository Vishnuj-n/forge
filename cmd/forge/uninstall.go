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
		installDir = filepath.Join("C:", string(filepath.Separator), "Program Files", "Forge")
		exePath = filepath.Join(installDir, "forge.exe")
	} else {
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
	globalTemplatesDir := ""
	configPath := ""
	if userProfile != "" {
		globalTemplatesDir = filepath.Join(userProfile, ".forge", "templates")
		configPath = filepath.Join(userProfile, ".forge", "config.yaml")
		if _, err := os.Stat(globalTemplatesDir); err == nil {
			fmt.Printf("\nRemoving global templates directory: %s\n", globalTemplatesDir)
			if err := os.RemoveAll(globalTemplatesDir); err != nil {
				fmt.Printf("⚠ Warning: Could not remove templates directory: %v\n", err)
			} else {
				fmt.Println("✓ Removed templates directory")
			}
		}
		// Remove config file
		if _, err := os.Stat(configPath); err == nil {
			fmt.Printf("\nRemoving config file: %s\n", configPath)
			if err := os.Remove(configPath); err != nil {
				fmt.Printf("⚠ Warning: Could not remove config file: %v\n", err)
			} else {
				fmt.Println("✓ Removed config file")
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

	// Remove FORGE_TEMPLATES environment variable
	removeForgeTemplatesEnv()

	// Create and spawn cleanup script
	fmt.Println("\nRemoving forge.exe...")
	if err := createAndSpawnCleanupScript(exePath); err != nil {
		fmt.Printf("⚠ Warning: Could not spawn cleanup script: %v\n", err)
		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Println("Forge has been partially uninstalled.")
		fmt.Println(strings.Repeat("=", 60))
		fmt.Printf("\nPlease delete the executable manually:\n  %s\n", exePath)
		fmt.Println("\nYou can delete it using:")
		fmt.Printf("  del \"%s\"\n", exePath)
		return
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Forge has been uninstalled successfully!")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\nCleanup in progress...")
	fmt.Println("The executable will be removed automatically.")
	fmt.Println("\nTo reinstall Forge later, download and run: forge install")
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

func removeForgeTemplatesEnv() {
	// Remove FORGE_TEMPLATES environment variable
	psCmd := `[Environment]::SetEnvironmentVariable("FORGE_TEMPLATES", $null, "User")`
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCmd)
	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠ Warning: Could not remove FORGE_TEMPLATES variable: %v\n", err)
	} else {
		fmt.Println("✓ Removed FORGE_TEMPLATES environment variable")
	}
}

func createAndSpawnCleanupScript(exePath string) error {
	// Create cleanup script in TEMP directory
	tempDir := os.Getenv("TEMP")
	if tempDir == "" {
		tempDir = os.TempDir()
	}

	// Generate unique script name
	scriptPath := filepath.Join(tempDir, fmt.Sprintf("forge-cleanup-%d.ps1", os.Getpid()))

	// PowerShell script that:
	// 1. Waits for forge.exe to exit
	// 2. Deletes forge.exe
	// 3. Deletes itself
	cleanupScript := fmt.Sprintf(`
# Forge cleanup script
# Wait for forge.exe to exit
Start-Sleep -Seconds 2

# Attempt to delete forge.exe
$targetFile = "%s"
$maxAttempts = 10
$attempt = 0

while ($attempt -lt $maxAttempts) {
    try {
        if (Test-Path $targetFile) {
            Remove-Item -Path $targetFile -Force -ErrorAction Stop
            Write-Host "✓ Deleted forge.exe"
            break
        } else {
            Write-Host "✓ forge.exe already removed"
            break
        }
    }
    catch {
        $attempt++
        if ($attempt -lt $maxAttempts) {
            Start-Sleep -Milliseconds 500
        } else {
            Write-Host "⚠ Could not delete forge.exe: $_"
            Write-Host "Please delete manually: $targetFile"
        }
    }
}

# Delete this cleanup script
Start-Sleep -Seconds 1
Remove-Item -Path $PSCommandPath -Force -ErrorAction SilentlyContinue
`, strings.ReplaceAll(exePath, `\`, `\\`))

	// Write cleanup script
	if err := os.WriteFile(scriptPath, []byte(cleanupScript), 0644); err != nil {
		return fmt.Errorf("failed to create cleanup script: %w", err)
	}

	// Spawn PowerShell to run cleanup script in background
	// -WindowStyle Hidden keeps it invisible
	// -ExecutionPolicy Bypass allows script to run
	cmd := exec.Command("powershell",
		"-WindowStyle", "Hidden",
		"-ExecutionPolicy", "Bypass",
		"-NoProfile",
		"-File", scriptPath)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to spawn cleanup script: %w", err)
	}

	// Don't wait for the command - let it run in background
	// The cleanup script will delete forge.exe after we exit

	return nil
}
