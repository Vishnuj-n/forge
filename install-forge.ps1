# Forge - Install/Uninstall Script
# Downloads Forge and installs it to Program Files with PATH integration
# Run as Administrator

param(
    [switch]$Uninstall = $false
)

# Configuration
$InstallDir = "C:\Program Files\Forge"
$ExeName = "forge.exe"
$ExePath = Join-Path $InstallDir $ExeName
$TemplatesDir = Join-Path $InstallDir "templates"

function Test-Admin {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Get-ForgeVersion {
    if (Test-Path $ExePath) {
        try {
            $version = & $ExePath --version 2>$null
            return $version
        } catch {
            return $null
        }
    }
    return $null
}

function Test-ForgeInstalled {
    return Test-Path $ExePath
}

function Add-ToPath {
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$InstallDir*") {
        $newPath = "$userPath;$InstallDir"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        Write-Host "✓ Added $InstallDir to User PATH" -ForegroundColor Green
        Write-Host "  (Restart your terminal to use 'forge' command globally)" -ForegroundColor Yellow
    } else {
        Write-Host "✓ $InstallDir already in PATH" -ForegroundColor Green
    }
}

function Remove-FromPath {
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -like "*$InstallDir*") {
        # Remove the path entry carefully
        $pathArray = $userPath -split ";"
        $newPathArray = @($pathArray | Where-Object { $_ -ne $InstallDir })
        $newPath = $newPathArray -join ";"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        Write-Host "✓ Removed $InstallDir from PATH" -ForegroundColor Green
    }
}

function Uninstall-Forge {
    if (-not (Test-ForgeInstalled)) {
        Write-Host "Forge is not installed." -ForegroundColor Yellow
        return $true
    }

    $version = Get-ForgeVersion
    Write-Host ""
    Write-Host "Forge is currently installed at: $ExePath" -ForegroundColor Cyan
    if ($version) {
        Write-Host "Version: $version" -ForegroundColor Cyan
    }
    Write-Host ""

    # Ask for confirmation
    $response = Read-Host "Do you want to uninstall Forge? (yes/no)"
    if ($response -ne "yes" -and $response -ne "y") {
        Write-Host "Uninstall cancelled." -ForegroundColor Yellow
        return $false
    }

    try {
        # Close any running forge processes
        Get-Process forge -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
        Start-Sleep -Milliseconds 500

        # Remove from PATH first
        Remove-FromPath

        # Remove installation directory
        if (Test-Path $InstallDir) {
            Remove-Item -Recurse -Force $InstallDir
            Write-Host "✓ Uninstalled Forge from $InstallDir" -ForegroundColor Green
        }

        Write-Host ""
        Write-Host "✓ Forge has been successfully uninstalled" -ForegroundColor Green
        return $true
    } catch {
        Write-Host "✗ Error during uninstall: $_" -ForegroundColor Red
        return $false
    }
}

function Install-Forge {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "  Forge Installation Script" -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host ""

    # Check if already installed
    if (Test-ForgeInstalled) {
        $version = Get-ForgeVersion
        Write-Host "Forge is already installed at: $ExePath" -ForegroundColor Yellow
        if ($version) {
            Write-Host "Version: $version" -ForegroundColor Yellow
        }
        Write-Host ""

        $response = Read-Host "Do you want to reinstall/update? (yes/no)"
        if ($response -ne "yes" -and $response -ne "y") {
            Write-Host "Installation cancelled." -ForegroundColor Yellow
            return $false
        }

        Write-Host "Uninstalling previous version..." -ForegroundColor Cyan
        Uninstall-Forge | Out-Null
        Write-Host ""
    }

    # Download forge.exe
    Write-Host "Downloading Forge..." -ForegroundColor Cyan
    
    # Create installation directory
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        Write-Host "✓ Created $InstallDir" -ForegroundColor Green
    }

    # Check if forge.exe exists in current directory
    $localForge = Join-Path $PSScriptRoot "forge.exe"
    if (Test-Path $localForge) {
        Write-Host "Found local forge.exe, using that..." -ForegroundColor Cyan
        Copy-Item $localForge $ExePath
        Write-Host "✓ Copied forge.exe to $InstallDir" -ForegroundColor Green
    } else {
        # Try to download from GitHub (if releases are available)
        Write-Host "No local forge.exe found." -ForegroundColor Yellow
        Write-Host "Please build forge.exe first:" -ForegroundColor Yellow
        Write-Host "  go build -o forge.exe" -ForegroundColor Gray
        Write-Host ""
        Write-Host "Then run this script again from the same directory." -ForegroundColor Yellow
        return $false
    }

    # Verify installation
    if (-not (Test-Path $ExePath)) {
        Write-Host "✗ Installation failed: forge.exe not found at $ExePath" -ForegroundColor Red
        return $false
    }

    # Add to PATH
    Add-ToPath

    # Verify forge command works
    Write-Host ""
    Write-Host "Verifying installation..." -ForegroundColor Cyan
    try {
        $testOutput = & $ExePath --help 2>&1
        if ($LASTEXITCODE -eq 0 -or $testOutput -like "*Forge*") {
            Write-Host "✓ Forge is working correctly" -ForegroundColor Green
        }
    } catch {
        Write-Host "⚠ Could not verify installation: $_" -ForegroundColor Yellow
    }

    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "✓ Installation Complete!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Usage:" -ForegroundColor Cyan
    Write-Host "  forge init <template>          Initialize a project" -ForegroundColor Gray
    Write-Host "  forge test <template>          Test a template" -ForegroundColor Gray
    Write-Host "  forge new <template-name>      Create a new template" -ForegroundColor Gray
    Write-Host "  forge list                     List available templates" -ForegroundColor Gray
    Write-Host "  forge --help                   Show help" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Important:" -ForegroundColor Yellow
    Write-Host "  Close and reopen your terminal to use 'forge' command globally" -ForegroundColor Yellow
    Write-Host ""

    return $true
}

# Main script execution
if (-not (Test-Admin)) {
    Write-Host "✗ This script must be run as Administrator" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please run PowerShell as Administrator and try again:" -ForegroundColor Yellow
    Write-Host "  Right-click PowerShell > Run as administrator" -ForegroundColor Gray
    exit 1
}

if ($Uninstall) {
    Uninstall-Forge
} else {
    Install-Forge
}
