# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-01-08

### Added
- Initial release of Factory CLI
- **INTAKE Mode**: Guided interview to capture project vision and generate specifications
- **REVIEW Mode**: Analyze code against specifications for drift detection
- **RESCUE Mode**: Reverse-engineer existing codebases to generate specifications
- **CHANGE_ORDER Mode**: Track and manage specification changes over time
- **LLM Integration**: Support for Ollama, OpenAI, and Anthropic providers
- **Auto-detection**: Automatic detection of available LLM providers
- **GitHub Integration**: OAuth device flow authentication
- **GitHub API**: Repository management, branch creation, PR creation
- **TUI Framework**: Beautiful terminal UI with Bubble Tea
- **Configuration**: TOML-based configuration with sensible defaults
- **Secret Storage**: Secure storage via OS keyring or encrypted files
- **CI/CD**: GitHub Actions for testing, linting, and releases
- **Cross-platform**: Builds for Linux, macOS, and Windows

### Security
- Secrets stored in OS keyring (preferred) or AES-256 encrypted files
- No plaintext credential storage
- Argon2 key derivation for file-based encryption

[Unreleased]: https://github.com/ssdajoker/Code-Factory/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/ssdajoker/Code-Factory/releases/tag/v0.1.0
