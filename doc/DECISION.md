# DECISION.md

This file records major technical and process decisions for the Forge CLI project.

---

## 2026-02-06: Template Execution Model
- `forge init` runs commands directly in the target directory (not a temp workspace).
- Interactive commands (e.g., `npm init`) work natively.
- Partial state is allowed for `forge init` (like native tools).
- `forge test` uses a temporary workspace, never prompts, and uses `test_cmd` or skips interactive steps.
- Two-phase commit is only for `forge test`.

## 2026-02-06: Template Metadata
- Templates support `description` and `version` fields in `template.yaml`.
- These fields are optional and used for documentation and version tracking.

## 2026-02-06: Global Template Directory Protection
- During install/reinstall, user is prompted before overwriting the global templates directory.
- User can choose to preserve or replace existing templates.

## 2026-02-06: `forge pull` Command
- Added `forge pull` to download/update templates from the official repo.
- Supports pulling a single template or all templates at once.
- Existing templates are replaced, no duplicates.

## 2026-02-06: Documentation Simplification
- All documentation rewritten to be concise and action-focused.
- `ARCHITECTURE.md` kept detailed for maintainers.

## 2026-02-07: Release Automation
- Releases are fully automated using GitHub Actions.
- Releases trigger on push to `main` with `[release]` in the merge commit message.
- Tags are not created manually; PATCH version bumps only (e.g., `v0.1.3` → `v0.1.4`).
- Release notes are taken from the final merge commit message.
- Windows `.exe` is cross-compiled on Ubuntu and attached to the release.

## 2026-02-07: Branching and Merging
- All development happens on feature branches.
- `main` is always releasable.
- Feature branches are merged into `main` using squash merge.

---

For rationale and context, see commit messages and `MAINTAINER.md`.
