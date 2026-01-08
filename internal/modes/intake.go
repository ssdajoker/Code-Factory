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

// IntakeStep represents a step in the intake interview
type IntakeStep int

const (
	StepProjectName IntakeStep = iota
	StepDescription
	StepTargetUsers
	StepCoreFeatures
	StepTechnicalConstraints
	StepSuccessCriteria
	StepPreview
	StepComplete
)

// IntakeData holds all collected information
type IntakeData struct {
	ProjectName          string
	Description          string
	TargetUsers          string
	CoreFeatures         string
	TechnicalConstraints string
	SuccessCriteria      string
	GeneratedSpec        string
}

// IntakeMode handles the INTAKE workflow
type IntakeMode struct {
	provider    llm.Provider
	data        IntakeData
	currentStep IntakeStep
	contractsDir string
}

// NewIntakeMode creates a new intake mode
func NewIntakeMode(provider llm.Provider, contractsDir string) *IntakeMode {
	if contractsDir == "" {
		contractsDir = "contracts"
	}
	return &IntakeMode{
		provider:     provider,
		contractsDir: contractsDir,
	}
}

// StepCount returns total number of steps
func (m *IntakeMode) StepCount() int {
	return int(StepPreview)
}

// CurrentStep returns the current step
func (m *IntakeMode) CurrentStep() IntakeStep {
	return m.currentStep
}

// StepTitle returns the title for a step
func (m *IntakeMode) StepTitle(step IntakeStep) string {
	titles := map[IntakeStep]string{
		StepProjectName:          "Project Name",
		StepDescription:          "Project Description",
		StepTargetUsers:          "Target Users",
		StepCoreFeatures:         "Core Features",
		StepTechnicalConstraints: "Technical Constraints",
		StepSuccessCriteria:      "Success Criteria",
		StepPreview:              "Preview Specification",
	}
	return titles[step]
}

// StepPrompt returns the prompt for a step
func (m *IntakeMode) StepPrompt(step IntakeStep) string {
	prompts := map[IntakeStep]string{
		StepProjectName:          "What is the name of your project?",
		StepDescription:          "Describe your project in a few sentences. What problem does it solve?",
		StepTargetUsers:          "Who are the target users? What are their needs?",
		StepCoreFeatures:         "List the core features (one per line):",
		StepTechnicalConstraints: "Any technical constraints? (language, framework, integrations, etc.)",
		StepSuccessCriteria:      "How will you measure success? What are the acceptance criteria?",
	}
	return prompts[step]
}

// SetStepValue sets the value for the current step
func (m *IntakeMode) SetStepValue(value string) {
	switch m.currentStep {
	case StepProjectName:
		m.data.ProjectName = value
	case StepDescription:
		m.data.Description = value
	case StepTargetUsers:
		m.data.TargetUsers = value
	case StepCoreFeatures:
		m.data.CoreFeatures = value
	case StepTechnicalConstraints:
		m.data.TechnicalConstraints = value
	case StepSuccessCriteria:
		m.data.SuccessCriteria = value
	}
}

// GetStepValue gets the current value for a step
func (m *IntakeMode) GetStepValue(step IntakeStep) string {
	switch step {
	case StepProjectName:
		return m.data.ProjectName
	case StepDescription:
		return m.data.Description
	case StepTargetUsers:
		return m.data.TargetUsers
	case StepCoreFeatures:
		return m.data.CoreFeatures
	case StepTechnicalConstraints:
		return m.data.TechnicalConstraints
	case StepSuccessCriteria:
		return m.data.SuccessCriteria
	default:
		return ""
	}
}

// NextStep advances to the next step
func (m *IntakeMode) NextStep() {
	if m.currentStep < StepComplete {
		m.currentStep++
	}
}

// PrevStep goes back to the previous step
func (m *IntakeMode) PrevStep() {
	if m.currentStep > StepProjectName {
		m.currentStep--
	}
}

// GenerateSpec uses LLM to expand inputs into a full specification
func (m *IntakeMode) GenerateSpec(ctx context.Context) (string, error) {
	if m.provider == nil {
		// Fallback to template-based generation
		return m.generateTemplateSpec(), nil
	}

	prompt := fmt.Sprintf(`Based on the following project information, generate a comprehensive software specification document in Markdown format.

Project Name: %s
Description: %s
Target Users: %s
Core Features:
%s
Technical Constraints: %s
Success Criteria: %s

Generate a professional specification document with the following sections:
1. Executive Summary
2. Problem Statement
3. Target Audience
4. Functional Requirements (expand the core features into detailed requirements)
5. Non-Functional Requirements
6. Technical Architecture (based on constraints)
7. Success Metrics
8. Out of Scope
9. Risks and Mitigations

Be specific and actionable. Use clear, professional language.`,
		m.data.ProjectName,
		m.data.Description,
		m.data.TargetUsers,
		m.data.CoreFeatures,
		m.data.TechnicalConstraints,
		m.data.SuccessCriteria,
	)

	opts := llm.DefaultOptions()
	opts.SystemPrompt = "You are a senior software architect creating detailed technical specifications. Be thorough but concise."
	opts.MaxTokens = 4096

	spec, err := m.provider.Complete(ctx, prompt, opts)
	if err != nil {
		// Fallback to template
		return m.generateTemplateSpec(), nil
	}

	m.data.GeneratedSpec = spec
	return spec, nil
}

func (m *IntakeMode) generateTemplateSpec() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s - Vision Specification\n\n", m.data.ProjectName))
	sb.WriteString(fmt.Sprintf("*Generated: %s*\n\n", time.Now().Format("2006-01-02 15:04")))

	sb.WriteString("## Executive Summary\n\n")
	sb.WriteString(m.data.Description + "\n\n")

	sb.WriteString("## Target Audience\n\n")
	sb.WriteString(m.data.TargetUsers + "\n\n")

	sb.WriteString("## Core Features\n\n")
	for _, feature := range strings.Split(m.data.CoreFeatures, "\n") {
		feature = strings.TrimSpace(feature)
		if feature != "" {
			sb.WriteString(fmt.Sprintf("- %s\n", feature))
		}
	}
	sb.WriteString("\n")

	sb.WriteString("## Technical Constraints\n\n")
	sb.WriteString(m.data.TechnicalConstraints + "\n\n")

	sb.WriteString("## Success Criteria\n\n")
	sb.WriteString(m.data.SuccessCriteria + "\n\n")

	m.data.GeneratedSpec = sb.String()
	return m.data.GeneratedSpec
}

// SaveSpec saves the specification to the contracts directory
func (m *IntakeMode) SaveSpec() (string, error) {
	if m.data.GeneratedSpec == "" {
		return "", fmt.Errorf("no specification generated")
	}

	if err := os.MkdirAll(m.contractsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create contracts directory: %w", err)
	}

	filename := filepath.Join(m.contractsDir, "vision_spec.md")
	if err := os.WriteFile(filename, []byte(m.data.GeneratedSpec), 0644); err != nil {
		return "", fmt.Errorf("failed to write spec: %w", err)
	}

	return filename, nil
}

// Data returns the collected intake data
func (m *IntakeMode) Data() IntakeData {
	return m.data
}
