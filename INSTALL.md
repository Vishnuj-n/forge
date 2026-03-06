# Forge Installation Guide

## Automated / Silent Installation

If you are scripting the installation or installing via an automated tool like WinGet, you can use the `--yes` (or `-y`) flag to run the installation non-interactively.

```powershell
forge install --yes
```

When the `--yes` flag is provided, Forge will:
* Skip all prompts
* Automatically accept setting up the global template directory
* Use the default configuration

## Manual Installation

1. Download **`forge.exe`** from GitHub Releases:
   [https://github.com/Vishnuj-n/forge/releases](https://github.com/Vishnuj-n/forge/releases)
2. Open PowerShell in the folder where `forge.exe` was downloaded.
3. Run:
   ```powershell
   .\forge.exe install
   ```
4. Follow the interactive prompts to set up your environment.
5. Close and reopen PowerShell.
