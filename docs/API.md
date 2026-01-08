# Factory API Documentation

## Overview

Factory provides both a CLI interface and programmatic Go API for spec-driven development workflows.

## CLI Commands

### Global Flags

```
--help, -h    Show help for any command
--version     Show version information
```

### Commands

#### `factory init`

Initialize Factory in the current project.

```bash
factory init
```

Creates:
- `.factory/` directory
- `contracts/` directory for specifications
- Default configuration

#### `factory intake`

Start INTAKE mode to capture project vision.

```bash
factory intake [flags]
```

Flags:
- `--output, -o` - Output directory for generated spec
- `--no-llm` - Skip LLM enhancement, use template only

#### `factory review`

Start REVIEW mode to check code against specifications.

```bash
factory review [flags]
```

Flags:
- `--spec, -s` - Path to specification file
- `--output, -o` - Output path for review report

#### `factory rescue`

Start RESCUE mode to reverse-engineer a codebase.

```bash
factory rescue [flags]
```

Flags:
- `--path, -p` - Path to codebase (default: current directory)
- `--output, -o` - Output path for generated spec

#### `factory change-order`

Start CHANGE_ORDER mode to track specification drift.

```bash
factory change-order [flags]
```

Flags:
- `--spec, -s` - Path to specification file
- `--watch` - Watch for changes continuously

#### `factory llm`

Manage LLM providers.

```bash
factory llm status          # Show current LLM status
factory llm detect          # Auto-detect available providers
factory llm set <provider>  # Set active provider
```

#### `factory github`

GitHub integration commands.

```bash
factory github auth         # Authenticate with GitHub
factory github status       # Show authentication status
factory github repos        # List accessible repositories
```

## Go API

### Config Package

```go
import "github.com/ssdajoker/Code-Factory/internal/config"

// Load configuration
cfg, err := config.Load()

// Get defaults
cfg := config.GetDefault()
```

### LLM Package

```go
import "github.com/ssdajoker/Code-Factory/internal/llm"

// Create provider
provider, err := llm.NewProvider(llm.Config{
    Type:   llm.ProviderOllama,
    Model:  "llama3.2",
})

// Generate completion
response, err := provider.Complete(ctx, prompt, llm.DefaultOptions())

// Auto-detect providers
detector := llm.NewDetector(ollamaURL, openAIKey, anthropicKey)
result := detector.Detect(ctx)
```

### Modes Package

```go
import "github.com/ssdajoker/Code-Factory/internal/modes"

// Create intake mode
intake := modes.NewIntakeMode(provider, "contracts")

// Navigate steps
intake.NextStep()
intake.SetStepValue("My Project")

// Generate specification
spec, err := intake.GenerateSpec(ctx)
```

### GitHub Package

```go
import "github.com/ssdajoker/Code-Factory/internal/github"

// OAuth flow
flow := github.NewOAuthFlow(clientID)
deviceCode, err := flow.InitiateDeviceFlow()
token, err := flow.PollForToken(deviceCode.DeviceCode, deviceCode.Interval)

// API client
client := github.NewClient(token)
user, err := client.GetUser()
repos, err := client.ListRepos()
```

### Store Package

```go
import "github.com/ssdajoker/Code-Factory/internal/store"

// Keyring store (preferred)
keyring := store.NewKeyringStore()
if keyring.IsAvailable() {
    keyring.Set("api_key", "secret")
    value, err := keyring.Get("api_key")
}

// File store (fallback)
fileStore, err := store.NewFileStore("password")
fileStore.Set("api_key", "secret")
```

## Error Handling

Common errors:

- `ErrNoProvider` - No LLM provider available
- `ErrProviderFailed` - LLM request failed
- `ErrSecretNotFound` - Secret not found in store
- `ErrKeyringUnavailable` - OS keyring not available
