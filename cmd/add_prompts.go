package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/backsofangels/grimoire/internal/logging"
	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/validator"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// runAddInteractive runs a short TUI to collect fields for adding a resource.
func runAddInteractive(cmd *cobra.Command, provider providers.Provider, kind string) error {
	// Prefill from flags if present
	pkgFlag, _ := cmd.Flags().GetString("package")
	moduleFlag, _ := cmd.Flags().GetString("module")
	langFlag, _ := cmd.Flags().GetString("lang")
	layoutFlag, _ := cmd.Flags().GetString("layout")
	uiFlag, _ := cmd.Flags().GetString("ui")
	noUIFlag, _ := cmd.Flags().GetBool("no-ui")
	overrideFlag, _ := cmd.Flags().GetBool("override")

	// Determine default module (prefer provided flag, else 'app' if exists, else '.')
	module := moduleFlag
	if module == "" {
		if _, err := os.Stat("app"); err == nil {
			module = "app"
		} else {
			module = "."
		}
	}

	// Detect package from manifest if not provided
	pkg := pkgFlag
	if pkg == "" {
		if p := detectPackage(module); p != "" {
			pkg = p
		}
	}

	// Detect Compose support
	composeOk := detectCompose(module)

	// defaults
	name := ""
	lang := langFlag
	if lang == "" {
		lang = "kotlin"
	}
	ui := uiFlag
	if noUIFlag {
		ui = "none"
	}
	if ui == "" {
		ui = "xml"
	}
	layout := layoutFlag
	includeVM := false
	di := "none"
	addNav := false
	override := overrideFlag

	// Build theme
	purple := lipgloss.Color("135")
	white := lipgloss.Color("255")
	dimWhite := lipgloss.Color("240")
	th := huh.ThemeBase()
	th.Focused.Base = th.Focused.Base.BorderForeground(purple)
	th.Focused.TextInput.Cursor = th.Focused.TextInput.Cursor.Foreground(purple)
	th.Focused.SelectSelector = th.Focused.SelectSelector.Foreground(purple)
	th.Focused.Option = th.Focused.Option.Foreground(purple).Bold(true)
	th.Focused.UnselectedOption = th.Focused.UnselectedOption.Foreground(dimWhite)
	th.Blurred.UnselectedOption = th.Blurred.UnselectedOption.Foreground(dimWhite)
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(white)
	sepStyle := lipgloss.NewStyle().Foreground(purple)
	th.Focused.Title = titleStyle
	th.Group.Title = titleStyle
	th.Group.Description = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	renderGroupTitle := func(s string) string {
		rendered := titleStyle.Render(s)
		width := len(s)
		if width < 6 {
			width = 6
		}
		if width > 36 {
			width = 36
		}
		sep := strings.Repeat("─", width)
		return rendered + "\n" + sepStyle.Render(sep)
	}

	headerBanner := fmt.Sprintf("🔮 grimoire — add %s", kind)
	headerStyle := lipgloss.NewStyle().Border(lipgloss.ThickBorder()).BorderLeft(true).BorderLeftForeground(purple).PaddingLeft(1).Bold(true).Foreground(white)
	subtitleStyle := lipgloss.NewStyle().Foreground(dimWhite)
	fmt.Println(headerStyle.Render(headerBanner))
	fmt.Println(subtitleStyle.Render("Fill fields — Enter to confirm · Ctrl+C to cancel"))
	fmt.Println()

	// Build form
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Name").Description("PascalCase (e.g. MyActivity)").Value(&name),
		).Title(renderGroupTitle("Step 1 — Name")),

		huh.NewGroup(
			huh.NewInput().Title("Package").Description("Target package (detected from manifest)").Value(&pkg),
		).Title(renderGroupTitle("Step 2 — Package")),

		huh.NewGroup(
			huh.NewInput().Title("Module").Description("Target module folder").Value(&module),
		).Title(renderGroupTitle("Step 3 — Module")),

		huh.NewGroup(
			huh.NewSelect[string]().Title("Language").Options(
				huh.NewOption("Kotlin", "kotlin"),
				huh.NewOption("Java", "java"),
			).Value(&lang),
		).Title(renderGroupTitle("Step 4 — Language")),

		huh.NewGroup(
			func() huh.Field {
				sel := huh.NewSelect[string]().Title("UI Type").Value(&ui)
				sel.Options(huh.NewOption("XML layout", "xml"))
				sel.Options(huh.NewOption("None", "none"))
				if composeOk {
					sel.Options(huh.NewOption("Jetpack Compose", "compose"))
				}
				return sel
			}(),
		).Title(renderGroupTitle("Step 5 — UI")),

		huh.NewGroup(
			huh.NewConfirm().Title("Include ViewModel?").Value(&includeVM),
		).Title(renderGroupTitle("Step 6 — ViewModel")),

		huh.NewGroup(
			huh.NewSelect[string]().Title("Dependency injection").Options(
				huh.NewOption("None", "none"),
				huh.NewOption("Hilt", "hilt"),
				huh.NewOption("Koin", "koin"),
			).Value(&di),
		).Title(renderGroupTitle("Step 7 — DI")),

		huh.NewGroup(
			huh.NewConfirm().Title("Add navigation entry?").Value(&addNav),
		).Title(renderGroupTitle("Step 8 — Navigation")),

		huh.NewGroup(
			huh.NewConfirm().Title("Overwrite existing files if any?").Value(&override),
		).Title(renderGroupTitle("Step 9 — Conflict handling")),

		huh.NewGroup(
			huh.NewNote().Title("✓ Ready").DescriptionFunc(func() string {
				return fmt.Sprintf("  Kind: %s\n  Name: %s\n  Package: %s\n  Module: %s\n  Lang: %s\n  UI: %s\n  ViewModel: %v\n  DI: %s\n  Navigation: %v\n  Overwrite: %v",
					kind, name, pkg, module, lang, ui, includeVM, di, addNav, override)
			}, nil),
		).Title(renderGroupTitle("Step 10 — Confirm")),
	).WithTheme(th)

	if err := form.Run(); err != nil {
		return err
	}

	// Build cfg and call provider.Add
	if ui == "none" {
		// don't pass a layout name when user chooses no UI
		if layout != "" {
			logging.Info("Ignoring layout because UI is 'none'")
		}
		layout = ""
	}

	// Validate inputs
	// Only validate explicit UI values; interactive selection may set 'none'
	if strings.ToLower(ui) != "none" {
		if err := validator.ValidateUI(ui); err != nil {
			return err
		}
	}
	if err := validator.ValidateLanguage(lang); err != nil {
		return err
	}
	if err := validator.ValidateDI(di); err != nil {
		return err
	}

	cfg := providers.ProviderConfig{
		"Kind":        kind,
		"Name":        name,
		"PackageName": pkg,
		"Module":      module,
		"Lang":        lang,
		"Layout":      layout,
		"UI":          ui,
		"Override":    override,
		"ViewModel":   includeVM,
		"DI":          di,
		"Nav":         addNav,
	}

	if err := provider.Add(cfg); err != nil {
		return err
	}

	logging.Success(fmt.Sprintf("Added %s %s", kind, name))
	return nil
}

// detectPackage tries to read package from module/src/main/AndroidManifest.xml
func detectPackage(module string) string {
	// 1) AndroidManifest.xml package attribute
	man := filepath.Join(module, "src", "main", "AndroidManifest.xml")
	if b, err := os.ReadFile(man); err == nil {
		re := regexp.MustCompile(`package\s*=\s*"([^"]+)"`)
		if m := re.FindStringSubmatch(string(b)); len(m) > 1 {
			return m[1]
		}
	}

	// helper to parse namespace or applicationId from Gradle files
	parseFromGradle := func(content string) string {
		reNs := regexp.MustCompile(`namespace\s*(?:=)?\s*["']([^"']+)["']`)
		if m := reNs.FindStringSubmatch(content); len(m) > 1 {
			return m[1]
		}
		reApp := regexp.MustCompile(`applicationId\s*(?:=)?\s*["']([^"']+)["']`)
		if m := reApp.FindStringSubmatch(content); len(m) > 1 {
			return m[1]
		}
		return ""
	}

	// 2) module build.gradle
	build := filepath.Join(module, "build.gradle")
	if b, err := os.ReadFile(build); err == nil {
		if p := parseFromGradle(string(b)); p != "" {
			return p
		}
	}
	// 3) module build.gradle.kts
	buildKts := filepath.Join(module, "build.gradle.kts")
	if b, err := os.ReadFile(buildKts); err == nil {
		if p := parseFromGradle(string(b)); p != "" {
			return p
		}
	}

	// 4) scan source files for package declaration
	roots := []string{
		filepath.Join(module, "src", "main", "java"),
		filepath.Join(module, "src", "main", "kotlin"),
	}
	stopErr := errors.New("__pkg_found__")
	var found string
	for _, root := range roots {
		_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".kt") || strings.HasSuffix(path, ".java") {
				if b, err := os.ReadFile(path); err == nil {
					re := regexp.MustCompile(`(?m)^\s*package\s+([a-zA-Z0-9_.]+)`)
					if m := re.FindStringSubmatch(string(b)); len(m) > 1 {
						found = m[1]
						return stopErr
					}
				}
			}
			return nil
		})
		if found != "" {
			return found
		}
	}
	return ""
}

// detectCompose does a simple check for compose usage in module build files
func detectCompose(module string) bool {
	build := filepath.Join(module, "build.gradle")
	b, err := os.ReadFile(build)
	if err == nil {
		s := string(b)
		if strings.Contains(s, "compose true") || strings.Contains(s, "androidx.compose") {
			return true
		}
	}
	// also check module gradle.kts variant
	buildKts := filepath.Join(module, "build.gradle.kts")
	if b, err := os.ReadFile(buildKts); err == nil {
		s := string(b)
		if strings.Contains(s, "compose") || strings.Contains(s, "androidx.compose") {
			return true
		}
	}
	return false
}
