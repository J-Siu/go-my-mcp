/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"fj.vms.local/jsds/go-my-mcp/global"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   global.AppName,
	Short: "MCP for SearXNG",
	// Resolve the effective configuration once for every sub command, after flag parsing
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := global.Conf.Load(); err != nil {
			// a broken config is not a usage error
			cmd.SilenceUsage = true
			return err
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Bind a cobra flag to a viper key, so viper resolves the effective value
func bindFlag(name, key string) {
	if err := global.Vip.BindPFlag(key, rootCmd.PersistentFlags().Lookup(name)); err != nil {
		// only reachable on a flag/key mismatch in the code below
		panic("cmd.bindFlag: " + err.Error())
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file path (default \""+global.DefaultConfigFile+"\")")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug logging")
	rootCmd.PersistentFlags().StringP("listen", "l", global.DefaultListenAddress, "HTTP listen address")
	rootCmd.PersistentFlags().StringP("SearXngUrl", "u", global.DefaultSearXngUrl, "SearXNG URL")
	rootCmd.PersistentFlags().Int("RedirectMax", global.DefaultRedirectMax, "maximum redirects followed by web_fetch")
	rootCmd.PersistentFlags().Int64("FetchBodySizeMax", global.DefaultFetchBodySizeMax, "web_fetch response body limit in bytes")
	rootCmd.PersistentFlags().Int64("SearchBodySizeMax", global.DefaultSearchBodySizeMax, "searxng_search response body limit in bytes")

	bindFlag("config", "FileConf")
	bindFlag("debug", "Debug")
	bindFlag("listen", "ListenAddress")
	bindFlag("SearXngUrl", "SearXngUrl")
	bindFlag("RedirectMax", "RedirectMax")
	bindFlag("FetchBodySizeMax", "FetchBodySizeMax")
	bindFlag("SearchBodySizeMax", "SearchBodySizeMax")
}
