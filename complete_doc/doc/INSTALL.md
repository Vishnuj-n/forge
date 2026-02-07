
# Installation Guide
## Downloading Templates from the Official Repository

Forge now supports downloading and updating templates directly from the official Forge templates repository using the `forge pull` command.

### Download a Single Template

```powershell
forge pull <template-name>
# Example:
forge pull git
```

### Download All Templates
# Installation (Short)

User install (recommended):

1. Download `forge.exe` from Releases
2. Run: `.orge.exe install`
3. Answer the setup prompt to create a global templates directory (optional)
4. Close and reopen PowerShell

Developer build:

```powershell
git clone https://github.com/Vishnuj-n/forge.git
cd forge
go build -o forge.exe
.\forge.exe install
```

Notes:
- Default user install places `forge.exe` in `%USERPROFILE%\bin` and sets `FORGE_TEMPLATES` to `%USERPROFILE%\.forge\templates`.
- System install (`--system`) requires admin and installs to `C:\Program Files\Forge`.

Uninstall:

```powershell
forge uninstall
del "$env:USERPROFILE\bin\forge.exe"  # if left behind
```

Quick verification:

```powershell
forge --version
forge list
```

If you need help, see `README.md` and `doc/TEMPLATE-GUIDE.md`.

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
