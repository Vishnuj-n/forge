# Forge — Windows Project Bootstrapper CLI

Forge is a Windows command-line tool that helps you start new projects using reusable templates.

It runs **real tools** like `git`, `python`, and `uv` for you, so you don’t have to repeat the same setup steps every time.

👉 **Official templates:** [https://github.com/Vishnuj-n/forge-templates](https://github.com/Vishnuj-n/forge-templates)

---

## Why Forge?

Starting a new project usually means doing the same things again and again:

* Initialize version control
* Set up a language environment
* Create common files like `README.md` and `.gitignore`
* Follow the same structure every time
* Fix mistakes if something fails halfway

Forge lets you save these steps as a template and reuse them with one command.

It does **not replace your tools**.
It simply runs them for you in a repeatable way.

---

## Example: Python Project

### Without Forge

```powershell
git init
python -m venv .venv
.venv\Scripts\activate

# create README.md
# create .gitignore
# repeat for every project
```

You do this manually for every new Python project.

---

### With Forge

```powershell
forge init python
```

Forge sets up the project **in the current directory**, runs all required commands, and creates standard files in one step.

---

## When Should You Use Forge?

Forge is useful if you:

* Create many projects with similar setup steps
* Want the same structure every time
* Work mainly on Windows
* Prefer automation over manual setup
* Want to use real tools instead of generators

---

## Features

* Create projects from templates
* Run real ecosystem commands (`git`, `python`, `uv`, `npm`)
* Support for global and local templates
* Safe and predictable execution
* Single Go binary (no dependencies)
* Windows-friendly install and uninstall

---

## Installation

1. Download **`forge.exe`** from GitHub Releases
   [https://github.com/Vishnuj-n/forge/releases](https://github.com/Vishnuj-n/forge/releases)

2. Open **PowerShell in the folder where `forge.exe` was downloaded**

   * Shift + Right Click → **Open PowerShell here**
   * Or open PowerShell and `cd` into the folder

3. Run:

```powershell
.\forge.exe install
```

4. Close and reopen PowerShell

5. Verify:

```powershell
forge --version
```

---

## Basic Usage

````markdown
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

Initializes the project in the **current working directory** which should be empty.

---

### Other Commands

```powershell
forge list          # show available templates
forge pull python   # download a template
forge new my-temp   # create a new template
forge test my-temp  # test a template safely
```

---

## Where Templates Are Stored

Forge looks for templates in this order:

1. `./templates` (project-local)
2. `$FORGE_TEMPLATES` (if set)
3. `%USERPROFILE%\.forge\templates` (global)

Global templates are stored in:

```
%USERPROFILE%\.forge\templates
```

---

## How Forge Works

* Forge runs commands exactly as written in the template
* It does not modify or replace your tools
* Interactive tools work normally during `forge init`
* Templates define **commands and files**, not generators
* Projects are never created in global directories

---

## Uninstall

```powershell
forge uninstall
```

This removes:

* Forge executable
* Global templates
* Forge configuration

If removal fails, you may need to delete `forge.exe` manually from your bin folder.

---

## Documentation

* `doc\INSTALL.md` — Installation details
* `doc\TEMPLATE-GUIDE.md` — Writing templates
* `doc\ARCHITECTURE.md` — Internal design
* `CONTRIBUTING.md` — Contribution guide

---

## License

MIT License

---

**Forge helps you start projects faster, with fewer mistakes and less repetition.**