package views

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ssdajoker/Code-Factory/internal/llm"
	"github.com/ssdajoker/Code-Factory/internal/modes"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle.Copy()
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

// IntakeView is the TUI view for INTAKE mode
type IntakeView struct {
	intake      *modes.IntakeMode
	textInput   textinput.Model
	textArea    textarea.Model
	spinner     spinner.Model
	width       int
	height      int
	generating  bool
	previewMode bool
	saved       bool
	savePath    string
	err         error
	llmStatus   string
}

// specGeneratedMsg is sent when spec generation completes
type specGeneratedMsg struct {
	spec string
	err  error
}

// specSavedMsg is sent when spec is saved
type specSavedMsg struct {
	path string
	err  error
}

// NewIntakeView creates a new intake view
func NewIntakeView(provider llm.Provider, contractsDir string) IntakeView {
	ti := textinput.New()
	ti.Placeholder = "Enter value..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60

	ta := textarea.New()
	ta.Placeholder = "Enter multiple lines..."
	ta.SetWidth(60)
	ta.SetHeight(6)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	llmStatus := "No LLM"
	if provider != nil {
		llmStatus = provider.Name()
	}

	return IntakeView{
		intake:    modes.NewIntakeMode(provider, contractsDir),
		textInput: ti,
		textArea:  ta,
		spinner:   s,
		llmStatus: llmStatus,
	}
}

// Init implements tea.Model
func (v IntakeView) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (v IntakeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
		return v, nil

	case spinner.TickMsg:
		if v.generating {
			var cmd tea.Cmd
			v.spinner, cmd = v.spinner.Update(msg)
			return v, cmd
		}

	case specGeneratedMsg:
		v.generating = false
		if msg.err != nil {
			v.err = msg.err
		} else {
			v.previewMode = true
		}
		return v, nil

	case specSavedMsg:
		if msg.err != nil {
			v.err = msg.err
		} else {
			v.saved = true
			v.savePath = msg.path
		}
		return v, nil

	case tea.KeyMsg:
		return v.handleKey(msg)
	}

	// Update text input or textarea
	if v.useTextArea() {
		var cmd tea.Cmd
		v.textArea, cmd = v.textArea.Update(msg)
		cmds = append(cmds, cmd)
	} else if !v.previewMode && !v.generating {
		var cmd tea.Cmd
		v.textInput, cmd = v.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return v, tea.Batch(cmds...)
}

func (v IntakeView) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return v, tea.Quit

	case "esc":
		if v.previewMode && !v.saved {
			v.previewMode = false
			v.intake.PrevStep()
			return v, nil
		}
		return v, nil

	case "enter":
		if v.generating {
			return v, nil
		}

		if v.saved {
			return v, tea.Quit
		}

		if v.previewMode {
			// Save the spec
			return v, v.saveSpec()
		}

		// For textarea, only advance on ctrl+enter
		if v.useTextArea() {
			return v, nil
		}

		return v.advanceStep()

	case "ctrl+enter":
		if v.useTextArea() && !v.generating && !v.previewMode {
			return v.advanceStep()
		}

	case "ctrl+b", "left":
		if !v.generating && !v.previewMode && v.intake.CurrentStep() > modes.StepProjectName {
			v.intake.PrevStep()
			v.loadCurrentValue()
		}
	}

	return v, nil
}

func (v *IntakeView) advanceStep() (tea.Model, tea.Cmd) {
	// Save current value
	var value string
	if v.useTextArea() {
		value = v.textArea.Value()
	} else {
		value = v.textInput.Value()
	}
	v.intake.SetStepValue(value)

	// Move to next step
	v.intake.NextStep()

	// Check if we need to generate spec
	if v.intake.CurrentStep() == modes.StepPreview {
		v.generating = true
		return *v, tea.Batch(v.spinner.Tick, v.generateSpec())
	}

	// Reset input for next step
	v.loadCurrentValue()
	return *v, nil
}

func (v *IntakeView) loadCurrentValue() {
	value := v.intake.GetStepValue(v.intake.CurrentStep())
	if v.useTextArea() {
		v.textArea.SetValue(value)
	} else {
		v.textInput.SetValue(value)
	}
}

func (v IntakeView) useTextArea() bool {
	step := v.intake.CurrentStep()
	return step == modes.StepCoreFeatures || step == modes.StepTechnicalConstraints || step == modes.StepSuccessCriteria
}

func (v IntakeView) generateSpec() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		spec, err := v.intake.GenerateSpec(ctx)
		return specGeneratedMsg{spec: spec, err: err}
	}
}

func (v IntakeView) saveSpec() tea.Cmd {
	return func() tea.Msg {
		path, err := v.intake.SaveSpec()
		return specSavedMsg{path: path, err: err}
	}
}

// View implements tea.Model
func (v IntakeView) View() string {
	var sb strings.Builder

	// Header
	sb.WriteString(titleStyle.Render("📝 INTAKE MODE - Capture Your Vision"))
	sb.WriteString("\n")
	sb.WriteString(blurredStyle.Render("LLM: " + v.llmStatus))
	sb.WriteString("\n\n")

	if v.saved {
		sb.WriteString(successStyle.Render("✓ Specification saved!"))
		sb.WriteString("\n")
		sb.WriteString("Path: " + v.savePath)
		sb.WriteString("\n\n")
		sb.WriteString(helpStyle.Render("Press Enter to exit"))
		return sb.String()
	}

	if v.err != nil {
		sb.WriteString(errorStyle.Render("Error: " + v.err.Error()))
		sb.WriteString("\n")
	}

	if v.generating {
		sb.WriteString(v.spinner.View())
		sb.WriteString(" Generating specification...")
		return sb.String()
	}

	if v.previewMode {
		return v.viewPreview()
	}

	// Progress indicator
	step := int(v.intake.CurrentStep()) + 1
	total := v.intake.StepCount()
	sb.WriteString(focusedStyle.Render(strings.Repeat("●", step)))
	sb.WriteString(blurredStyle.Render(strings.Repeat("○", total-step)))
	sb.WriteString(blurredStyle.Render(" Step " + string(rune('0'+step)) + "/" + string(rune('0'+total))))
	sb.WriteString("\n\n")

	// Current step title and prompt
	sb.WriteString(titleStyle.Render(v.intake.StepTitle(v.intake.CurrentStep())))
	sb.WriteString("\n")
	sb.WriteString(v.intake.StepPrompt(v.intake.CurrentStep()))
	sb.WriteString("\n\n")

	// Input
	if v.useTextArea() {
		sb.WriteString(v.textArea.View())
		sb.WriteString("\n\n")
		sb.WriteString(helpStyle.Render("Ctrl+Enter: next • Ctrl+B: back • Ctrl+C: quit"))
	} else {
		sb.WriteString(v.textInput.View())
		sb.WriteString("\n\n")
		sb.WriteString(helpStyle.Render("Enter: next • ←: back • Ctrl+C: quit"))
	}

	return sb.String()
}

func (v IntakeView) viewPreview() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("📄 Preview Specification"))
	sb.WriteString("\n\n")

	// Show truncated preview
	spec := v.intake.Data().GeneratedSpec
	lines := strings.Split(spec, "\n")
	maxLines := 20
	if len(lines) > maxLines {
		for _, line := range lines[:maxLines] {
			sb.WriteString(line + "\n")
		}
		sb.WriteString(blurredStyle.Render("... (truncated)"))
	} else {
		sb.WriteString(spec)
	}

	sb.WriteString("\n\n")
	sb.WriteString(helpStyle.Render("Enter: save • Esc: edit • Ctrl+C: quit"))

	return sb.String()
}
