# Code-Factory 🏭

**Spec-Driven Software Development with AI**

---

## What is Code-Factory?

Code-Factory is a revolutionary development tool that bridges the gap between specifications and implementation. It transforms natural language requirements into production-ready code through AI-powered workflows, while keeping humans firmly in control.

### Key Features

✨ **Beautiful Terminal UI** - Canvas-style interface built with Charm.sh/Bubble Tea  
🤖 **AI-Powered** - Works with local Ollama or cloud LLMs (OpenAI, Claude)  
📝 **Spec-Driven** - Everything starts and ends with clear specifications  
🔄 **Git-Native** - All artifacts are plain text files in git  
🚀 **Zero Dependencies** - Single binary, works anywhere  
🔒 **Privacy-First** - Run completely offline with local LLM  
🌐 **GitHub Integration** - Seamless OAuth + PR workflows  

---

## Quick Start

### Installation

**macOS / Linux:**
```bash
curl -sSL https://factory.dev/install.sh | bash
```

**Homebrew (macOS):**
```bash
brew install code-factory
```

**From Source:**
```bash
git clone https://github.com/ssdajoker/Code-Factory.git
cd Code-Factory
go build -o factory ./cmd/factory
sudo mv factory /usr/local/bin/
```

### First Run

```bash
# Initialize Code-Factory
factory init

# Follow the interactive setup:
# 1. Configure LLM (Ollama auto-detected)
# 2. Optional: Connect to GitHub
# 3. Create project structure
```

---

## The Four Modes

### 🎯 INTAKE - Capture Requirements

Transform ideas into structured specifications:

```bash
factory intake
```

**What it does:**
- Interactive requirement gathering
- AI-powered clarifying questions
- Generate comprehensive specifications
- Human review and editing
- Git commit and sync

**Example:**
```bash
$ factory intake --prompt "Build a REST API for user authentication"
```

---

### 🔍 REVIEW - Analyze Code

Compare code against specifications:

```bash
factory review src/auth
```

**What it does:**
- Scan codebase and match specs
- Check compliance, security, performance
- Generate detailed reports with actionable recommendations
- Track metrics and test coverage

**Example:**
```bash
$ factory review --spec contracts/specs/user-auth.md src/auth
```

---

### ⚙️ CHANGE_ORDER - Implement Changes

Generate code from specifications:

```bash
factory change-order
```

**What it does:**
- Analyze change request and impact
- Plan implementation steps
- Generate code with AI
- Show diffs for human review
- Create feature branch + PR

**Example:**
```bash
$ factory change-order --description "Add password reset functionality"
```

---

### 🚑 RESCUE - Debug & Fix

Diagnose and solve problems:

```bash
factory rescue
```

**What it does:**
- Gather context (code, tests, logs, specs)
- AI-powered root cause analysis
- Generate solutions
- Apply fixes and verify

**Example:**
```bash
$ factory rescue --problem "Tests failing after adding password reset"
```

---

## Project Structure

Code-Factory organizes projects with simple, git-friendly structure:

```
my-project/
├── contracts/              # All specifications
│   ├── specs/              # Feature specs
│   │   ├── user-auth.md
│   │   └── api-gateway.md
│   ├── architecture/       # System design docs
│   │   └── system-design.md
│   └── decisions/          # Architecture Decision Records
│       └── adr-001-use-jwt.md
├── reports/                # Generated analysis reports
│   └── review-2026-01-07.md
└── .factory/               # Factory metadata (gitignored)
    ├── cache/
    └── config.yaml
```

---

## Configuration

Code-Factory stores configuration in `~/.config/factory/config.yaml`:

```yaml
version: "1.0.0"

llm:
  provider: ollama          # or openai, anthropic, google
  endpoint: http://localhost:11434
  model: llama3.2:latest
  
github:
  enabled: true
  username: yourusername
  
project:
  name: my-awesome-app
  path: /path/to/project
```

### Reconfigure Anytime

```bash
factory config              # View current config
factory config llm          # Change LLM provider
factory config github       # GitHub integration
factory config reset        # Reset to defaults
```

---

## LLM Providers

### Ollama (Recommended)

**Free, private, local AI:**

```bash
# Install Ollama
curl -sSL https://ollama.ai/install.sh | bash

# Pull a model
ollama pull llama3.2

# Code-Factory auto-detects Ollama
factory init
```

**Recommended models:**
- `llama3.2:latest` - Best overall
- `codellama:latest` - Optimized for code
- `deepseek-coder:latest` - Strong coding performance

### BYOK (Bring Your Own Key)

**Cloud LLM options:**

- **OpenAI** (GPT-4, GPT-3.5)
- **Anthropic** (Claude 3)
- **Google** (Gemini)
- **Azure OpenAI**
- **Custom endpoint** (any OpenAI-compatible API)

Configure during `factory init` or via:
```bash
factory config llm
```

---

## GitHub Integration

### Features

- **OAuth Device Flow** - Secure authentication
- **GitHub App** - Fine-grained permissions
- **PR Creation** - Auto-generate from CHANGE_ORDER
- **Spec Sync** - Push specs to repository
- **Issue Linking** - Connect specs to issues

### Setup

```bash
factory config github

# Follow prompts:
# 1. Open https://github.com/login/device
# 2. Enter provided code
# 3. Install Code-Factory GitHub App
# 4. Grant repository access
```

---

## Examples

### Full Workflow

```bash
# 1. Create specification
factory intake --prompt "User authentication with JWT tokens"

# 2. Review existing auth code
factory review src/auth

# 3. Implement missing features from spec
factory change-order --spec contracts/specs/user-auth.md

# 4. Fix any issues
factory rescue --problem "Login tests failing"
```

### Advanced Usage

```bash
# INTAKE: From file
factory intake --file requirements.txt

# REVIEW: Focus on security
factory review --focus security src/

# REVIEW: Export as JSON
factory review --format json --output report.json src/

# CHANGE_ORDER: Dry run (plan only)
factory change-order --dry-run --description "Add rate limiting"

# CHANGE_ORDER: No PR creation
factory change-order --no-pr --description "Update dependencies"

# RESCUE: From error log
factory rescue --error "$(cat error.log)"
```

---

## Architecture

### Design Principles

1. **Simplicity Over Complexity** - Single binary, sensible defaults
2. **Git-Native Storage** - Plain text files, no databases
3. **Terminal-First** - Beautiful TUI, keyboard-driven
4. **Privacy by Default** - Local LLM support, secure secrets
5. **Human-in-the-Loop** - AI assists, humans validate

### Technology Stack

- **Language:** Go 1.21+
- **TUI:** Charm.sh (Bubble Tea, Lipgloss, Bubbles)
- **LLM:** Ollama, OpenAI, Anthropic, Google
- **Git:** go-git
- **GitHub:** go-github + OAuth2
- **Secrets:** OS keyring (Keychain, Secret Service, Credential Manager)

---

## Development

### Building from Source

```bash
# Clone repository
git clone https://github.com/ssdajoker/Code-Factory.git
cd Code-Factory

# Install dependencies
go mod download

# Build
go build -o factory ./cmd/factory

# Run tests
go test ./...

# Run with local changes
./factory --help
```

### Project Structure

```
Code-Factory/
├── cmd/factory/           # Main entry point
├── internal/
│   ├── tui/               # Terminal UI
│   ├── llm/               # LLM integration
│   ├── github/            # GitHub integration
│   ├── modes/             # INTAKE, REVIEW, CHANGE_ORDER, RESCUE
│   └── core/              # Core utilities
├── contracts/             # Project specifications
├── docs/                  # Documentation
└── README.md
```

### Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

## Roadmap

### v1.0 (Current)
- ✅ Core four modes (INTAKE, REVIEW, CHANGE_ORDER, RESCUE)
- ✅ Ollama + BYOK support
- ✅ GitHub OAuth integration
- ✅ Beautiful TUI
- ✅ Git-native storage

### v1.1 (Next)
- [ ] Web UI mirror (optional)
- [ ] Plugin system
- [ ] Multi-project management
- [ ] Team collaboration features

### v2.0 (Future)
- [ ] IDE extensions (VSCode, JetBrains)
- [ ] CI/CD integration
- [ ] Real-time co-editing
- [ ] Knowledge graph

---

## FAQ

### Q: Does Code-Factory write code without human review?
**A:** No. All code changes require explicit human approval. Code-Factory shows diffs, explains changes, and waits for your confirmation.

### Q: Can I use Code-Factory offline?
**A:** Yes! With Ollama, Code-Factory works completely offline. Only GitHub integration requires internet.

### Q: How much does it cost?
**A:** Code-Factory is free and open source (MIT license). Ollama is free. Cloud LLMs have their own pricing.

### Q: Is my code sent to external servers?
**A:** Only if you choose cloud LLMs (OpenAI, Claude). With Ollama, everything stays local.

### Q: What languages does it support?
**A:** Code-Factory is language-agnostic. It works with any programming language or framework.

### Q: Can I use it with existing projects?
**A:** Absolutely! Run `factory init` in any project directory.

---

## Support

- **Documentation:** [https://docs.code-factory.dev](https://docs.code-factory.dev)
- **Issues:** [GitHub Issues](https://github.com/ssdajoker/Code-Factory/issues)
- **Discussions:** [GitHub Discussions](https://github.com/ssdajoker/Code-Factory/discussions)
- **Discord:** [Join our community](https://discord.gg/code-factory)

---

## License

MIT License - see [LICENSE](LICENSE) for details.

---

## Acknowledgments

Code-Factory is built on the shoulders of giants:

- [Charm.sh](https://charm.sh/) - Beautiful TUI framework
- [Ollama](https://ollama.ai/) - Local LLM runtime
- [Go](https://go.dev/) - Fast, reliable language
- All the open source contributors who make this possible

---

**Built with ❤️ by developers, for developers**

Transform your development workflow. Start with specs, end with code.

```bash
factory init
```
