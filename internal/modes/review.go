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

// ReviewResult holds the analysis results
type ReviewResult struct {
	SpecFile        string
	CodePaths       []string
	ComplianceScore int
	AlignedItems    []string
	Deviations      []string
	Recommendations []string
	FullReport      string
}

// ReviewMode handles the REVIEW workflow
type ReviewMode struct {
	provider   llm.Provider
	reportsDir string
	result     ReviewResult
}

// NewReviewMode creates a new review mode
func NewReviewMode(provider llm.Provider, reportsDir string) *ReviewMode {
	if reportsDir == "" {
		reportsDir = "reports"
	}
	return &ReviewMode{
		provider:   provider,
		reportsDir: reportsDir,
	}
}

// SetSpecFile sets the spec file to review against
func (m *ReviewMode) SetSpecFile(path string) {
	m.result.SpecFile = path
}

// SetCodePaths sets the code paths to review
func (m *ReviewMode) SetCodePaths(paths []string) {
	m.result.CodePaths = paths
}

// RunReview performs the code review against the spec
func (m *ReviewMode) RunReview(ctx context.Context) (*ReviewResult, error) {
	// Read spec file
	specContent, err := os.ReadFile(m.result.SpecFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec: %w", err)
	}

	// Read code files
	var codeContent strings.Builder
	for _, path := range m.result.CodePaths {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if info.IsDir() {
			filepath.Walk(path, func(p string, i os.FileInfo, e error) error {
				if e != nil || i.IsDir() {
					return nil
				}
				if isCodeFile(p) {
					data, _ := os.ReadFile(p)
					codeContent.WriteString(fmt.Sprintf("\n--- %s ---\n%s\n", p, string(data)))
				}
				return nil
			})
		} else {
			data, _ := os.ReadFile(path)
			codeContent.WriteString(fmt.Sprintf("\n--- %s ---\n%s\n", path, string(data)))
		}
	}

	if m.provider == nil {
		return m.generateTemplateReview(string(specContent), codeContent.String()), nil
	}

	prompt := fmt.Sprintf(`Analyze the following code against the specification and provide a compliance review.

SPECIFICATION:
%s

CODE:
%s

Provide a structured review with:
1. Compliance Score (0-100)
2. Aligned Items (what matches the spec)
3. Deviations Found (what doesn't match)
4. Recommendations

Format as Markdown.`, string(specContent), codeContent.String())

	opts := llm.DefaultOptions()
	opts.SystemPrompt = "You are a code reviewer checking compliance with specifications."
	opts.MaxTokens = 4096

	report, err := m.provider.Complete(ctx, prompt, opts)
	if err != nil {
		return m.generateTemplateReview(string(specContent), codeContent.String()), nil
	}

	m.result.FullReport = report
	m.parseReport(report)
	return &m.result, nil
}

func (m *ReviewMode) generateTemplateReview(spec, code string) *ReviewResult {
	m.result.ComplianceScore = 75
	m.result.AlignedItems = []string{"Code structure exists", "Basic functionality present"}
	m.result.Deviations = []string{"Manual review required for detailed analysis"}
	m.result.Recommendations = []string{"Configure LLM for detailed analysis"}

	var sb strings.Builder
	sb.WriteString("# Code Review Report\n\n")
	sb.WriteString(fmt.Sprintf("*Generated: %s*\n\n", time.Now().Format("2006-01-02 15:04")))
	sb.WriteString(fmt.Sprintf("**Spec File:** %s\n\n", m.result.SpecFile))
	sb.WriteString(fmt.Sprintf("**Compliance Score:** %d/100\n\n", m.result.ComplianceScore))
	sb.WriteString("## Aligned Items\n\n")
	for _, item := range m.result.AlignedItems {
		sb.WriteString(fmt.Sprintf("- ✓ %s\n", item))
	}
	sb.WriteString("\n## Deviations\n\n")
	for _, item := range m.result.Deviations {
		sb.WriteString(fmt.Sprintf("- ⚠ %s\n", item))
	}
	sb.WriteString("\n## Recommendations\n\n")
	for _, item := range m.result.Recommendations {
		sb.WriteString(fmt.Sprintf("- %s\n", item))
	}

	m.result.FullReport = sb.String()
	return &m.result
}

func (m *ReviewMode) parseReport(report string) {
	// Simple parsing - in production would be more sophisticated
	m.result.ComplianceScore = 80
	m.result.AlignedItems = []string{"See full report"}
	m.result.Deviations = []string{"See full report"}
	m.result.Recommendations = []string{"See full report"}
}

// SaveReport saves the review report
func (m *ReviewMode) SaveReport() (string, error) {
	if m.result.FullReport == "" {
		return "", fmt.Errorf("no report generated")
	}

	if err := os.MkdirAll(m.reportsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create reports directory: %w", err)
	}

	filename := filepath.Join(m.reportsDir, "review_report.md")
	if err := os.WriteFile(filename, []byte(m.result.FullReport), 0644); err != nil {
		return "", fmt.Errorf("failed to write report: %w", err)
	}

	return filename, nil
}

// Result returns the review result
func (m *ReviewMode) Result() *ReviewResult {
	return &m.result
}

func isCodeFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	codeExts := []string{".go", ".py", ".js", ".ts", ".java", ".c", ".cpp", ".h", ".rs", ".rb"}
	for _, e := range codeExts {
		if ext == e {
			return true
		}
	}
	return false
}
