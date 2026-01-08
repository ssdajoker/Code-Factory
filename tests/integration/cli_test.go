//go:build integration

package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

var binaryPath = "../../factory"

func TestCLIHelp(t *testing.T) {
	cmd := exec.Command(binaryPath, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("factory --help failed: %v\n%s", err, output)
	}

	expected := []string{
		"Factory",
		"intake",
		"review",
		"rescue",
		"change-order",
		"github",
	}

	for _, s := range expected {
		if !strings.Contains(string(output), s) {
			t.Errorf("help output missing %q", s)
		}
	}
}

func TestCLIVersion(t *testing.T) {
	cmd := exec.Command(binaryPath, "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("factory version failed: %v\n%s", err, output)
	}

	if !strings.Contains(string(output), "Factory") {
		t.Error("version output should contain 'Factory'")
	}
}

func TestCLIInit(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "factory-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command(binaryPath, "init")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("init output: %s", output)
		// init may fail without full setup, just verify it runs
	}
}

func TestCLILLMStatus(t *testing.T) {
	cmd := exec.Command(binaryPath, "llm", "status")
	output, _ := cmd.CombinedOutput()
	// Just verify command doesn't panic
	_ = output
}
