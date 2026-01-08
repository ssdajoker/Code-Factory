package views

import (
	"github.com/ssdajoker/Code-Factory/internal/tui/components"
)

// HomeView represents the home screen
type HomeView struct {
	header *components.Header
	menu   *components.Menu
	width  int
	height int
}

// NewHomeView creates a new home view
func NewHomeView() *HomeView {
	menuItems := []components.MenuItem{
		{Label: "Initialize Project", Icon: "🚀", Description: "Set up Factory in your project"},
		{Label: "INTAKE - Capture Vision", Icon: "📝", Description: "Create specifications from ideas"},
		{Label: "REVIEW - Check Code", Icon: "🔍", Description: "Verify code against specs"},
		{Label: "RESCUE - Reverse Engineer", Icon: "🆘", Description: "Generate specs from existing code"},
		{Label: "CHANGE_ORDER - Track Drift", Icon: "📋", Description: "Manage specification changes"},
		{Label: "Settings", Icon: "⚙️", Description: "Configure Factory"},
		{Label: "Quit", Icon: "❌", Description: "Exit Factory"},
	}

	return &HomeView{
		header: components.NewHeader("SPEC-DRIVEN SOFTWARE FACTORY"),
		menu:   components.NewMenu(menuItems),
	}
}

// SetSize sets the view dimensions
func (v *HomeView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.header.SetWidth(min(width-4, 60))
	v.menu.Width = min(width-4, 50)
}

// Up moves menu selection up
func (v *HomeView) Up() {
	v.menu.Up()
}

// Down moves menu selection down
func (v *HomeView) Down() {
	v.menu.Down()
}

// Selected returns the selected menu index
func (v *HomeView) Selected() int {
	return v.menu.Selected
}

// Render renders the home view
func (v *HomeView) Render() string {
	s := "\n"
	s += v.header.Render()
	s += "\n\n"
	s += v.menu.Render()
	return s
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
