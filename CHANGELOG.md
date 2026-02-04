# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2023-10-27

### Added
-   **Core Commands:**
    -   `forge init`: Initialize a new project from a template.
    -   `forge test`: Test a template in a temporary workspace without committing.
    -   `forge new`: Scaffold a new template directory.
    -   `forge list`: List available templates.
    -   `forge install`: Self-installation helper.
-   **Template System:**
    -   YAML-based configuration (`template.yaml`).
    -   Support for executing external commands (token arrays).
    -   File copying (`files/` directory).
    -   Append-only file patching (`patches/` directory).
-   **Safety Features:**
    -   **Workspace Isolation:** All operations occur in `%TEMP%` first.
    -   **Two-Phase Commit:** Atomic (or best-effort) move to target directory only upon success.
    -   **Cross-Volume Detection:** Protects against partial moves across drives.
-   **Global Templates:**
    -   Support for `~/.forge/templates` and `FORGE_TEMPLATES` environment variable.

### Notes
-   Initial release based on the V1 design plan.
-   Windows-only support.
