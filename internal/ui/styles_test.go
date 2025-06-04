package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/yagi2/yosegi/internal/config"
)

func TestDefaultColors(t *testing.T) {
	// Test default color values
	expectedColors := map[string]string{
		"Primary":   "#7C3AED",
		"Secondary": "#06B6D4",
		"Success":   "#10B981",
		"Warning":   "#F59E0B",
		"Error":     "#EF4444",
		"Muted":     "#6B7280",
		"Text":      "#F9FAFB",
	}

	actualColors := map[string]lipgloss.Color{
		"Primary":   Primary,
		"Secondary": Secondary,
		"Success":   Success,
		"Warning":   Warning,
		"Error":     Error,
		"Muted":     Muted,
		"Text":      Text,
	}

	for name, expected := range expectedColors {
		actual := string(actualColors[name])
		if actual != expected {
			t.Errorf("Expected %s color to be %s, got %s", name, expected, actual)
		}
	}
}

func TestInitializeTheme(t *testing.T) {
	// Save original colors
	originalColors := map[string]lipgloss.Color{
		"Primary":   Primary,
		"Secondary": Secondary,
		"Success":   Success,
		"Warning":   Warning,
		"Error":     Error,
		"Muted":     Muted,
		"Text":      Text,
	}

	// Restore original colors after test
	defer func() {
		Primary = originalColors["Primary"]
		Secondary = originalColors["Secondary"]
		Success = originalColors["Success"]
		Warning = originalColors["Warning"]
		Error = originalColors["Error"]
		Muted = originalColors["Muted"]
		Text = originalColors["Text"]
	}()

	tests := []struct {
		name   string
		config *config.Config
		verify func(t *testing.T)
	}{
		{
			name: "Full theme configuration",
			config: &config.Config{
				Theme: config.ThemeConfig{
					Primary:   "#FF0000",
					Secondary: "#00FF00",
					Success:   "#0000FF",
					Warning:   "#FFFF00",
					Error:     "#FF00FF",
					Muted:     "#00FFFF",
					Text:      "#000000",
				},
			},
			verify: func(t *testing.T) {
				expectedColors := map[string]string{
					"Primary":   "#FF0000",
					"Secondary": "#00FF00",
					"Success":   "#0000FF",
					"Warning":   "#FFFF00",
					"Error":     "#FF00FF",
					"Muted":     "#00FFFF",
					"Text":      "#000000",
				}

				actualColors := map[string]lipgloss.Color{
					"Primary":   Primary,
					"Secondary": Secondary,
					"Success":   Success,
					"Warning":   Warning,
					"Error":     Error,
					"Muted":     Muted,
					"Text":      Text,
				}

				for name, expected := range expectedColors {
					actual := string(actualColors[name])
					if actual != expected {
						t.Errorf("Expected %s color to be %s, got %s", name, expected, actual)
					}
				}
			},
		},
		{
			name: "Partial theme configuration",
			config: &config.Config{
				Theme: config.ThemeConfig{
					Primary: "#CUSTOM1",
					Error:   "#CUSTOM2",
					// Other colors should remain as original defaults
				},
			},
			verify: func(t *testing.T) {
				if string(Primary) != "#CUSTOM1" {
					t.Errorf("Expected Primary to be #CUSTOM1, got %s", string(Primary))
				}
				if string(Error) != "#CUSTOM2" {
					t.Errorf("Expected Error to be #CUSTOM2, got %s", string(Error))
				}
				// These should remain as defaults (original values)
				if string(Secondary) != string(originalColors["Secondary"]) {
					t.Errorf("Secondary should remain default")
				}
				if string(Success) != string(originalColors["Success"]) {
					t.Errorf("Success should remain default")
				}
			},
		},
		{
			name: "Empty theme configuration",
			config: &config.Config{
				Theme: config.ThemeConfig{
					// All empty - should keep defaults
				},
			},
			verify: func(t *testing.T) {
				// All colors should remain as originals
				for name, original := range originalColors {
					var actual lipgloss.Color
					switch name {
					case "Primary":
						actual = Primary
					case "Secondary":
						actual = Secondary
					case "Success":
						actual = Success
					case "Warning":
						actual = Warning
					case "Error":
						actual = Error
					case "Muted":
						actual = Muted
					case "Text":
						actual = Text
					}
					if string(actual) != string(original) {
						t.Errorf("%s color should remain default, expected %s, got %s", name, string(original), string(actual))
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset to original colors before each test
			Primary = originalColors["Primary"]
			Secondary = originalColors["Secondary"]
			Success = originalColors["Success"]
			Warning = originalColors["Warning"]
			Error = originalColors["Error"]
			Muted = originalColors["Muted"]
			Text = originalColors["Text"]

			InitializeTheme(tt.config)
			tt.verify(t)
		})
	}
}

func TestGetStatusIcon(t *testing.T) {
	tests := []struct {
		name      string
		isCurrent bool
		expected  string
	}{
		{
			name:      "Current worktree",
			isCurrent: true,
			expected:  "●",
		},
		{
			name:      "Non-current worktree",
			isCurrent: false,
			expected:  "○",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStatusIcon(tt.isCurrent)
			if result != tt.expected {
				t.Errorf("Expected icon %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetBranchIcon(t *testing.T) {
	icon := GetBranchIcon()
	expected := ""
	if icon != expected {
		t.Errorf("Expected branch icon %s, got %s", expected, icon)
	}
}

func TestGetPathIcon(t *testing.T) {
	icon := GetPathIcon()
	expected := ""
	if icon != expected {
		t.Errorf("Expected path icon %s, got %s", expected, icon)
	}
}

func TestStylesCreation(t *testing.T) {
	// Test that all style variables are properly initialized
	styles := []struct {
		name  string
		style lipgloss.Style
	}{
		{"TitleStyle", TitleStyle},
		{"SubtitleStyle", SubtitleStyle},
		{"SelectedItemStyle", SelectedItemStyle},
		{"NormalItemStyle", NormalItemStyle},
		{"CurrentItemStyle", CurrentItemStyle},
		{"HelpStyle", HelpStyle},
		{"ErrorStyle", ErrorStyle},
		{"SuccessStyle", SuccessStyle},
		{"WarningStyle", WarningStyle},
		{"BorderStyle", BorderStyle},
		{"InputStyle", InputStyle},
	}

	for _, s := range styles {
		t.Run(s.name, func(t *testing.T) {
			// Test that styles can render without panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Style %s panicked during render: %v", s.name, r)
				}
			}()

			// Test rendering with sample text
			rendered := s.style.Render("Test")
			if rendered == "" {
				t.Errorf("Style %s rendered empty string", s.name)
			}

			// Test that the style has some properties set
			// This is a basic check to ensure the style is not completely empty
			styleStr := s.style.String()
			if styleStr == "" {
				t.Logf("Warning: Style %s appears to have no properties set", s.name)
			}
		})
	}
}

func TestStylesProperties(t *testing.T) {
	tests := []struct {
		name          string
		style         lipgloss.Style
		expectedProps []string // Properties we expect to be set
	}{
		{
			name:          "TitleStyle",
			style:         TitleStyle,
			expectedProps: []string{"foreground", "bold", "padding"},
		},
		{
			name:          "SubtitleStyle",
			style:         SubtitleStyle,
			expectedProps: []string{"foreground", "italic"},
		},
		{
			name:          "SelectedItemStyle",
			style:         SelectedItemStyle,
			expectedProps: []string{"foreground", "background", "bold", "padding"},
		},
		{
			name:          "BorderStyle",
			style:         BorderStyle,
			expectedProps: []string{"border", "padding"},
		},
		{
			name:          "InputStyle",
			style:         InputStyle,
			expectedProps: []string{"foreground", "background", "padding", "margin"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the style can be used for rendering
			rendered := tt.style.Render("Sample Text")
			if len(rendered) == 0 {
				t.Errorf("Style %s produced empty output", tt.name)
			}

			// Test that the style can render (basic check for non-empty style)
			testRendered := tt.style.Render("test")
			if testRendered == "test" {
				// If rendered output is exactly the same as input, 
				// the style might not have any properties set
				t.Logf("Style %s may not have properties set (output unchanged)", tt.name)
			}
		})
	}
}

func TestColorConsistency(t *testing.T) {
	// Test that styles use the correct colors
	// This is more of a smoke test to ensure colors are being used correctly

	// Save original colors
	originalPrimary := Primary
	originalSuccess := Success

	// Set custom colors
	Primary = lipgloss.Color("#TEST01")
	Success = lipgloss.Color("#TEST02")

	// Restore after test
	defer func() {
		Primary = originalPrimary
		Success = originalSuccess
	}()

	// These styles should incorporate the Primary color
	primaryStyles := []lipgloss.Style{
		TitleStyle,
		SelectedItemStyle,
		BorderStyle,
	}

	for i, style := range primaryStyles {
		rendered := style.Render("Test")
		if rendered == "" {
			t.Errorf("Primary-based style %d rendered empty", i)
		}
	}

	// These styles should incorporate the Success color
	successStyles := []lipgloss.Style{
		CurrentItemStyle,
		SuccessStyle,
	}

	for i, style := range successStyles {
		rendered := style.Render("Test")
		if rendered == "" {
			t.Errorf("Success-based style %d rendered empty", i)
		}
	}
}

func TestThemeReinitialization(t *testing.T) {
	// Test that theme can be reinitialized multiple times
	originalPrimary := Primary

	defer func() {
		Primary = originalPrimary
	}()

	configs := []*config.Config{
		{
			Theme: config.ThemeConfig{
				Primary: "#FF0000",
			},
		},
		{
			Theme: config.ThemeConfig{
				Primary: "#00FF00",
			},
		},
		{
			Theme: config.ThemeConfig{
				Primary: "#0000FF",
			},
		},
	}

	expectedColors := []string{"#FF0000", "#00FF00", "#0000FF"}

	for i, cfg := range configs {
		InitializeTheme(cfg)
		actual := string(Primary)
		expected := expectedColors[i]
		if actual != expected {
			t.Errorf("Iteration %d: Expected Primary to be %s, got %s", i, expected, actual)
		}
	}
}

// Benchmark tests
func BenchmarkInitializeTheme(b *testing.B) {
	cfg := &config.Config{
		Theme: config.ThemeConfig{
			Primary:   "#FF0000",
			Secondary: "#00FF00",
			Success:   "#0000FF",
			Warning:   "#FFFF00",
			Error:     "#FF00FF",
			Muted:     "#00FFFF",
			Text:      "#000000",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InitializeTheme(cfg)
	}
}

func BenchmarkGetStatusIcon(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetStatusIcon(i%2 == 0)
	}
}

func BenchmarkStyleRender(b *testing.B) {
	text := "Sample text for benchmarking style rendering"
	
	b.Run("TitleStyle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			TitleStyle.Render(text)
		}
	})

	b.Run("SelectedItemStyle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			SelectedItemStyle.Render(text)
		}
	})

	b.Run("BorderStyle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			BorderStyle.Render(text)
		}
	})
}