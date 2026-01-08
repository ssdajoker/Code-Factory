# Code-Factory 🏭

**Spec-Driven Software Factory** - Turn specifications into reality

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)]()

---

## 🎯 Vision

Code-Factory is a **single-binary, cross-platform tool** that helps developers maintain perfect alignment between specifications and code. It operates in four distinct modes (INTAKE, REVIEW, CHANGE_ORDER, RESCUE) and provides a beautiful terminal UI for an exceptional developer experience.

**Key Features:**
- 🚀 **Zero Dependencies** - Single Go binary, works anywhere
- 🎨 **Beautiful TUI** - Canvas-style interface using Charm.sh/Bubble Tea
- 🤖 **LLM Integration** - Ollama (local) or BYOK (OpenAI, Claude, etc.)
- 🔗 **GitHub Integration** - Seamless OAuth and GitHub App setup
- 📁 **Git-Native Storage** - Flat files, no database required
- 🌐 **Cross-Platform** - Linux, macOS, Windows (WSL/Git Bash)
- 👥 **Team-Friendly** - One person sets up, others clone and go

---

## 🚀 Quick Start

### Installation

**Recommended (Safer Two-Step):**
```bash
# Download and inspect the script first
curl -sSLO https://raw.githubusercontent.com/ssdajoker/Code-Factory/main/scripts/install.sh
less install.sh  # Review the script
sh install.sh
```

**Quick Install (One-liner):**
```bash
# ⚠️ Only use if you trust the source
curl -sSL https://raw.githubusercontent.com/ssdajoker/Code-Factory/main/scripts/install.sh | sh
```

### Initialize

```bash
cd /path/to/your/project
factory init
```

### Start Using

```bash
factory          # Start TUI
factory intake   # Create specification
factory review   # Review code against specs
```

---

## 📖 What is Spec-Driven Development?

Spec-Driven Development is a methodology where **specifications are the source of truth** for your software. Instead of code drifting away from documentation, Factory ensures they stay in sync.

**The Problem:**
- Documentation gets outdated
- Code doesn't match requirements
- Technical debt accumulates
- New team members struggle to understand the system

**The Solution:**
- Write specs first (or generate from existing code)
- Factory checks code against specs automatically
- Track and manage specification drift
- Keep documentation always up-to-date

---

## 🎭 Four Modes

### 1. INTAKE Mode 📝
**Capture your vision and create specifications**

Describe what you want to build, and Factory's LLM generates a comprehensive, actionable specification document.

```bash
factory intake
```

### 2. REVIEW Mode 🔍
**Check code against specifications**

Factory analyzes your codebase and identifies compliance issues, deviations, and areas for improvement.

```bash
factory review
```

### 3. CHANGE_ORDER Mode 📊
**Track specification drift over time**

Detect code changes that deviate from specs, and decide whether to update the spec or revert the change.

```bash
factory change-order
```

### 4. RESCUE Mode 🆘
**Reverse-engineer existing codebase**

Generate specifications from existing code - perfect for legacy projects without documentation.

```bash
factory rescue
```

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     Factory CLI                         │
│                  (Single Go Binary)                     │
└─────────────────────────────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
        ▼                 ▼                 ▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│  Beautiful  │  │     LLM     │  │   GitHub    │
│     TUI     │  │ Integration │  │ Integration │
│  (Bubble    │  │  (Ollama/   │  │  (OAuth +   │
│    Tea)     │  │   OpenAI)   │  │  App API)   │
└─────────────┘  └─────────────┘  └─────────────┘
        │                 │                 │
        └─────────────────┼─────────────────┘
                          ▼
                ┌─────────────────┐
                │  Git-Native     │
                │    Storage      │
                │  /contracts/    │
                │  /reports/      │
                └─────────────────┘
```

**Key Components:**
- **TUI Layer** - Beautiful terminal interface (Charm.sh/Bubble Tea)
- **Mode Orchestrator** - INTAKE, REVIEW, CHANGE_ORDER, RESCUE
- **LLM Layer** - Ollama (local) or cloud providers (OpenAI, Claude, etc.)
- **GitHub Layer** - OAuth, GitHub App, API integration
- **Storage Layer** - Markdown files in `/contracts/` and `/reports/`

---

## 📁 Project Structure

```
your-project/
├── contracts/              # Specification documents
│   ├── system_architecture.md
│   ├── user_authentication.md
│   └── api_endpoints.md
├── reports/                # Generated reports
│   ├── review_2026-01-07.md
│   └── change_orders/
│       └── co_001.md
└── .factory/               # Project-specific config
    └── config.toml
```

---

## 🛠️ Configuration

### Global Config (`~/.factory/config.toml`)

```toml
[user]
name = "Your Name"
email = "you@example.com"

[github]
token_storage = "keyring"  # or "file", "env"

[llm]
provider = "ollama"
model = "codellama:7b"
endpoint = "http://localhost:11434"

[ui]
theme = "auto"  # auto, light, dark
```

### Project Config (`{project}/.factory/config.toml`)

```toml
[project]
name = "My Project"
repository = "owner/repo"

[contracts]
directory = "contracts"

[reports]
directory = "reports"
```

---

## 🤝 Team Collaboration

**First Team Member (Admin):**
1. Run `factory init` (full setup)
2. Commit `.factory/config.toml` to repository
3. Share repository with team

**Additional Team Members:**
1. Clone repository
2. Run `factory init --team`
3. Authenticate with GitHub (personal account)
4. Configure LLM (personal preference)
5. Start using Factory!

---

## 📚 Documentation

- **[User Guide](docs/USER_GUIDE.md)** - Complete user documentation
- **[Architecture](docs/ARCHITECTURE.md)** - Technical architecture for developers
- **[Bootstrap Spec](contracts/factory_bootstrap_spec.md)** - Onboarding flow specification
- **[System Architecture](contracts/system_architecture.md)** - System design
- **[Mode Specifications](contracts/mode_specifications.md)** - Detailed mode specs

---

## 🎯 Roadmap

### Phase 1: Foundation (Current)
- ✅ Core architecture design
- ✅ Specification documents
- 🚧 TUI implementation
- 🚧 LLM integration
- 🚧 GitHub integration

### Phase 2: Core Modes (Q1 2026)
- 🔜 INTAKE mode
- 🔜 REVIEW mode
- 🔜 CHANGE_ORDER mode
- 🔜 RESCUE mode

### Phase 3: Polish & Distribution (Q2 2026)
- 🔜 Installation scripts
- 🔜 Package managers (Homebrew, Winget, etc.)
- 🔜 Documentation site
- 🔜 Example projects

### Phase 4: Advanced Features (Q3 2026)
- 🔜 Web UI (localhost:3333)
- 🔜 Plugin system
- 🔜 CI/CD integration
- 🔜 Cloud sync (optional)

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

**Ways to Contribute:**
- 🐛 Report bugs
- 💡 Suggest features
- 📝 Improve documentation
- 🔧 Submit pull requests
- ⭐ Star the project

---

## 📄 License

Code-Factory is licensed under the [GNU General Public License v3.0](LICENSE).

**TL;DR:**
- ✅ Use for any purpose
- ✅ Modify and distribute
- ✅ Use in commercial projects
- ⚠️ Must disclose source
- ⚠️ Must use same license

---

## 🙏 Acknowledgments

Built with amazing open-source tools:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Ollama](https://ollama.ai) - Local LLM runtime
- [Go](https://go.dev) - Programming language

---

## 📞 Contact

- **GitHub:** [@ssdajoker](https://github.com/ssdajoker)
- **Issues:** [GitHub Issues](https://github.com/ssdajoker/Code-Factory/issues)
- **Discussions:** [GitHub Discussions](https://github.com/ssdajoker/Code-Factory/discussions)

---

**Happy Building! 🏭**
