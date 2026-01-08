package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// MenuItem represents a single menu item
type MenuItem struct {
	Label       string
	Description string
	Icon        string
	Disabled    bool
}

// Menu is a selectable menu component
type Menu struct {
	Items    []MenuItem
	Selected int
	Width    int
}

// NewMenu creates a new menu
func NewMenu(items []MenuItem) *Menu {
	return &Menu{
		Items:    items,
		Selected: 0,
		Width:    40,
	}
}

// Up moves selection up
func (m *Menu) Up() {
	if m.Selected > 0 {
		m.Selected--
		// Skip disabled items
		if m.Items[m.Selected].Disabled {
			m.Up()
		}
	}
}

// Down moves selection down
func (m *Menu) Down() {
	if m.Selected < len(m.Items)-1 {
		m.Selected++
		// Skip disabled items
		if m.Items[m.Selected].Disabled {
			m.Down()
		}
	}
}

// SelectedItem returns the currently selected item
func (m *Menu) SelectedItem() MenuItem {
	if m.Selected >= 0 && m.Selected < len(m.Items) {
		return m.Items[m.Selected]
	}
	return MenuItem{}
}

// Render renders the menu
func (m *Menu) Render() string {
	var b strings.Builder

	normalStyle := lipgloss.NewStyle().
		Padding(0, 2)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7C3AED")).
		Background(lipgloss.Color("#1F2937")).
		Bold(true).
		Padding(0, 2)

	disabledStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4B5563")).
		Padding(0, 2)

	for i, item := range m.Items {
		label := item.Label
		if item.Icon != "" {
			label = item.Icon + " " + label
		}

		var line string
		if item.Disabled {
			line = disabledStyle.Render("  " + label)
		} else if i == m.Selected {
			line = selectedStyle.Render("▸ " + label)
		} else {
			line = normalStyle.Render("  " + label)
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}
