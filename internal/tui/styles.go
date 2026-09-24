package tui

import "github.com/charmbracelet/lipgloss"

var (
	cream = lipgloss.Color("#fff1b8")
	gold  = lipgloss.Color("#ffd75f")
	muted = lipgloss.Color("#7a88b8")
	blue  = lipgloss.Color("#8eb0ff")
	miss  = lipgloss.Color("#ff8b8b")
	okc   = lipgloss.Color("#9dffb0")
	ink   = lipgloss.Color("#1a1f33")

	frameStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6b8cff")).
			Padding(1, 2)

	titleStyle  = lipgloss.NewStyle().Foreground(cream).Bold(true)
	goldStyle   = lipgloss.NewStyle().Foreground(gold).Bold(true)
	mutedStyle  = lipgloss.NewStyle().Foreground(muted)
	bodyStyle   = lipgloss.NewStyle().Foreground(cream)
	missStyle   = lipgloss.NewStyle().Foreground(miss).Bold(true)
	okStyle     = lipgloss.NewStyle().Foreground(okc).Bold(true)
	letterStyle = lipgloss.NewStyle().Foreground(blue).Bold(true)
	selStyle    = lipgloss.NewStyle().Foreground(ink).Background(gold).Bold(true)
)
