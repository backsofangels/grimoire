package android

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/backsofangels/grimoire/internal/config"
	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/validator"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Prompt displays an interactive wizard for the Android provider.
// It implements the 4-step form plus a confirmation screen described in the
// project plan. All fields are prefilled from `~/.grimoire/config.json` when
// available; otherwise sensible defaults are used.
func (p *AndroidProvider) Prompt() (providers.ProviderConfig, error) {
	// Load defaults from user config (silently fall back to defaults).
	cfg, _ := config.Load()

	// initialize form values from config
	appName, packageName, outputDir, lang, minSdkStr, template, initGit, vscode := initialState(cfg)

	// Build a custom theme that matches Grimoire's identity. We keep the
	// terminal background (transparent) and tune accent/focus colors.
	purple := lipgloss.Color("135")
	white := lipgloss.Color("255")
	dim := lipgloss.Color("#666666")
	// dimWhite used for subtitle and other subtle text
	dimWhite := lipgloss.Color("240")

	th := huh.ThemeBase()
	// Focused border accent
	th.Focused.Base = th.Focused.Base.BorderForeground(purple)
	// Input cursor color
	th.Focused.TextInput.Cursor = th.Focused.TextInput.Cursor.Foreground(purple)
	// Active selection (option) purple + bold
	th.Focused.SelectSelector = th.Focused.SelectSelector.Foreground(purple)
	th.Focused.Option = th.Focused.Option.Foreground(purple).Bold(true)
	// Inactive selection: dimmed white
	th.Focused.UnselectedOption = th.Focused.UnselectedOption.Foreground(dimWhite)
	th.Blurred.UnselectedOption = th.Blurred.UnselectedOption.Foreground(dimWhite)
	th.Blurred.Option = th.Blurred.Option.Foreground(dimWhite)
	// Titles and descriptions
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(white)
	sepStyle := lipgloss.NewStyle().Foreground(purple)
	th.Focused.Title = titleStyle
	th.Group.Title = titleStyle
	th.Group.Description = lipgloss.NewStyle().Foreground(dim)

	// Helper to render a styled group title with a short separator underneath.
	renderGroupTitle := func(s string) string {
		// Render the bold white title
		rendered := titleStyle.Render(s)
		// Separator length matches title rune length (clamped)
		width := utf8.RuneCountInString(s)
		if width < 6 {
			width = 6
		}
		if width > 36 {
			width = 36
		}
		sep := strings.Repeat("─", width)
		return rendered + "\n" + sepStyle.Render(sep)
	}

	// Branded header printed once before the form runs (skip if already printed)
	if os.Getenv("GRIMOIRE_HEADER_PRINTED") == "" {
		headerBanner := "🔮 grimoire — new project"
		headerStyle := lipgloss.NewStyle().Border(lipgloss.ThickBorder()).BorderLeft(true).BorderLeftForeground(purple).PaddingLeft(1).Bold(true).Foreground(white)
		subtitleStyle := lipgloss.NewStyle().Foreground(dimWhite)
		fmt.Println(headerStyle.Render(headerBanner))
		fmt.Println(subtitleStyle.Render("Use arrow keys · Enter to confirm · Ctrl+C to cancel"))
		fmt.Println()
	}

	// Build the form with one field per step so the TUI shows a single
	// prompt at a time rather than multiple prompts on the same screen.

	// Template select is dynamic: show 'compose' only when language == kotlin.
	templateSelect := huh.NewSelect[string]().
		Title("Project template").
		Value(&template)
	templateSelect.OptionsFunc(func() []huh.Option[string] {
		opts := []huh.Option[string]{
			huh.NewOption("basic  — Activity + layout XML + ViewModel", "basic"),
		}
		if lang == "kotlin" {
			opts = append(opts, huh.NewOption("compose — Activity + Compose layout + ViewModel", "compose"))
		}
		opts = append(opts, huh.NewOption("empty  — Bare MainActivity, no layout", "empty"))
		return opts
	}, &lang).Height(0)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("App name").
				Description("PascalCase recommended (e.g. MyApp)").
				Value(&appName).
				Validate(func(s string) error { return validator.ValidateAppName(s) }),
		).Title(renderGroupTitle("Step 1 — App name")),

		huh.NewGroup(
			huh.NewInput().
				Title("Package name").
				Description("Reverse domain notation (e.g. com.example.myapp)").
				Value(&packageName).
				Validate(func(s string) error { return validator.ValidatePackageName(s) }),
		).Title(renderGroupTitle("Step 2 — Package name")),

		huh.NewGroup(
			huh.NewInput().
				Title("Output directory").
				Description("Leave empty to use ./<app-name>").
				Value(&outputDir),
		).Title(renderGroupTitle("Step 3 — Output directory")),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Language").
				Options(
					huh.NewOption("Kotlin", "kotlin"),
					huh.NewOption("Java", "java"),
				).
				Value(&lang),
		).Title(renderGroupTitle("Step 4 — Language")),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Minimum SDK").
				Options(
					huh.NewOption("API 21 — Android 5.0 Lollipop", "21"),
					huh.NewOption("API 24 — Android 7.0 Nougat", "24"),
					huh.NewOption("API 26 — Android 8.0 Oreo", "26"),
					huh.NewOption("API 28 — Android 9.0 Pie", "28"),
					huh.NewOption("API 33 — Android 13", "33"),
					huh.NewOption("API 35 — Android 15", "35"),
				).
				Value(&minSdkStr),
		).Title(renderGroupTitle("Step 5 — Minimum SDK")),

		huh.NewGroup(templateSelect).Title(renderGroupTitle("Step 6 — Template")),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Initialize git repository?").
				Value(&initGit),
		).Title(renderGroupTitle("Step 7 — Initialize git")),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Generate .vscode/ config?").
				Value(&vscode),
		).Title(renderGroupTitle("Step 8 — VSCode config")),

		huh.NewGroup(
			huh.NewNote().
				Title("✓ Ready to create").
				DescriptionFunc(func() string {
					minSdkInt := mustAtoi(minSdkStr)
					return fmt.Sprintf(
						"  App:       %s\n  Package:   %s\n  Language:  %s\n  Min SDK:   %s (%s)\n  Template:  %s\n  Git:       %s\n  VSCode:    %s",
						appName,
						packageName,
						lang,
						minSdkStr, validator.SdkVersionLabel(minSdkInt),
						template,
						boolLabel(initGit),
						boolLabel(vscode),
					)
				}, nil),
		).Title(renderGroupTitle("Step 9 — Confirm")),
	).WithTheme(th)

	// Run the interactive form; return any error (e.g. user aborted with Ctrl+C).
	if err := form.Run(); err != nil {
		return nil, err
	}

	if outputDir == "" {
		outputDir = "./" + appName
	}
	minSdk, _ := strconv.Atoi(minSdkStr)

	return providers.ProviderConfig{
		"AppName":     appName,
		"PackageName": packageName,
		"OutputDir":   outputDir,
		"Lang":        lang,
		"MinSdk":      minSdk,
		"TargetSdk":   35,
		"Template":    template,
		"Git":         initGit,
		"Vscode":      vscode,
		"NoWrapper":   false,
	}, nil
}

func boolLabel(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

func mustAtoi(s string) int {
	// Min return to 26 to ensure minimum SDK for android
	if s == "" {
		return 26
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 26
	}
	return v
}

// initialState returns the initial wizard values populated from the provided
// Config. Extracted to allow unit testing of prefill behavior without running
// the interactive form.
func initialState(cfg config.Config) (appName, packageName, outputDir, lang, minSdkStr, template string, initGit, vscode bool) {
	appName = ""
	if cfg.DefaultPackage == "" {
		cfg.DefaultPackage = "com.example"
	}
	packageName = cfg.DefaultPackage + ".myapp"
	outputDir = ""
	if cfg.DefaultLang == "" {
		lang = "kotlin"
	} else {
		lang = cfg.DefaultLang
	}
	if cfg.DefaultMinSdk == 0 {
		cfg.DefaultMinSdk = 26
	}
	minSdkStr = strconv.Itoa(cfg.DefaultMinSdk)
	if cfg.DefaultTemplate == "" {
		template = "basic"
	} else {
		template = cfg.DefaultTemplate
	}
	initGit = cfg.Git
	vscode = cfg.VSCode
	return
}
