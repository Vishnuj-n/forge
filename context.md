# AI-FRIENDLY PROJECT CONTEXT — FORGE (PASTE VERBATIM)

**READ THIS CAREFULLY AND FOLLOW IT STRICTLY.**

You are assisting with **Forge**, a **Windows-only CLI developer tool** written in Go.

---

## What Forge Is

Forge is a **workflow-aware project bootstrapper**.

It:

* Runs ecosystem-native initialization commands (e.g. `uv init`, `npm init`)
* Executes everything inside an **OS-provided temporary workspace**
* Applies template file overlays
* Applies **append-only** file patches
* Commits the final result to the user directory **only at the end**

Forge **orchestrates existing tools**.
It does **not** replace them.

---

## Core Safety Model (NON-NEGOTIABLE)

* All commands and file operations run in a **temporary directory**
* User project directory is untouched until commit
* No partial writes
* No silent behavior
* Failure aborts safely

This is a **transactional execution model**.

---

## Interactivity Rules (UPDATED)

* Forge is **non-interactive by default**
* Child process `stdin` is closed unless explicitly enabled
* Interactive tools fail fast with a clear error message

Interactive mode is enabled **only by the user** via:

```
--interactive
```

When enabled:

* stdin/stdout/stderr are forwarded
* User may respond to prompts
* Execution still occurs in temp workspace
* Commit semantics remain unchanged

Templates must **never force interactivity**.

---

## Hard Constraints (DO NOT VIOLATE)

* Windows only
* YAML is the only configuration format
* Commands are token arrays, e.g. `["uv","init"]`
* **No shell execution** (no `cmd.exe`, no PowerShell, no shell strings)
* No file merging
* No structured file editing (JSON/TOML/XML)
* Append-only patches only
* No inference, guessing, or auto-detection
* No refactors unless explicitly requested
* No new features beyond what is asked

If a request conflicts with these rules, **do not implement it**.

---

## Execution Phases (FIXED ORDER)

1. Create temp workspace
2. Run commands
3. Copy template files (`files/`)
4. Apply append-only patches (`patches/`)
5. Optional inspection (`forge test`)
6. Commit to user directory

Order is not configurable.

---

## Template Model

Each template has:

```
template/
├─ template.yaml
├─ files/     # full files copied verbatim
└─ patches/   # partial content appended to existing files
```

* `files/` → creates or replaces files
* `patches/` → appends only, target must exist
* Command-generated files are authoritative

---

## Execution Modes

* `forge init <template>`

  * Default
  * Non-interactive
  * Runs full workflow
  * Commits automatically

* `forge test <template>`

  * Runs full workflow in temp
  * Does NOT commit
  * Prints commands executed and files created
  * Supports `--interactive` only if user enables it

---

## Non-Goals (V1)

Forge does **not**:

* Manage dependencies
* Install tools
* Support plugins
* Support cross-platform
* Merge or edit structured configs
* Provide interactive prompts by default
* Optimize for performance over correctness

---

## Your Role as AI

* Implement **only** what is requested
* Keep solutions minimal and explicit
* Prefer correctness over elegance
* Ask for clarification **only if absolutely required**
* Do not introduce abstractions, patterns, or features not asked for

---

## Guiding Principle

> **Forge favors determinism, safety, and clarity over flexibility.**
