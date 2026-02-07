# MAINTAINER.md

## Project Overview

- Go CLI tool
- Go version: 1.25
- All development on feature branches
- `main` branch is always releasable

## Branching and Commit Rules

- Feature branches: small, incremental commits allowed
- Merge to `main` via squash merge only
- `main` must always pass build and tests

## Release Process

- Releases are fully automated via GitHub Actions
- No manual tag creation
- Tags follow `v0.1.x` (PATCH only)
- Release triggers on push to `main` with `[release]` in the final merge commit message
- Workflow:
  - Finds latest tag
  - Builds from latest commit on `main`
  - Cross-compiles Windows `.exe` on Ubuntu
  - Creates GitHub release
  - Uses final merge commit message as release notes

## Versioning

- Only PATCH version bumps (e.g., `v0.1.3` to `v0.1.4`)
- No manual tag management

## Release Notes

- Release notes are taken from the final merge commit message

## Intentional Releases

- Only push to `main` with `[release]` in the merge commit message when a release is intended

## Maintenance

- Keep `main` clean and always releasable
- Use feature branches for all changes
- Review and test before merging to `main`
- Monitor GitHub Actions for build and release status
