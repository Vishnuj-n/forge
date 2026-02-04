# Installation Guide

## Quick Install

### For Users (From Release)

1. **Download** `forge.exe` from [GitHub Releases](https://github.com/Vishnuj-n/forge/releases)

2. **Run the installer:**
   ```powershell
   .\forge.exe install
   ```

3. **Answer the setup question:**
   ```
   Would you like to set up a global templates directory? (yes/no): yes
   ```

4. **Close and reopen PowerShell**, then verify:
   ```powershell
   forge --version
   forge --help
   ```

---

## For Developers (From Source)

### Prerequisites
- Windows 10 or later
- Go 1.19 or later
- Git

### Installation Steps

```powershell
# Clone the repository
git clone https://github.com/Vishnuj-n/forge.git
cd forge

# Build the binary
go build -o forge.exe

# Run the installer
.\forge.exe install

# Close and reopen PowerShell
```

---

## What Installation Does

### User Installation (Default - No Admin Needed)

```powershell
forge install
```

This command:
- ✅ Copies `forge.exe` to `%USERPROFILE%\bin`
- ✅ Adds `%USERPROFILE%\bin` to your User PATH
- ✅ Creates global templates directory: `~\.forge\templates`
- ✅ Sets `FORGE_TEMPLATES` environment variable
- ✅ No administrator privileges required

### System Installation (Admin Required)

```powershell
forge install --system
```

This command:
- ✅ Copies `forge.exe` to `C:\Program Files\Forge`
- ✅ Adds to system PATH (affects all users)
- ❌ Requires administrator privileges

---

## Global Templates Directory Setup

During installation, Forge asks:
```
Would you like to set up a global templates directory? (yes/no):
```

**Choose `yes` if you want to:**
- Create templates once and use them in multiple projects
- Share templates with team members
- Keep templates organized in one place

**The setup:**
- Creates: `C:\Users\YourName\.forge\templates`
- Sets environment variable: `FORGE_TEMPLATES`
- Forge will use this location for all `forge new` and `forge init` commands

**If you choose `no`:**
- Templates stored locally in `./templates` of each project
- You can set it up later manually:
  ```powershell
  mkdir "$env:USERPROFILE\.forge\templates"
  [Environment]::SetEnvironmentVariable("FORGE_TEMPLATES", "$env:USERPROFILE\.forge\templates", "User")
  ```

---

## Uninstall

```powershell
# Uninstall Forge
forge uninstall

# You'll see:
# Forge has been uninstalled.
# Please delete the executable manually:
#   C:\Users\YourName\bin\forge.exe
```

The uninstall command:
- ✅ Removes from User PATH
- ✅ Removes global templates directory
- ⚠️ Leaves `forge.exe` (delete manually)

**Delete manually:**
```powershell
del "$env:USERPROFILE\bin\forge.exe"
```

---

## Troubleshooting

### "forge: command not found"

**Solution:** PowerShell needs to reload the PATH
```powershell
# Close PowerShell completely (don't just open a new tab)
# Then reopen PowerShell and try again
forge --version
```

### Install Failed - Permission Denied

**Solution:** Using default user installation (doesn't need admin):
```powershell
# Don't use --system flag
forge install

# If you really need system installation:
# Right-click PowerShell → "Run as Administrator"
# Then run: forge install --system
```

### Can't Create Templates

**Solution:** Check if global templates directory exists:
```powershell
# Check the environment variable
$env:FORGE_TEMPLATES

# If empty, create manually:
mkdir "$env:USERPROFILE\.forge\templates"
[Environment]::SetEnvironmentVariable("FORGE_TEMPLATES", "$env:USERPROFILE\.forge\templates", "User")

# Close and reopen PowerShell
```

### PATH Still Shows Old Value

**Solution:** PowerShell caches the PATH variable
```powershell
# Reload PATH in current session:
$env:Path = [Environment]::GetEnvironmentVariable("Path", "User")

# Or close and reopen PowerShell
```

---

## Verify Installation

After installation, verify everything works:

```powershell
# Check version
forge --version
# Output: forge version 0.1.0

# Check help
forge --help
# Output: Shows all available commands

# List templates
forge list
# Output: Shows templates from global directory

# Check environment variables
$env:FORGE_TEMPLATES
# Output: C:\Users\YourName\.forge\templates
```

---

## Environment Variables

### FORGE_TEMPLATES

**Purpose:** Set custom location for templates

**Default:** `C:\Users\YourName\.forge\templates`

**Set custom location:**
```powershell
[Environment]::SetEnvironmentVariable("FORGE_TEMPLATES", "C:\my\custom\templates", "User")
```

**Reload in current session:**
```powershell
$env:FORGE_TEMPLATES = [Environment]::GetEnvironmentVariable("FORGE_TEMPLATES", "User")
```

---

## Getting Help

- **README:** See [README.md](./README.md) for overview and usage
- **Template Guide:** See [TEMPLATE-GUIDE.md](doc/TEMPLATE-GUIDE.md) for creating templates
- **Architecture:** See [ARCHITECTURE.md](doc/ARCHITECTURE.md) for technical details
- **Issues:** [GitHub Issues](https://github.com/Vishnuj-n/forge/issues)
