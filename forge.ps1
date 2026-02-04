param(
    [ValidateSet("install", "uninstall")]
    [string]$Action = "install"
)

$InstallDir = Join-Path $env:LOCALAPPDATA "Forge"
$ExeName = "forge.exe"
$ExePath = Join-Path $InstallDir $ExeName

function Add-ToPath {
    $path = [Environment]::GetEnvironmentVariable("Path", "User") -split ";"
    if ($path -notcontains $InstallDir) {
        $path += $InstallDir
        [Environment]::SetEnvironmentVariable("Path", ($path -join ";"), "User")
        Write-Host "[OK] Added Forge to PATH"
    }
}

function Remove-FromPath {
    $path = [Environment]::GetEnvironmentVariable("Path", "User") -split ";"
    if ($path -contains $InstallDir) {
        $path = $path | Where-Object { $_ -ne $InstallDir }
        [Environment]::SetEnvironmentVariable("Path", ($path -join ";"), "User")
        Write-Host "[OK] Removed Forge from PATH"
    }
}

function Install-Forge {
    Write-Host "Installing Forge..."

    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null

    $localForge = Join-Path $PSScriptRoot $ExeName
    if (-not (Test-Path $localForge)) {
        Write-Host "[ERROR] forge.exe not found in current directory"
        Write-Host "Build it first:"
        Write-Host "  go build -o forge.exe"
        exit 1
    }

    Copy-Item $localForge $ExePath -Force
    Add-ToPath

    Write-Host "[OK] Forge installed at $ExePath"
    Write-Host "Restart your terminal to use 'forge'"
}

function Uninstall-Forge {
    Write-Host "Uninstalling Forge..."

    Get-Process forge -ErrorAction SilentlyContinue | Stop-Process -Force

    Remove-FromPath

    if (Test-Path $InstallDir) {
        Remove-Item -Recurse -Force $InstallDir
        Write-Host "[OK] Forge removed"
    } else {
        Write-Host "Forge is not installed"
    }
}

if ($Action -eq "install") {
    Install-Forge
} else {
    Uninstall-Forge
}
