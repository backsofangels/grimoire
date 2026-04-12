package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/backsofangels/grimoire/internal/logging"
	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a resource to an existing project (activity|fragment|viewmodel)",
	Long:  "Add a resource to an existing project (activity|fragment|viewmodel). Use --ui (xml|compose|none) to select UI type; --no-ui is a shortcut for --ui none.",
	Run: func(cmd *cobra.Command, args []string) {
		providerName, _ := cmd.Flags().GetString("provider")
		provider, err := providers.Get(providerName)
		if err != nil {
			logging.Error("✗ Provider not found", "provider", providerName, "error", err)
			return
		}

		// selection TUI
		kind := "activity"
		th := huh.ThemeBase()
		sel := huh.NewSelect[string]().Title("What would you like to add?")
		sel.Options(
			huh.NewOption("Activity", "activity"),
			huh.NewOption("Fragment", "fragment"),
			huh.NewOption("ViewModel", "viewmodel"),
		).Value(&kind)

		form := huh.NewForm(
			huh.NewGroup(sel).Title("Choose resource"),
		).WithTheme(th)

		fmt.Println("🔮 grimoire — add")
		if err := form.Run(); err != nil {
			logging.Info("Interrupted by user")
			return
		}

		if err := runAddInteractive(cmd, provider, kind); err != nil {
			if isUserAbort(err) {
				logging.Info("Interrupted by user")
				return
			}
			logging.Error("✗ Add failed", "error", err)
			return
		}
	},
}

var addActivityCmd = &cobra.Command{
	Use:   "activity [name]",
	Short: "Add an Activity",
	Long:  "Add an Activity. Use --ui (xml|compose|none) to control UI generation or --no-ui to skip UI generation.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		providerName, _ := cmd.Flags().GetString("provider")
		provider, err := providers.Get(providerName)
		if err != nil {
			return fmt.Errorf("provider not found: %s: %w", providerName, err)
		}

		// allow --name flag as alternative to positional arg; interactive when both absent
		nameFlag, _ := cmd.Flags().GetString("name")
		if len(args) == 0 && nameFlag == "" {
			// write a quick marker so we can detect whether this Run executed
			_ = os.WriteFile(".add_run_marker", []byte("run"), 0o644)

			if err := runAddInteractive(cmd, provider, "activity"); err != nil {
				if isUserAbort(err) {
					logging.Info("Interrupted by user")
					return nil
				}
				return fmt.Errorf("add failed: %w", err)
			}
			return nil
		}

		var name string
		if nameFlag != "" {
			name = nameFlag
		} else {
			name = args[0]
		}
		pkg, _ := cmd.Flags().GetString("package")
		module, _ := cmd.Flags().GetString("module")
		lang, _ := cmd.Flags().GetString("lang")
		layout, _ := cmd.Flags().GetString("layout")
		ui, _ := cmd.Flags().GetString("ui")
		// support --no-ui alias which forces UI to "none"
		noUI, _ := cmd.Flags().GetBool("no-ui")
		if noUI {
			ui = "none"
		}

		// normalize/validate inputs
		if strings.ToLower(ui) == "none" && layout != "" {
			logging.Info("Ignoring --layout because --ui is 'none'")
			layout = ""
		}
		if err := validateUI(ui); err != nil {
			return fmt.Errorf("invalid --ui: %w", err)
		}
		if err := validateLang(lang); err != nil {
			return fmt.Errorf("invalid --lang: %w", err)
		}
		override, _ := cmd.Flags().GetBool("override")

		di, _ := cmd.Flags().GetString("di")
		if err := validateDI(di); err != nil {
			return fmt.Errorf("invalid --di: %w", err)
		}
		includeVM, _ := cmd.Flags().GetBool("viewmodel")
		addNav, _ := cmd.Flags().GetBool("nav")

		cfg := providers.ProviderConfig{
			"Kind":        "activity",
			"Name":        name,
			"PackageName": pkg,
			"Module":      module,
			"Lang":        lang,
			"Layout":      layout,
			"UI":          ui,
			"Override":    override,
			"DI":          di,
			"ViewModel":   includeVM,
			"Nav":         addNav,
		}

		if err := provider.Add(cfg); err != nil {
			return fmt.Errorf("add failed: %w", err)
		}
		logging.Success(fmt.Sprintf("Added activity %s", name))
		return nil
	},
}

var addFragmentCmd = &cobra.Command{
	Use:   "fragment [name]",
	Short: "Add a Fragment",
	Long:  "Add a Fragment. Use --ui (xml|compose|none) to control UI generation or --no-ui to skip UI generation.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		providerName, _ := cmd.Flags().GetString("provider")
		provider, err := providers.Get(providerName)
		if err != nil {
			return fmt.Errorf("provider not found: %s: %w", providerName, err)
		}

		// allow --name flag as alternative to positional arg; interactive when both absent
		nameFlag, _ := cmd.Flags().GetString("name")
		if len(args) == 0 && nameFlag == "" {
			if err := runAddInteractive(cmd, provider, "fragment"); err != nil {
				if isUserAbort(err) {
					logging.Info("Interrupted by user")
					return nil
				}
				return fmt.Errorf("add failed: %w", err)
			}
			return nil
		}

		var name string
		if nameFlag != "" {
			name = nameFlag
		} else {
			name = args[0]
		}
		pkg, _ := cmd.Flags().GetString("package")
		module, _ := cmd.Flags().GetString("module")
		lang, _ := cmd.Flags().GetString("lang")
		layout, _ := cmd.Flags().GetString("layout")
		ui, _ := cmd.Flags().GetString("ui")
		// support --no-ui alias which forces UI to "none"
		noUI, _ := cmd.Flags().GetBool("no-ui")
		if noUI {
			ui = "none"
		}

		// normalize/validate inputs
		if strings.ToLower(ui) == "none" && layout != "" {
			logging.Info("Ignoring --layout because --ui is 'none'")
			layout = ""
		}
		if err := validateUI(ui); err != nil {
			return fmt.Errorf("invalid --ui: %w", err)
		}
		if err := validateLang(lang); err != nil {
			return fmt.Errorf("invalid --lang: %w", err)
		}
		override, _ := cmd.Flags().GetBool("override")

		di, _ := cmd.Flags().GetString("di")
		if err := validateDI(di); err != nil {
			return fmt.Errorf("invalid --di: %w", err)
		}
		includeVM, _ := cmd.Flags().GetBool("viewmodel")
		addNav, _ := cmd.Flags().GetBool("nav")

		cfg := providers.ProviderConfig{
			"Kind":        "fragment",
			"Name":        name,
			"PackageName": pkg,
			"Module":      module,
			"Lang":        lang,
			"Layout":      layout,
			"UI":          ui,
			"Override":    override,
			"DI":          di,
			"ViewModel":   includeVM,
			"Nav":         addNav,
		}

		if err := provider.Add(cfg); err != nil {
			return fmt.Errorf("add failed: %w", err)
		}
		logging.Success(fmt.Sprintf("Added fragment %s", name))
		return nil
	},
}

var addViewModelCmd = &cobra.Command{
	Use:   "viewmodel [name]",
	Short: "Add a ViewModel",
	Long:  "Add a ViewModel. Use --ui (xml|compose|none) or --no-ui to indicate no UI should be generated (viewmodels typically have no UI).",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		providerName, _ := cmd.Flags().GetString("provider")
		provider, err := providers.Get(providerName)
		if err != nil {
			return fmt.Errorf("provider not found: %s: %w", providerName, err)
		}

		// allow --name flag as alternative to positional arg; interactive when both absent
		nameFlag, _ := cmd.Flags().GetString("name")
		if len(args) == 0 && nameFlag == "" {
			if err := runAddInteractive(cmd, provider, "viewmodel"); err != nil {
				if isUserAbort(err) {
					logging.Info("Interrupted by user")
					return nil
				}
				return fmt.Errorf("add failed: %w", err)
			}
			return nil
		}

		var name string
		if nameFlag != "" {
			name = nameFlag
		} else {
			name = args[0]
		}
		pkg, _ := cmd.Flags().GetString("package")
		module, _ := cmd.Flags().GetString("module")
		lang, _ := cmd.Flags().GetString("lang")
		ui, _ := cmd.Flags().GetString("ui")
		// support --no-ui alias which forces UI to "none"
		noUI, _ := cmd.Flags().GetBool("no-ui")
		if noUI {
			ui = "none"
		}
		override, _ := cmd.Flags().GetBool("override")

		// validate inputs for viewmodel command
		if err := validateUI(ui); err != nil {
			return fmt.Errorf("invalid --ui: %w", err)
		}
		if err := validateLang(lang); err != nil {
			return fmt.Errorf("invalid --lang: %w", err)
		}

		cfg := providers.ProviderConfig{
			"Kind":        "viewmodel",
			"Name":        name,
			"PackageName": pkg,
			"Module":      module,
			"Lang":        lang,
			"UI":          ui,
			"Override":    override,
		}

		if err := provider.Add(cfg); err != nil {
			return fmt.Errorf("add failed: %w", err)
		}
		logging.Success(fmt.Sprintf("Added viewmodel %s", name))
		return nil
	},
}

func init() {
	// Register add command and subcommands
	rootCmd.AddCommand(addCmd)
	addCmd.AddCommand(addActivityCmd)
	addCmd.AddCommand(addFragmentCmd)
	addCmd.AddCommand(addViewModelCmd)

	// common persistent flags on top-level `add` so `grimoire add` can prefill
	addCmd.PersistentFlags().String("name", "", "Resource name (alternative to positional)")
	addCmd.PersistentFlags().String("di", "none", "Dependency injection (none|hilt|koin)")
	addCmd.PersistentFlags().Bool("viewmodel", false, "Generate associated ViewModel")
	addCmd.PersistentFlags().Bool("nav", false, "Add navigation entry")

	addCmd.PersistentFlags().StringP("package", "p", "", "Target package (e.g. com.example.app)")
	addCmd.PersistentFlags().StringP("module", "m", "app", "Target module (default: app)")
	addCmd.PersistentFlags().StringP("lang", "l", "kotlin", "Language (kotlin|java)")
	addCmd.PersistentFlags().StringP("layout", "", "", "Layout resource name (optional)")
	addCmd.PersistentFlags().StringP("ui", "", "xml", "UI type (xml|compose|none). Use --no-ui as a shortcut for --ui none")
	addCmd.PersistentFlags().Bool("no-ui", false, "Shortcut for --ui none (equivalent to --ui none)")
	addCmd.PersistentFlags().BoolP("override", "", false, "Overwrite existing files")
}

// validateUI ensures the provided UI type is one of the allowed values.
func validateUI(ui string) error {
	s := strings.ToLower(strings.TrimSpace(ui))
	if s == "" {
		return nil
	}
	switch s {
	case "xml", "compose", "none":
		return nil
	default:
		return fmt.Errorf("invalid UI type: %s (allowed: xml|compose|none)", ui)
	}
}

// validateLang ensures language is kotlin or java.
func validateLang(lang string) error {
	s := strings.ToLower(strings.TrimSpace(lang))
	if s == "" {
		return nil
	}
	switch s {
	case "kotlin", "java":
		return nil
	default:
		return fmt.Errorf("invalid language: %s (allowed: kotlin|java)", lang)
	}
}

// validateDI ensures DI selection is valid.
func validateDI(di string) error {
	s := strings.ToLower(strings.TrimSpace(di))
	if s == "" {
		return nil
	}
	switch s {
	case "none", "hilt", "koin":
		return nil
	default:
		return fmt.Errorf("invalid DI option: %s (allowed: none|hilt|koin)", di)
	}
}

// isUserAbort attempts to heuristically detect errors caused by user
// cancellation (Ctrl+C) from interactive forms. Returns true when the
// error message contains common abort/cancel keywords.
func isUserAbort(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	if strings.Contains(s, "user aborted") || strings.Contains(s, "aborted") || strings.Contains(s, "cancel") || strings.Contains(s, "interrupt") {
		return true
	}
	return false
}
