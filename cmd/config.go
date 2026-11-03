/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/J-Siu/go-my-mcp/global"
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
		// stdout, so the output can be piped; cmd.Println would use stderr
		_, err = fmt.Fprintln(cmd.OutOrStdout(), string(out))
		return err
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
