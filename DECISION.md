# Decisions Record

### Decision: build-time version injection
**Reason:** prevent manual version drift
We inject the version during the GitHub Actions build process using `ldflags` to ensure the version string in the binary always matches the exact git tag.

### Decision: `--yes` flag for non-interactive installs
**Reason:** required for WinGet automation
WinGet installers need to be able to run silently without blocking on user prompts. The `--yes` flag automates the install process by bypassing all interactive prompts and assigning default values.
