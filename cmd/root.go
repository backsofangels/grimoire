package cmd

import (
	"fmt"

	"github.com/backsofangels/grimoire/internal/logging"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	version     string
	showVersion bool
)

var rootCmd = &cobra.Command{
	Use:   "grimoire",
	Short: "Grimoire CLI scaffolding tool",
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			// Print version only when explicitly requested.
			fmt.Println("grimoire", version)
			return
		}
		_ = cmd.Help()
	},
}

// Execute runs the root command.
func Execute(v string) {
	version = v
	// initialize logger (centralized)
	logging.Init()
	// Configure global logger appearance: remove timestamps, hide caller, and
	// default to Warn level so INFO logs are not shown in normal usage.
	log.SetTimeFormat("")
	log.SetReportCaller(false)
	log.SetLevel(log.WarnLevel)

	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Print version")
	// default provider
	rootCmd.PersistentFlags().String("provider", "android", "Scaffolding provider to use")
	// Provider flags are registered in each command's init(); avoid duplicate registration here.

	if err := rootCmd.Execute(); err != nil {
		logging.Error("✗ Command failed", "error", err)
	}
}
