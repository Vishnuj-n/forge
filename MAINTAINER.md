# Maintainer Guidelines

## Versioning Rules

Versions are injected automatically from Git tags.
Do not manually edit `root.go` version. 

Keep `var Version = "development"` in `cmd/forge/root.go` and let the CI/CD pipeline inject the correct version at build time.
