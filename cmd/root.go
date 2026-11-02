/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"fj.vms.local/jsds/go-my-mcp/global"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   global.AppName,
	Short: "MCP for SearXNG",
	// Resolve the effective configuration once for every sub command, after flag parsing
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		trackFlags(cmd)
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

// flagKeys - command line flags and the viper keys they feed
var flagKeys = []struct{ name, key string }{
	{"config", global.KeyFileConf},
	{"debug", global.KeyDebug},
	{"listen", global.KeyListenAddress},
	{"SearXngUrl", global.KeySearXngUrl},
	{"RedirectMax", global.KeyRedirectMax},
	{"FetchBodySizeMax", global.KeyFetchBodySizeMax},
	{"SearchBodySizeMax", global.KeySearchBodySizeMax},
}

// record which flags were given on the command line, for `config` to report
func trackFlags(cmd *cobra.Command) {
	setFlags := map[string]string{}
	cmd.Flags().Visit(func(f *pflag.Flag) {
		for _, fk := range flagKeys {
			if fk.name == f.Name {
				setFlags[fk.key] = f.Value.String()
			}
		}
	})
	global.Vip.Set(global.KeyConfigFlag, setFlags)
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

	for _, fk := range flagKeys {
		bindFlag(fk.name, fk.key)
	}
}
