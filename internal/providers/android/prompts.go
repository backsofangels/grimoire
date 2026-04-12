package android

import (
	"fmt"
	"strconv"

	"github.com/backsofangels/grimoire/internal/config"
	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/tui"
	"github.com/backsofangels/grimoire/internal/validator"
	"github.com/charmbracelet/huh"
)

// Prompt displays an interactive wizard for the Android provider.
// It implements the 4-step form plus a confirmation screen described in the
// project plan. All fields are prefilled from `~/.grimoire/config.json` when
// available; otherwise sensible defaults are used.
func (p *AndroidProvider) Prompt() (providers.ProviderConfig, error) {
	// Load defaults from user config (silently fall back to defaults).
	cfg, _ := config.Load()

	// initialize form values from config
	appName, packageName, outputDir, lang, minSdkStr, template, useWrapper, initGit, vscode := initialState(cfg)

	// Create theme
	theme := tui.NewTheme()
	theme.PrintHeader()

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
		).Title(theme.RenderGroupTitle("Step 1 — App name")),

		huh.NewGroup(
			huh.NewInput().
				Title("Package name").
				Description("Reverse domain notation (e.g. com.example.myapp)").
				Value(&packageName).
				Validate(func(s string) error { return validator.ValidatePackageName(s) }),
		).Title(theme.RenderGroupTitle("Step 2 — Package name")),

		huh.NewGroup(
			huh.NewInput().
				Title("Output directory").
				Description("Leave empty to use ./<app-name>").
				Value(&outputDir),
		).Title(theme.RenderGroupTitle("Step 3 — Output directory")),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Language").
				Options(
					huh.NewOption("Kotlin", "kotlin"),
					huh.NewOption("Java", "java"),
				).
				Value(&lang),
		).Title(theme.RenderGroupTitle("Step 4 — Language")),

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
		).Title(theme.RenderGroupTitle("Step 5 — Minimum SDK")),

		huh.NewGroup(templateSelect).Title(theme.RenderGroupTitle("Step 6 — Template")),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Use Gradle wrapper (gradlew)?").
				Description("Builds will be reproducible across machines. Wrapper (9.4.1) included. Choose 'no' to use system gradle.").
				Value(&useWrapper),
		).Title(theme.RenderGroupTitle("Step 7 — Gradle wrapper")),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Initialize git repository?").
				Value(&initGit),
		).Title(theme.RenderGroupTitle("Step 8 — Initialize git")),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Generate .vscode/ config?").
				Value(&vscode),
		).Title(theme.RenderGroupTitle("Step 9 — VSCode config")),

		huh.NewGroup(
			huh.NewNote().
				Title("✓ Ready to create").
				DescriptionFunc(func() string {
					minSdkInt := mustAtoi(minSdkStr)
					return fmt.Sprintf(
						"  App:       %s\n  Package:   %s\n  Language:  %s\n  Min SDK:   %s (%s)\n  Template:  %s\n  Wrapper:   %s\n  Git:       %s\n  VSCode:    %s",
						appName,
						packageName,
						lang,
						minSdkStr, validator.SdkVersionLabel(minSdkInt),
						template,
						boolLabel(useWrapper),
						boolLabel(initGit),
						boolLabel(vscode),
					)
				}, nil),
		).Title(theme.RenderGroupTitle("Step 10 — Confirm")),
	).WithTheme(theme.HuhTheme)

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
		"Wrapper":     useWrapper,
		"Git":         initGit,
		"Vscode":      vscode,
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
func initialState(cfg config.Config) (appName, packageName, outputDir, lang, minSdkStr, template string, useWrapper, initGit, vscode bool) {
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
	useWrapper = true // Default to using Gradle wrapper for reproducible builds
	initGit = cfg.Git
	vscode = cfg.VSCode
	return
}
