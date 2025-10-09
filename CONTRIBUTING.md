# Contributing to Forkspacer CLI

Thank you for your interest in contributing to the Forkspacer CLI! We welcome contributions from the community and are excited to see what you'll bring to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Submitting Changes](#submitting-changes)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Community](#community)

---

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors. We expect everyone to:

- Be respectful and considerate
- Welcome newcomers and help them get started
- Accept constructive criticism gracefully
- Focus on what is best for the community
- Show empathy towards other community members

---

## Getting Started

### Prerequisites

- **Go 1.24 or later** - [Install Go](https://golang.org/doc/install)
- **Git** - [Install Git](https://git-scm.com/downloads)
- **kubectl** - [Install kubectl](https://kubernetes.io/docs/tasks/tools/)
- **Kubernetes cluster** - For testing (kind, minikube, or any K8s cluster)
- **Forkspacer Operator** - [Install the operator](https://github.com/forkspacer/forkspacer)

### Fork and Clone

1. **Fork the repository** on GitHub by clicking the "Fork" button

2. **Clone your fork locally:**

```bash
git clone https://github.com/YOUR_USERNAME/cli.git
cd cli
```

3. **Add upstream remote:**

```bash
git remote add upstream https://github.com/forkspacer/cli.git
```

4. **Install dependencies:**

```bash
go mod download
```

5. **Build the CLI:**

```bash
go build -o forkspacer .
```

6. **Verify the build:**

```bash
./forkspacer version
```

---

## Development Workflow

### 1. Create a Feature Branch

Always create a new branch for your work. Never commit directly to `main`.

```bash
git checkout -b feature/my-new-feature
```

**Branch naming conventions:**
- `feature/` - New features (e.g., `feature/module-commands`)
- `fix/` - Bug fixes (e.g., `fix/validation-error`)
- `docs/` - Documentation updates (e.g., `docs/update-readme`)
- `chore/` - Maintenance tasks (e.g., `chore/update-dependencies`)
- `refactor/` - Code refactoring (e.g., `refactor/printer-package`)

### 2. Make Your Changes

Keep your changes focused and atomic. Each PR should address a single concern.

```bash
# Edit files
vim cmd/workspace/create.go

# Run the CLI locally
./forkspacer workspace list

# Test with verbose mode
./forkspacer workspace create test-env --verbose
```

### 3. Follow Coding Standards

- **Format your code:**
  ```bash
  go fmt ./...
  ```

- **Run the linter:**
  ```bash
  go vet ./...
  ```

- **Run static analysis:**
  ```bash
  go install honnef.co/go/tools/cmd/staticcheck@latest
  staticcheck ./...
  ```

### 4. Write Tests

Add tests for new functionality:

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out
```

### 5. Commit Your Changes

Write clear, descriptive commit messages:

```bash
git add .
git commit -m "feat: add module list command

- Implement module list command with table output
- Add filtering by workspace
- Add tests for list functionality

Closes #123"
```

**Commit message format:**
- Use present tense ("add feature" not "added feature")
- Start with a type: `feat:`, `fix:`, `docs:`, `chore:`, `refactor:`, `test:`
- Reference issues with `Closes #123` or `Fixes #123`
- Keep first line under 72 characters
- Add detailed description if needed

### 6. Keep Your Branch Updated

```bash
# Fetch latest changes from upstream
git fetch upstream

# Rebase your branch onto upstream main
git rebase upstream/main

# Force push if needed (only for your feature branch!)
git push origin feature/my-new-feature --force
```

---

## Submitting Changes

### Creating a Pull Request

1. **Push your branch to your fork:**

```bash
git push origin feature/my-new-feature
```

2. **Open a Pull Request** on GitHub from your fork to `forkspacer/cli:main`

3. **Fill out the PR template** with:
   - Clear description of changes
   - Link to related issues
   - Screenshots/examples if applicable
   - Checklist of completed items

4. **Wait for CI checks** to pass:
   - ✅ Linting
   - ✅ Tests
   - ✅ Build for multiple platforms

5. **Respond to review feedback** promptly and professionally

### PR Review Process

- **All PRs require at least one approval** from a maintainer
- **CI must pass** before merging
- **Squash and merge** is preferred for a clean history
- **Reviews may take a few days** - please be patient

---

## Coding Standards

### Go Style Guide

We follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

**Key principles:**

- Use `gofmt` for formatting
- Use meaningful variable and function names
- Keep functions small and focused
- Comment exported functions and packages
- Handle errors explicitly (don't ignore them)
- Use cobra patterns for CLI commands

### Project Structure

```
cli/
├── cmd/                    # Command implementations
│   ├── root.go            # Root command setup
│   └── workspace/         # Workspace subcommands
├── pkg/                   # Reusable packages
│   ├── k8s/              # Kubernetes client
│   ├── printer/          # Output formatting
│   ├── styles/           # Terminal styling
│   └── validation/       # Input validation
└── main.go               # Entry point
```

**When adding new commands:**

1. Create file in appropriate `cmd/` subdirectory
2. Use cobra.Command struct
3. Add validation in RunE function
4. Use printer package for output
5. Handle errors gracefully

**Example command structure:**

```go
var myCmd = &cobra.Command{
    Use:   "mycommand [name]",
    Short: "Short description",
    Long:  "Detailed description...",
    Args:  cobra.ExactArgs(1),
    RunE:  runMyCommand,
}

func runMyCommand(cmd *cobra.Command, args []string) error {
    // Validation
    // Business logic
    // Output formatting
    return nil
}
```

### Error Handling

- Return errors, don't panic
- Wrap errors with context: `fmt.Errorf("failed to create workspace: %w", err)`
- Use helpful error messages
- Show suggested fixes when possible

**Good error message:**
```go
return fmt.Errorf("workspace name %q is invalid: must be DNS-1123 compliant (lowercase letters, numbers, hyphens)\nExamples: dev-env, staging, prod-2", name)
```

### Output Formatting

- Use `pkg/styles` for consistent colors and formatting
- Use `pkg/printer` for tables and spinners
- Always show progress for long operations
- Provide helpful next steps after success

---

## Testing

### Unit Tests

Write unit tests for:
- Validation logic
- Data transformation
- Helper functions

```go
func TestValidateDNS(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid", "dev-env", false},
        {"invalid uppercase", "Dev-Env", true},
        {"invalid chars", "dev@env", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateDNS(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Tests

Test against a real Kubernetes cluster when possible:

```bash
# Start a kind cluster
kind create cluster --name forkspacer-test

# Install operator
kubectl apply -f https://github.com/forkspacer/forkspacer/releases/latest/download/install.yaml

# Run CLI
./forkspacer workspace create test-env
./forkspacer workspace list
./forkspacer workspace delete test-env --force
```

---

## Documentation

### Code Documentation

- Document all exported functions, types, and packages
- Use GoDoc format
- Include examples where helpful

```go
// ValidateDNS1123Subdomain validates that a name follows DNS-1123 subdomain rules.
// It must consist of lowercase alphanumeric characters, '-' or '.', and must
// start and end with an alphanumeric character.
//
// Example:
//   err := ValidateDNS1123Subdomain("my-workspace") // valid
//   err := ValidateDNS1123Subdomain("My-Workspace") // invalid - uppercase
func ValidateDNS1123Subdomain(name string) error {
    // ...
}
```

### README Updates

Update the README.md when you:
- Add new commands
- Change existing behavior
- Add new features
- Update installation steps

### Examples

Add examples to the README showing:
- Real command usage
- Expected output
- Common use cases

---

## Community

### Getting Help

- 💬 **GitHub Discussions** - Ask questions and share ideas
- 🐛 **GitHub Issues** - Report bugs or request features
- 📖 **Documentation** - Check the README and docs

### Communication

- Be kind and respectful
- Ask questions - there are no stupid questions
- Help others when you can
- Share your use cases and feedback

### Recognition

Contributors are recognized in:
- Git commit history
- Release notes
- README acknowledgments

---

## Release Process

Maintainers will handle releases, but contributors should know:

1. Version follows [Semantic Versioning](https://semver.org/)
2. Releases are triggered by pushing git tags (`v*`)
3. CI automatically builds binaries for all platforms
4. Release notes are generated from PR titles
5. Binaries are attached to GitHub releases

---

## Questions?

If you have questions about contributing:

1. Check this guide and the README
2. Search existing issues and discussions
3. Open a new discussion or issue
4. Reach out to maintainers

**Thank you for contributing to Forkspacer CLI! 🎉**
