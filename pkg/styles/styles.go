package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Color palette inspired by Claude Code
var (
	// Primary colors
	Primary      = lipgloss.Color("#7C3AED") // Purple
	SuccessColor = lipgloss.Color("#10B981") // Green
	WarningColor = lipgloss.Color("#F59E0B") // Amber
	ErrorColor   = lipgloss.Color("#EF4444") // Red
	InfoColor    = lipgloss.Color("#3B82F6") // Blue
	Muted        = lipgloss.Color("#9CA3AF") // Gray

	// Text colors
	TextPrimary   = lipgloss.Color("#F3F4F6") // Light gray
	TextSecondary = lipgloss.Color("#D1D5DB") // Medium gray
	TextMuted     = lipgloss.Color("#6B7280") // Dark gray
)

// Detect if terminal supports colors
func HasColorSupport() bool {
	return termenv.DefaultOutput().Profile != termenv.Ascii
}

// Style definitions
var (
	// Headers and titles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(TextSecondary).
			Italic(true)

	// Status styles
	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(SuccessColor)

	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ErrorColor)

	WarningStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(WarningColor)

	InfoStyle = lipgloss.NewStyle().
			Foreground(InfoColor)

	// Content styles
	KeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary)

	ValueStyle = lipgloss.NewStyle().
			Foreground(TextPrimary)

	MutedStyle = lipgloss.NewStyle().
			Foreground(TextMuted)

	// Box styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	// Code/command style
	CodeStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1F2937")).
			Foreground(lipgloss.Color("#F9FAFB")).
			Padding(0, 1)

	// Divider
	DividerStyle = lipgloss.NewStyle().
			Foreground(Muted)
)

// Common symbols
const (
	SymbolSuccess   = "✓"
	SymbolError     = "✗"
	SymbolWarning   = "⚠"
	SymbolInfo      = "ℹ"
	SymbolSpinner   = "⏳"
	SymbolArrow     = "→"
	SymbolBullet    = "•"
	SymbolCheckbox  = "☐"
	SymbolChecked   = "☑"
	SymbolSparkles  = "✨"
)

// Helper functions for styled output
func Success(message string) string {
	return SuccessStyle.Render(SymbolSuccess + " " + message)
}

func Error(message string) string {
	return ErrorStyle.Render(SymbolError + " " + message)
}

func Warning(message string) string {
	return WarningStyle.Render(SymbolWarning + " " + message)
}

func Info(message string) string {
	return InfoStyle.Render(SymbolInfo + " " + message)
}

func Loading(message string) string {
	return MutedStyle.Render(SymbolSpinner + " " + message + "...")
}

func Code(text string) string {
	return CodeStyle.Render(text)
}

func Key(text string) string {
	return KeyStyle.Render(text)
}

func Value(text string) string {
	return ValueStyle.Render(text)
}

func Divider() string {
	return DividerStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
