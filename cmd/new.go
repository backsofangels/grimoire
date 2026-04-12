package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/backsofangels/grimoire/internal/logging"
	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/validator"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var newCmd = &cobra.Command{
	Use:   "new [app-name]",
	Short: "Create a new project",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		providerName, _ := cmd.Flags().GetString("provider")
		provider, err := providers.Get(providerName)
		if err != nil {
			logging.Error("✗ Provider not found", "provider", providerName, "error", err)
			return
		}

		var cfg providers.ProviderConfig

		// Interactive wizard: no positional args AND no flags changed
		if len(args) == 0 && !anyFlagChanged(cmd) {
			var perr error
			cfg, perr = provider.Prompt()
			if perr != nil {
				// user aborted or prompt error
				logging.Info("✗ Aborted: project creation cancelled by user")
				return
			}
		} else {
			// Non-interactive mode: positional arg (app name) required
			if len(args) == 0 {
				logging.Error("✗ Error: app name required")
				cmd.Usage()
				return
			}
			appName := args[0]

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
	// Register flags for known providers dynamically (android exists by default)
	if p, err := providers.Get("android"); err == nil {
		for _, f := range p.Flags() {
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
