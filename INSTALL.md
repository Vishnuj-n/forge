# Installation Guide

## Quick Install (Windows)

### Prerequisites
- Windows 7 or later
- PowerShell 5.0 or later
- Administrator access

### Step 1: Build Forge

```powershell
cd C:\Users\vishn\PROJECT\CLI - GO
go build -o forge.exe
```

### Step 2: Run Installation Script

**Open PowerShell as Administrator**, then:

```powershell
cd C:\Users\vishn\PROJECT\CLI - GO
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
.\install-forge.ps1
```

### Step 3: Verify Installation

Close and reopen PowerShell, then:

```powershell
forge --version
forge --help
```

---

## Uninstall

**Open PowerShell as Administrator**, then:

```powershell
.\install-forge.ps1 -Uninstall
```

The script will:
1. Detect the current installation
2. Ask for confirmation
3. Remove Forge from `Program Files`
4. Remove from system PATH

---

## Manual Installation (Alternative)

If you prefer manual installation:

1. **Build Forge:**
   ```powershell
   go build -o forge.exe
   ```

2. **Create Program Files directory:**
   ```powershell
   New-Item -ItemType Directory -Path "C:\Program Files\Forge" -Force
   ```

3. **Copy forge.exe:**
   ```powershell
   Copy-Item forge.exe "C:\Program Files\Forge\forge.exe"
   ```

4. **Add to PATH:**
   - Open Settings → System → About → Advanced system settings
   - Click "Environment Variables"
   - Under "User variables", edit "Path"
   - Add: `C:\Program Files\Forge`
   - Click OK and restart PowerShell

5. **Verify:**
   ```powershell
   forge --version
   ```

---

## Troubleshooting

### "Forge is not installed"
- Run `forge --version` to verify
- Check if `C:\Program Files\Forge\forge.exe` exists

### "forge command not found"
- Restart your PowerShell terminal (close and reopen)
- Check PATH: `$env:Path | Select-String "Program Files\\Forge"`

### "Script execution is disabled"
- Run: `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser`
- This allows running local scripts

### "This script must be run as Administrator"
- Right-click PowerShell → "Run as administrator"
- Try the installation again

---

## Reinstall / Update

To reinstall or update Forge:

```powershell
# Rebuild
go build -o forge.exe

# Run installer (it will detect existing installation)
.\install-forge.ps1
```

The script will ask if you want to reinstall.

---

## What Gets Installed

```
C:\Program Files\Forge\
├── forge.exe              (Main executable)
└── templates/             (Default templates directory)
```

The installer also:
- Adds `C:\Program Files\Forge` to your User PATH
- Does NOT require Administrator after installation
- Can be uninstalled cleanly via the script

---

## Next: Using Forge

Once installed, see:
- `forge --help` for available commands
- `README.md` for usage examples
- `TEMPLATE-GUIDE.md` for creating custom templates
