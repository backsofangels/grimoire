package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/validator"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <app-name>",
	Short: "Create a new project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]

		providerName, _ := cmd.Flags().GetString("provider")
		provider, err := providers.Get(providerName)
		if err != nil {
			log.Error("provider not found", "provider", providerName, "error", err)
			return
		}

		// Build config from flags
		cfg := providers.ProviderConfig{}
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
			cfg["AppNameLower"] = strings.ToLower(appName)
		}
		if _, ok := cfg["SanitizedAppName"]; !ok {
			cfg["SanitizedAppName"] = validator.SanitizeAppName(appName)
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

		// Validate
		if err := provider.Validate(cfg); err != nil {
			log.Error("validation failed", "error", err)
			return
		}

		if err := provider.Generate(cfg); err != nil {
			log.Error("generate failed", "error", err)
			return
		}

		// Success summary (minimal)
		fmt.Printf("✓ Project \"%s\" created\n", appName)
	},
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
