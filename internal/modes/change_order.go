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

// ChangeItem represents a detected change
type ChangeItem struct {
	ID          string
	Description string
	SpecSection string
	CodePath    string
	Status      string // "pending", "approved", "rejected", "deferred"
	Reason      string
}

// ChangeOrderResult holds the change order analysis
type ChangeOrderResult struct {
	SpecFile     string
	CodebasePath string
	Changes      []ChangeItem
	FullReport   string
}

// ChangeOrderMode handles the CHANGE_ORDER workflow
type ChangeOrderMode struct {
	provider     llm.Provider
	contractsDir string
	result       ChangeOrderResult
}

// NewChangeOrderMode creates a new change order mode
func NewChangeOrderMode(provider llm.Provider, contractsDir string) *ChangeOrderMode {
	if contractsDir == "" {
		contractsDir = "contracts"
	}
	return &ChangeOrderMode{
		provider:     provider,
		contractsDir: contractsDir,
	}
}

// SetSpecFile sets the spec file to compare against
func (m *ChangeOrderMode) SetSpecFile(path string) {
	m.result.SpecFile = path
}

// SetCodebasePath sets the codebase path
func (m *ChangeOrderMode) SetCodebasePath(path string) {
	m.result.CodebasePath = path
}

// DetectDrift analyzes spec vs code for intentional drift
func (m *ChangeOrderMode) DetectDrift(ctx context.Context) (*ChangeOrderResult, error) {
	specContent, err := os.ReadFile(m.result.SpecFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec: %w", err)
	}

	var codeContent strings.Builder
	filepath.Walk(m.result.CodebasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.Contains(path, "/.git/") || strings.Contains(path, "/vendor/") {
			return nil
		}
		if isCodeFile(path) {
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to read file %s: %v\n", path, err)
				return nil
			}
			if len(data) < 10000 {
				codeContent.WriteString(fmt.Sprintf("\n--- %s ---\n%s\n", path, string(data)))
			}
		}
		return nil
	})

	if m.provider == nil {
		return m.generateTemplateChangeOrder(string(specContent)), nil
	}

	prompt := fmt.Sprintf(`Compare this specification with the codebase and identify intentional drift (changes that deviate from spec).

SPECIFICATION:
%s

CODEBASE:
%s

List each deviation as:
- ID: CO-XXX
- Description: What changed
- Spec Section: Which part of spec it affects
- Code Path: Where in code

Format as Markdown.`, string(specContent), codeContent.String())

	opts := llm.DefaultOptions()
	opts.SystemPrompt = "You are analyzing code drift from specifications."
	opts.MaxTokens = 4096

	report, err := m.provider.Complete(ctx, prompt, opts)
	if err != nil {
		return m.generateTemplateChangeOrder(string(specContent)), nil
	}

	m.result.FullReport = report
	m.parseChanges(report)
	return &m.result, nil
}

func (m *ChangeOrderMode) generateTemplateChangeOrder(spec string) *ChangeOrderResult {
	m.result.Changes = []ChangeItem{
		{
			ID:          "CO-001",
			Description: "Manual review required",
			SpecSection: "All",
			CodePath:    m.result.CodebasePath,
			Status:      "pending",
		},
	}

	var sb strings.Builder
	sb.WriteString("# Change Order Report\n\n")
	sb.WriteString(fmt.Sprintf("*Generated: %s*\n\n", time.Now().Format("2006-01-02 15:04")))
	sb.WriteString(fmt.Sprintf("**Spec:** %s\n", m.result.SpecFile))
	sb.WriteString(fmt.Sprintf("**Codebase:** %s\n\n", m.result.CodebasePath))
	sb.WriteString("## Detected Changes\n\n")
	for _, c := range m.result.Changes {
		sb.WriteString(fmt.Sprintf("### %s\n", c.ID))
		sb.WriteString(fmt.Sprintf("- **Description:** %s\n", c.Description))
		sb.WriteString(fmt.Sprintf("- **Spec Section:** %s\n", c.SpecSection))
		sb.WriteString(fmt.Sprintf("- **Status:** %s\n\n", c.Status))
	}
	sb.WriteString("\n*Configure LLM for detailed drift analysis*\n")

	m.result.FullReport = sb.String()
	return &m.result
}

func (m *ChangeOrderMode) parseChanges(report string) {
	// Simple parsing - would be more sophisticated in production
	m.result.Changes = []ChangeItem{
		{ID: "CO-001", Description: "See full report", Status: "pending"},
	}
}

// ApproveChange approves a change
func (m *ChangeOrderMode) ApproveChange(id, reason string) {
	for i := range m.result.Changes {
		if m.result.Changes[i].ID == id {
			m.result.Changes[i].Status = "approved"
			m.result.Changes[i].Reason = reason
		}
	}
}

// RejectChange rejects a change
func (m *ChangeOrderMode) RejectChange(id, reason string) {
	for i := range m.result.Changes {
		if m.result.Changes[i].ID == id {
			m.result.Changes[i].Status = "rejected"
			m.result.Changes[i].Reason = reason
		}
	}
}

// DeferChange defers a change
func (m *ChangeOrderMode) DeferChange(id, reason string) {
	for i := range m.result.Changes {
		if m.result.Changes[i].ID == id {
			m.result.Changes[i].Status = "deferred"
			m.result.Changes[i].Reason = reason
		}
	}
}

// SaveChangeOrder saves or appends to change_order.md
func (m *ChangeOrderMode) SaveChangeOrder() (string, error) {
	if err := os.MkdirAll(m.contractsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create contracts directory: %w", err)
	}

	filename := filepath.Join(m.contractsDir, "change_order.md")

	// Append if exists
	var content string
	if existing, err := os.ReadFile(filename); err == nil {
		content = string(existing) + "\n\n---\n\n" + m.result.FullReport
	} else {
		content = m.result.FullReport
	}

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write change order: %w", err)
	}

	return filename, nil
}

// Result returns the change order result
func (m *ChangeOrderMode) Result() *ChangeOrderResult {
	return &m.result
}
