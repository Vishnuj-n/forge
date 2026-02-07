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

- Releases are **automatic**
- Triggered by pushing to `main`
- Tags are **NOT created manually**
- Version bump: PATCH only (v0.1.x)
- Tag is created by GitHub Actions
- Release notes come from the final merge commit

To release:
- Use squash merge
- Include `[release]` in commit message

---

## CI/CD Behavior

- Go version: 1.25
- Runner: ubuntu-latest
- Cross-compiles Windows `.exe`
- CGO disabled
- Output: `dist/forge.exe`
- Releases created via GitHub Actions

---

## Commit Conventions

- Small commits allowed in feature branches
- Release commit must be clean and descriptive
- Avoid “merge branch …” commits

---

## Things I Must Remember

- Do NOT manually create tags
- Do NOT create releases from GitHub UI
- Do NOT push directly to `main`
- Keep release notes human-readable

---

## Future Improvements (Ideas)
- Winget automation
- Chocolatey packaging
- Code signing
