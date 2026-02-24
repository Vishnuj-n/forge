# Maintainer Journal — Forge CLI

## Purpose
This document exists so future-me understands:
- how this repo is structured
- how releases work
- what rules must not be broken

---

## Branch Strategy

### main
- Always releasable
- Every release is built from `main`
- No direct commits except via PR or squash merge

### feature/*
- All development happens here
- Many small commits allowed
- Never released directly

---

## Release Rules

- Releases are **automatic** when `[release]` keyword is present
- Triggered only by pushing to `main` with `[release]` in commit message
- Tags are **NOT created manually**
- Version bump: PATCH only (v0.1.x)
- Tag is created by GitHub Actions
- Release notes come from the final merge commit

To release:
- Use squash merge into `main`
- Include `[release]` keyword anywhere in the merge commit message
- Example: `Merge pull request #42: Add new feature [release]`
- The workflow will:
  1. Build `forge.exe`
  2. Increment version (v0.1.5 → v0.1.6)
  3. Create git tag
  4. Push tag to GitHub
  5. Create release with binary attached

**Do NOT:**
- Create tags manually
- Create releases manually via GitHub UI
- Push directly to `main` (always use PRs)
- Use `[release]` on non-main branches (only works on `main`)

Notes:
- `forge install` now persists configuration at `%USERPROFILE%/.forge/config.yaml` to avoid re-running setup on reinstall
- Use `forge install --force` to re-run setup prompts, or `forge install --bin-only` to only replace the executable

## CI/CD Workflow Behavior

### Triggers
| Event | Build | Release |
|-------|-------|---------|
| Push to `main` without `[release]` | ❌ Skipped | ❌ Skipped |
| Push to `main` with `[release]` | ✓ Yes | ✓ Yes |
| PR → `main` (opened/updated) | ✓ Yes | ❌ Skipped |
| PR on feature branches | ❌ Skipped | ❌ Skipped |

### Build Details
- Go version: 1.25
- Runner: ubuntu-latest
- Cross-compiles Windows `.exe`
- CGO disabled
- Output: `dist/forge.exe`
- All tags fetched (`fetch-depth: 0`) for accurate version detection

### Concurrency
- Queue group: `release-{{ branch }}`
- Prevents simultaneous releases on the same branch
- If two `[release]` commits are pushed rapidly, the second waits for the first to finish, then increments correctly

## Commit Conventions

- Small commits allowed in feature branches
- Release commit must be clean and descriptive
- Avoid “merge branch …” commits

---

## Things I Must Remember

- Do NOT manually create tags
- Do NOT create releases from GitHub UI
- Do NOT push directly to `main`
- Do NOT include `[release]` in commits to feature branches
- Keep release notes human-readable
- If a release fails, check GitHub Actions logs before retrying
- The `[release]` keyword is case-sensitive and can appear anywhere in the commit message

## Edge Cases & Troubleshooting

### Release created but version jumped unexpectedly
**Cause:** Multiple tags with non-semver format (e.g., `v1`, `release-1.0`)
**Fix:** Delete malformed tags:
```bash
git tag -d <malformed-tag>
git push origin --delete <malformed-tag>
```
Then retry the release.

### Release failed to push tag
**Cause:** GitHub Actions token lacks write permissions
**Fix:** Verify repository Settings → Actions → General → "Workflow permissions" is set to "Read and write permissions"

### Two rapid [release] commits created duplicate tags
**Cause:** Concurrency group didn't queue in time (network lag)
**Fix:** Delete the duplicate tag:
```bash
git tag -d v0.1.X
git push origin --delete v0.1.X
```
Then manually set the correct version for the next release.

### Release created with wrong version
**Cause:** Old version tag left unpushed locally
**Fix:** Verify tags are in sync:
```bash
git tag -l
git ls-remote --tags origin
```

### [release] in PR title didn't trigger release
**Cause:** Only the **merge commit message** is checked, not PR title
**Fix:** Ensure `[release]` is in the commit message when merging

## Version Management

To manually change the base version for future releases:

### Jump to v0.2.0 (or any version)
```bash
git tag v0.2.0
git push origin v0.2.0
```
The next `[release]` commit will detect `v0.2.0` and create `v0.2.1`.

### Check current version
```bash
git describe --tags --abbrev=0
```
