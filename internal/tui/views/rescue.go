package views

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ssdajoker/Code-Factory/internal/llm"
	"github.com/ssdajoker/Code-Factory/internal/modes"
)

// RescueStep represents steps in rescue flow
type RescueStep int

const (
	RescueStepSelectDir RescueStep = iota
	RescueStepScanning
	RescueStepPreview
	RescueStepSaved
)

// RescueView is the TUI view for RESCUE mode
type RescueView struct {
	rescue     *modes.RescueMode
	filePicker filepicker.Model
	spinner    spinner.Model
	step       RescueStep
	codeDir    string
	width      int
	height     int
	saved      bool
	specPath   string
	reportPath string
	err        error
	llmStatus  string
}

type rescueDoneMsg struct {
	result *modes.RescueResult
	err    error
}

type rescueSavedMsg struct {
	specPath   string
	reportPath string
	err        error
}

// NewRescueView creates a new rescue view
func NewRescueView(provider llm.Provider, contractsDir, reportsDir string) RescueView {
	fp := filepicker.New()
	fp.CurrentDirectory, _ = os.Getwd()
	fp.DirAllowed = true
	fp.FileAllowed = false

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	llmStatus := "No LLM"
	if provider != nil {
		llmStatus = provider.Name()
	}

	return RescueView{
		rescue:    modes.NewRescueMode(provider, contractsDir, reportsDir),
		filePicker: fp,
		spinner:   s,
		step:      RescueStepSelectDir,
		llmStatus: llmStatus,
	}
}

// Init implements tea.Model
func (v RescueView) Init() tea.Cmd {
	return v.filePicker.Init()
}

// Update implements tea.Model
func (v RescueView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
		return v, nil

	case spinner.TickMsg:
		if v.step == RescueStepScanning {
			var cmd tea.Cmd
			v.spinner, cmd = v.spinner.Update(msg)
			return v, cmd
		}

	case rescueDoneMsg:
		if msg.err != nil {
			v.err = msg.err
			v.step = RescueStepSelectDir
		} else {
			v.step = RescueStepPreview
		}
		return v, nil

	case rescueSavedMsg:
		if msg.err != nil {
			v.err = msg.err
		} else {
			v.saved = true
			v.specPath = msg.specPath
			v.reportPath = msg.reportPath
			v.step = RescueStepSaved
		}
		return v, nil

	case tea.KeyMsg:
		return v.handleKey(msg)
	}

	if v.step == RescueStepSelectDir {
		var cmd tea.Cmd
		v.filePicker, cmd = v.filePicker.Update(msg)
		cmds = append(cmds, cmd)

		if didSelect, path := v.filePicker.DidSelectFile(msg); didSelect {
			v.codeDir = path
			v.rescue.SetCodebasePath(path)
			v.step = RescueStepScanning
			return v, tea.Batch(v.spinner.Tick, v.runRescue())
		}
	}

	return v, tea.Batch(cmds...)
}

func (v RescueView) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return v, tea.Quit
	case "enter":
		if v.step == RescueStepPreview {
			return v, v.saveResults()
		}
		if v.step == RescueStepSaved {
			return v, tea.Quit
		}
	case "esc":
		if v.step == RescueStepPreview {
			v.step = RescueStepSelectDir
		}
	}
	return v, nil
}

func (v RescueView) runRescue() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		result, err := v.rescue.ScanCodebase(ctx)
		return rescueDoneMsg{result: result, err: err}
	}
}

func (v RescueView) saveResults() tea.Cmd {
	return func() tea.Msg {
		specPath, reportPath, err := v.rescue.SaveResults()
		return rescueSavedMsg{specPath: specPath, reportPath: reportPath, err: err}
	}
}

// View implements tea.Model
func (v RescueView) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("🆘 RESCUE MODE - Reverse Engineer Spec"))
	sb.WriteString("\n")
	sb.WriteString(blurredStyle.Render("LLM: " + v.llmStatus))
	sb.WriteString("\n\n")

	if v.err != nil {
		sb.WriteString(errorStyle.Render("Error: " + v.err.Error()))
		sb.WriteString("\n\n")
	}

	switch v.step {
	case RescueStepSelectDir:
		sb.WriteString("Select codebase directory to analyze:\n\n")
		sb.WriteString(v.filePicker.View())

	case RescueStepScanning:
		sb.WriteString(v.spinner.View())
		result := v.rescue.Result()
		sb.WriteString(fmt.Sprintf(" Scanning codebase... (%d files)", result.FilesScanned))

	case RescueStepPreview:
		result := v.rescue.Result()
		sb.WriteString(successStyle.Render("Scan Complete!"))
		sb.WriteString("\n\n")
		sb.WriteString(fmt.Sprintf("Files scanned: %d\n\n", result.FilesScanned))

		sb.WriteString(focusedStyle.Render("Inferred Specification Preview:"))
		sb.WriteString("\n")
		lines := strings.Split(result.InferredSpec, "\n")
		maxLines := 12
		if len(lines) > maxLines {
			for _, line := range lines[:maxLines] {
				sb.WriteString(line + "\n")
			}
			sb.WriteString(blurredStyle.Render("... (truncated)"))
		} else {
			sb.WriteString(result.InferredSpec)
		}
		sb.WriteString("\n\n")
		sb.WriteString(helpStyle.Render("Enter: save • Esc: back"))

	case RescueStepSaved:
		sb.WriteString(successStyle.Render("✓ Files saved!"))
		sb.WriteString("\n")
		sb.WriteString("Spec: " + v.specPath + "\n")
		sb.WriteString("Report: " + v.reportPath + "\n\n")
		sb.WriteString(helpStyle.Render("Press Enter to exit"))
	}

	return sb.String()
}
