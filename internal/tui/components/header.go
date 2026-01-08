package components

import (
	"github.com/charmbracelet/lipgloss"
)

// Header renders a header component with title and optional status
type Header struct {
	Title  string
	Status string
	Width  int
}

// NewHeader creates a new header
func NewHeader(title string) *Header {
	return &Header{
		Title: title,
		Width: 60,
	}
}

// SetStatus sets the status text
func (h *Header) SetStatus(status string) {
	h.Status = status
}

// SetWidth sets the header width
func (h *Header) SetWidth(width int) {
	h.Width = width
}

// Render renders the header
func (h *Header) Render() string {
	headerStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		Padding(0, 2).
		Align(lipgloss.Center).
		Width(h.Width)

	title := "🏭 " + h.Title + " 🏭"

	if h.Status != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).
			Italic(true)
		title += "\n" + statusStyle.Render(h.Status)
	}

	return headerStyle.Render(title)
}
