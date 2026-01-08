# Contributing to Factory

Thank you for your interest in contributing to Factory! This document provides guidelines and instructions for contributing.

## Getting Started

### Prerequisites

- Go 1.22 or later
- Git
- Make (optional, but recommended)

### Setting Up Development Environment

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/Code-Factory.git
   cd Code-Factory
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/ssdajoker/Code-Factory.git
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```
5. Build:
   ```bash
   make build
   ```

## Development Workflow

### Branching Strategy

- `main` - stable release branch
- `develop` - development branch
- `feature/*` - feature branches
- `fix/*` - bug fix branches

### Making Changes

1. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```
2. Make your changes
3. Run tests:
   ```bash
   make test
   ```
4. Run linters:
   ```bash
   make lint
   ```
5. Commit your changes:
   ```bash
   git commit -m "feat: add your feature description"
   ```

### Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - new feature
- `fix:` - bug fix
- `docs:` - documentation changes
- `test:` - adding or updating tests
- `refactor:` - code refactoring
- `chore:` - maintenance tasks

### Pull Requests

1. Push your branch:
   ```bash
   git push origin feature/your-feature-name
   ```
2. Open a Pull Request against `main`
3. Fill out the PR template
4. Wait for review

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Use meaningful variable and function names
- Add comments for exported functions
- Write table-driven tests

## Testing

- Write unit tests for new functionality
- Maintain test coverage above 70%
- Use mocks for external dependencies
- Run `make test` before submitting PR

## Documentation

- Update README.md for user-facing changes
- Add godoc comments to exported functions
- Update relevant docs in `docs/` directory

## Questions?

Feel free to open an issue for questions or discussions.
