# System Architecture Specification

**Version:** 1.0.0  
**Status:** Draft  
**Last Updated:** 2026-01-07  
**Owner:** Code-Factory Core Team

---

## 1. Overview

The Spec-Driven Software Factory is a single-binary, cross-platform tool that helps developers maintain alignment between specifications and code. It operates in four distinct modes (INTAKE, REVIEW, CHANGE_ORDER, RESCUE) and provides both a beautiful TUI and optional web interface.

### 1.1 Design Principles

1. **Zero Dependencies:** Single Go binary, no external dependencies required
2. **Beautiful UX:** Canvas-style TUI using Charm.sh/Bubble Tea
3. **Local-First:** All data stored locally, git-native, no database
4. **Graceful Degradation:** Works offline, without LLM, without GitHub
5. **Cross-Platform:** Linux, macOS, Windows (WSL/Git Bash)
6. **Team-Friendly:** One person sets up, others clone and go

---

## 2. System Components

### 2.1 Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         Factory CLI                         │
│                      (cmd/factory)                          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ├─────────────────────────────┐
                              │                             │
                              ▼                             ▼
┌─────────────────────────────────────┐   ┌─────────────────────────────────┐
│          TUI Layer                  │   │       HTTP Server               │
│      (internal/tui)                 │   │    (internal/server)            │
│                                     │   │                                 │
│  • Bubble Tea Application           │   │  • localhost:3333               │
│  • Lipgloss Styling                 │   │  • Web Mirror of TUI            │
│  • Canvas-Style Interface           │   │  • REST API                     │
│  • Keyboard Navigation              │   │  • WebSocket for Live Updates   │
└─────────────────────────────────────┘   └─────────────────────────────────┘
                              │
                              ├─────────────────────────────┐
                              │                             │
                              ▼                             ▼
┌─────────────────────────────────────┐   ┌─────────────────────────────────┐
│         Mode Orchestrator           │   │      Configuration Manager      │
│      (internal/modes)               │   │      (internal/config)          │
│                                     │   │                                 │
│  • INTAKE Mode                      │   │  • TOML Parser                  │
│  • REVIEW Mode                      │   │  • Secret Management            │
│  • CHANGE_ORDER Mode                │   │  • Environment Variables        │
│  • RESCUE Mode                      │   │  • Validation                   │
└─────────────────────────────────────┘   └─────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   LLM Layer     │  │  GitHub Layer   │  │  Storage Layer  │
│ (internal/llm)  │  │(internal/github)│  │(internal/store) │
│                 │  │                 │  │                 │
│ • Ollama        │  │ • OAuth Flow    │  │ • File I/O      │
│ • OpenAI        │  │ • GitHub API    │  │ • Git Ops       │
│ • Anthropic     │  │ • App Install   │  │ • Markdown      │
│ • Azure         │  │ • Webhooks      │  │ • Templates     │
│ • Gemini        │  │                 │  │                 │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

---

## 3. Core Components

### 3.1 CLI Entry Point (`cmd/factory`)

**Responsibilities:**
- Parse command-line arguments
- Initialize configuration
- Route to appropriate mode or TUI
- Handle global flags (--version, --help, --debug)

**Key Files:**
- `main.go` - Entry point
- `commands.go` - Command definitions
- `flags.go` - Flag parsing

**Example Usage:**
```go
package main

import (
    "github.com/ssdajoker/Code-Factory/internal/config"
    "github.com/ssdajoker/Code-Factory/internal/tui"
    "github.com/ssdajoker/Code-Factory/internal/modes"
)

func main() {
    cfg := config.Load()
    
    switch cmd {
    case "init":
        modes.RunInit(cfg)
    case "intake":
        tui.StartTUI(cfg, modes.IntakeMode)
    case "review":
        tui.StartTUI(cfg, modes.ReviewMode)
    default:
        tui.StartTUI(cfg, modes.DefaultMode)
    }
}
```

### 3.2 TUI Layer (`internal/tui`)

**Responsibilities:**
- Render beautiful terminal interface
- Handle keyboard/mouse input
- Manage application state
- Coordinate with mode implementations

**Key Components:**
- `app.go` - Main Bubble Tea application
- `styles.go` - Lipgloss styling definitions
- `components/` - Reusable UI components
  - `banner.go` - Header banners
  - `menu.go` - Navigation menus
  - `editor.go` - Text editor component
  - `list.go` - Scrollable lists
  - `progress.go` - Progress indicators
  - `dialog.go` - Modal dialogs

**Bubble Tea Model:**
```go
type Model struct {
    mode        modes.Mode
    currentView string
    width       int
    height      int
    keymap      KeyMap
    styles      Styles
    
    // Mode-specific state
    modeState   interface{}
}

func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m Model) View() string
```

### 3.3 Mode Orchestrator (`internal/modes`)

**Responsibilities:**
- Implement four core modes
- Coordinate between LLM, GitHub, and Storage layers
- Manage mode-specific state and workflows

**Mode Interface:**
```go
type Mode interface {
    Name() string
    Description() string
    Init(cfg *config.Config) error
    Run(ctx context.Context) error
    HandleInput(input string) (Response, error)
    Cleanup() error
}
```

**Modes:**

#### 3.3.1 INTAKE Mode
- **Purpose:** Capture project vision and create specifications
- **Workflow:**
  1. Prompt user for project description
  2. Use LLM to generate structured spec
  3. Allow user to edit and refine
  4. Save to `/contracts/` directory
  5. Optionally commit to git

#### 3.3.2 REVIEW Mode
- **Purpose:** Check code against specifications
- **Workflow:**
  1. Load specs from `/contracts/`
  2. Scan codebase for relevant files
  3. Use LLM to compare code vs spec
  4. Generate compliance report
  5. Highlight deviations and suggest fixes
  6. Save report to `/reports/`

#### 3.3.3 CHANGE_ORDER Mode
- **Purpose:** Track specification drift over time
- **Workflow:**
  1. Detect changes in codebase (git diff)
  2. Compare changes against specs
  3. Identify intentional vs unintentional drift
  4. Create change order document
  5. Optionally create GitHub issue
  6. Update specs if change is approved

#### 3.3.4 RESCUE Mode
- **Purpose:** Reverse-engineer existing codebase into specs
- **Workflow:**
  1. Scan entire codebase
  2. Use LLM to infer architecture and patterns
  3. Generate specification documents
  4. Create dependency diagrams
  5. Identify technical debt
  6. Save generated specs to `/contracts/`

### 3.4 LLM Layer (`internal/llm`)

**Responsibilities:**
- Abstract LLM provider differences
- Handle API calls and streaming
- Manage context windows and token limits
- Implement fallback strategies

**Provider Interface:**
```go
type Provider interface {
    Name() string
    IsAvailable() bool
    Complete(ctx context.Context, prompt string, opts Options) (string, error)
    Stream(ctx context.Context, prompt string, opts Options) (<-chan string, error)
    Models() []string
}
```

**Implementations:**
- `ollama.go` - Local Ollama integration
- `openai.go` - OpenAI API
- `anthropic.go` - Claude API
- `azure.go` - Azure OpenAI
- `gemini.go` - Google Gemini

**Prompt Templates:**
- `prompts/intake.tmpl` - Spec generation prompts
- `prompts/review.tmpl` - Code review prompts
- `prompts/change_order.tmpl` - Change detection prompts
- `prompts/rescue.tmpl` - Reverse engineering prompts

### 3.5 GitHub Layer (`internal/github`)

**Responsibilities:**
- OAuth authentication flow
- GitHub App installation
- API interactions (repos, issues, PRs)
- Webhook handling (optional)

**Key Components:**
- `oauth.go` - OAuth flow implementation
- `app.go` - GitHub App management
- `api.go` - GitHub API client wrapper
- `webhooks.go` - Webhook server (optional)

**API Client:**
```go
type Client struct {
    token      string
    httpClient *http.Client
    github     *github.Client
}

func (c *Client) GetRepo(owner, repo string) (*Repository, error)
func (c *Client) CreateIssue(owner, repo string, issue *Issue) error
func (c *Client) ListPRs(owner, repo string) ([]*PullRequest, error)
func (c *Client) CommentOnPR(owner, repo string, number int, comment string) error
```

### 3.6 Storage Layer (`internal/store`)

**Responsibilities:**
- Read/write specifications
- Generate reports
- Manage git operations
- Handle file templates

**Key Components:**
- `contracts.go` - Spec file management
- `reports.go` - Report generation
- `git.go` - Git operations wrapper
- `templates.go` - Template rendering

**Contract Store:**
```go
type ContractStore struct {
    basePath string
}

func (s *ContractStore) List() ([]*Contract, error)
func (s *ContractStore) Get(name string) (*Contract, error)
func (s *ContractStore) Save(contract *Contract) error
func (s *ContractStore) Delete(name string) error
```

### 3.7 Configuration Manager (`internal/config`)

**Responsibilities:**
- Load configuration from files and environment
- Validate configuration
- Manage secrets securely
- Provide defaults

**Configuration Structure:**
```go
type Config struct {
    User     UserConfig
    GitHub   GitHubConfig
    LLM      LLMConfig
    UI       UIConfig
    Project  ProjectConfig
}

type UserConfig struct {
    Name  string
    Email string
}

type GitHubConfig struct {
    Token         string
    TokenStorage  string // "keyring", "file", "env"
    DefaultOrg    string
}

type LLMConfig struct {
    Provider    string
    Model       string
    Endpoint    string
    APIKey      string
    Temperature float64
    MaxTokens   int
    Fallback    *LLMConfig
}

type UIConfig struct {
    Theme   string // "auto", "light", "dark"
    Editor  string
    Browser string
}

type ProjectConfig struct {
    Name              string
    Repository        string
    ContractsDir      string
    ReportsDir        string
    GitHubAppInstalled bool
    InstallationID    int64
}
```

---

## 4. Data Flow

### 4.1 INTAKE Mode Flow

```
User Input (Vision)
    │
    ▼
TUI (Capture Input)
    │
    ▼
LLM Layer (Generate Spec)
    │
    ▼
TUI (Display Draft)
    │
    ▼
User (Edit & Refine)
    │
    ▼
Storage Layer (Save to /contracts/)
    │
    ▼
Git Layer (Optional Commit)
    │
    ▼
GitHub Layer (Optional Push)
```

### 4.2 REVIEW Mode Flow

```
Storage Layer (Load Specs)
    │
    ▼
Git Layer (Get Changed Files)
    │
    ▼
Storage Layer (Load Code Files)
    │
    ▼
LLM Layer (Compare Code vs Spec)
    │
    ▼
TUI (Display Results)
    │
    ▼
Storage Layer (Generate Report)
    │
    ▼
GitHub Layer (Optional PR Comment)
```

### 4.3 CHANGE_ORDER Mode Flow

```
Git Layer (Detect Changes)
    │
    ▼
Storage Layer (Load Relevant Specs)
    │
    ▼
LLM Layer (Analyze Drift)
    │
    ▼
TUI (Display Change Order)
    │
    ▼
User (Approve/Reject)
    │
    ├─ Approve ─▶ Storage Layer (Update Spec)
    │
    └─ Reject ──▶ GitHub Layer (Create Issue)
```

### 4.4 RESCUE Mode Flow

```
Git Layer (Scan Codebase)
    │
    ▼
Storage Layer (Load All Files)
    │
    ▼
LLM Layer (Infer Architecture)
    │
    ▼
LLM Layer (Generate Specs)
    │
    ▼
TUI (Display Generated Specs)
    │
    ▼
User (Review & Edit)
    │
    ▼
Storage Layer (Save to /contracts/)
```

---

## 5. File System Layout

### 5.1 Global Configuration

```
~/.factory/
├── config.toml              # User preferences
├── github_token             # OAuth token (encrypted)
├── logs/
│   ├── github.log           # GitHub API calls
│   ├── llm.log              # LLM queries
│   └── factory.log          # General logs
├── cache/
│   ├── ollama_models/       # Cached Ollama models
│   └── github_repos/        # Cached repo metadata
└── templates/
    ├── spec_template.md     # Default spec template
    └── report_template.md   # Default report template
```

### 5.2 Project Structure

```
project-root/
├── .factory/
│   ├── config.toml          # Project-specific config
│   ├── cache/               # Local cache (gitignored)
│   └── temp/                # Temporary files (gitignored)
├── contracts/
│   ├── README.md
│   ├── system_architecture.md
│   ├── mode_specifications.md
│   └── features/
│       ├── feature_a.md
│       └── feature_b.md
└── reports/
    ├── README.md
    ├── review_2026-01-07.md
    └── change_orders/
        └── co_001.md
```

---

## 6. Security Architecture

### 6.1 Secret Management

**Token Storage Hierarchy:**
1. **System Keyring (Preferred)**
   - macOS: Keychain Access
   - Windows: Credential Manager
   - Linux: Secret Service API

2. **Encrypted File (Fallback)**
   - AES-256-GCM encryption
   - Key derived from machine ID + user ID

3. **Plain File (Last Resort)**
   - Warn user about security implications
   - Recommend environment variable instead

**Implementation:**
```go
type SecretStore interface {
    Store(key, value string) error
    Retrieve(key string) (string, error)
    Delete(key string) error
}

// Implementations
type KeyringStore struct { ... }
type EncryptedFileStore struct { ... }
type PlainFileStore struct { ... }
```

### 6.2 API Security

**GitHub Token:**
- Minimum required scopes
- Token rotation every 90 days
- Revocation on logout

**LLM API Keys:**
- Never logged or displayed
- Stored in keyring when possible
- Option to use environment variables

### 6.3 Data Privacy

**Local-First:**
- All specs stored locally
- Reports generated locally
- No telemetry by default

**Cloud LLM Privacy:**
- Warn user when using cloud LLMs
- Option to redact sensitive data
- Option to use local Ollama

---

## 7. Performance Considerations

### 7.1 Startup Time

**Target:** < 100ms cold start

**Optimizations:**
- Lazy load LLM providers
- Defer GitHub API calls
- Cache configuration
- Minimal dependencies

### 7.2 LLM Response Time

**Target:** < 5 seconds for typical queries

**Optimizations:**
- Stream responses for immediate feedback
- Use smaller models for simple tasks
- Cache common prompts
- Parallel processing where possible

### 7.3 File I/O

**Target:** < 1 second for typical projects

**Optimizations:**
- Incremental file scanning
- Git-based change detection
- Parallel file reading
- Memory-mapped files for large codebases

---

## 8. Error Handling

### 8.1 Error Categories

1. **User Errors:** Invalid input, missing configuration
2. **Network Errors:** GitHub API failures, LLM timeouts
3. **System Errors:** File I/O failures, permission issues
4. **Logic Errors:** Unexpected state, assertion failures

### 8.2 Error Recovery

**Graceful Degradation:**
- Offline mode when GitHub unavailable
- Manual mode when LLM unavailable
- Read-only mode when permissions insufficient

**User Feedback:**
- Clear error messages
- Actionable recovery steps
- Links to documentation
- Option to report bugs

### 8.3 Logging

**Log Levels:**
- DEBUG: Detailed diagnostic information
- INFO: General informational messages
- WARN: Warning messages (non-critical)
- ERROR: Error messages (recoverable)
- FATAL: Fatal errors (unrecoverable)

**Log Rotation:**
- Daily rotation
- Keep 30 days of logs
- Compress old logs
- User can clear logs

---

## 9. Testing Strategy

### 9.1 Unit Tests

**Coverage Target:** > 80%

**Key Areas:**
- Configuration parsing
- LLM prompt generation
- GitHub API interactions
- File I/O operations
- Template rendering

### 9.2 Integration Tests

**Test Scenarios:**
- End-to-end mode workflows
- OAuth flow (mocked)
- LLM integration (mocked)
- Git operations (real)

### 9.3 E2E Tests

**Test Scenarios:**
- Installation on different platforms
- First-time setup flow
- Team setup flow
- Offline mode
- Error recovery

---

## 10. Deployment Architecture

### 10.1 Binary Distribution

**Platforms:**
- Linux: amd64, arm64
- macOS: amd64 (Intel), arm64 (Apple Silicon)
- Windows: amd64 (via WSL or Git Bash)

**Distribution Channels:**
- GitHub Releases (primary)
- Homebrew (macOS/Linux)
- Winget (Windows)
- Docker (cross-platform)
- Nix (NixOS)

### 10.2 Release Process

1. **Version Bump:** Update version in code
2. **Changelog:** Generate from git commits
3. **Build:** Cross-compile for all platforms
4. **Test:** Run E2E tests on each platform
5. **Sign:** Code sign binaries (macOS, Windows)
6. **Checksum:** Generate SHA256 checksums
7. **Release:** Create GitHub release
8. **Publish:** Update package managers

### 10.3 Update Mechanism

**Auto-Update (Future):**
- Check for updates on startup (opt-in)
- Download and verify new binary
- Replace old binary atomically
- Restart with new version

**Manual Update:**
```bash
factory update
```

---

## 11. Extensibility

### 11.1 Plugin System (Future)

**Plugin Interface:**
```go
type Plugin interface {
    Name() string
    Version() string
    Init(cfg *config.Config) error
    RegisterCommands() []Command
    RegisterModes() []Mode
    Shutdown() error
}
```

**Plugin Discovery:**
- Load from `~/.factory/plugins/`
- Verify plugin signature
- Sandbox plugin execution

### 11.2 Custom Modes

**User-Defined Modes:**
- Define in `~/.factory/modes/`
- Use Lua or JavaScript for scripting
- Access to core APIs (LLM, GitHub, Storage)

### 11.3 Custom Templates

**Template System:**
- Jinja2-style templates
- User templates in `~/.factory/templates/`
- Project templates in `.factory/templates/`

---

## 12. Monitoring & Observability

### 12.1 Metrics (Opt-In)

**Collected Metrics:**
- Command usage frequency
- Mode usage distribution
- LLM provider usage
- Error rates by category
- Performance metrics (latency, throughput)

**Storage:**
- Local: `~/.factory/metrics.log`
- Remote: (opt-in) Anonymous telemetry

### 12.2 Health Checks

**`factory doctor` Command:**
- Check GitHub connectivity
- Check LLM availability
- Verify configuration
- Test file permissions
- Diagnose common issues

---

## 13. Future Enhancements

### 13.1 Phase 2 (Q2 2026)

- Web UI (localhost:3333)
- Plugin system
- Auto-update mechanism
- Cloud sync (optional)

### 13.2 Phase 3 (Q3 2026)

- CI/CD integration (GitHub Actions, GitLab CI)
- Slack/Discord notifications
- Multi-repo support
- Spec marketplace

### 13.3 Phase 4 (Q4 2026)

- VS Code extension
- JetBrains plugin
- Real-time collaboration
- Enterprise features (SSO, audit logs)

---

**End of Specification**
