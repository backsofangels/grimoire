package springboot

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/backsofangels/grimoire/internal/config"
	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/validator"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// RunPrompt displays a simple interactive wizard for Java/Spring projects.
func RunPrompt() (providers.ProviderConfig, error) {
	cfg, _ := config.Load()

	appName, group, artifact, packageName, outputDir, template, buildSystem, initGit := initialState(cfg)

	// Minimal theme
	th := huh.ThemeBase()
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255"))
	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("135"))
	renderGroupTitle := func(s string) string {
		rendered := titleStyle.Render(s)
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

	// Header (skip if top-level wizard already printed it)
	if os.Getenv("GRIMOIRE_HEADER_PRINTED") == "" {
		headerBanner := "🔮 grimoire — new project"
		headerStyle := lipgloss.NewStyle().Border(lipgloss.ThickBorder()).BorderLeft(true).BorderLeftForeground(lipgloss.Color("135")).PaddingLeft(1).Bold(true).Foreground(lipgloss.Color("255"))
		subtitleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		fmt.Println(headerStyle.Render(headerBanner))
		fmt.Println(subtitleStyle.Render("Use arrow keys · Enter to confirm · Ctrl+C to cancel"))
		fmt.Println()
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("App name").
				Description("PascalCase recommended (e.g. MyApp)").
				Value(&appName).
				Validate(func(s string) error { return validator.ValidateAppName(s) }),
		).Title(renderGroupTitle("Step 2 — App name")),

		huh.NewGroup(
			huh.NewInput().
				Title("Group ID").
				Description("Reverse domain (e.g. com.example)").
				Value(&group),
		).Title(renderGroupTitle("Step 3 — Group ID")),

		huh.NewGroup(
			huh.NewInput().
				Title("Artifact ID").
				Description("Module/artifact name (e.g. myapp)").
				Value(&artifact),
		).Title(renderGroupTitle("Step 4 — Artifact ID")),

		huh.NewGroup(
			huh.NewInput().
				Title("Package name").
				Description("Leave empty to use <group>.<artifact>").
				Value(&packageName),
		).Title(renderGroupTitle("Step 5 — Package name")),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Project template").
				Options(
					huh.NewOption("Spring Boot — Spring Boot application", "springboot"),
					huh.NewOption("Plain Java — Standard jar application", "plain"),
				).
				Value(&template).
				Height(0),
		).Title(renderGroupTitle("Step 6 — Framework")),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Build system").
				Options(
					huh.NewOption("Gradle — Gradle build system", "gradle"),
					huh.NewOption("Maven (recommended) — Maven build system", "maven"),
				).
				Value(&buildSystem).
				Height(0),
		).Title(renderGroupTitle("Step 7 — Build system")),

		huh.NewGroup(
			huh.NewInput().
				Title("Output directory").
				Description("Leave empty to use ./<artifact>").
				Value(&outputDir),
		).Title(renderGroupTitle("Step 7 — Output directory")),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Initialize git repository?").
				Value(&initGit),
		).Title(renderGroupTitle("Step 8 — Initialize git")),

		huh.NewGroup(
			huh.NewNote().
				Title("✓ Ready to create").
				DescriptionFunc(func() string {
					pkg := packageName
					if pkg == "" {
						pkg = fmt.Sprintf("%s.%s", group, artifact)
					}
					return fmt.Sprintf("  App:      %s\n  Group:    %s\n  Artifact: %s\n  Package:  %s\n  Template: %s\n  Git:      %s",
						appName, group, artifact, pkg, template, boolLabel(initGit))
				}, nil),
		).Title(renderGroupTitle("Step 10 — Confirm")),
	).WithTheme(th)

	if err := form.Run(); err != nil {
		return nil, err
	}

	if packageName == "" {
		packageName = fmt.Sprintf("%s.%s", group, artifact)
	}
	if outputDir == "" {
		outputDir = "./" + artifact
	}

	return providers.ProviderConfig{
		"AppName":     appName,
		"Group":       group,
		"Artifact":    artifact,
		"PackageName": packageName,
		"OutputDir":   outputDir,
		"Template":    template,
		"BuildSystem": buildSystem,
		"Git":         initGit,
	}, nil
}

func boolLabel(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

func initialState(cfg config.Config) (appName, group, artifact, packageName, outputDir, template, buildSystem string, initGit bool) {
	appName = ""
	if cfg.DefaultPackage == "" {
		cfg.DefaultPackage = "com.example"
	}
	group = cfg.DefaultPackage
	artifact = "app"
	packageName = ""
	outputDir = ""
	if cfg.DefaultTemplate == "" {
		template = "springboot"
	} else {
		template = cfg.DefaultTemplate
	}
	// Default build system for Java projects
	buildSystem = "gradle"
	initGit = cfg.Git
	return
}
