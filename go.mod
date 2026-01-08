module github.com/ssdajoker/Code-Factory

go 1.21

require (
	// TUI Framework
	github.com/charmbracelet/bubbletea v0.25.0
	github.com/charmbracelet/bubbles v0.18.0
	github.com/charmbracelet/lipgloss v0.9.1
	
	// CLI Framework
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2
	
	// Git Operations
	github.com/go-git/go-git/v5 v5.11.0
	
	// GitHub API
	github.com/google/go-github/v57 v57.0.0
	golang.org/x/oauth2 v0.15.0
	
	// LLM Clients
	github.com/sashabaranov/go-openai v1.19.2
	
	// Configuration
	gopkg.in/yaml.v3 v3.0.1
	
	// Secret Management
	github.com/zalando/go-keyring v0.2.3
	
	// Markdown
	github.com/yuin/goldmark v1.6.0
	
	// Utilities
	github.com/google/uuid v1.5.0
)
