package tui

import (
        "fmt"

        tea "github.com/charmbracelet/bubbletea"
        "github.com/ssdajoker/Code-Factory/internal/tui/views"
)

// View represents different screens in the TUI
type View int

const (
        ViewHome View = iota
        ViewInit
        ViewIntake
        ViewReview
        ViewRescue
        ViewChangeOrder
        ViewSettings
)

// Model is the main Bubble Tea model
type Model struct {
        currentView     View
        menuIndex       int
        quitting        bool
        width           int
        height          int
        err             error
        intakeView      *views.IntakeView
        reviewView      *views.ReviewView
        rescueView      *views.RescueView
        changeOrderView *views.ChangeOrderView
}

// menuItems for the home screen
var menuItems = []string{
        "🚀 Initialize Project",
        "📝 INTAKE - Capture Vision",
        "🔍 REVIEW - Check Code",
        "🆘 RESCUE - Reverse Engineer",
        "📋 CHANGE_ORDER - Track Drift",
        "⚙️  Settings",
        "❌ Quit",
}

// New creates a new TUI model
func New() Model {
        return Model{
                currentView: ViewHome,
                menuIndex:   0,
        }
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
        return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
        // Delegate to active view
        switch m.currentView {
        case ViewIntake:
                if m.intakeView != nil {
                        updatedModel, cmd := m.intakeView.Update(msg)
                        if iv, ok := updatedModel.(views.IntakeView); ok {
                                m.intakeView = &iv
                        }
                        return m, cmd
                }
        case ViewReview:
                if m.reviewView != nil {
                        updatedModel, cmd := m.reviewView.Update(msg)
                        if rv, ok := updatedModel.(views.ReviewView); ok {
                                m.reviewView = &rv
                        }
                        return m, cmd
                }
        case ViewRescue:
                if m.rescueView != nil {
                        updatedModel, cmd := m.rescueView.Update(msg)
                        if rv, ok := updatedModel.(views.RescueView); ok {
                                m.rescueView = &rv
                        }
                        return m, cmd
                }
        case ViewChangeOrder:
                if m.changeOrderView != nil {
                        updatedModel, cmd := m.changeOrderView.Update(msg)
                        if cv, ok := updatedModel.(views.ChangeOrderView); ok {
                                m.changeOrderView = &cv
                        }
                        return m, cmd
                }
        }

        switch msg := msg.(type) {
        case tea.KeyMsg:
                return m.handleKey(msg)
        case tea.WindowSizeMsg:
                m.width = msg.Width
                m.height = msg.Height
                return m, nil
        }
        return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
        switch msg.String() {
        case "ctrl+c", "q":
                m.quitting = true
                return m, tea.Quit
        case "up", "k":
                if m.menuIndex > 0 {
                        m.menuIndex--
                }
        case "down", "j":
                if m.menuIndex < len(menuItems)-1 {
                        m.menuIndex++
                }
        case "enter", " ":
                return m.selectMenuItem()
        case "esc":
                if m.currentView != ViewHome {
                        m.currentView = ViewHome
                        m.intakeView = nil
                }
        }
        return m, nil
}

func (m Model) selectMenuItem() (tea.Model, tea.Cmd) {
        switch m.menuIndex {
        case 0:
                m.currentView = ViewInit
        case 1:
                m.currentView = ViewIntake
                iv := views.NewIntakeView(nil, "contracts")
                m.intakeView = &iv
                return m, m.intakeView.Init()
        case 2:
                m.currentView = ViewReview
                rv := views.NewReviewView(nil, "contracts", "reports")
                m.reviewView = &rv
                return m, m.reviewView.Init()
        case 3:
                m.currentView = ViewRescue
                rv := views.NewRescueView(nil, "contracts", "reports")
                m.rescueView = &rv
                return m, m.rescueView.Init()
        case 4:
                m.currentView = ViewChangeOrder
                cv := views.NewChangeOrderView(nil, "contracts")
                m.changeOrderView = &cv
                return m, m.changeOrderView.Init()
        case 5:
                m.currentView = ViewSettings
        case 6:
                m.quitting = true
                return m, tea.Quit
        }
        return m, nil
}

// View implements tea.Model
func (m Model) View() string {
        if m.quitting {
                return "Goodbye! 🏭\n"
        }

        switch m.currentView {
        case ViewHome:
                return m.viewHome()
        case ViewIntake:
                if m.intakeView != nil {
                        return m.intakeView.View()
                }
        case ViewReview:
                if m.reviewView != nil {
                        return m.reviewView.View()
                }
        case ViewRescue:
                if m.rescueView != nil {
                        return m.rescueView.View()
                }
        case ViewChangeOrder:
                if m.changeOrderView != nil {
                        return m.changeOrderView.View()
                }
        }
        return m.viewPlaceholder()
}

func (m Model) viewHome() string {
        s := "\n"
        s += RenderHeader("SPEC-DRIVEN SOFTWARE FACTORY", m.width)
        s += "\n\n"
        s += RenderMenu(menuItems, m.menuIndex)
        s += "\n\n"
        s += StyleSubtle.Render("↑/↓: navigate • enter: select • q: quit")
        s += "\n"
        return s
}

func (m Model) viewPlaceholder() string {
        viewNames := map[View]string{
                ViewInit:        "Initialize Project",
                ViewIntake:      "INTAKE Mode",
                ViewReview:      "REVIEW Mode",
                ViewRescue:      "RESCUE Mode",
                ViewChangeOrder: "CHANGE_ORDER Mode",
                ViewSettings:    "Settings",
        }

        s := "\n"
        s += RenderHeader(viewNames[m.currentView], m.width)
        s += "\n\n"
        s += StyleWarning.Render("  Coming soon...")
        s += "\n\n"
        s += StyleSubtle.Render("  Press ESC to go back")
        s += "\n"
        return s
}

// Run starts the TUI application
func Run() error {
        p := tea.NewProgram(New(), tea.WithAltScreen())
        _, err := p.Run()
        if err != nil {
                return fmt.Errorf("error running TUI: %w", err)
        }
        return nil
}

// RunIntake starts the INTAKE mode TUI directly
func RunIntake() error {
        iv := views.NewIntakeView(nil, "contracts")
        p := tea.NewProgram(iv, tea.WithAltScreen())
        _, err := p.Run()
        if err != nil {
                return fmt.Errorf("error running intake: %w", err)
        }
        return nil
}
