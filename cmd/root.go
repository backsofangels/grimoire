package cmd

import (
	"fmt"
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
			fmt.Printf("grimoire %s\n", version)
			return
		}
		_ = cmd.Help()
	},
}

// Execute runs the root command.
func Execute(v string) {
	version = v
	// initialize logger (lightweight)
	log.Info("logger initialized", "version", version)

	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Print version")
	// default provider
	rootCmd.PersistentFlags().String("provider", "android", "Scaffolding provider to use")
	// Provider flags are registered in each command's init(); avoid duplicate registration here.

	if err := rootCmd.Execute(); err != nil {
		log.Error("command failed", "error", err)
	}
}
