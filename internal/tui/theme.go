package tui

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Theme holds Grimoire's consistent TUI styling.
type Theme struct {
	Purple    lipgloss.Color
	White     lipgloss.Color
	Dim       lipgloss.Color
	DimWhite  lipgloss.Color
	HuhTheme  *huh.Theme
	TitleStyle   lipgloss.Style
	SepStyle     lipgloss.Style
	HeaderStyle  lipgloss.Style
	SubtitleStyle lipgloss.Style
}

// NewTheme creates a new Grimoire-branded theme with purple/violet accents.
func NewTheme() *Theme {
	purple := lipgloss.Color("135")    // Violet
	white := lipgloss.Color("255")     // Bright white
	dim := lipgloss.Color("#666666")   // Gray
	dimWhite := lipgloss.Color("240")  // Dim white

	// Build base huh theme with Grimoire colors
	th := huh.ThemeBase()
	
	// Focused state: purple accent
	th.Focused.Base = th.Focused.Base.BorderForeground(purple)
	th.Focused.TextInput.Cursor = th.Focused.TextInput.Cursor.Foreground(purple)
	th.Focused.SelectSelector = th.Focused.SelectSelector.Foreground(purple)
	th.Focused.Option = th.Focused.Option.Foreground(purple).Bold(true)
	
	// Unfocused state: dimmed
	th.Focused.UnselectedOption = th.Focused.UnselectedOption.Foreground(dimWhite)
	th.Blurred.UnselectedOption = th.Blurred.UnselectedOption.Foreground(dimWhite)
	th.Blurred.Option = th.Blurred.Option.Foreground(dimWhite)
	
	// Text styles
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(white)
	sepStyle := lipgloss.NewStyle().Foreground(purple)
	th.Focused.Title = titleStyle
	th.Group.Title = titleStyle
	th.Group.Description = lipgloss.NewStyle().Foreground(dim)

	headerStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderLeft(true).
		BorderLeftForeground(purple).
		PaddingLeft(1).
		Bold(true).
		Foreground(white)

	subtitleStyle := lipgloss.NewStyle().Foreground(dimWhite)

	return &Theme{
		Purple:        purple,
		White:         white,
		Dim:           dim,
		DimWhite:      dimWhite,
		HuhTheme:      th,
		TitleStyle:    titleStyle,
		SepStyle:      sepStyle,
		HeaderStyle:   headerStyle,
		SubtitleStyle: subtitleStyle,
	}
}

// RenderGroupTitle renders a step title with a purple separator line underneath.
// Example: "Step 1 — App name" with a line of dashes.
func (t *Theme) RenderGroupTitle(s string) string {
	rendered := t.TitleStyle.Render(s)
	width := utf8.RuneCountInString(s)
	if width < 6 {
		width = 6
	}
	if width > 40 {
		width = 40
	}
	sep := strings.Repeat("─", width)
	return rendered + "\n" + t.SepStyle.Render(sep)
}

// PrintHeader prints the Grimoire branded header (called once per session).
// It's safe to call multiple times; subsequent calls are no-ops.
func (t *Theme) PrintHeader() {
	if os.Getenv("GRIMOIRE_HEADER_PRINTED") == "" {
		headerBanner := "🔮 grimoire — new project"
		fmt.Println(t.HeaderStyle.Render(headerBanner))
		fmt.Println(t.SubtitleStyle.Render("Use arrow keys · Enter to confirm · Ctrl+C to cancel"))
		fmt.Println()
		os.Setenv("GRIMOIRE_HEADER_PRINTED", "1")
	}
}

// ConfirmationStyle returns a styled confirmation note with the provided message.
func (t *Theme) ConfirmationStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.White).
		MarginLeft(1)
}
