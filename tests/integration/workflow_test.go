//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFullWorkflow(t *testing.T) {
	// Create temp project directory
	tmpDir, err := os.MkdirTemp("", "factory-workflow")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a minimal project structure
	contractsDir := filepath.Join(tmpDir, "contracts")
	os.MkdirAll(contractsDir, 0755)

	// Write a sample spec
	specContent := `# Test Project Spec

## Overview
A test project for integration testing.

## Features
- Feature 1
- Feature 2
`
	os.WriteFile(filepath.Join(contractsDir, "spec.md"), []byte(specContent), 0644)

	// Verify spec file exists
	if _, err := os.Stat(filepath.Join(contractsDir, "spec.md")); os.IsNotExist(err) {
		t.Error("spec file was not created")
	}

	// TODO: Add full workflow tests when modes are fully implemented
	// - Run intake mode
	// - Run review mode
	// - Run change-order mode
}

func TestFixturesExist(t *testing.T) {
	fixtures := []string{
		"../fixtures/sample_spec.md",
		"../fixtures/sample_config.toml",
	}

	for _, f := range fixtures {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Logf("fixture %s not found (expected for initial setup)", f)
		}
	}
}
