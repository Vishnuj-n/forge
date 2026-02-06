# Forge Architecture Guide

This document provides a technical overview of Forge's internal design, module responsibilities, and execution flow. It is intended for contributors and architects who want to understand how Forge works under the hood.

## 🏗️ System Design Overview

Forge is a **transactional, workflow-aware project bootstrapper**. It allows users to create projects using declarative templates that mix ecosystem-native commands (like `npm init` or `git init`) with file operations, all within a safe, isolated environment.

### Core Design Principles

1.  **Safety by Design:** All operations occur in a temporary workspace. User data is never modified until all steps succeed.
2.  **Deterministic Execution:** Templates are declarative (YAML) and sequential. There is no hidden logic or "magic".
3.  **Two-Phase Commit:** The "prepare" phase happens in a temp directory. The "commit" phase moves the result to the target directory.
4.  **Windows Atomicity:** Forge is designed with Windows filesystem constraints in mind, detecting cross-volume operations and handling them safely.

---

## 🧩 Module Breakdown

Forge is organized into distinct modules, each with a single responsibility.

### 1. CLI Layer (`cmd/forge/`)
-   **Role:** Handles argument parsing, flag validation, and command routing.
-   **Key Components:** `root.go`, `init.go`, `test.go`, `new.go`.
-   **Responsibility:** It validates user input and hands off control to the core logic. It does *not* contain business logic.

### 2. Template Module (`internal/template/`)
-   **Role:** Parses and validates `template.yaml` files.
-   **Key Components:** `template.go`
-   **Responsibility:**
    -   Reads YAML configuration.
    -   Validates structure (commands, files, patches).
    -   Ensures no forbidden operations are requested.
    -   Supports optional metadata fields: `description` and `version`.
    -   Recognizes `interactive` and `test_cmd` on commands to support deterministic `forge test` runs.

### 3. Workspace Module (`internal/workspace/`)
-   **Role:** Manages the temporary execution environment.
-   **Key Components:** `workspace.go`
-   **Responsibility:**
    -   Creates a secure temporary directory for execution.
    -   Ensures isolation from the host system.
    -   Cleans up resources after execution (or failure).

### 4. Executor Module (`internal/executor/`)
-   **Role:** Runs external commands defined in the template.
-   **Key Components:** `executor.go`
-   **Responsibility:**
    -   Executes commands (e.g., `git`, `npm`) within the workspace context.
    -   Handles `stdin`/`stdout` streams (suppressed by default, attached in `--interactive` mode).
    -   Enforces "fail-fast" behavior on command errors.
    -   Supports a *test mode* that replaces or skips interactive commands: `test_cmd` is used when present; otherwise interactive steps are skipped with a clear log message.

### 5. File Operations (`internal/fileops/`)
-   **Role:** Handles file copying and patching.
-   **Key Components:** `fileops.go`
-   **Responsibility:**
    -   **Copy:** recursively copies files from the template `files/` directory to the workspace.
    -   **Append:** applies append-only patches from `patches/` to existing files.
    -   **Verify:** ensures targets for patches exist.

### 6. Commit Module (`internal/commit/`)
-   **Role:** Finalizes the transaction.
-   **Key Components:** `commit.go`
-   **Responsibility:**
    -   Moves the contents of the temporary workspace to the final user destination.
    -   **Safety Check:** Detects if the temp dir and target dir are on different volumes.
        -   **Atomic Move:** Uses `MoveFileEx` where possible, or falls back to a safe copy-and-delete strategy with rollback capabilities if a move isn't possible.

### 7. Remote / Pull Module (`internal/remote/`)
- **Role:** Download and install templates from the official templates repository.
- **Key Components:** `download.go` (ZIP download, prefix detection, extraction helpers)
- **Responsibility:**
    - Download repository ZIP from GitHub into a temp file and detect the dynamic top-level prefix (e.g., `forge-templates-main/`).
    - Enumerate and validate top-level directories; install a single template or all templates into `%USERPROFILE%\\.forge\\templates`.
    - Replace existing templates atomically (remove then extract) and provide clear errors for network, extraction, or validation failures.

---

## 🔄 Execution Flow

When a user runs `forge init <template> <target>`, the following sequence occurs:

1.  **Validation:**
    -   CLI validates arguments.
    -   Template loader finds and parses `template.yaml`.

2.  **Workspace Creation (Phase 1 Start):**
    -   A temporary directory is created (e.g., `%TEMP%\forge-xxxx`).

3.  **Command Execution:**
    -   Commands defined in `template.yaml` are executed sequentially in the temp workspace.
    -   If any command fails (non-zero exit code), the process aborts, and the temp dir is cleaned up.

4.  **File Operations:**
    -   **Copy:** Files from `template/files/` are copied over the workspace.
    -   **Patch:** Content from `template/patches/` is appended to target files in the workspace.

5.  **Commit (Phase 2):**
    -   If all previous steps succeed, the Commit module takes over.
    -   It verifies the target directory doesn't already exist (or is empty).
    -   It moves the workspace content to the target path.

6.  **Cleanup:**
    -   The temporary directory wrapper is removed.

---

## 🛡️ Safety Model

### Workspace Isolation
By running everything in `%TEMP%`, we ensure that:
-   Partial failures don't leave "half-baked" projects in your working directory.
-   Scripts inside templates cannot accidentally delete or modify your actual files during the generation phase.

### Two-Phase Commit
Forge separates **Generation** (Phase 1) from **Commit** (Phase 2).
-   **Generation:** Can fail at any point without side effects.
-   **Commit:** Is the only point where the user's filesystem is permanently altered.

### Append-Only Patching
Forge intentionally limits file modifications to **append-only**.
-   **Why?** Merging structured files (JSON, YAML, XML) is complex and error-prone without specific parsers.
-   **Behavior:** Forge simply adds content to the end of a file. This is safe for `.gitignore`, `.env`, and many config formats.
-   **Conflict:** If a complex merge is needed, the template should provide the full file instead.

---

## 📐 Design Decisions

### No Shell Execution
Forge executes commands as token arrays (`["npm", "install"]`), not shell strings (`"npm install"`).
-   **Reason:** Prevents shell injection attacks and cross-platform shell incompatibilities (PowerShell vs CMD vs Bash).

### No "Magic" Logic
Templates cannot contain conditionals (`if`, `loop`).
-   **Reason:** Keeps templates readable and deterministic. If complex logic is needed, it should be wrapped in a script included in the template.

### Windows-First
Forge is optimized for Windows.
-   **Reason:** Windows filesystem semantics (locking, volume boundaries) are often overlooked in cross-platform tools. Forge explicitly handles these cases.
