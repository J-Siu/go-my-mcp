/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fj.vms.local/jsds/go-my-mcp/global"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Show the version",
	Args:    cobra.NoArgs,
	// Skip the configuration load inherited from the root command, so the
	// version is shown even when the configuration file is broken.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("%s version %s\n", global.AppName, global.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
