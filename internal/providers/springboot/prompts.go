package springboot

import (
	"fmt"

	"github.com/backsofangels/grimoire/internal/config"
	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/tui"
	"github.com/backsofangels/grimoire/internal/validator"
	"github.com/charmbracelet/huh"
)

// RunPrompt displays a simple interactive wizard for Java/Spring projects.
func RunPrompt() (providers.ProviderConfig, error) {
	cfg, _ := config.Load()

	appName, group, artifact, packageName, outputDir, template, buildSystem, initGit := initialState(cfg)

	// Create theme
	theme := tui.NewTheme()
	theme.PrintHeader()

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
				Title("Group ID").
				Description("Reverse domain (e.g. com.example)").
				Value(&group),
		).Title(theme.RenderGroupTitle("Step 2 — Group ID")),

		huh.NewGroup(
			huh.NewInput().
				Title("Artifact ID").
				Description("Module/artifact name (e.g. myapp)").
				Value(&artifact),
		).Title(theme.RenderGroupTitle("Step 3 — Artifact ID")),

		huh.NewGroup(
			huh.NewInput().
				Title("Package name").
				Description("Leave empty to use <group>.<artifact>").
				Value(&packageName),
		).Title(theme.RenderGroupTitle("Step 4 — Package name")),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Project template").
				Options(
					huh.NewOption("Spring Boot — Spring Boot application", "springboot"),
					huh.NewOption("Plain Java — Standard jar application", "plain"),
				).
				Value(&template).
				Height(0),
		).Title(theme.RenderGroupTitle("Step 5 — Framework")),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Build system").
				Options(
					huh.NewOption("Gradle — Gradle build system", "gradle"),
					huh.NewOption("Maven (recommended) — Maven build system", "maven"),
				).
				Value(&buildSystem).
				Height(0),
		).Title(theme.RenderGroupTitle("Step 6 — Build system")),

		huh.NewGroup(
			huh.NewInput().
				Title("Output directory").
				Description("Leave empty to use ./<artifact>").
				Value(&outputDir),
		).Title(theme.RenderGroupTitle("Step 7 — Output directory")),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Initialize git repository?").
				Value(&initGit),
		).Title(theme.RenderGroupTitle("Step 8 — Initialize git")),

		huh.NewGroup(
			huh.NewNote().
				Title("✓ Ready to create").
				DescriptionFunc(func() string {
					pkg := packageName
					if pkg == "" {
						pkg = fmt.Sprintf("%s.%s", group, artifact)
					}
					return fmt.Sprintf("  App:      %s\n  Group:    %s\n  Artifact: %s\n  Package:  %s\n  Template: %s\n  Build:    %s\n  Git:      %s",
						appName, group, artifact, pkg, template, buildSystem, boolLabel(initGit))
				}, nil),
		).Title(theme.RenderGroupTitle("Step 9 — Confirm")),
	).WithTheme(theme.HuhTheme)

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
