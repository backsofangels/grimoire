package cmd

import (
	"fmt"
	"strings"

	"github.com/backsofangels/grimoire/internal/logging"
	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/backsofangels/grimoire/internal/validator"
	"github.com/spf13/cobra"
)

// AddResourceInput consolidates all inputs for add command handlers.
type AddResourceInput struct {
	Kind       string
	Name       string
	Package    string
	Module     string
	Language   string
	Layout     string
	UI         string
	Override   bool
	DI         string
	ViewModel  bool
	Navigation bool
}

// ErrNeedsInteractive signals that interactive mode should be entered.
var ErrNeedsInteractive = fmt.Errorf("interactive mode required")

// extractAddFlags consolidates flag extraction for all add subcommands.
// Returns ErrNeedsInteractive if both name and --name flag are absent,
// signaling the caller to enter interactive mode.
func extractAddFlags(cmd *cobra.Command, args []string, kind string) (*AddResourceInput, error) {
	nameFlag, _ := cmd.Flags().GetString("name")
	var name string
	if len(args) == 0 && nameFlag == "" {
		return nil, ErrNeedsInteractive
	}
	if nameFlag != "" {
		name = nameFlag
	} else {
		name = args[0]
	}

	// Extract all flags
	pkg, _ := cmd.Flags().GetString("package")
	module, _ := cmd.Flags().GetString("module")
	lang, _ := cmd.Flags().GetString("lang")
	layout, _ := cmd.Flags().GetString("layout")
	ui, _ := cmd.Flags().GetString("ui")
	noUI, _ := cmd.Flags().GetBool("no-ui")
	override, _ := cmd.Flags().GetBool("override")
	di, _ := cmd.Flags().GetString("di")
	includeVM, _ := cmd.Flags().GetBool("viewmodel")
	addNav, _ := cmd.Flags().GetBool("nav")

	if noUI {
		ui = "none"
	}

	// Normalize/validate inputs
	if strings.ToLower(ui) == "none" && layout != "" {
		logging.Info("Ignoring --layout because --ui is 'none'")
		layout = ""
	}

	// Only validate explicit --ui values ("none" is set via --no-ui)
	if strings.ToLower(ui) != "none" {
		if err := validator.ValidateUI(ui); err != nil {
			return nil, fmt.Errorf("invalid --ui: %w", err)
		}
	}

	if err := validator.ValidateLanguage(lang); err != nil {
		return nil, fmt.Errorf("invalid --lang: %w", err)
	}

	if err := validator.ValidateDI(di); err != nil {
		return nil, fmt.Errorf("invalid --di: %w", err)
	}

	return &AddResourceInput{
		Kind:       kind,
		Name:       name,
		Package:    pkg,
		Module:     module,
		Language:   lang,
		Layout:     layout,
		UI:         ui,
		Override:   override,
		DI:         di,
		ViewModel:  includeVM,
		Navigation: addNav,
	}, nil
}

// inputToProviderConfig converts AddResourceInput to ProviderConfig.
func inputToProviderConfig(input *AddResourceInput) providers.ProviderConfig {
	cfg := providers.ProviderConfig{
		"Kind":        input.Kind,
		"Name":        input.Name,
		"PackageName": input.Package,
		"Module":      input.Module,
		"Lang":        input.Language,
		"Layout":      input.Layout,
		"UI":          input.UI,
		"Override":    input.Override,
		"DI":          input.DI,
	}

	// Add optional fields depending on resource type
	if input.Kind == "activity" || input.Kind == "fragment" {
		cfg["ViewModel"] = input.ViewModel
		cfg["Nav"] = input.Navigation
	}

	return cfg
}

// runAddResourceCommand orchestrates the add command flow for a specific resource type.
func runAddResourceCommand(cmd *cobra.Command, provider providers.Provider, kind string, args []string) error {
	input, err := extractAddFlags(cmd, args, kind)
	if err == ErrNeedsInteractive {
		return runAddInteractive(cmd, provider, kind)
	}
	if err != nil {
		return err
	}

	cfg := inputToProviderConfig(input)
	if err := provider.Add(cfg); err != nil {
		return fmt.Errorf("add %s failed: %w", kind, err)
	}

	logging.Success(fmt.Sprintf("Added %s %s to %s", kind, input.Name, input.Module))
	return nil
}
