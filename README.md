# Forge — Windows Project Bootstrapper CLI

Forge is a Windows CLI for creating new projects from reusable templates.
It runs real tools (like `git`, `uv`, `npm`), copies files, and applies safe patches.

👉 **Official templates:** [https://github.com/Vishnuj-n/forge-templates](https://github.com/Vishnuj-n/forge-templates)

---

## Features

* Create projects from templates
* Run real ecosystem commands (git, uv, npm)
* Global and local templates
* Safe, non-destructive execution
* Single Go binary (no dependencies)
* Windows-friendly install and uninstall

---

## Installation (Recommended)

1. **Download `forge.exe`** from GitHub Releases
   [https://github.com/Vishnuj-n/forge/releases](https://github.com/Vishnuj-n/forge/releases)

2. **Open PowerShell in the folder where `forge.exe` was downloaded**

   * Shift + Right Click → **Open PowerShell here**
   * Or open PowerShell and `cd` into the folder

3. **Run the install command**

```powershell
.\forge.exe install
```

4. **Close and reopen PowerShell**

5. **Verify installation**

```powershell
forge --version
```

---

## Basic Usage

### Pull templates

```powershell
forge pull git
forge pull --all
```

Templates are stored in:

```
%USERPROFILE%\.forge\templates
```

---

### Create a project

```powershell
forge init git ./my-project
```

* Creates (or uses) `./my-project`
* Initializes the project inside that directory

---

### Initialize in current directory

```powershell
forge init git
```

Initializes the project in the **current working directory**.

---

### Other useful commands

```powershell
forge list          # list templates
forge new my-temp   # create a template
forge test my-temp  # test template safely
```

---

## `forge init` behavior

```powershell
forge init <template> [project-dir]
```

* If `project-dir` is provided → Forge creates or uses that directory
* If not provided → Forge uses the current directory
* Forge **never creates projects in global directories**
* Global directories are used **only for templates**

---

## Template Locations (Priority Order)

1. `./templates` (project-local)
2. `$FORGE_TEMPLATES` (if set)
3. `~/.forge/templates` (global)

---

## Documentation

* INSTALL.md — Installation details
* TEMPLATE-GUIDE.md — Writing templates
* ARCHITECTURE.md — Internals
* CONTRIBUTING.md — Contributions

---

## Uninstall

```powershell
forge uninstall
```

If needed, delete manually:

```
%USERPROFILE%\bin\forge.exe
```

---

## License

MIT License

---

**Forge focuses on speed, safety, and predictable project setup.**
