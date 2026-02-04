# Documentation Workflow

## Agent Roles & Responsibilities

---

### 🎯 Agent 1: Documentation Planner (Current)

**Role:** Creates specifications and requirements for all documentation

**Responsibilities:**
- Define what documentation is needed
- Specify content requirements for each doc
- Create and maintain this checklist
- Review completed docs (optional)

**Process:**
1. Analyze project structure and purpose
2. List required documentation below
3. Add checkboxes for tracking
4. Provide brief/summary for each item

---

### ✍️ Agent 2: Documentation Writer

**Role:** Creates the actual documentation files

**Responsibilities:**
- Read the briefs below
- Write comprehensive documentation
- Check off items when complete
- Ensure consistency and quality

---

## 📋 Documentation Checklist

### Required Documentation

- [ ] **README.md**
  - **Brief:** Main project documentation covering: project overview, what Forge does, installation instructions, basic usage examples (`forge init`, `forge test`), template structure explanation, and quick start guide
  - **Target Audience:** End users and developers wanting to use Forge
  - **Key Sections:** Overview, Installation, Quick Start, Usage Examples, Template Format, Safety Features
  - **Tone:** Clear, concise, practical

- [ ] **ARCHITECTURE.md**
  - **Brief:** Technical architecture documentation covering: system design, module breakdown (template, workspace, executor, fileops, commit), execution flow, safety model (two-phase commit, workspace isolation), and design decisions
  - **Target Audience:** Contributors and architects evaluating the codebase
  - **Key Sections:** Overview, Module Architecture, Execution Flow, Safety Model, Design Principles
  - **Tone:** Technical, detailed, architectural

- [ ] **TEMPLATE-GUIDE.md**
  - **Brief:** Complete guide for creating templates covering: template.yaml structure, command execution rules, file operations (copy/append), example templates, best practices, and common patterns
  - **Target Audience:** Template authors
  - **Key Sections:** Template Structure, Command Syntax, File Operations, Complete Examples, Best Practices, Troubleshooting
  - **Tone:** Tutorial-style, practical examples

- [ ] **CONTRIBUTING.md**
  - **Brief:** Contribution guidelines covering: development setup, code structure, testing requirements, pull request process, coding standards, and how to add new features
  - **Target Audience:** Contributors
  - **Key Sections:** Getting Started, Development Workflow, Testing, Code Standards, PR Guidelines
  - **Tone:** Welcoming, clear instructions

- [ ] **CHANGELOG.md**
  - **Brief:** Version history and changes. Start with v0.1.0 (initial release) covering: core features implemented (init/test commands, template system, workspace isolation, two-phase commit), and note that this is the initial implementation based on the design plan
  - **Target Audience:** Users tracking versions
  - **Format:** Standard keepachangelog.com format
  - **Tone:** Factual, organized by version

---

## 📝 Notes for Agent 2

### Project Context
- **Project Name:** Forge
- **Purpose:** Safety-first Windows CLI tool for project bootstrapping
- **Language:** Go
- **Key Feature:** Transactional workspace with two-phase commit
- **Philosophy:** Safety over convenience, deterministic execution

### Key Implementation Details
- Commands: `forge init <template>` and `forge test <template>`
- Flag: `--interactive` / `-i` for interactive mode
- Template format: YAML with commands (token arrays), files.copy, files.append
- Execution: Temp workspace → Execute commands → Apply files → Commit
- Safety: Cross-volume detection, append-only patches, fail-fast errors

### Source Files to Reference
- `cmd/forge/*.go` - CLI commands
- `internal/template/template.go` - Template parsing
- `internal/workspace/workspace.go` - Workspace management
- `internal/executor/executor.go` - Command execution
- `internal/fileops/fileops.go` - File operations
- `internal/commit/commit.go` - Two-phase commit
- `plan.md` - Complete design philosophy
- `example.yaml` - Template example

### Documentation Standards
- Use clear headings and structure
- Include code examples with syntax highlighting
- Add practical examples over theory
- Cross-reference between docs where helpful
- Keep README concise, move details to specialized docs

---

## ✅ Completion Criteria

Each document is considered complete when it:
1. Covers all points in the brief
2. Is technically accurate
3. Includes relevant code examples
4. Follows consistent formatting
5. Has been checked off above

---

## 📦 New Addition: Template Creation Command (forge new)

### CLI Command for Template Generation

- [ ] **forge new <template-name>**
  - **Purpose:** Generate a new template directory structure with all necessary files
  - **Behavior:** 
    - Creates `templates/<template-name>/` directory
    - Generates `template.yaml` skeleton with comments
    - Creates empty `files/` and `patches/` subdirectories
    - Generates `README.md` for template documentation
  - **Output:** 
    ```
    templates/my-template/
    ├── template.yaml      (with example structure + comments)
    ├── README.md          (template documentation)
    ├── files/             (empty, for files to copy)
    └── patches/           (empty, for append files)
    ```
  - **Usage:** `forge new my-awesome-template`
  - **Target User:** Template creators

### Template Creation Details

The `forge new` command should:

1. **Validate template name** — alphanumeric + hyphens only
2. **Check for duplicates** — warn if template already exists
3. **Generate template.yaml** with:
   - Comments explaining each section
   - Example commands (git init, etc.)
   - Commented-out file operations
   - Clear structure for copy/append
4. **Generate README.md** explaining:
   - What this template does
   - How to customize it
   - Prerequisites (tools needed)
5. **Create empty directories** — files/ and patches/
6. **Output success message** with next steps

### Example Generated `template.yaml`

```yaml
# Template: my-awesome-template
# Generated by: forge scaffold

name: my-awesome-template

# Commands to execute (in order)
commands:
  - cmd: ["git", "init"]
  # - cmd: ["git", "config", "user.name", "Your Name"]

# File operations (copy and append)
files:
  # Copy files/directories from this template
  copy:
    - files/README.md
    # - files/config.json
  
  # Append content to existing files (created by commands or copies)
  append:
    - target: ".gitignore"
      source: "patches/gitignore.append"
    # - target: "package.json"
    #   source: "patches/package.json.append"
```

---

**Status:** Agent 2 should create documentation files + implement forge new command

---

## 🔧 Implementation Task for forge new

When implementing, follow these guidelines:

### Code Structure
- Create `cmd/forge/new.go` for the CLI command
- Add helper in `internal/scaffold/` for template generation
- Write tests in `internal/scaffold/new_test.go`

### Error Handling
- Validate template name (no special chars, length limits)
- Check if target directory already exists
- Clear error messages for failures

### Generated Content
- Use heredocs or embedded templates for scaffolding
- Ensure YAML is valid and properly indented
- Make generated README practical and actionable

### User Experience
- Show success message with template location
- Provide next steps: "Edit files/ and patches/, then use: forge test templates/my-template"
- Make it easy to get started immediately
