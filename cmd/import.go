// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package cmd

import (
	"fmt"

	"github.com/sapk/go-genesys/tool/check"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	importUsername string
	importPassword string
)

func init() {
	importCmd.Flags().StringVarP(&importUsername, "user", "u", "default", "GAX user name")
	importCmd.Flags().StringVarP(&importPassword, "pass", "p", "password", "GAX user password")
	RootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "[WIP] Connect to a GAX server to import object from dump",
	Long: `[WIP] Use GAX APIs to load  objects from dump of previous configuration.
	Ex:  genesys-tools import hostb:8080 Application/*.md`,
	Args: func(cmd *cobra.Command, args []string) error {
		logrus.Debug("Checking args for list cmd: ", args)
		if len(args) < 2 {
			return fmt.Errorf("requires at least one GAX server and one file to import")
		}
		if !check.IsValidClientArg(args[0]) {
			return fmt.Errorf("invalid gax host argument specified (ex: gax_host:8080): %s", args[0])
		}
		for _, arg := range args[1:] {
			if !check.IsValidFileArg(arg) {
				return fmt.Errorf("invalid file argument specified: %s", arg)
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}
