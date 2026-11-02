/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fj.vms.local/jsds/go-my-mcp/global"
	"fj.vms.local/jsds/go-my-mcp/server"
	"fj.vms.local/jsds/go-my-mcp/tool/searxng"
	"fj.vms.local/jsds/go-my-mcp/tool/webfetch"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"s"},
	Short:   "Start the SearXNG MCP server over HTTP",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Printf("MCP server listening on %s\n", global.Conf.ListenAddress)
		var (
			mcpSrv = new(server.TypeServer)
			sx     = new(searxng.TypeSearXng)
			wf     = new(webfetch.TypeWebFetch)
			err    error
		)
		sx, err = sx.New(global.Conf.SearXngUrl, global.Conf.SearchBodySizeMax)
		if err == nil {
			wf, err = wf.New(global.Conf.FetchBodySizeMax, global.Conf.RedirectMax)
		}
		if err == nil {
			if global.Conf.Debug {
				mcpSrv.Debug = true
				sx.Debug = true
				wf.Debug = true
			}
			sx.TlsCheck = false
			err = mcpSrv.
				New(global.Conf.ListenAddress).
				Register(sx).
				Register(wf).
				Run()
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
