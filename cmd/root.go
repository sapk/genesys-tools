// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var appVerbose bool

var (
	//Version version of app set by build flag
	//Version = "v0.0.4"
	Version = "latest"
	//Branch git branch of app set by build flag
	Branch = "master"
	//Commit git commit of app set by build flag
	Commit string
	//BuildTime build time of app set by build flag
	BuildTime string
)

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
	RootCmd.Long = RootCmd.Short + fmt.Sprintf("\nVersion: %s - Branch: %s - Commit: %s - BuildTime: %s\n\n", Version, Branch, Commit, BuildTime)
	RootCmd.PersistentFlags().BoolVarP(&appVerbose, "verbose", "v", false, "Turn on verbose logging")
	RootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Display current version and build date",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\nVersion: %s - Branch: %s - Commit: %s - BuildTime: %s\n\n", Version, Branch, Commit, BuildTime)
		},
	})
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
