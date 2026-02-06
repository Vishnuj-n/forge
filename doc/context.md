# Forge — Short Context

Forge is a Windows-only CLI tool that bootstraps projects from declarative templates.

Key points:
- Runs template commands inside a temporary workspace (safe, isolated)
- Uses a two-phase commit: prepare in temp, then commit to target only on success
- Non-interactive by default; enable interactive commands with `--interactive`
- Templates are YAML; authors declare all test behavior (no guessing)

Keep templates small, documented, and testable.

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
