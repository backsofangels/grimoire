package cmd

import (
	"fmt"

	"github.com/backsofangels/grimoire/internal/logging"
	"github.com/backsofangels/grimoire/internal/providers"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a resource to an existing project (activity|fragment|viewmodel)",
}

var addActivityCmd = &cobra.Command{
	Use:   "activity [name]",
	Short: "Add an Activity",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		providerName, _ := cmd.Flags().GetString("provider")
		provider, err := providers.Get(providerName)
		if err != nil {
			logging.Error("✗ Provider not found", "provider", providerName, "error", err)
			return
		}

		// allow --name flag as alternative to positional arg; interactive when both absent
		nameFlag, _ := cmd.Flags().GetString("name")
		if len(args) == 0 && nameFlag == "" {
			if err := runAddInteractive(cmd, provider, "activity"); err != nil {
				logging.Error("✗ Add failed", "error", err)
			}
			return
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
		override, _ := cmd.Flags().GetBool("override")

		di, _ := cmd.Flags().GetString("di")
		includeVM, _ := cmd.Flags().GetBool("viewmodel")
		addNav, _ := cmd.Flags().GetBool("nav")

		cfg := providers.ProviderConfig{
			"Kind":        "activity",
			"Name":        name,
			"PackageName": pkg,
			"Module":      module,
			"Lang":        lang,
			"Layout":      layout,
			"Override":    override,
			"DI":          di,
			"ViewModel":   includeVM,
			"Nav":         addNav,
		}

		if err := provider.Add(cfg); err != nil {
			logging.Error("✗ Add failed", "error", err)
			return
		}
		logging.Success(fmt.Sprintf("Added activity %s", name))
	},
}

var addFragmentCmd = &cobra.Command{
	Use:   "fragment [name]",
	Short: "Add a Fragment",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		providerName, _ := cmd.Flags().GetString("provider")
		provider, err := providers.Get(providerName)
		if err != nil {
			logging.Error("✗ Provider not found", "provider", providerName, "error", err)
			return
		}

		// allow --name flag as alternative to positional arg; interactive when both absent
		nameFlag, _ := cmd.Flags().GetString("name")
		if len(args) == 0 && nameFlag == "" {
			if err := runAddInteractive(cmd, provider, "fragment"); err != nil {
				logging.Error("✗ Add failed", "error", err)
			}
			return
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
		override, _ := cmd.Flags().GetBool("override")

		di, _ := cmd.Flags().GetString("di")
		includeVM, _ := cmd.Flags().GetBool("viewmodel")
		addNav, _ := cmd.Flags().GetBool("nav")

		cfg := providers.ProviderConfig{
			"Kind":        "fragment",
			"Name":        name,
			"PackageName": pkg,
			"Module":      module,
			"Lang":        lang,
			"Layout":      layout,
			"Override":    override,
			"DI":          di,
			"ViewModel":   includeVM,
			"Nav":         addNav,
		}

		if err := provider.Add(cfg); err != nil {
			logging.Error("✗ Add failed", "error", err)
			return
		}
		logging.Success(fmt.Sprintf("Added fragment %s", name))
	},
}

var addViewModelCmd = &cobra.Command{
	Use:   "viewmodel [name]",
	Short: "Add a ViewModel",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		providerName, _ := cmd.Flags().GetString("provider")
		provider, err := providers.Get(providerName)
		if err != nil {
			logging.Error("✗ Provider not found", "provider", providerName, "error", err)
			return
		}

		// allow --name flag as alternative to positional arg; interactive when both absent
		nameFlag, _ := cmd.Flags().GetString("name")
		if len(args) == 0 && nameFlag == "" {
			if err := runAddInteractive(cmd, provider, "viewmodel"); err != nil {
				logging.Error("✗ Add failed", "error", err)
			}
			return
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
		override, _ := cmd.Flags().GetBool("override")

		cfg := providers.ProviderConfig{
			"Kind":        "viewmodel",
			"Name":        name,
			"PackageName": pkg,
			"Module":      module,
			"Lang":        lang,
			"Override":    override,
		}

		if err := provider.Add(cfg); err != nil {
			logging.Error("✗ Add failed", "error", err)
			return
		}
		logging.Success(fmt.Sprintf("Added viewmodel %s", name))
	},
}

func init() {
	// Register add command and subcommands
	rootCmd.AddCommand(addCmd)
	addCmd.AddCommand(addActivityCmd)
	addCmd.AddCommand(addFragmentCmd)
	addCmd.AddCommand(addViewModelCmd)

	// common flags
	// allow --name as alternative to positional
	addActivityCmd.Flags().String("name", "", "Resource name (alternative to positional)")
	addActivityCmd.Flags().String("di", "none", "Dependency injection (none|hilt|koin)")
	addActivityCmd.Flags().Bool("viewmodel", false, "Generate associated ViewModel")
	addActivityCmd.Flags().Bool("nav", false, "Add navigation entry")

	addActivityCmd.Flags().StringP("package", "p", "", "Target package (e.g. com.example.app)")
	addActivityCmd.Flags().StringP("module", "m", "app", "Target module (default: app)")
	addActivityCmd.Flags().StringP("lang", "l", "kotlin", "Language (kotlin|java)")
	addActivityCmd.Flags().StringP("layout", "", "", "Layout resource name (optional)")
	addActivityCmd.Flags().BoolP("override", "", false, "Overwrite existing files")

	addFragmentCmd.Flags().String("name", "", "Resource name (alternative to positional)")
	addFragmentCmd.Flags().String("di", "none", "Dependency injection (none|hilt|koin)")
	addFragmentCmd.Flags().Bool("viewmodel", false, "Generate associated ViewModel")
	addFragmentCmd.Flags().Bool("nav", false, "Add navigation entry")
	addFragmentCmd.Flags().StringP("package", "p", "", "Target package (e.g. com.example.app)")
	addFragmentCmd.Flags().StringP("module", "m", "app", "Target module (default: app)")
	addFragmentCmd.Flags().StringP("lang", "l", "kotlin", "Language (kotlin|java)")
	addFragmentCmd.Flags().StringP("layout", "", "", "Layout resource name (optional)")
	addFragmentCmd.Flags().BoolP("override", "", false, "Overwrite existing files")

	addViewModelCmd.Flags().String("name", "", "Resource name (alternative to positional)")
	addViewModelCmd.Flags().StringP("package", "p", "", "Target package (e.g. com.example.app)")
	addViewModelCmd.Flags().StringP("module", "m", "app", "Target module (default: app)")
	addViewModelCmd.Flags().StringP("lang", "l", "kotlin", "Language (kotlin|java)")
	addViewModelCmd.Flags().BoolP("override", "", false, "Overwrite existing files")
}
