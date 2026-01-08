package modes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ssdajoker/Code-Factory/internal/llm"
)

// RescueResult holds the rescue analysis results
type RescueResult struct {
	CodebasePath    string
	FilesScanned    int
	InferredSpec    string
	AlignmentReport string
	Architecture    string
	Patterns        []string
	Dependencies    []string
}

// RescueMode handles the RESCUE workflow
type RescueMode struct {
	provider     llm.Provider
	contractsDir string
	reportsDir   string
	result       RescueResult
}

// NewRescueMode creates a new rescue mode
func NewRescueMode(provider llm.Provider, contractsDir, reportsDir string) *RescueMode {
	if contractsDir == "" {
		contractsDir = "contracts"
	}
	if reportsDir == "" {
		reportsDir = "reports"
	}
	return &RescueMode{
		provider:     provider,
		contractsDir: contractsDir,
		reportsDir:   reportsDir,
	}
}

// SetCodebasePath sets the codebase to analyze
func (m *RescueMode) SetCodebasePath(path string) {
	m.result.CodebasePath = path
}

// ScanCodebase scans and analyzes the codebase
func (m *RescueMode) ScanCodebase(ctx context.Context) (*RescueResult, error) {
	var codeContent strings.Builder
	var fileList []string

	err := filepath.Walk(m.result.CodebasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		// Skip hidden and vendor directories
		if strings.Contains(path, "/.git/") || strings.Contains(path, "/vendor/") || strings.Contains(path, "/node_modules/") {
			return nil
		}
		if isCodeFile(path) || isConfigFile(path) {
			fileList = append(fileList, path)
			data, _ := os.ReadFile(path)
			// Limit file size
			if len(data) < 10000 {
				codeContent.WriteString(fmt.Sprintf("\n--- %s ---\n%s\n", path, string(data)))
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan codebase: %w", err)
	}

	m.result.FilesScanned = len(fileList)

	if m.provider == nil {
		return m.generateTemplateRescue(fileList), nil
	}

	prompt := fmt.Sprintf(`Analyze this codebase and reverse-engineer a specification document.

CODEBASE:
%s

Generate:
1. A comprehensive specification document inferring the project's purpose, architecture, and requirements
2. An alignment report showing what was discovered

Format both as Markdown, separated by "---ALIGNMENT---"`, codeContent.String())

	opts := llm.DefaultOptions()
	opts.SystemPrompt = "You are a software architect reverse-engineering specifications from code."
	opts.MaxTokens = 8192

	response, err := m.provider.Complete(ctx, prompt, opts)
	if err != nil {
		return m.generateTemplateRescue(fileList), nil
	}

	parts := strings.Split(response, "---ALIGNMENT---")
	m.result.InferredSpec = parts[0]
	if len(parts) > 1 {
		m.result.AlignmentReport = parts[1]
	} else {
		m.result.AlignmentReport = m.generateAlignmentReport(fileList)
	}

	return &m.result, nil
}

func (m *RescueMode) generateTemplateRescue(files []string) *RescueResult {
	var sb strings.Builder
	sb.WriteString("# Inferred Specification\n\n")
	sb.WriteString(fmt.Sprintf("*Generated: %s*\n\n", time.Now().Format("2006-01-02 15:04")))
	sb.WriteString(fmt.Sprintf("**Codebase:** %s\n\n", m.result.CodebasePath))
	sb.WriteString(fmt.Sprintf("**Files Analyzed:** %d\n\n", len(files)))
	sb.WriteString("## Discovered Structure\n\n")
	for _, f := range files {
		sb.WriteString(fmt.Sprintf("- %s\n", f))
	}
	sb.WriteString("\n## Architecture\n\n")
	sb.WriteString("*Configure LLM for detailed architecture analysis*\n\n")
	sb.WriteString("## Inferred Requirements\n\n")
	sb.WriteString("*Configure LLM for requirement inference*\n")

	m.result.InferredSpec = sb.String()
	m.result.AlignmentReport = m.generateAlignmentReport(files)
	return &m.result
}

func (m *RescueMode) generateAlignmentReport(files []string) string {
	var sb strings.Builder
	sb.WriteString("# Alignment Report\n\n")
	sb.WriteString(fmt.Sprintf("*Generated: %s*\n\n", time.Now().Format("2006-01-02 15:04")))
	sb.WriteString("## Discovery Summary\n\n")
	sb.WriteString(fmt.Sprintf("- Total files scanned: %d\n", len(files)))
	sb.WriteString("- Code patterns detected: Manual review needed\n")
	sb.WriteString("- Dependencies found: Check go.mod/package.json\n\n")
	sb.WriteString("## Recommendations\n\n")
	sb.WriteString("- Review inferred spec for accuracy\n")
	sb.WriteString("- Add missing requirements\n")
	sb.WriteString("- Configure LLM for deeper analysis\n")
	return sb.String()
}

// SaveResults saves both the spec and alignment report
func (m *RescueMode) SaveResults() (specPath, reportPath string, err error) {
	if m.result.InferredSpec == "" {
		return "", "", fmt.Errorf("no spec generated")
	}

	if err := os.MkdirAll(m.contractsDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create contracts directory: %w", err)
	}
	if err := os.MkdirAll(m.reportsDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create reports directory: %w", err)
	}

	specPath = filepath.Join(m.contractsDir, "current_spec.md")
	if err := os.WriteFile(specPath, []byte(m.result.InferredSpec), 0644); err != nil {
		return "", "", fmt.Errorf("failed to write spec: %w", err)
	}

	reportPath = filepath.Join(m.reportsDir, "alignment_report.md")
	if err := os.WriteFile(reportPath, []byte(m.result.AlignmentReport), 0644); err != nil {
		return specPath, "", fmt.Errorf("failed to write report: %w", err)
	}

	return specPath, reportPath, nil
}

// Result returns the rescue result
func (m *RescueMode) Result() *RescueResult {
	return &m.result
}

func isConfigFile(path string) bool {
	base := filepath.Base(path)
	configs := []string{"go.mod", "go.sum", "package.json", "Cargo.toml", "requirements.txt", "Makefile", "Dockerfile"}
	for _, c := range configs {
		if base == c {
			return true
		}
	}
	return false
}
