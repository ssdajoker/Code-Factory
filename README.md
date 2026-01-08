# Factory

[![CI](https://github.com/ssdajoker/Code-Factory/actions/workflows/ci.yml/badge.svg)](https://github.com/ssdajoker/Code-Factory/actions/workflows/ci.yml)
[![Go Report Card](https://i.ytimg.com/vi/P8xMWWfbvR4/sddefault.jpg)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)

**A Spec-Driven Software Factory CLI**

Factory helps you capture project vision, generate specifications, and maintain alignment between code and contracts through an intelligent, LLM-powered workflow.

## Features

- **INTAKE Mode** - Guided interview to capture project vision and generate comprehensive specifications
- **REVIEW Mode** - Analyze code against specifications to detect drift and inconsistencies
- **RESCUE Mode** - Reverse-engineer existing codebases to generate specifications
- **CHANGE_ORDER Mode** - Track and manage specification changes over time
- **LLM Integration** - Support for Ollama (local), OpenAI, and Anthropic
- **GitHub Integration** - OAuth authentication, repository management, PR creation
- **Beautiful TUI** - Terminal UI built with Bubble Tea
- **Secure Storage** - OS keyring or encrypted file storage for secrets

## Installation

### Using Go

```bash
go install github.com/ssdajoker/Code-Factory/cmd/factory@latest
```

### From Source

```bash
git clone https://github.com/ssdajoker/Code-Factory.git
cd Code-Factory
make build
sudo make install-local
```

### Using Install Script

```bash
curl -sSL https://raw.githubusercontent.com/ssdajoker/Code-Factory/main/scripts/install.sh | bash
```

### Pre-built Binaries

Download from [Releases](https://github.com/ssdajoker/Code-Factory/releases).

## Quick Start

```bash
# Initialize Factory in your project
factory init

# Start the interactive TUI
factory

# Or use specific modes
factory intake      # Capture project vision
factory review      # Review code against specs
factory rescue      # Generate specs from existing code
factory change-order # Track spec changes
```

## LLM Setup

Factory auto-detects available LLM providers:

### Ollama (Recommended - Local & Free)

```bash
# Install Ollama
curl -fsSL https://ollama.com/install.sh | sh

# Pull a model
ollama pull llama3.2

# Factory will auto-detect Ollama
factory llm status
```

### OpenAI / Anthropic

```bash
# Set API key
factory llm set openai
# Follow prompts to enter API key
```

## GitHub Integration

```bash
# Authenticate with GitHub
factory github auth

# Check status
factory github status

# List repositories
factory github repos
```

## Configuration

Configuration is stored in `~/.factory/config.toml`:

```toml
[llm]
provider = "ollama"
model = "llama3.2"
base_url = "http://localhost:11434"

[github]
token_storage = "keyring"

[ui]
theme = "dark"
animations = true

[paths]
specs_dir = ".factory/specs"
reports_dir = ".factory/reports"
```

## Documentation

- [User Guide](docs/USER_GUIDE.md)
- [Architecture](docs/ARCHITECTURE.md)
- [API Reference](docs/API.md)
- [Contributing](docs/CONTRIBUTING.md)
- [Changelog](docs/CHANGELOG.md)

## Project Structure

```
Code-Factory/
├── cmd/factory/        # CLI entry point
├── internal/
│   ├── config/         # Configuration management
│   ├── github/         # GitHub OAuth & API
│   ├── llm/            # LLM providers (Ollama, OpenAI, Anthropic)
│   ├── modes/          # INTAKE, REVIEW, RESCUE, CHANGE_ORDER
│   ├── store/          # Secret storage (keyring, encrypted files)
│   └── tui/            # Terminal UI components
├── contracts/          # Specification documents
├── docs/               # Documentation
└── tests/              # Integration tests
```

## Development

```bash
# Run tests
make test

# Run linters
make lint

# Build
make build

# Run integration tests
make test-integration
```

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please read our [Contributing Guide](docs/CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md).
