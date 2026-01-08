package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileStore(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "factory-secrets-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Override HOME
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	store, err := NewFileStore("test-password-123")
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	tests := []struct {
		name  string
		key   string
		value string
	}{
		{"simple key", "api_key", "sk-test123"},
		{"key with special chars", "oauth_token", "ghp_abc123XYZ!@#"},
		{"empty value", "empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set
			if err := store.Set(tt.key, tt.value); err != nil {
				t.Fatalf("Set() error = %v", err)
			}

			// Get
			got, err := store.Get(tt.key)
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if got != tt.value {
				t.Errorf("Get() = %q, want %q", got, tt.value)
			}

			// Delete
			if err := store.Delete(tt.key); err != nil {
				t.Fatalf("Delete() error = %v", err)
			}

			// Verify deleted
			_, err = store.Get(tt.key)
			if err != ErrSecretNotFound {
				t.Errorf("expected ErrSecretNotFound after delete, got %v", err)
			}
		})
	}
}

func TestFileStoreNotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "factory-secrets-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	store, _ := NewFileStore("password")
	_, err = store.Get("nonexistent")
	if err != ErrSecretNotFound {
		t.Errorf("expected ErrSecretNotFound, got %v", err)
	}
}

func TestEncryption(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "factory-secrets-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	store, _ := NewFileStore("password123")
	secret := "super-secret-api-key"
	store.Set("test_key", secret)

	// Read raw file and verify it's not plaintext
	secretsDir := filepath.Join(tmpDir, ".factory", "secrets")
	files, _ := os.ReadDir(secretsDir)
	if len(files) == 0 {
		t.Fatal("no secret files created")
	}

	data, _ := os.ReadFile(filepath.Join(secretsDir, files[0].Name()))
	if string(data) == secret {
		t.Error("secret stored in plaintext!")
	}
}

func TestKeyringStoreInterface(t *testing.T) {
	// Just verify KeyringStore implements SecretStore
	var _ SecretStore = &KeyringStore{}
	var _ SecretStore = &FileStore{}
}
