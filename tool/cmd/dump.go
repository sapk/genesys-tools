// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package cmd

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sapk/genesys/api/client"
	"github.com/sapk/genesys/tool/check"
)

//TODO add flag for username and pass

// listCmd represents the list command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Connect to a GAX server to dump his state",
	Args: func(cmd *cobra.Command, args []string) error {
		logrus.Debug("Checking args for list cmd: ", args)
		if len(args) > 1 {
			return fmt.Errorf("requires at least one GAX server")
		}
		for _, arg := range args {
			if !check.IsValidClientArg(arg) {
				return fmt.Errorf("invalid argument specified (ex: gax_host:8080): %s", arg)
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		for _, gax := range args {
			if !strings.Contains(gax, ":") {
				//By default use port 8080
				gax += "8080"
			}
			c := client.NewClient(gax)
			user, err := c.Login("default", "password")
			if err != nil {
				logrus.Panicf("Login failed : %v", err)
			}
			logrus.Println(user)

			hosts, err := c.ListHost()
			if err != nil {
				logrus.Panicf("ListHost failed : %v", err)
			}
			for _, host := range hosts {
				logrus.Infof("Host: %s (%s)", host.Name, host.DBID)
			}
			apps, err := c.ListApplication()
			if err != nil {
				logrus.Panicf("ListApplication failed : %v", err)
			}
			for _, app := range apps {
				logrus.Infof("App: %s (%s) @ %s", app.Name, app.DBID, app.WorkDirectory)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(dumpCmd)
}
