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
	"github.com/J-Siu/go-my-mcp/global"
	"github.com/J-Siu/go-my-mcp/server"
	"github.com/J-Siu/go-my-mcp/tool/searxng"
	"github.com/J-Siu/go-my-mcp/tool/webfetch"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"s"},
	Short:   "Start the SearXNG MCP server over HTTP",
	RunE: func(cmd *cobra.Command, args []string) error {
		sx, err := new(searxng.TypeSearXng).New(global.Conf.SearXngUrl, global.Conf.SearchBodySizeMax)
		if err == nil {
			wf, err := new(webfetch.TypeWebFetch).New(global.Conf.FetchBodySizeMax, global.Conf.RedirectMax)
			if err == nil {
				err = new(server.TypeServer).
					New(global.Conf.ListenAddress, global.Conf.ApiKey).
					Register(sx).
					Register(wf).
					Run()
			}
		}
		if err != nil {
			// a bind failure is an environment problem, not a usage error
			cmd.SilenceUsage = true
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
