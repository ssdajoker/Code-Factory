package tui

import (
	"github.com/charmbracelet/bubbletea"
)

// App represents the main TUI application
type App struct {
	program *tea.Program
}

// NewApp creates a new TUI application
func NewApp() *App {
	return &App{}
}

// Run starts the TUI application
func (a *App) Run() error {
	// TODO: Initialize with appropriate model based on mode
	// model := NewIntakeModel()
	// a.program = tea.NewProgram(model, tea.WithAltScreen())
	// return a.program.Start()
	return nil
}
