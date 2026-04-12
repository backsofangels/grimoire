package cmd

import (
	"fmt"
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
		return runAddResourceCommand(cmd, provider, "activity", args)
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
		return runAddResourceCommand(cmd, provider, "fragment", args)
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
		return runAddResourceCommand(cmd, provider, "viewmodel", args)
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
	addCmd.PersistentFlags().StringP("module", "m", "app", "Target module folder (default: app)")
	addCmd.PersistentFlags().StringP("lang", "l", "kotlin", "Language (kotlin|java)")
	addCmd.PersistentFlags().StringP("layout", "", "", "Layout resource name (optional)")
	addCmd.PersistentFlags().StringP("ui", "", "xml", "UI type (xml|compose). Use --no-ui to disable UI generation")
	addCmd.PersistentFlags().Bool("no-ui", false, "Shortcut for --ui none (equivalent to --ui none)")
	addCmd.PersistentFlags().BoolP("override", "", false, "Overwrite existing files")
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
