// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package cmd

import (
	"fmt"

	"github.com/sapk/genesys-tools/tool/check"
	"github.com/sapk/genesys-tools/tool/render"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	servPort string
)

func init() {
	serveCmd.Flags().StringVarP(&servPort, "port", "p", "8080", "Listening port")
	RootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Render markdown file of a dump on a local web server",
	Long: `Render markdown file of a dump on a local web server.
	Based on github.com/lithammer/go-wiki to render markdown file and display git history if files.
	Ex:  genesys-tools serve hostb-8080`,
	Args: func(cmd *cobra.Command, args []string) error {
		logrus.Debug("Checking args for serve cmd: ", args)
		if len(args) != 1 {
			return fmt.Errorf("Requires one folder to display")
		}
		if !check.IsValidFolderArg(args[0]) {
			return fmt.Errorf("Invalid folder argument specified. Must contain an index.md file. (ex: gax_host-8080): %s", args[0])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		//TODO start browser
		render.Serve(args[0], servPort)
	},
}
