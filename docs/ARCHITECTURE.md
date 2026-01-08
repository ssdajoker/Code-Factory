# Code-Factory Technical Architecture

**Version:** 1.0.0  
**Last Updated:** 2026-01-07

---

## Table of Contents

1. [Overview](#overview)
2. [Technology Stack](#technology-stack)
3. [Project Structure](#project-structure)
4. [Core Components](#core-components)
5. [Data Flow](#data-flow)
6. [Build & Deployment](#build--deployment)
7. [Development Guide](#development-guide)
8. [Testing Strategy](#testing-strategy)
9. [Contributing](#contributing)

---

## Overview

Code-Factory is a single-binary, cross-platform tool for maintaining alignment 
between software specifications and code. It's built in Go with a beautiful 
terminal UI using Charm.sh's Bubble Tea framework.

**Design Goals:**
- Zero dependencies (single binary)
- Beautiful, intuitive TUI
- Cross-platform (Linux, macOS, Windows)
- Local-first (git-native storage)
- Graceful degradation (works offline, without LLM, without GitHub)

---

## Technology Stack

### Core
- **Language:** Go 1.21+
- **TUI Framework:** Bubble Tea (charmbracelet/bubbletea)
- **Styling:** Lipgloss (charmbracelet/lipgloss)
- **Configuration:** TOML (BurntSushi/toml)

### Integrations
- **GitHub:** go-github/v57 + golang.org/x/oauth2
- **LLM:** 
  - Ollama (local)
  - OpenAI API
  - Anthropic API
  - Azure OpenAI
  - Google Gemini

### Storage
- **Specs:** Markdown files in `/contracts/`
- **Reports:** Markdown files in `/reports/`
- **Config:** TOML files in `~/.factory/` and `.factory/`
- **Secrets:** System keyring or encrypted files

---

## Project Structure

```
Code-Factory/
├── cmd/
│   └── factory/              # Main entry point
│       ├── main.go           # CLI initialization
│       ├── commands.go       # Command definitions
│       └── flags.go          # Flag parsing
│
├── internal/
│   ├── tui/                  # Terminal UI
│   │   ├── app.go            # Main Bubble Tea app
│   │   ├── styles.go         # Lipgloss styles
│   │   ├── keymap.go         # Keyboard shortcuts
│   │   └── components/       # Reusable UI components
│   │       ├── banner.go
│   │       ├── menu.go
│   │       ├── editor.go
│   │       ├── list.go
│   │       ├── progress.go
│   │       └── dialog.go
│   │
│   ├── modes/                # Four core modes
│   │   ├── mode.go           # Mode interface
│   │   ├── intake.go         # INTAKE mode
│   │   ├── review.go         # REVIEW mode
│   │   ├── change_order.go   # CHANGE_ORDER mode
│   │   └── rescue.go         # RESCUE mode
│   │
│   ├── llm/                  # LLM integration
│   │   ├── provider.go       # Provider interface
│   │   ├── ollama.go         # Ollama implementation
│   │   ├── openai.go         # OpenAI implementation
│   │   ├── anthropic.go      # Anthropic implementation
│   │   ├── azure.go          # Azure OpenAI
│   │   ├── gemini.go         # Google Gemini
│   │   └── prompts/          # Prompt templates
│   │       ├── intake.tmpl
│   │       ├── review.tmpl
│   │       ├── change_order.tmpl
│   │       └── rescue.tmpl
│   │
│   ├── github/               # GitHub integration
│   │   ├── oauth.go          # OAuth flow
│   │   ├── app.go            # GitHub App
│   │   ├── api.go            # API client
│   │   └── webhooks.go       # Webhook server (optional)
│   │
│   ├── store/                # Storage layer
│   │   ├── contracts.go      # Spec file management
│   │   ├── reports.go        # Report generation
│   │   ├── git.go            # Git operations
│   │   └── templates.go      # Template rendering
│   │
│   ├── config/               # Configuration
│   │   ├── config.go         # Config loading
│   │   ├── secrets.go        # Secret management
│   │   └── validation.go     # Config validation
│   │
│   └── server/               # HTTP server (optional)
│       ├── server.go         # HTTP server
│       ├── handlers.go       # HTTP handlers
│       └── websocket.go      # WebSocket for live updates
│
├── contracts/                # Specification documents
│   ├── factory_bootstrap_spec.md
│   ├── system_architecture.md
│   └── mode_specifications.md
│
├── reports/                  # Generated reports (examples)
│   └── README.md
│
├── docs/                     # Documentation
│   ├── ARCHITECTURE.md       # This file
│   └── USER_GUIDE.md         # User documentation
│
├── scripts/                  # Installation scripts
│   └── install.sh            # Curl-able install script
│
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
├── README.md                 # Project README
└── LICENSE                   # GPL v3.0 license
```

---

## Core Components

### 1. CLI Entry Point (`cmd/factory`)

**Responsibilities:**
- Parse command-line arguments
- Load configuration
- Initialize components
- Route to appropriate mode or TUI

**Key Functions:**
```go
func main() {
    // Parse flags
    cmd, args := parseArgs()
    
    // Load config
    cfg := config.Load()
    
    // Route to command
    switch cmd {
    case "init":
        runInit(cfg)
    case "intake":
        runMode(cfg, modes.IntakeMode)
    case "review":
        runMode(cfg, modes.ReviewMode)
    case "change-order":
        runMode(cfg, modes.ChangeOrderMode)
    case "rescue":
        runMode(cfg, modes.RescueMode)
    default:
        runTUI(cfg)
    }
}
```

### 2. TUI Layer (`internal/tui`)

**Responsibilities:**
- Render terminal interface
- Handle user input
- Manage application state
- Coordinate with modes

**Bubble Tea Model:**
```go
type Model struct {
    // Core state
    mode        modes.Mode
    currentView string
    width       int
    height      int
    
    // UI components
    keymap      KeyMap
    styles      Styles
    
    // Mode-specific state
    modeState   interface{}
}

func (m Model) Init() tea.Cmd {
    return m.mode.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        return m, nil
    default:
        return m.mode.Update(msg)
    }
}

func (m Model) View() string {
    return m.mode.View(m.width, m.height)
}
```

### 3. Mode Implementations (`internal/modes`)

**Mode Interface:**
```go
type Mode interface {
    Name() string
    Description() string
    Init() tea.Cmd
    Update(msg tea.Msg) (tea.Model, tea.Cmd)
    View(width, height int) string
    Cleanup() error
}
```

**Example: INTAKE Mode:**
```go
type IntakeMode struct {
    cfg         *config.Config
    llm         llm.Provider
    store       *store.ContractStore
    
    // State
    input       string
    generating  bool
    draft       string
    editing     bool
}

func (m *IntakeMode) Init() tea.Cmd {
    return nil
}

func (m *IntakeMode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "ctrl+enter" {
            return m, m.generateSpec()
        }
    }
    return m, nil
}

func (m *IntakeMode) generateSpec() tea.Cmd {
    return func() tea.Msg {
        prompt := loadPrompt("intake", m.input)
        spec, err := m.llm.Complete(context.Background(), prompt)
        if err != nil {
            return errorMsg{err}
        }
        return specGeneratedMsg{spec}
    }
}
```

### 4. LLM Layer (`internal/llm`)

**Provider Interface:**
```go
type Provider interface {
    Name() string
    IsAvailable() bool
    Complete(ctx context.Context, prompt string, opts Options) (string, error)
    Stream(ctx context.Context, prompt string, opts Options) (<-chan string, error)
    Models() []string
}

type Options struct {
    Model       string
    Temperature float64
    MaxTokens   int
    Stop        []string
}
```

**Ollama Implementation:**
```go
type OllamaProvider struct {
    endpoint string
    model    string
    client   *http.Client
}

func (p *OllamaProvider) IsAvailable() bool {
    resp, err := p.client.Get(p.endpoint + "/api/tags")
    return err == nil && resp.StatusCode == 200
}

func (p *OllamaProvider) Complete(ctx context.Context, prompt string, opts Options) (string, error) {
    req := ollamaRequest{
        Model:  opts.Model,
        Prompt: prompt,
        Options: map[string]interface{}{
            "temperature": opts.Temperature,
            "num_predict": opts.MaxTokens,
        },
    }
    
    resp, err := p.client.Post(p.endpoint+"/api/generate", "application/json", toJSON(req))
    if err != nil {
        return "", err
    }
    
    var result ollamaResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return result.Response, nil
}
```

### 5. GitHub Layer (`internal/github`)

**OAuth Flow:**
```go
type OAuthFlow struct {
    clientID     string
    clientSecret string
    redirectURL  string
    server       *http.Server
}

func (f *OAuthFlow) Start() (string, error) {
    // Generate state for CSRF protection
    state := generateRandomState()
    
    // Build authorization URL
    authURL := fmt.Sprintf(
        "https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s&scope=repo,read:user",
        f.clientID, f.redirectURL, state,
    )
    
    // Start local callback server
    tokenChan := make(chan string)
    f.server = f.startCallbackServer(state, tokenChan)
    
    // Open browser
    openBrowser(authURL)
    
    // Wait for callback
    token := <-tokenChan
    return token, nil
}

func (f *OAuthFlow) startCallbackServer(expectedState string, tokenChan chan string) *http.Server {
    mux := http.NewServeMux()
    mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
        code := r.URL.Query().Get("code")
        state := r.URL.Query().Get("state")
        
        if state != expectedState {
            http.Error(w, "Invalid state", http.StatusBadRequest)
            return
        }
        
        token, err := f.exchangeCode(code)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        tokenChan <- token
        fmt.Fprintf(w, "Success! You can close this window.")
    })
    
    server := &http.Server{Addr: ":8765", Handler: mux}
    go server.ListenAndServe()
    return server
}
```

### 6. Storage Layer (`internal/store`)

**Contract Store:**
```go
type ContractStore struct {
    basePath string
}

func (s *ContractStore) List() ([]*Contract, error) {
    files, err := filepath.Glob(filepath.Join(s.basePath, "*.md"))
    if err != nil {
        return nil, err
    }
    
    contracts := make([]*Contract, 0, len(files))
    for _, file := range files {
        contract, err := s.Get(filepath.Base(file))
        if err != nil {
            continue
        }
        contracts = append(contracts, contract)
    }
    
    return contracts, nil
}

func (s *ContractStore) Get(name string) (*Contract, error) {
    path := filepath.Join(s.basePath, name)
    content, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    return &Contract{
        Name:    name,
        Content: string(content),
        Path:    path,
    }, nil
}

func (s *ContractStore) Save(contract *Contract) error {
    path := filepath.Join(s.basePath, contract.Name)
    return os.WriteFile(path, []byte(contract.Content), 0644)
}
```

---

## Data Flow

### INTAKE Mode Flow

```
User Input
    │
    ▼
TUI (Capture)
    │
    ▼
LLM Provider (Generate)
    │
    ▼
TUI (Display Draft)
    │
    ▼
User (Edit)
    │
    ▼
Contract Store (Save)
    │
    ▼
Git (Commit)
```

### REVIEW Mode Flow

```
Contract Store (Load Specs)
    │
    ▼
Git (Get Changed Files)
    │
    ▼
File System (Load Code)
    │
    ▼
LLM Provider (Analyze)
    │
    ▼
TUI (Display Results)
    │
    ▼
Report Store (Generate Report)
```

---

## Build & Deployment

### Building

```bash
# Build for current platform
go build -o factory ./cmd/factory

# Build for all platforms
GOOS=linux GOARCH=amd64 go build -o factory-linux-amd64 ./cmd/factory
GOOS=darwin GOARCH=amd64 go build -o factory-darwin-amd64 ./cmd/factory
GOOS=darwin GOARCH=arm64 go build -o factory-darwin-arm64 ./cmd/factory
GOOS=windows GOARCH=amd64 go build -o factory-windows-amd64.exe ./cmd/factory
```

### Release Process

1. **Version Bump:** Update version in `cmd/factory/main.go`
2. **Changelog:** Generate from git commits
3. **Build:** Cross-compile for all platforms
4. **Test:** Run E2E tests on each platform
5. **Sign:** Code sign binaries (macOS, Windows)
6. **Checksum:** Generate SHA256 checksums
7. **Release:** Create GitHub release with binaries
8. **Publish:** Update Homebrew formula, Winget manifest

---

## Development Guide

### Prerequisites

- Go 1.21+
- Git
- (Optional) Ollama for local LLM testing

### Setup

```bash
# Clone repository
git clone https://github.com/ssdajoker/Code-Factory.git
cd Code-Factory

# Install dependencies
go mod download

# Build
go build -o factory ./cmd/factory

# Run
./factory
```

### Development Workflow

1. **Create Feature Branch:** `git checkout -b feature/my-feature`
2. **Make Changes:** Edit code
3. **Test:** `go test ./...`
4. **Lint:** `golangci-lint run`
5. **Commit:** `git commit -m "feat: my feature"`
6. **Push:** `git push origin feature/my-feature`
7. **PR:** Create pull request on GitHub

### Code Style

- Follow Go conventions (gofmt, golint)
- Use meaningful variable names
- Add comments for exported functions
- Keep functions small and focused
- Write tests for new features

---

## Testing Strategy

### Unit Tests

```bash
go test ./...
```

**Coverage Target:** > 80%

**Key Areas:**
- Configuration parsing
- LLM prompt generation
- GitHub API interactions
- File I/O operations

### Integration Tests

```bash
go test -tags=integration ./...
```

**Test Scenarios:**
- End-to-end mode workflows
- OAuth flow (mocked)
- LLM integration (mocked)
- Git operations (real)

### E2E Tests

```bash
./scripts/e2e_test.sh
```

**Test Scenarios:**
- Installation on different platforms
- First-time setup flow
- Team setup flow
- Offline mode

---

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for detailed contribution guidelines.

**Quick Start:**
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Write tests
5. Submit a pull request

---

**Last Updated:** 2026-01-07  
**Maintainer:** ssdajoker
