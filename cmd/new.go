package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/backsofangels/grimoire/internal/logging"
	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/validator"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var newCmd = &cobra.Command{
	Use:   "new [app-name]",
	Short: "Create a new project",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Determine default provider from flag (may be overridden by interactive selection)
		providerName, _ := cmd.Flags().GetString("provider")

		var cfg providers.ProviderConfig
		var provider providers.Provider
		var err error

		// Interactive wizard: no positional args AND no flags changed
		if len(args) == 0 && !anyFlagChanged(cmd) {
			// Top-level TUI: choose project type first
			projectType := "android"
			th := huh.ThemeBase()
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Project type").
						Options(
							huh.NewOption("Android — Android project (Kotlin/Java)", "android"),
							huh.NewOption("Java — Plain Java or Spring Boot", "java"),
						).
						Value(&projectType).
						Height(0),
				).Title("Step 1 — Project type"),
			).WithTheme(th)

			// Branded header for top-level wizard (match provider prompts)
			purple := lipgloss.Color("135")
			headerBanner := "🔮 grimoire — new project"
			headerStyle := lipgloss.NewStyle().Border(lipgloss.ThickBorder()).BorderLeft(true).BorderLeftForeground(purple).PaddingLeft(1).Bold(true).Foreground(lipgloss.Color("255"))
			subtitleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
			fmt.Println(headerStyle.Render(headerBanner))
			fmt.Println(subtitleStyle.Render("Use arrow keys · Enter to confirm · Ctrl+C to cancel"))
			fmt.Println()

			if err := form.Run(); err != nil {
				logging.Info("✗ Aborted: project creation cancelled by user")
				return
			}

			// Map selection to provider name
			if projectType == "java" {
				providerName = "springboot"
			} else {
				providerName = "android"
			}

			provider, err = providers.Get(providerName)
			if err != nil {
				logging.Error("✗ Provider not found", "provider", providerName, "error", err)
				return
			}

			// We've already rendered the branded header for the top-level wizard;
			// set a flag so provider prompts don't re-render it.
			_ = os.Setenv("GRIMOIRE_HEADER_PRINTED", "1")

			cfg, err = provider.Prompt()
			if err != nil {
				logging.Info("✗ Aborted: project creation cancelled by user")
				return
			}
		} else {
			// Non-interactive mode: accept either positional arg or --app-name flag
			provider, err = providers.Get(providerName)
			if err != nil {
				logging.Error("✗ Provider not found", "provider", providerName, "error", err)
				return
			}

			appNameFlag, _ := cmd.Flags().GetString("app-name")
			var appName string
			if appNameFlag != "" {
				appName = appNameFlag
			} else if len(args) > 0 {
				appName = args[0]
			} else {
				logging.Error("✗ Error: app name required")
				cmd.Usage()
				return
			}

			// Build config from flags
			cfg = providers.ProviderConfig{}
			cfg["AppName"] = appName

			// Pull provider-specific flags
			for _, f := range provider.Flags() {
				switch def := f.Default.(type) {
				case bool:
					if v, err := cmd.Flags().GetBool(f.Name); err == nil {
						cfg[f.Name] = v
					} else {
						cfg[f.Name] = def
					}
				case int:
					if v, err := cmd.Flags().GetInt(f.Name); err == nil {
						cfg[f.Name] = v
					} else {
						cfg[f.Name] = def
					}
				default:
					if v, err := cmd.Flags().GetString(f.Name); err == nil {
						cfg[f.Name] = v
					} else {
						// fallback to default printed value
						cfg[f.Name] = fmt.Sprintf("%v", def)
					}
				}
			}

			// Derived values
			if _, ok := cfg["AppNameLower"]; !ok {
				if an, _ := cfg["AppName"].(string); an != "" {
					cfg["AppNameLower"] = strings.ToLower(an)
				}
			}
			if _, ok := cfg["SanitizedAppName"]; !ok {
				if an, _ := cfg["AppName"].(string); an != "" {
					cfg["SanitizedAppName"] = validator.SanitizeAppName(an)
				}
			}
			if pkg, ok := cfg["package"].(string); ok && pkg != "" {
				cfg["PackageName"] = pkg
				cfg["PackagePath"] = validator.PackageToPath(pkg)
			}

			// Cast some well-known flags into canonical keys
			if v, ok := cfg["min-sdk"].(int); ok {
				cfg["MinSdk"] = v
			} else if s, ok := cfg["min-sdk"].(string); ok {
				if n, err := strconv.Atoi(s); err == nil {
					cfg["MinSdk"] = n
				}
			}
			if v, ok := cfg["target-sdk"].(int); ok {
				cfg["TargetSdk"] = v
			} else if s, ok := cfg["target-sdk"].(string); ok {
				if n, err := strconv.Atoi(s); err == nil {
					cfg["TargetSdk"] = n
				}
			}
		}

		// Validate
		if err := provider.Validate(cfg); err != nil {
			logging.Error("✗ Validation failed", "error", err)
			return
		}

		if err := provider.Generate(cfg); err != nil {
			logging.Error("✗ Generation failed", "error", err)
			return
		}

		// Render a branded success box to the terminal.
		printSuccess(cfg)
	},
}

// printSuccess renders a lipgloss-styled success box summarizing the created project.
func printSuccess(cfg providers.ProviderConfig) {
	an, _ := cfg["AppName"].(string)
	if an == "" {
		if s, _ := cfg["app-name"].(string); s != "" {
			an = s
		}
	}
	pkg, _ := cfg["PackageName"].(string)
	if pkg == "" {
		if s, _ := cfg["package"].(string); s != "" {
			pkg = s
		}
	}
	lang, _ := cfg["Lang"].(string)
	if lang == "" {
		if s, _ := cfg["lang"].(string); s != "" {
			lang = s
		}
	}
	minSdk := 0
	if v, ok := cfg["MinSdk"].(int); ok {
		minSdk = v
	} else if s, ok := cfg["min-sdk"].(string); ok {
		if n, err := strconv.Atoi(s); err == nil {
			minSdk = n
		}
	}
	templateKind, _ := cfg["Template"].(string)
	if templateKind == "" {
		if s, _ := cfg["template"].(string); s != "" {
			templateKind = s
		}
	}

	purple := lipgloss.NewStyle().Foreground(lipgloss.Color("135"))
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	bold := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255"))

	// Build content lines
	lines := []string{}
	lines = append(lines, bold.Render("✓  Project created"))
	lines = append(lines, "")
	lines = append(lines, purple.Copy().Bold(true).Render("  →  "+an))
	lines = append(lines, dim.Render("     "+pkg))
	lines = append(lines, dim.Render(fmt.Sprintf("     %s · API %d · %s template", lang, minSdk, templateKind)))
	lines = append(lines, "")
	lines = append(lines, dim.Render("  Next step:"))
	lines = append(lines, purple.Render(fmt.Sprintf("     cd %s && ./gradlew assembleDebug", an)))

	content := strings.Join(lines, "\n")

	box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("135")).Padding(1, 3)
	fmt.Println(box.Render(content))
}

func init() {
	// support explicit --app-name flag in addition to positional arg
	newCmd.Flags().String("app-name", "", "Application name (alternative to positional app-name)")
	// Register flags for all registered providers (deduplicate by flag name)
	seen := map[string]bool{}
	for _, p := range providers.All() {
		for _, f := range p.Flags() {
			if seen[f.Name] {
				continue
			}
			seen[f.Name] = true
			switch def := f.Default.(type) {
			case bool:
				newCmd.Flags().BoolP(f.Name, f.Short, def, f.Usage)
			case int:
				newCmd.Flags().IntP(f.Name, f.Short, def, f.Usage)
			default:
				newCmd.Flags().StringP(f.Name, f.Short, fmt.Sprintf("%v", def), f.Usage)
			}
		}
	}
	rootCmd.AddCommand(newCmd)
}

// anyFlagChanged returns true if the user explicitly set any flag on the command.
func anyFlagChanged(cmd *cobra.Command) bool {
	changed := false
	cmd.Flags().Visit(func(f *pflag.Flag) { changed = true })
	return changed
}
