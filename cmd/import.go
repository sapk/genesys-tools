// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

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
		gax := args[0]
		logrus.Debugln(gax)
		for _, file := range args[1:] {
			o := getObj(file)
			logrus.Debugln(o)
		}
	},
}

func getObj(file string) map[string]interface{} {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		logrus.Fatalf("Read file %s failed : %v", file, err)
	}
	fileStr := string(b)

	pos := strings.LastIndex(fileStr, "[//]: # ({")
	if pos == -1 {
		logrus.Fatalf("Fail to found raw dump in file %s : %v", file, err)
	}
	jsonStr := fileStr[pos+9:]

	//TODO regex
	pos = strings.Index(jsonStr, "})\n")
	if pos == -1 {
		logrus.Fatalf("Fail to found raw dump in file %s : %v", file, err)
	}
	jsonStr = jsonStr[:pos+1]
	logrus.Debugf("Parsing JSON : %s", jsonStr)

	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		logrus.Fatalf("Fail failed to parse %s : %v", jsonStr, err)
	}
	return data
}
