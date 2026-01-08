package config

import (
        "os"
        "path/filepath"
        "testing"
)

func TestGetDefault(t *testing.T) {
        cfg := GetDefault()

        tests := []struct {
                name     string
                got      string
                expected string
        }{
                {"LLM Provider", cfg.LLM.Provider, "ollama"},
                {"LLM Model", cfg.LLM.Model, "llama3.2"},
                {"LLM APIKeyStore", cfg.LLM.APIKeyStore, "keyring"},
                {"LLM BaseURL", cfg.LLM.BaseURL, "http://localhost:11434"},
                {"GitHub TokenStorage", cfg.GitHub.TokenStorage, "keyring"},
                {"UI Theme", cfg.UI.Theme, "dark"},
                {"Paths SpecsDir", cfg.Paths.SpecsDir, ".factory/specs"},
                {"Paths ReportsDir", cfg.Paths.ReportsDir, ".factory/reports"},
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        if tt.got != tt.expected {
                                t.Errorf("got %q, want %q", tt.got, tt.expected)
                        }
                })
        }

        if !cfg.UI.Animations {
                t.Error("expected UI.Animations to be true by default")
        }
}

func TestLoadSave(t *testing.T) {
        // Create temp directory for test
        tmpDir, err := os.MkdirTemp("", "factory-config-test")
        if err != nil {
                t.Fatalf("failed to create temp dir: %v", err)
        }
        defer os.RemoveAll(tmpDir)

        // Override home for test
        oldHome := os.Getenv("HOME")
        os.Setenv("HOME", tmpDir)
        defer os.Setenv("HOME", oldHome)

        // Test Load returns default when no file exists
        cfg, err := Load()
        if err != nil {
                t.Fatalf("Load() error = %v", err)
        }
        if cfg.LLM.Provider != "ollama" {
                t.Errorf("expected default provider 'ollama', got %q", cfg.LLM.Provider)
        }

        // Verify default values loaded correctly
        if cfg.LLM.Model != "llama3.2" {
                t.Errorf("expected default model 'llama3.2', got %q", cfg.LLM.Model)
        }
}

func TestConfigPath(t *testing.T) {
        path, err := configPath()
        if err != nil {
                t.Fatalf("configPath() error = %v", err)
        }
        if !filepath.IsAbs(path) {
                t.Errorf("expected absolute path, got %q", path)
        }
        if filepath.Base(path) != "config.toml" {
                t.Errorf("expected config.toml, got %q", filepath.Base(path))
        }
}
