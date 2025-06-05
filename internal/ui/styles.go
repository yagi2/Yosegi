package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/yagi2/yosegi/internal/config"
)

// Color palette (default values, can be overridden by config)
var (
	Primary   = lipgloss.Color("#7C3AED") // Purple
	Secondary = lipgloss.Color("#06B6D4") // Cyan
	Success   = lipgloss.Color("#10B981") // Green
	Warning   = lipgloss.Color("#F59E0B") // Yellow
	Error     = lipgloss.Color("#EF4444") // Red
	Muted     = lipgloss.Color("#6B7280") // Gray
	Text      = lipgloss.Color("#F9FAFB") // Light
)

// InitializeTheme initializes the theme from config
func InitializeTheme(cfg *config.Config) {
	if cfg.Theme.Primary != "" {
		Primary = lipgloss.Color(cfg.Theme.Primary)
	}
	if cfg.Theme.Secondary != "" {
		Secondary = lipgloss.Color(cfg.Theme.Secondary)
	}
	if cfg.Theme.Success != "" {
		Success = lipgloss.Color(cfg.Theme.Success)
	}
	if cfg.Theme.Warning != "" {
		Warning = lipgloss.Color(cfg.Theme.Warning)
	}
	if cfg.Theme.Error != "" {
		Error = lipgloss.Color(cfg.Theme.Error)
	}
	if cfg.Theme.Muted != "" {
		Muted = lipgloss.Color(cfg.Theme.Muted)
	}
	if cfg.Theme.Text != "" {
		Text = lipgloss.Color(cfg.Theme.Text)
	}
}

// Base styles
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Italic(true)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Text).
				Background(Primary).
				Bold(true).
				Padding(0, 1)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(Text).
			Padding(0, 1)

	CurrentItemStyle = lipgloss.NewStyle().
				Foreground(Success).
				Bold(true).
				Padding(0, 1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Italic(true).
			Margin(1, 0)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true).
			Padding(0, 1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true).
			Padding(0, 1)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true).
			Padding(0, 1)

	NormalStyle = lipgloss.NewStyle().
			Foreground(Text)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	InputStyle = lipgloss.NewStyle().
			Foreground(Text).
			Background(lipgloss.Color("#374151")).
			Padding(0, 1).
			Margin(0, 1)
)

// GetStatusIcon returns an icon based on status
func GetStatusIcon(isCurrent bool) string {
	if isCurrent {
		return "●"
	}
	return "○"
}

// GetBranchIcon returns an icon for branch
func GetBranchIcon() string {
	return ""
}

// GetPathIcon returns an icon for path
func GetPathIcon() string {
	return ""
}
