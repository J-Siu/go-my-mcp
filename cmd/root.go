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
	"fmt"
	"os"

	"github.com/J-Siu/go-my-mcp/conf"
	"github.com/J-Siu/go-my-mcp/helper"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     conf.AppName,
	Short:   "MCP for SearXNG",
	Version: helper.Version(),
	// Resolve the effective configuration once for every sub command, after flag parsing
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := conf.Conf.Load(); err != nil {
			// a broken config is not a usage error
			cmd.SilenceUsage = true
		}
		if conf.Conf.Debug {
			fmt.Println("Version:", helper.Version())
			fmt.Println(helper.Struct2String("Conf", &conf.Conf))
		}
		return err
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&conf.Conf.Debug, "Debug", "d", conf.Conf.Debug, "debug logging")
	rootCmd.PersistentFlags().BoolVarP(&conf.Conf.SkipTlsVerify, "SkipTlsVerify", "", conf.Conf.SkipTlsVerify, "skip TLS verify")
	rootCmd.PersistentFlags().Int64VarP(&conf.Conf.FetchBodySizeMax, "FetchBodySizeMax", "", conf.Conf.FetchBodySizeMax, "web_fetch response body limit in bytes")
	rootCmd.PersistentFlags().Int64VarP(&conf.Conf.SearchBodySizeMax, "SearchBodySizeMax", "", conf.Conf.SearchBodySizeMax, "web_search response body limit in bytes")
	rootCmd.PersistentFlags().IntVarP(&conf.Conf.RedirectMax, "RedirectMax", "", conf.Conf.RedirectMax, "maximum redirects followed by web_fetch")
	rootCmd.PersistentFlags().StringVarP(&conf.Conf.FileConf, "FileConf", "c", conf.Conf.FileConf, "config file path (default: "+conf.DefaultConfigFile+")")
	rootCmd.PersistentFlags().StringVarP(&conf.Conf.ListenAddress, "ListenAddress", "l", conf.Conf.ListenAddress, "HTTP listen address")
	rootCmd.PersistentFlags().StringVarP(&conf.Conf.SearXngUrl, "SearXngUrl", "u", conf.Conf.SearXngUrl, "SearXNG URL")

	if err := conf.Vip.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		panic("cmd.init: " + err.Error())
	}
}
