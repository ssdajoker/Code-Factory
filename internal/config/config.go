package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Config holds all Factory configuration
type Config struct {
	// LLM settings
	LLM LLMConfig `toml:"llm"`

	// GitHub settings
	GitHub GitHubConfig `toml:"github"`

	// UI settings
	UI UIConfig `toml:"ui"`

	// Paths
	Paths PathsConfig `toml:"paths"`
}

// LLMConfig holds LLM provider settings
type LLMConfig struct {
	Provider    string `toml:"provider"`     // "openai", "anthropic", "ollama", "openrouter"
	Model       string `toml:"model"`        // Model name
	APIKeyStore string `toml:"api_key_store"` // "keyring", "env", "file"
	BaseURL     string `toml:"base_url"`     // Custom API endpoint (for Ollama, etc.)
}

// GitHubConfig holds GitHub integration settings
type GitHubConfig struct {
	TokenStorage string `toml:"token_storage"` // "keyring", "file", "env"
	DefaultOrg   string `toml:"default_org"`   // Default organization
}

// UIConfig holds UI preferences
type UIConfig struct {
	Theme       string `toml:"theme"`        // "dark", "light", "auto"
	ColorScheme string `toml:"color_scheme"` // Color scheme name
	Animations  bool   `toml:"animations"`   // Enable animations
}

// PathsConfig holds path settings
type PathsConfig struct {
	SpecsDir    string `toml:"specs_dir"`    // Specifications directory
	ReportsDir  string `toml:"reports_dir"`  // Reports output directory
	TemplateDir string `toml:"template_dir"` // Custom templates
}

// configDir returns the Factory config directory
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".factory"), nil
}

// configPath returns the path to config.toml
func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.toml"), nil
}

// GetDefault returns a Config with sensible defaults
func GetDefault() *Config {
	return &Config{
		LLM: LLMConfig{
			Provider:    "ollama",
			Model:       "llama3.2",
			APIKeyStore: "keyring",
			BaseURL:     "http://localhost:11434",
		},
		GitHub: GitHubConfig{
			TokenStorage: "keyring",
			DefaultOrg:   "",
		},
		UI: UIConfig{
			Theme:       "dark",
			ColorScheme: "default",
			Animations:  true,
		},
		Paths: PathsConfig{
			SpecsDir:    ".factory/specs",
			ReportsDir:  ".factory/reports",
			TemplateDir: "",
		},
	}
}

// Load reads configuration from ~/.factory/config.toml
// Returns default config if file doesn't exist
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return GetDefault(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return GetDefault(), nil
		}
		return nil, err
	}

	cfg := GetDefault()
	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Save writes configuration to ~/.factory/config.toml
func (c *Config) Save() error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	// Create config directory if needed
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := toml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// ConfigDir returns the Factory config directory path
func ConfigDir() (string, error) {
	return configDir()
}
