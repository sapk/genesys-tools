// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package main

import "github.com/sapk/go-genesys/cmd"

var (
	//Version version of app set by build flag
	Version string
	//Branch git branch of app set by build flag
	Branch string
	//Commit git commit of app set by build flag
	Commit string
	//BuildTime build time of app set by build flag
	BuildTime string
)

func main() {
	if Version != "" {
		cmd.Version = Version
	}
	if Branch != "" {
		cmd.Branch = Branch
	}
	cmd.Commit = Commit
	cmd.BuildTime = BuildTime
	cmd.Execute()
}
