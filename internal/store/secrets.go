package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/argon2"
)

const (
	serviceName = "code-factory"
	argon2Time  = 1
	argon2Mem   = 64 * 1024
	argon2Threads = 4
	argon2KeyLen  = 32
)

var (
	ErrSecretNotFound = errors.New("secret not found")
	ErrKeyringUnavailable = errors.New("keyring unavailable")
)

// SecretStore defines the interface for secret storage
type SecretStore interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

// KeyringStore uses OS keyring for secret storage (Tier 1)
type KeyringStore struct{}

// NewKeyringStore creates a new keyring-based store
func NewKeyringStore() *KeyringStore {
	return &KeyringStore{}
}

// Get retrieves a secret from the keyring
func (k *KeyringStore) Get(key string) (string, error) {
	val, err := keyring.Get(serviceName, key)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", ErrSecretNotFound
		}
		return "", err
	}
	return val, nil
}

// Set stores a secret in the keyring
func (k *KeyringStore) Set(key, value string) error {
	return keyring.Set(serviceName, key, value)
}

// Delete removes a secret from the keyring
func (k *KeyringStore) Delete(key string) error {
	err := keyring.Delete(serviceName, key)
	if err == keyring.ErrNotFound {
		return ErrSecretNotFound
	}
	return err
}

// IsAvailable checks if keyring is available on this system
func (k *KeyringStore) IsAvailable() bool {
	testKey := "__factory_test__"
	err := keyring.Set(serviceName, testKey, "test")
	if err != nil {
		return false
	}
	keyring.Delete(serviceName, testKey)
	return true
}

// FileStore uses encrypted files for secret storage (Tier 2)
type FileStore struct {
	dir      string
	password string // User-provided password for key derivation
}

// NewFileStore creates a new file-based encrypted store
// password should be user-provided or a high-entropy random key
func NewFileStore(password string) (*FileStore, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(home, ".factory", "secrets")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return &FileStore{dir: dir, password: password}, nil
}

// deriveKey uses Argon2id to derive an encryption key from password
func (f *FileStore) deriveKey(salt []byte) []byte {
	return argon2.IDKey(
		[]byte(f.password),
		salt,
		argon2Time,
		argon2Mem,
		argon2Threads,
		argon2KeyLen,
	)
}

// Get retrieves and decrypts a secret from file
func (f *FileStore) Get(key string) (string, error) {
	path := filepath.Join(f.dir, key+".enc")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrSecretNotFound
		}
		return "", err
	}

	// Decode base64
	ciphertext, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return "", err
	}

	// Extract salt (first 16 bytes) and nonce (next 12 bytes)
	if len(ciphertext) < 28 {
		return "", errors.New("invalid encrypted data")
	}
	salt := ciphertext[:16]
	nonce := ciphertext[16:28]
	encrypted := ciphertext[28:]

	// Derive key and decrypt
	key32 := f.deriveKey(salt)
	block, err := aes.NewCipher(key32)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Set encrypts and stores a secret to file
func (f *FileStore) Set(key, value string) error {
	// Generate random salt and nonce
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	// Derive key and encrypt
	key32 := f.deriveKey(salt)
	block, err := aes.NewCipher(key32)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	encrypted := gcm.Seal(nil, nonce, []byte(value), nil)

	// Combine salt + nonce + ciphertext
	ciphertext := append(salt, nonce...)
	ciphertext = append(ciphertext, encrypted...)

	// Encode and write
	path := filepath.Join(f.dir, key+".enc")
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return os.WriteFile(path, []byte(encoded), 0600)
}

// Delete removes a secret file
func (f *FileStore) Delete(key string) error {
	path := filepath.Join(f.dir, key+".enc")
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return ErrSecretNotFound
	}
	return err
}

// AutoStore automatically selects the best available storage tier
type AutoStore struct {
	primary   SecretStore
	fallback  SecretStore
	usingFallback bool
}

// NewAutoStore creates a store that uses keyring if available, else encrypted file
func NewAutoStore(filePassword string) (*AutoStore, error) {
	kr := NewKeyringStore()
	if kr.IsAvailable() {
		return &AutoStore{primary: kr, usingFallback: false}, nil
	}

	fs, err := NewFileStore(filePassword)
	if err != nil {
		return nil, err
	}
	return &AutoStore{primary: fs, usingFallback: true}, nil
}

// Get retrieves a secret
func (a *AutoStore) Get(key string) (string, error) {
	return a.primary.Get(key)
}

// Set stores a secret
func (a *AutoStore) Set(key, value string) error {
	return a.primary.Set(key, value)
}

// Delete removes a secret
func (a *AutoStore) Delete(key string) error {
	return a.primary.Delete(key)
}

// UsingFallback returns true if using file storage instead of keyring
func (a *AutoStore) UsingFallback() bool {
	return a.usingFallback
}
