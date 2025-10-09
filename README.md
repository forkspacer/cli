# Forkspacer CLI

A beautiful command-line interface for managing the Forkspacer Kubernetes operator. Create, manage, and hibernate ephemeral development environments with style.

[![License](https://img.shields.io/github/license/forkspacer/cli)](LICENSE)
[![CI](https://github.com/forkspacer/cli/workflows/CI/badge.svg)](https://github.com/forkspacer/cli/actions)
[![Release](https://img.shields.io/github/v/release/forkspacer/cli)](https://github.com/forkspacer/cli/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/forkspacer/cli)](go.mod)

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
  - [Using Install Script](#using-install-script)
  - [Manual Installation](#manual-installation)
  - [From Source](#from-source)
- [Quick Start](#quick-start)
- [Commands](#commands)
  - [Workspace Management](#workspace-management)
- [Examples](#examples)
- [Shell Completion](#shell-completion)
- [Configuration](#configuration)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

The Forkspacer CLI is part of the Forkspacer ecosystem, which consists of:

1. **Forkspacer Operator** - Core Kubernetes operator managing custom resources
2. **API Server** - Backend API service (optional)
3. **Operator UI** - Web-based dashboard (optional)
4. **CLI** (this project) - Command-line interface

The CLI provides direct Kubernetes integration, beautiful terminal output, and fast client-side validation.

```
┌─────────────────────────────────────────────┐
│           Kubernetes Cluster                │
│                                             │
│  ┌──────────────┐      ┌──────────────┐     │
│  │  Forkspacer  │◄────►│  Workspaces  │     │
│  │   Operator   │      │   & Modules  │     │
│  │  (CRD Watch) │      │   (CRDs)     │     │
│  └──────────────┘      └──────────────┘     │
│         ▲                                   │
│         │                                   │
│         │ kubectl-style access              │
└─────────┼───────────────────────────────────┘
          │
    ┌─────▼──────┐
    │ Forkspacer │◄─── You
    │    CLI     │
    └────────────┘
```

---

## Features

- 🎨 **Beautiful Output** - Styled terminal output with colors, spinners, and progress indicators
- ⚡ **Fast Validation** - Client-side validation for instant feedback
- 🚀 **Easy to Use** - Intuitive commands that feel natural
- 🔧 **Direct K8s Access** - No API server required, uses your kubeconfig
- 🌍 **Cross-Platform** - Works on macOS, Linux, and Windows
- 📝 **Shell Completion** - Tab completion for bash, zsh, fish, and powershell
- 🔄 **Workspace Lifecycle** - Create, hibernate, wake, and manage workspaces
- ⏰ **Auto-Hibernation** - Schedule workspaces to sleep and wake automatically

---

## Installation

### Prerequisites

- Kubernetes cluster (1.24+)
- [Forkspacer Operator](https://github.com/forkspacer/forkspacer) installed in your cluster
- `kubectl` configured with cluster access

### Using Install Script

The quickest way to install (macOS & Linux):

```bash
curl -sSL https://raw.githubusercontent.com/forkspacer/cli/main/scripts/install.sh | bash
```

### Manual Installation

#### macOS (Apple Silicon)
```bash
curl -sSL https://github.com/forkspacer/cli/releases/latest/download/forkspacer-darwin-arm64.tar.gz | tar xz
sudo mv forkspacer /usr/local/bin/
```

#### macOS (Intel)
```bash
curl -sSL https://github.com/forkspacer/cli/releases/latest/download/forkspacer-darwin-amd64.tar.gz | tar xz
sudo mv forkspacer /usr/local/bin/
```

#### Linux (amd64)
```bash
curl -sSL https://github.com/forkspacer/cli/releases/latest/download/forkspacer-linux-amd64.tar.gz | tar xz
sudo mv forkspacer /usr/local/bin/
```

#### Linux (arm64)
```bash
curl -sSL https://github.com/forkspacer/cli/releases/latest/download/forkspacer-linux-arm64.tar.gz | tar xz
sudo mv forkspacer /usr/local/bin/
```

#### Windows
Download `forkspacer-windows-amd64.zip` from the [releases page](https://github.com/forkspacer/cli/releases), extract, and add to your PATH.

### From Source

```bash
git clone https://github.com/forkspacer/cli.git
cd cli
go build -o forkspacer .
sudo mv forkspacer /usr/local/bin/
```

### Verify Installation

```bash
forkspacer version
```

---

## Quick Start

```bash
# List existing workspaces
forkspacer workspace list

# Create a simple workspace
forkspacer workspace create dev-env

# Create with auto-hibernation (sleep at 6 PM, wake at 8 AM)
forkspacer workspace create dev-env \
  --hibernation-schedule "0 18 * * *" \
  --wake-schedule "0 8 * * *"

# Get detailed workspace information
forkspacer workspace get dev-env

# Hibernate a workspace manually
forkspacer workspace hibernate dev-env

# Wake up a hibernated workspace
forkspacer workspace wake dev-env

# Delete a workspace
forkspacer workspace delete dev-env --force

# Use short alias
forkspacer ws list
```

---

## Commands

### Workspace Management

```bash
# Create
forkspacer workspace create <name> [flags]
  --hibernation-schedule string   Cron schedule for auto-hibernation
  --wake-schedule string          Cron schedule for auto-wake
  --connection string             Connection type (default "in-cluster")
  --from string                   Fork from existing workspace
  --wait                          Wait for workspace to be ready

# List
forkspacer workspace list [flags]
  --all-namespaces, -A            List workspaces across all namespaces

# Get
forkspacer workspace get <name>

# Delete
forkspacer workspace delete <name> [flags]
  --force                         Skip confirmation prompt

# Hibernate
forkspacer workspace hibernate <name>

# Wake
forkspacer workspace wake <name>

# Alias
forkspacer ws <command>  # Short form for workspace commands
```

### Global Flags

```bash
-n, --namespace string   Kubernetes namespace (default "default")
-o, --output string      Output format: table|json|yaml (default "table")
-v, --verbose            Enable verbose output
-h, --help               Help for any command
```

### Utility Commands

```bash
# Show version
forkspacer version

# Generate shell completion
forkspacer completion [bash|zsh|fish|powershell]
```

---

## Examples

### Basic Workspace Creation

```bash
$ forkspacer workspace create my-dev-env

✨ Creating workspace my-dev-env

✓ Workspace name is valid
✓ Connected to cluster (context: minikube)
✓ Forkspacer operator is installed
✓ Workspace name is available
✓ Workspace resource created

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Name:  my-dev-env
Namespace:  default
Type:  kubernetes
Hibernation:  disabled

Next steps:
  →  forkspacer workspace get my-dev-env
  →  forkspacer module deploy redis --workspace my-dev-env

Documentation: https://forkspacer.com/docs/workspaces
```

### Auto-Hibernation Schedule

```bash
$ forkspacer workspace create staging \
    --hibernation-schedule "0 18 * * 1-5" \
    --wake-schedule "0 8 * * 1-5"

# Workspace will automatically:
# - Hibernate at 6 PM on weekdays
# - Wake up at 8 AM on weekdays
```

### Listing Workspaces

```bash
$ forkspacer workspace list

┌──────────┬───────────┬────────────┬───────┬────────────┬─────────────────────┐
│   NAME   │ NAMESPACE │   PHASE    │ READY │ HIBERNATED │    LAST ACTIVITY    │
├──────────┼───────────┼────────────┼───────┼────────────┼─────────────────────┤
│ dev-env  │ default   │ ready      │ true  │ false      │ 2025-10-09 18:40:59 │
│ staging  │ default   │ hibernated │ true  │ true       │ 2025-10-09 12:30:00 │
└──────────┴───────────┴────────────┴───────┴────────────┴─────────────────────┘

Total: 2 workspace(s)
```

### Working with Multiple Namespaces

```bash
# Create workspace in specific namespace
forkspacer workspace create prod-env -n production

# List all workspaces across namespaces
forkspacer workspace list --all-namespaces
```

### Cron Schedule Examples

Common hibernation schedules:

```bash
# Every day at 6 PM
--hibernation-schedule "0 18 * * *"

# Weekdays at 6 PM
--hibernation-schedule "0 18 * * 1-5"

# Every Monday at 9 AM
--hibernation-schedule "0 9 * * 1"

# Every 15 minutes
--hibernation-schedule "*/15 * * * *"

# Every Sunday at midnight
--hibernation-schedule "0 0 * * 0"
```

---

## Shell Completion

Enable tab completion for faster command entry:

### Bash

```bash
# Linux
forkspacer completion bash > /etc/bash_completion.d/forkspacer

# macOS
forkspacer completion bash > /usr/local/etc/bash_completion.d/forkspacer
```

### Zsh

```bash
# macOS (Homebrew)
forkspacer completion zsh > $(brew --prefix)/share/zsh/site-functions/_forkspacer

# Linux
forkspacer completion zsh > /usr/local/share/zsh/site-functions/_forkspacer

# Then restart your shell
exec zsh
```

### Fish

```bash
forkspacer completion fish > ~/.config/fish/completions/forkspacer.fish
```

### PowerShell

```powershell
forkspacer completion powershell | Out-String | Invoke-Expression
```

After enabling completion, you can:
```bash
forkspacer work<TAB>           # → workspace
forkspacer workspace cr<TAB>   # → create
forkspacer ws li<TAB>          # → list
```

---

## Configuration

The CLI uses your `~/.kube/config` for Kubernetes authentication. Configure your context:

```bash
# View current context
kubectl config current-context

# Switch context
kubectl config use-context my-cluster

# List available contexts
kubectl config get-contexts
```

---

## Development

### Prerequisites

- Go 1.24 or later
- Kubernetes cluster for testing
- [Forkspacer Operator](https://github.com/forkspacer/forkspacer) installed

### Local Setup

```bash
# Clone the repository
git clone https://github.com/forkspacer/cli.git
cd cli

# Install dependencies
go mod download

# Build
go build -o forkspacer .

# Run
./forkspacer version
```

### Project Structure

```
cli/
├── cmd/                    # Command implementations
│   ├── root.go            # Root command & global flags
│   ├── version.go         # Version command
│   └── workspace/         # Workspace commands
│       ├── create.go
│       ├── list.go
│       ├── get.go
│       ├── delete.go
│       ├── hibernate.go
│       └── wake.go
├── pkg/                   # Shared packages
│   ├── k8s/              # Kubernetes client wrapper
│   ├── printer/          # Output formatting
│   ├── styles/           # Terminal styling
│   └── validation/       # Input validation
├── .github/              # GitHub workflows & templates
├── scripts/              # Install scripts
└── main.go               # Entry point
```

### Testing

```bash
# Run tests
go test -v ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage
go tool cover -html=coverage.out
```

### Building for Multiple Platforms

```bash
# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o forkspacer-darwin-arm64 .

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o forkspacer-darwin-amd64 .

# Linux
GOOS=linux GOARCH=amd64 go build -o forkspacer-linux-amd64 .

# Windows
GOOS=windows GOARCH=amd64 go build -o forkspacer-windows-amd64.exe .
```

---

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on:

- Code of Conduct
- Development workflow
- Submitting pull requests
- Coding standards

### Quick Contribution Guide

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

---

## Roadmap

### v0.1.0 (Current)
- ✅ Workspace CRUD operations
- ✅ Hibernation/Wake commands
- ✅ Beautiful terminal output
- ✅ Shell completion

### v0.2.0 (Planned)
- 🚧 Module management commands
- 🚧 Helm chart deployment
- 🚧 Resource filtering

### v0.3.0 (Future)
- 📋 Helm release discovery
- 📋 Batch import existing releases
- 📋 Interactive init wizard
- 📋 Workspace forking

### v1.0.0 (Future)
- 📋 Full feature parity with operator capabilities
- 📋 Advanced output formats (JSON, YAML)
- 📋 Plugin system
- 📋 Enhanced error handling

---

## Troubleshooting

### CLI can't connect to cluster

```bash
# Check kubectl connectivity
kubectl cluster-info

# Check current context
kubectl config current-context

# Verify Forkspacer operator is installed
kubectl get crd workspaces.batch.forkspacer.com
```

### Permission denied errors

```bash
# Check RBAC permissions
kubectl auth can-i create workspaces.batch.forkspacer.com

# View your permissions
kubectl auth can-i --list
```

### Workspace creation fails

```bash
# Use verbose mode to see detailed errors
forkspacer workspace create my-env --verbose

# Check operator logs
kubectl logs -n operator-system deployment/operator-controller-manager
```

### Getting help

- 📖 [Documentation](https://github.com/forkspacer/cli#documentation)
- 🐛 [Report a Bug](https://github.com/forkspacer/cli/issues/new?template=bug_report.yml)
- 💡 [Request a Feature](https://github.com/forkspacer/cli/issues/new?template=feature_request.yml)
- ❓ [Ask a Question](https://github.com/forkspacer/cli/issues/new?template=question.yml)

---

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss) for beautiful output
- Powered by [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime)

---

**Made with ❤️ by the Forkspacer Team**

⭐ Star us on GitHub if you find this project useful!
