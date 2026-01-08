package views

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ssdajoker/Code-Factory/internal/llm"
	"github.com/ssdajoker/Code-Factory/internal/modes"
)

// ChangeOrderStep represents steps in change order flow
type ChangeOrderStep int

const (
	COStepSelectSpec ChangeOrderStep = iota
	COStepSelectCode
	COStepAnalyzing
	COStepReview
	COStepReason
	COStepSaved
)

// ChangeOrderView is the TUI view for CHANGE_ORDER mode
type ChangeOrderView struct {
	changeOrder   *modes.ChangeOrderMode
	filePicker    filepicker.Model
	spinner       spinner.Model
	textInput     textinput.Model
	step          ChangeOrderStep
	specFile      string
	codeDir       string
	selectedIdx   int
	width         int
	height        int
	saved         bool
	savePath      string
	err           error
	llmStatus     string
}

type changeOrderDoneMsg struct {
	result *modes.ChangeOrderResult
	err    error
}

type changeOrderSavedMsg struct {
	path string
	err  error
}

// NewChangeOrderView creates a new change order view
func NewChangeOrderView(provider llm.Provider, contractsDir string) ChangeOrderView {
	fp := filepicker.New()
	fp.CurrentDirectory, _ = os.Getwd()
	fp.AllowedTypes = []string{".md", ".txt"}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ti := textinput.New()
	ti.Placeholder = "Enter reason for change..."
	ti.CharLimit = 256
	ti.Width = 60

	llmStatus := "No LLM"
	if provider != nil {
		llmStatus = provider.Name()
	}

	return ChangeOrderView{
		changeOrder: modes.NewChangeOrderMode(provider, contractsDir),
		filePicker:  fp,
		spinner:     s,
		textInput:   ti,
		step:        COStepSelectSpec,
		llmStatus:   llmStatus,
	}
}

// Init implements tea.Model
func (v ChangeOrderView) Init() tea.Cmd {
	return v.filePicker.Init()
}

// Update implements tea.Model
func (v ChangeOrderView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
		return v, nil

	case spinner.TickMsg:
		if v.step == COStepAnalyzing {
			var cmd tea.Cmd
			v.spinner, cmd = v.spinner.Update(msg)
			return v, cmd
		}

	case changeOrderDoneMsg:
		if msg.err != nil {
			v.err = msg.err
			v.step = COStepSelectSpec
		} else {
			v.step = COStepReview
		}
		return v, nil

	case changeOrderSavedMsg:
		if msg.err != nil {
			v.err = msg.err
		} else {
			v.saved = true
			v.savePath = msg.path
			v.step = COStepSaved
		}
		return v, nil

	case tea.KeyMsg:
		return v.handleKey(msg)
	}

	if v.step == COStepSelectSpec || v.step == COStepSelectCode {
		var cmd tea.Cmd
		v.filePicker, cmd = v.filePicker.Update(msg)
		cmds = append(cmds, cmd)

		if didSelect, path := v.filePicker.DidSelectFile(msg); didSelect {
			if v.step == COStepSelectSpec {
				v.specFile = path
				v.changeOrder.SetSpecFile(path)
				v.step = COStepSelectCode
				v.filePicker.AllowedTypes = nil
				v.filePicker.DirAllowed = true
				v.filePicker.FileAllowed = false
			} else {
				v.codeDir = path
				v.changeOrder.SetCodebasePath(path)
				v.step = COStepAnalyzing
				return v, tea.Batch(v.spinner.Tick, v.runDetection())
			}
		}
	}

	if v.step == COStepReason {
		var cmd tea.Cmd
		v.textInput, cmd = v.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return v, tea.Batch(cmds...)
}

func (v ChangeOrderView) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return v, tea.Quit

	case "up", "k":
		if v.step == COStepReview && v.selectedIdx > 0 {
			v.selectedIdx--
		}

	case "down", "j":
		if v.step == COStepReview {
			changes := v.changeOrder.Result().Changes
			if v.selectedIdx < len(changes)-1 {
				v.selectedIdx++
			}
		}

	case "a": // Approve
		if v.step == COStepReview {
			v.step = COStepReason
			v.textInput.Focus()
			return v, textinput.Blink
		}

	case "r": // Reject
		if v.step == COStepReview {
			changes := v.changeOrder.Result().Changes
			if len(changes) > v.selectedIdx {
				v.changeOrder.RejectChange(changes[v.selectedIdx].ID, "Rejected by user")
			}
		}

	case "d": // Defer
		if v.step == COStepReview {
			changes := v.changeOrder.Result().Changes
			if len(changes) > v.selectedIdx {
				v.changeOrder.DeferChange(changes[v.selectedIdx].ID, "Deferred for later")
			}
		}

	case "enter":
		if v.step == COStepReason {
			changes := v.changeOrder.Result().Changes
			if len(changes) > v.selectedIdx {
				v.changeOrder.ApproveChange(changes[v.selectedIdx].ID, v.textInput.Value())
			}
			v.textInput.SetValue("")
			v.step = COStepReview
		}
		if v.step == COStepSaved {
			return v, tea.Quit
		}

	case "s": // Save
		if v.step == COStepReview {
			return v, v.saveChangeOrder()
		}

	case "esc":
		if v.step == COStepReason {
			v.step = COStepReview
		}
	}
	return v, nil
}

func (v ChangeOrderView) runDetection() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		result, err := v.changeOrder.DetectDrift(ctx)
		return changeOrderDoneMsg{result: result, err: err}
	}
}

func (v ChangeOrderView) saveChangeOrder() tea.Cmd {
	return func() tea.Msg {
		path, err := v.changeOrder.SaveChangeOrder()
		return changeOrderSavedMsg{path: path, err: err}
	}
}

// View implements tea.Model
func (v ChangeOrderView) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("📋 CHANGE_ORDER MODE - Track Drift"))
	sb.WriteString("\n")
	sb.WriteString(blurredStyle.Render("LLM: " + v.llmStatus))
	sb.WriteString("\n\n")

	if v.err != nil {
		sb.WriteString(errorStyle.Render("Error: " + v.err.Error()))
		sb.WriteString("\n\n")
	}

	switch v.step {
	case COStepSelectSpec:
		sb.WriteString("Select specification file:\n\n")
		sb.WriteString(v.filePicker.View())

	case COStepSelectCode:
		sb.WriteString(focusedStyle.Render("Spec: " + v.specFile))
		sb.WriteString("\n\n")
		sb.WriteString("Select codebase directory:\n\n")
		sb.WriteString(v.filePicker.View())

	case COStepAnalyzing:
		sb.WriteString(v.spinner.View())
		sb.WriteString(" Detecting drift...")

	case COStepReview:
		result := v.changeOrder.Result()
		sb.WriteString(successStyle.Render("Drift Analysis Complete!"))
		sb.WriteString("\n\n")

		for i, change := range result.Changes {
			prefix := "  "
			if i == v.selectedIdx {
				prefix = "> "
			}
			status := change.Status
			statusStyle := blurredStyle
			switch status {
			case "approved":
				statusStyle = successStyle
			case "rejected":
				statusStyle = errorStyle
			}
			sb.WriteString(fmt.Sprintf("%s%s: %s [%s]\n", prefix, change.ID, change.Description, statusStyle.Render(status)))
		}

		sb.WriteString("\n")
		sb.WriteString(helpStyle.Render("↑/↓: select • a: approve • r: reject • d: defer • s: save"))

	case COStepReason:
		sb.WriteString("Enter reason for approval:\n\n")
		sb.WriteString(v.textInput.View())
		sb.WriteString("\n\n")
		sb.WriteString(helpStyle.Render("Enter: confirm • Esc: cancel"))

	case COStepSaved:
		sb.WriteString(successStyle.Render("✓ Change order saved!"))
		sb.WriteString("\n")
		sb.WriteString("Path: " + v.savePath)
		sb.WriteString("\n\n")
		sb.WriteString(helpStyle.Render("Press Enter to exit"))
	}

	return sb.String()
}
