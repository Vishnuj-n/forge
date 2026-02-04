# Forge - Windows Project Bootstrapper CLI

Forge is a powerful Windows CLI tool that automates project initialization using templates. Create, share, and reuse project templates with zero dependencies.

## ✨ Features

- 📦 **Template-based project generation** - Define projects in simple YAML
- 🔄 **Reusable templates** - Share templates across projects and teams
- 🎯 **Safe execution** - Isolated workspace with two-phase commit
- 📁 **Global templates** - Store templates once, use everywhere
- ⚡ **No dependencies** - Pure Go binary, runs anywhere on Windows
- 🛡️ **Non-destructive** - Fails safely if anything goes wrong

## 🚀 Quick Start (2 minutes)

### 1. Install Forge

**Option A: Download Release (Easiest)**
```powershell
# Download forge.exe from GitHub Releases
# https://github.com/Vishnuj-n/forge/releases

# Run installer
.\forge.exe install

# Close and reopen PowerShell
```

**Option B: Build from Source**
```powershell
git clone https://github.com/Vishnuj-n/forge.git
cd forge
go build -o forge.exe
.\forge.exe install
```

### 2. Set Up Global Templates (Optional but Recommended)
```powershell
# During installation, you'll be asked:
# "Would you like to set up a global templates directory? (yes/no)"

# Answer: yes
# This creates: C:\Users\YourName\.forge\templates
```

### 3. Use Forge

```powershell
# List available templates
forge list

# Create a new project from a template
forge init my-template ./my-new-project

# Create your own template
forge new my-awesome-template

# Test a template without committing
forge test my-template
```

---

## 📖 Installation Methods

### From GitHub Release (Recommended for Users)

1. **Download** the latest `forge.exe` from [Releases](https://github.com/Vishnuj-n/forge/releases)
2. **Run installer:**
   ```powershell
   .\forge.exe install
   ```
3. **Answer the setup question** about global templates
4. **Close and reopen** PowerShell
5. **Verify:**
   ```powershell
   forge --version
   ```

### From Source (For Developers)

**Prerequisites:**
- Windows 10 or later
- Go 1.19 or later
- Git

**Steps:**
```powershell
# Clone repository
git clone https://github.com/Vishnuj-n/forge.git
cd forge

# Build
go build -o forge.exe

# Install
.\forge.exe install

# Run setup
# Answer 'yes' when asked about global templates
```

### Manual Installation (Advanced)

```powershell
# 1. Create installation directory
New-Item -ItemType Directory -Path "$env:USERPROFILE\bin" -Force

# 2. Copy forge.exe
Copy-Item .\forge.exe "$env:USERPROFILE\bin\"

# 3. Add to PATH (PowerShell - already done by installer)
[Environment]::SetEnvironmentVariable(
    "Path",
    "$env:Path;$env:USERPROFILE\bin",
    "User"
)

# 4. Close and reopen PowerShell
```

---

## 🎓 Usage Examples

### Example 1: Create a Project from Template

```powershell
# Initialize a new project with git
forge init simple-git-project ./my-project

# Now you have a git-initialized project!
cd my-project
git log
```

### Example 2: Create Your Own Template

```powershell
# Create template scaffold
forge new my-web-template

# The template directory is created at:
# C:\Users\YourName\.forge\templates\my-web-template\

# Edit the template files:
# - template.yaml      (configuration)
# - README.md         (documentation)
# - files/            (files to copy)
# - patches/          (files to append to)
```

### Example 3: Test Before Committing

```powershell
# Test a template without creating the project
forge test my-web-template

# Inspect the test workspace, then delete it
# Useful for validating templates before sharing
```

### Example 4: List All Templates

```powershell
# Show all available templates
forge list

# Output:
# Templates in: C:\Users\YourName\.forge\templates
# 
# NAME                 COMMANDS  FILE OPS  PATH
# simple-git-project   2         1         C:\Users\YourName\.forge\templates\simple-git-project
# my-web-template      3         2         C:\Users\YourName\.forge\templates\my-web-template
```

---

## 📋 Template Structure

Templates are stored in `~/.forge/templates/template-name/`:

```
my-template/
├── template.yaml          # Template configuration
├── README.md             # Template documentation
├── files/                # Files to copy into project
│   └── src/
│       └── example.txt
└── patches/              # Files to append to existing files
    └── .gitignore
```

### template.yaml Example

```yaml
name: my-template
commands:
  - git init
  - npm init -y
files:
  copy:
    - files/src
  append:
    - target: .gitignore
      source: patches/.gitignore
```

For detailed template documentation, see [TEMPLATE-GUIDE.md](doc/TEMPLATE-GUIDE.md)

---

## ⚙️ Configuration

### Global Templates Directory

Forge looks for templates in this order:
1. `$FORGE_TEMPLATES` environment variable
2. `~/.forge/templates` (default global location)
3. `./templates` (project-local templates)

**Set custom location:**
```powershell
[Environment]::SetEnvironmentVariable("FORGE_TEMPLATES", "C:\my\templates", "User")
```

---

## 🔄 Uninstall

```powershell
# Remove from PATH and clean up
forge uninstall

# You'll see:
# "Please delete C:\Users\YourName\bin\forge.exe manually"

# Delete it:
del "$env:USERPROFILE\bin\forge.exe"

# Reopen PowerShell
```

---

## 🔍 Troubleshooting

### "forge: command not found"
- **Solution:** Close and reopen PowerShell completely (not just a new tab)
- Check PATH: `$env:Path -split ";" | Select-String "bin"`

### "Permission denied" on install
- **Solution:** The installer doesn't need admin for user install. If you get an error:
  - Run: `forge install` (default, user-based)
  - Not: `forge install --system` (requires admin)

### Global templates not found
- **Solution:** Make sure installation completed successfully:
  ```powershell
  $env:FORGE_TEMPLATES  # Should show: C:\Users\YourName\.forge\templates
  ```
- If empty, manually set it:
  ```powershell
  [Environment]::SetEnvironmentVariable("FORGE_TEMPLATES", "$env:USERPROFILE\.forge\templates", "User")
  ```

### Template won't initialize
- **Solution:** Test first to see detailed errors:
  ```powershell
  forge test my-template
  ```

---

## 📚 Documentation

- **[INSTALL.md](./INSTALL.md)** - Detailed installation guide
- **[TEMPLATE-GUIDE.md](doc/TEMPLATE-GUIDE.md)** - How to create and structure templates
- **[ARCHITECTURE.md](doc/ARCHITECTURE.md)** - Technical design and internals
- **[CONTRIBUTING.md](./CONTRIBUTING.md)** - How to contribute

---

## 💡 Common Workflows

### Share Templates with Team

1. Create template on your machine:
   ```powershell
   forge new team-template
   # Edit the template...
   ```

2. Share the template directory:
   ```powershell
   # Copy C:\Users\YourName\.forge\templates\team-template
   # To: \\shared-drive\templates\team-template
   ```

3. Team members use it:
   ```powershell
   set FORGE_TEMPLATES=\\shared-drive\templates
   forge list
   forge init team-template
   ```

### Use Project-Local Templates

```powershell
# Create templates in your project
mkdir templates
forge new my-project-template

# Others can use:
forge list templates
forge init my-project-template
```

---

## 🛡️ Safety Features

- **Workspace Isolation** - All operations happen in a temporary directory
- **Two-Phase Commit** - Atomic operations (or best-effort cross-volume)
- **Non-Destructive** - If anything fails, your project directory is untouched
- **Append-Only Patching** - Patches only append to files, never modify existing content

---

## 📄 License

MIT License - see [LICENSE](./LICENSE)

---

## 🤝 Contributing

Contributions welcome! See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

---

## 📞 Support

- **Issues:** [GitHub Issues](https://github.com/Vishnuj-n/forge/issues)
- **Discussions:** [GitHub Discussions](https://github.com/Vishnuj-n/forge/discussions)

---

## 🎯 Roadmap

- [x] Template-based project initialization
- [x] Global templates directory
- [x] Template scaffolding
- [x] User-based installation (no admin)
- [ ] Configuration file support (`forge.config.yaml`)
- [ ] Interactive mode enhancements
- [ ] Multi-platform support (macOS, Linux)
- [ ] Package managers (Chocolatey, Scoop)

---

**Made with ❤️ by the Forge community**
