# Contributing to Forge

Thank you for your interest in contributing to Forge! We welcome contributions that align with our core philosophy of **Safety**, **Determinism**, and **Simplicity**.

## 🚀 Getting Started

### Prerequisites
-   **Windows 10/11** (Forge is Windows-only)
-   **Go 1.19+**
-   **Git**

### Setup
1.  Fork the repository on GitHub.
2.  Clone your fork locally:
    ```powershell
    git clone https://github.com/YOUR_USERNAME/forge.git
    cd forge
    ```
3.  Build the project to ensure everything is working:
    ```powershell
    go build -o forge.exe
    .\forge.exe --version
    ```

---

## 🛠️ Development Workflow

1.  **Create a Branch:** Always work on a new branch for your feature or fix.
    ```powershell
    git checkout -b feature/my-awesome-feature
    ```

2.  **Make Changes:** Write your code.
    -   Keep changes focused and minimal.
    -   Follow existing code style (Go standard formatting).

3.  **Run Tests:**
    ```powershell
    go test ./...
    ```
    Ensure all tests pass before submitting.

---

## 🧪 Testing Requirements

Forge relies heavily on integration tests to ensure filesystem safety.

-   **Unit Tests:** Write unit tests for individual functions where possible.
-   **Safety Tests:** If you modify `fileops` or `workspace` code, you **must** verify that operations do not leak outside the temporary directory.
-   **Windows Paths:** Always use `filepath.Join` or `filepath.Clean` to handle Windows paths correctly.

---

## 📏 Code Standards

-   **Formatting:** Run `go fmt ./...` before committing.
-   **Linting:** We recommend `golangci-lint`.
-   **Error Handling:**
    -   Return errors, don't panic.
    -   Use descriptive error messages.
    -   Wrap errors when helpful: `fmt.Errorf("failed to create file: %w", err)`.
-   **Comments:** Document exported functions and complex logic.

---

## 📝 Pull Request Process

1.  **Descriptive Title:** Use a clear title (e.g., "fix: handle path spaces in command arguments").
2.  **Description:** Explain *what* you changed and *why*.
3.  **Link Issues:** If this PR fixes an issue, link it (e.g., "Fixes #123").
4.  **Checklist:**
    -   [ ] Tests passed locally
    -   [ ] Code formatted
    -   [ ] Documentation updated (if applicable)

---

## 🚫 What We Don't Accept

Forge has a strict scope for V1. Please do **not** submit PRs for:
-   Cross-platform support (Mac/Linux).
-   Interactive prompts (unless behind a flag).
-   Shell execution support.
-   Dependencies on external libraries (keep `go.mod` minimal).

See [ARCHITECTURE.md](doc/ARCHITECTURE.md) and `plan.md` for more on our design philosophy.
