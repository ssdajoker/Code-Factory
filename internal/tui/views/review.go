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

// ReviewStep represents steps in review flow
type ReviewStep int

const (
        ReviewStepSelectSpec ReviewStep = iota
        ReviewStepSelectCode
        ReviewStepAnalyzing
        ReviewStepPreview
        ReviewStepSaved
)

// ReviewView is the TUI view for REVIEW mode
type ReviewView struct {
        review       *modes.ReviewMode
        filePicker   filepicker.Model
        spinner      spinner.Model
        step         ReviewStep
        specFile     string
        codePaths    []string
        width        int
        height       int
        saved        bool
        savePath     string
        err          error
        llmStatus    string
        contractsDir string
}

type reviewDoneMsg struct {
        result *modes.ReviewResult
        err    error
}

type reviewSavedMsg struct {
        path string
        err  error
}

// NewReviewView creates a new review view
func NewReviewView(provider llm.Provider, contractsDir, reportsDir string) ReviewView {
        fp := filepicker.New()
        fp.CurrentDirectory, _ = os.Getwd()
        fp.AllowedTypes = []string{".md", ".txt"}

        s := spinner.New()
        s.Spinner = spinner.Dot
        s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

        llmStatus := "No LLM"
        if provider != nil {
                llmStatus = provider.Name()
        }

        if contractsDir == "" {
                contractsDir = "contracts"
        }

        return ReviewView{
                review:       modes.NewReviewMode(provider, reportsDir),
                filePicker:   fp,
                spinner:      s,
                step:         ReviewStepSelectSpec,
                llmStatus:    llmStatus,
                contractsDir: contractsDir,
        }
}

// Init implements tea.Model
func (v ReviewView) Init() tea.Cmd {
        return v.filePicker.Init()
}

// Update implements tea.Model
func (v ReviewView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        var cmds []tea.Cmd

        switch msg := msg.(type) {
        case tea.WindowSizeMsg:
                v.width = msg.Width
                v.height = msg.Height
                return v, nil

        case spinner.TickMsg:
                if v.step == ReviewStepAnalyzing {
                        var cmd tea.Cmd
                        v.spinner, cmd = v.spinner.Update(msg)
                        return v, cmd
                }

        case reviewDoneMsg:
                if msg.err != nil {
                        v.err = msg.err
                        v.step = ReviewStepSelectSpec
                } else {
                        v.step = ReviewStepPreview
                }
                return v, nil

        case reviewSavedMsg:
                if msg.err != nil {
                        v.err = msg.err
                } else {
                        v.saved = true
                        v.savePath = msg.path
                        v.step = ReviewStepSaved
                }
                return v, nil

        case tea.KeyMsg:
                return v.handleKey(msg)
        }

        if v.step == ReviewStepSelectSpec || v.step == ReviewStepSelectCode {
                var cmd tea.Cmd
                v.filePicker, cmd = v.filePicker.Update(msg)
                cmds = append(cmds, cmd)

                if didSelect, path := v.filePicker.DidSelectFile(msg); didSelect {
                        if v.step == ReviewStepSelectSpec {
                                v.specFile = path
                                v.review.SetSpecFile(path)
                                v.step = ReviewStepSelectCode
                                v.filePicker.AllowedTypes = nil // Allow all for code
                        } else {
                                v.codePaths = append(v.codePaths, path)
                        }
                }
        }

        return v, tea.Batch(cmds...)
}

func (v ReviewView) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "ctrl+c":
                return v, tea.Quit
        case "esc":
                if v.step == ReviewStepSelectCode && len(v.codePaths) > 0 {
                        // Remove last selected path
                        v.codePaths = v.codePaths[:len(v.codePaths)-1]
                }
                return v, nil
        case "enter":
                if v.step == ReviewStepSelectCode && len(v.codePaths) > 0 {
                        v.review.SetCodePaths(v.codePaths)
                        v.step = ReviewStepAnalyzing
                        return v, tea.Batch(v.spinner.Tick, v.runReview())
                }
                if v.step == ReviewStepPreview {
                        return v, v.saveReport()
                }
                if v.step == ReviewStepSaved {
                        return v, tea.Quit
                }
        }
        return v, nil
}

func (v ReviewView) runReview() tea.Cmd {
        return func() tea.Msg {
                ctx := context.Background()
                result, err := v.review.RunReview(ctx)
                return reviewDoneMsg{result: result, err: err}
        }
}

func (v ReviewView) saveReport() tea.Cmd {
        return func() tea.Msg {
                path, err := v.review.SaveReport()
                return reviewSavedMsg{path: path, err: err}
        }
}

// View implements tea.Model
func (v ReviewView) View() string {
        var sb strings.Builder

        sb.WriteString(titleStyle.Render("🔍 REVIEW MODE - Check Code Compliance"))
        sb.WriteString("\n")
        sb.WriteString(blurredStyle.Render("LLM: " + v.llmStatus))
        sb.WriteString("\n\n")

        if v.err != nil {
                sb.WriteString(errorStyle.Render("Error: " + v.err.Error()))
                sb.WriteString("\n\n")
        }

        switch v.step {
        case ReviewStepSelectSpec:
                sb.WriteString("Select specification file from contracts/:\n\n")
                sb.WriteString(v.filePicker.View())

        case ReviewStepSelectCode:
                sb.WriteString(focusedStyle.Render("Spec: " + v.specFile))
                sb.WriteString("\n\n")
                sb.WriteString("Select code files/directories to review:\n")
                if len(v.codePaths) > 0 {
                        sb.WriteString("\nSelected:\n")
                        for _, p := range v.codePaths {
                                sb.WriteString("  ✓ " + p + "\n")
                        }
                        sb.WriteString("\n")
                }
                sb.WriteString(v.filePicker.View())
                sb.WriteString("\n")
                sb.WriteString(helpStyle.Render("Enter: start review • Esc: remove last"))

        case ReviewStepAnalyzing:
                sb.WriteString(v.spinner.View())
                sb.WriteString(" Analyzing code compliance...")

        case ReviewStepPreview:
                result := v.review.Result()
                sb.WriteString(successStyle.Render("Analysis Complete!"))
                sb.WriteString("\n\n")
                sb.WriteString(focusedStyle.Render("Compliance Score: "))
                sb.WriteString(fmt.Sprintf("%d/100\n\n", result.ComplianceScore))

                // Show truncated report
                lines := strings.Split(result.FullReport, "\n")
                maxLines := 15
                if len(lines) > maxLines {
                        for _, line := range lines[:maxLines] {
                                sb.WriteString(line + "\n")
                        }
                        sb.WriteString(blurredStyle.Render("... (truncated)"))
                } else {
                        sb.WriteString(result.FullReport)
                }
                sb.WriteString("\n\n")
                sb.WriteString(helpStyle.Render("Enter: save report"))

        case ReviewStepSaved:
                sb.WriteString(successStyle.Render("✓ Report saved!"))
                sb.WriteString("\n")
                sb.WriteString("Path: " + v.savePath)
                sb.WriteString("\n\n")
                sb.WriteString(helpStyle.Render("Press Enter to exit"))
        }

        return sb.String()
}


