/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"

	"fj.vms.local/jsds/go-my-mcp/global"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"c"},
	Short:   "Show the effective configuration",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := json.MarshalIndent(global.Conf.Effective(), "", "  ")
		if err != nil {
			return err
		}
		cmd.Println(string(out))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
