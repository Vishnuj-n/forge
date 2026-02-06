# Forge Template Guide (Short)

Templates are folders that include a `template.yaml` and optional `files/`, `patches/`, and `README.md`.

Minimum required: `template.yaml` with `name`.

Short example:

```yaml
name: example
description: "Short description"
version: "1.0.0"

commands:
  - cmd: ["git", "init"]

files:
  copy:
    - files/README.md
  append:
    - target: ".gitignore"
      source: "patches/gitignore.append"
```

Quick rules:

- `name` is required.
- `cmd` is an array of tokens (no shell strings).
- Use `interactive: true` for commands that prompt; add `test_cmd` for non-interactive test runs.
- `files.copy` paths are relative to the template and must exist when used.
- `files.append.source` is relative to the template and `target` must exist in the project.

Testing and troubleshooting:

- `forge test <template>` runs commands in a temp workspace (non-interactive). Interactive steps are replaced by `test_cmd` or skipped.
- "target file not found" → ensure the file exists before appending.

Keep templates small, documented, and testable.
