// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var appVerbose bool

//var appCSVOutput bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "genesys-tools",
	Short: "A simple application to view and test some every day task on Genesys solution",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initVerbose)

	RootCmd.PersistentFlags().BoolVarP(&appVerbose, "verbose", "v", false, "Turn on verbose logging")
	//RootCmd.PersistentFlags().BoolVar(&appCSVOutput, "csv", false, "Turn on verbose logging output compatible with csv")
	//RootCmd.AddCommand(dump.DumpCmd)
	//RootCmd.AddCommand(check.CheckCmd)
}

func initVerbose() {
	if appVerbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	/*
		if !appCSVOutput {
			logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
			logrus.SetOutput(colorable.NewColorableStdout())
		}
	*/
}
