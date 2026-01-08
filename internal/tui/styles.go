package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Colors
var (
	ColorPrimary   = lipgloss.Color("#7C3AED") // Purple
	ColorSecondary = lipgloss.Color("#06B6D4") // Cyan
	ColorSuccess   = lipgloss.Color("#10B981") // Green
	ColorWarning   = lipgloss.Color("#F59E0B") // Amber
	ColorError     = lipgloss.Color("#EF4444") // Red
	ColorSubtle    = lipgloss.Color("#6B7280") // Gray
	ColorText      = lipgloss.Color("#F3F4F6") // Light gray
)

// Base styles
var (
	StyleTitle = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Padding(0, 1)

	StyleSubtitle = lipgloss.NewStyle().
		Foreground(ColorSecondary).
		Italic(true)

	StyleSuccess = lipgloss.NewStyle().
		Foreground(ColorSuccess)

	StyleWarning = lipgloss.NewStyle().
		Foreground(ColorWarning)

	StyleError = lipgloss.NewStyle().
		Foreground(ColorError).
		Bold(true)

	StyleSubtle = lipgloss.NewStyle().
		Foreground(ColorSubtle)

	StyleText = lipgloss.NewStyle().
		Foreground(ColorText)
)

// Menu styles
var (
	StyleMenuItemNormal = lipgloss.NewStyle().
		Padding(0, 2)

	StyleMenuItemSelected = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Background(lipgloss.Color("#1F2937")).
		Bold(true).
		Padding(0, 2)
)

// Box styles
var (
	StyleBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2)

	StyleHeader = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorPrimary).
		Padding(0, 2).
		Align(lipgloss.Center)
)

// RenderHeader renders a centered header box
func RenderHeader(title string, width int) string {
	if width < 40 {
		width = 60
	}
	headerWidth := min(width-4, 60)

	header := StyleHeader.
		Width(headerWidth).
		Render("🏭 " + title + " 🏭")

	return lipgloss.PlaceHorizontal(width, lipgloss.Center, header)
}

// RenderMenu renders a menu with the given items
func RenderMenu(items []string, selected int) string {
	var b strings.Builder

	for i, item := range items {
		if i == selected {
			b.WriteString(StyleMenuItemSelected.Render("▸ " + item))
		} else {
			b.WriteString(StyleMenuItemNormal.Render("  " + item))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
