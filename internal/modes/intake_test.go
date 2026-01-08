package modes

import (
	"context"
	"strings"
	"testing"
)

func TestNewIntakeMode(t *testing.T) {
	tests := []struct {
		name         string
		contractsDir string
		wantDir      string
	}{
		{"default dir", "", "contracts"},
		{"custom dir", "specs", "specs"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewIntakeMode(nil, tt.contractsDir)
			if m.contractsDir != tt.wantDir {
				t.Errorf("contractsDir = %v, want %v", m.contractsDir, tt.wantDir)
			}
		})
	}
}

func TestIntakeSteps(t *testing.T) {
	m := NewIntakeMode(nil, "")

	// Test step count
	if m.StepCount() != int(StepPreview) {
		t.Errorf("StepCount() = %d, want %d", m.StepCount(), int(StepPreview))
	}

	// Test initial step
	if m.CurrentStep() != StepProjectName {
		t.Errorf("initial step = %v, want StepProjectName", m.CurrentStep())
	}

	// Test step titles
	titles := []struct {
		step  IntakeStep
		title string
	}{
		{StepProjectName, "Project Name"},
		{StepDescription, "Project Description"},
		{StepTargetUsers, "Target Users"},
		{StepCoreFeatures, "Core Features"},
		{StepTechnicalConstraints, "Technical Constraints"},
		{StepSuccessCriteria, "Success Criteria"},
		{StepPreview, "Preview Specification"},
	}

	for _, tt := range titles {
		if got := m.StepTitle(tt.step); got != tt.title {
			t.Errorf("StepTitle(%v) = %q, want %q", tt.step, got, tt.title)
		}
	}
}

func TestIntakeNavigation(t *testing.T) {
	m := NewIntakeMode(nil, "")

	// Test NextStep
	m.NextStep()
	if m.CurrentStep() != StepDescription {
		t.Errorf("after NextStep, step = %v, want StepDescription", m.CurrentStep())
	}

	// Test PrevStep
	m.PrevStep()
	if m.CurrentStep() != StepProjectName {
		t.Errorf("after PrevStep, step = %v, want StepProjectName", m.CurrentStep())
	}

	// Test PrevStep at beginning (should stay)
	m.PrevStep()
	if m.CurrentStep() != StepProjectName {
		t.Errorf("PrevStep at start should stay at StepProjectName")
	}
}

func TestIntakeSetGetValue(t *testing.T) {
	m := NewIntakeMode(nil, "")

	tests := []struct {
		step  IntakeStep
		value string
	}{
		{StepProjectName, "MyProject"},
		{StepDescription, "A test project"},
		{StepTargetUsers, "Developers"},
		{StepCoreFeatures, "Feature 1\nFeature 2"},
		{StepTechnicalConstraints, "Go, PostgreSQL"},
		{StepSuccessCriteria, "100% test coverage"},
	}

	for _, tt := range tests {
		m.currentStep = tt.step
		m.SetStepValue(tt.value)
		if got := m.GetStepValue(tt.step); got != tt.value {
			t.Errorf("GetStepValue(%v) = %q, want %q", tt.step, got, tt.value)
		}
	}
}

func TestIntakeStepPrompts(t *testing.T) {
	m := NewIntakeMode(nil, "")

	steps := []IntakeStep{
		StepProjectName,
		StepDescription,
		StepTargetUsers,
		StepCoreFeatures,
		StepTechnicalConstraints,
		StepSuccessCriteria,
	}

	for _, step := range steps {
		prompt := m.StepPrompt(step)
		if prompt == "" {
			t.Errorf("StepPrompt(%v) returned empty string", step)
		}
	}
}

func TestGenerateSpecWithoutProvider(t *testing.T) {
	m := NewIntakeMode(nil, "")
	m.data = IntakeData{
		ProjectName:          "TestProject",
		Description:          "A test project",
		TargetUsers:          "Developers",
		CoreFeatures:         "Feature 1\nFeature 2",
		TechnicalConstraints: "Go",
		SuccessCriteria:      "Works",
	}

	spec, err := m.GenerateSpec(context.Background())
	if err != nil {
		t.Fatalf("GenerateSpec() error = %v", err)
	}

	// Should contain project name in template output
	if !strings.Contains(spec, "TestProject") {
		t.Error("spec should contain project name")
	}
}
