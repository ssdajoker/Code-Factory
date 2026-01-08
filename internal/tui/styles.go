package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	ColorPrimary   = lipgloss.Color("99")
	ColorSecondary = lipgloss.Color("205")
	ColorSuccess   = lipgloss.Color("42")
	ColorWarning   = lipgloss.Color("214")
	ColorError     = lipgloss.Color("196")
	ColorSubtle    = lipgloss.Color("240")

	// Styles
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	StyleSubtle = lipgloss.NewStyle().
			Foreground(ColorSubtle)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	StyleWarning = lipgloss.NewStyle().
			Foreground(ColorWarning)

	StyleError = lipgloss.NewStyle().
			Foreground(ColorError)

	StyleSelected = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true)

	StyleNormal = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))
)

// RenderHeader renders a centered header
func RenderHeader(title string, width int) string {
	if width < 40 {
		width = 80
	}

	border := strings.Repeat("═", width-4)
	header := StyleTitle.Render("╔" + border + "╗")
	header += "\n"

	padding := (width - 4 - len(title)) / 2
	if padding < 0 {
		padding = 0
	}
	titleLine := "║" + strings.Repeat(" ", padding) + title + strings.Repeat(" ", width-4-padding-len(title)) + "║"
	header += StyleTitle.Render(titleLine)
	header += "\n"
	header += StyleTitle.Render("╚" + border + "╝")

	return header
}

// RenderMenu renders a menu with selection
func RenderMenu(items []string, selected int) string {
	var sb strings.Builder
	for i, item := range items {
		if i == selected {
			sb.WriteString(StyleSelected.Render("  ▸ " + item))
		} else {
			sb.WriteString(StyleNormal.Render("    " + item))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// RenderLLMStatus renders the LLM status indicator
func RenderLLMStatus(provider string, available bool) string {
	if !available {
		return StyleWarning.Render("⚠ No LLM")
	}
	return StyleSuccess.Render("✓ " + provider)
}
