/*
Copyright © 2026 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"os"

	"github.com/J-Siu/go-my-mcp/global"
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

func init() {
	rootCmd.PersistentFlags().StringVarP(&global.Conf.FileConf, "FileConf", "c", global.Conf.FileConf, "config file path (default \""+global.DefaultConfigFile+"\")")
	rootCmd.PersistentFlags().BoolVarP(&global.Conf.Debug, "Debug", "d", global.Conf.Debug, "debug logging")
	rootCmd.PersistentFlags().StringVarP(&global.Conf.ListenAddress, "ListenAddress", "l", global.Conf.ListenAddress, "HTTP listen address")
	rootCmd.PersistentFlags().StringVarP(&global.Conf.SearXngUrl, "SearXngUrl", "u", global.Conf.SearXngUrl, "SearXNG URL")
	rootCmd.PersistentFlags().IntVar(&global.Conf.RedirectMax, "RedirectMax", global.Conf.RedirectMax, "maximum redirects followed by web_fetch")
	rootCmd.PersistentFlags().Int64Var(&global.Conf.FetchBodySizeMax, "FetchBodySizeMax", global.Conf.FetchBodySizeMax, "web_fetch response body limit in bytes")
	rootCmd.PersistentFlags().Int64Var(&global.Conf.SearchBodySizeMax, "SearchBodySizeMax", global.Conf.SearchBodySizeMax, "searxng_search response body limit in bytes")

	if err := global.Vip.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		panic("cmd.init: " + err.Error())
	}
}
