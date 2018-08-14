// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sapk/genesys/api/client"
	"github.com/sapk/genesys/api/object"
	"github.com/sapk/genesys/tool/check"
)

//TODO add flag for username and pass
//TODO add short option that only dump json doesn't format data
//TODO add from dump (from short flag
//TODO flag for no cleanup before export

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

			///Login
			c := client.NewClient(gax)
			user, err := c.Login("default", "password")
			if err != nil {
				logrus.Panicf("Login failed : %v", err)
			}
			logrus.Println(user)

			//Cleanup
			err = cleanAll()
			if err != nil {
				logrus.Panicf("Clean up failed : %v", err)
			}

			//Get DATA
			//Hosts
			hosts, err := c.ListHost()
			if err != nil {
				logrus.Panicf("ListHost failed : %v", err)
			}
			sort.Sort(object.CfgHostList(hosts)) //order data by name
			err = dumpToFile("./Hosts.json", hosts)
			if err != nil {
				logrus.Panicf("Dump failed : %v", err)
			}
			//Applications
			apps, err := c.ListApplication()
			if err != nil {
				logrus.Panicf("ListApplication failed : %v", err)
			}
			sort.Sort(object.CfgApplicationList(apps)) //order data by name
			err = dumpToFile("./Applications.json", apps)
			if err != nil {
				logrus.Panicf("Dump failed : %v", err)
			}

			//TODO format data
			err = os.Mkdir("Hosts", 0755)
			if err != nil {
				logrus.Panicf("Folder creation failed : %v", err)
			}
			err = os.Mkdir("Applications", 0755)
			if err != nil {
				logrus.Panicf("Folder creation failed : %v", err)
			}
			for _, host := range hosts {
				logrus.Infof("Host: %s (%s)", host.Name, host.Dbid)
				err = writeToFile("Hosts/"+host.Name+".md", formatHost(host, apps))
				if err != nil {
					logrus.Panicf("File creation failed : %v", err)
				}
			}
			for _, app := range apps {
				logrus.Infof("App: %s (%s)", app.Name, app.Dbid)
				err = writeToFile("Applications/"+app.Name+".md", formatApplication(app, apps, hosts))
				if err != nil {
					logrus.Panicf("File creation failed : %v", err)
				}
			}
		}
	},
}

//TODO order applications conn and port
func formatApplication(app object.CfgApplication, apps object.CfgApplicationList, hosts object.CfgHostList) string {
	ret := "# " + app.Name + "\n"
	ret += "\n"
	ret += "## Informations: \n"
	ret += " Dbid: " + app.Dbid + "\n"
	ret += " Name: " + app.Name + "\n"
	host := app.Hostdbid
	for _, h := range hosts {
		if app.Hostdbid == h.Dbid {
			host = h.Name
			break
		}
	}
	ret += " Host: " + host + "\n"
	ret += " Type: " + app.Type + "\n"
	ret += " Subtype: " + app.Subtype + "\n"
	ret += " Componenttype: " + app.Componenttype + "\n"
	ret += " Appprototypedbid: " + app.Appprototypedbid + "\n" //TODO
	ret += " Isserver: " + app.Isserver + "\n"
	ret += " Version: " + app.Version + "\n"
	ret += " State: " + app.State + "\n"
	ret += " Startuptype: " + app.Startuptype + "\n"
	ret += " Workdirectory: " + app.Workdirectory + "\n"
	ret += " Commandline: " + app.Commandline + "\n"
	ret += " Autorestart: " + app.Autorestart + "\n"
	ret += " Port principal: " + app.Port + "\n"
	ret += " Redundancytype: " + app.Redundancytype + "\n"
	ret += " Isprimary: " + app.Isprimary + "\n"
	backup := app.Backupserverdbid
	for _, a := range apps {
		if app.Backupserverdbid == a.Dbid {
			backup = a.Name
			break
		}
	}
	ret += " Backupserver: " + backup + "\n"
	ret += "\n"

	portList := ""
	for _, p := range app.Portinfos.Portinfo {
		portList += " - " + p.ID + " / " + p.Port + "\n"
	}
	ret += fmt.Sprintf("## Listening ports (%d): \n", strings.Count(portList, "\n"))
	ret += portList
	ret += "\n"

	connList := ""
	for _, c := range app.Appservers.Conninfo {
		appserv := c.Appserverdbid
		for _, a := range apps {
			if c.Appserverdbid == a.Dbid {
				appserv = a.Name
				break
			}
		}
		connList += " - " + appserv + " / " + c.ID + " / " + c.Mode + "\n"
	}
	ret += fmt.Sprintf("## Connections (%d): \n", strings.Count(connList, "\n"))
	ret += connList
	ret += "\n"

	//TODO format as ini
	optList := ""
	for _, o := range app.Options.Property {
		optList += " - [" + o.Section + "] / " + o.Key + " = " + o.Value + "\n"
	}
	ret += fmt.Sprintf("## Options (%d): \n", strings.Count(optList, "\n"))
	ret += optList
	ret += "\n"
	return ret
}

func formatHost(host object.CfgHost, apps object.CfgApplicationList) string {
	ret := "# " + host.Name + "\n"
	ret += "\n"

	ret += "## Informations: \n"
	ret += " Dbid: " + host.Dbid + "\n"
	ret += " Name: " + host.Name + "\n"
	ret += " Type: " + host.Type + "\n"
	ret += " Subtype: " + host.Subtype + "\n"
	ret += " OS: " + host.Ostype + "\n"
	ret += " State: " + host.State + "\n"
	ret += " IP: " + host.Ipaddress + "\n"
	ret += "\n"

	appList := ""
	for _, app := range apps {
		if app.Hostdbid == host.Dbid {
			appList += " - " + app.Name + "\n"
		}
	}
	ret += fmt.Sprintf("## Applications (%d): \n", strings.Count(appList, "\n"))
	ret += appList
	ret += "\n"

	ret += "## Listening ports (all applications): \n"
	ret += "TODO\n"
	ret += "\n"
	return ret
}
func writeToFile(file, data string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(data)
	if err != nil {
		return err
	}
	f.Sync()
	return nil
}

func cleanAll() error {
	err := os.RemoveAll("./Hosts.json")
	if err != nil {
		return err
	}
	err = os.RemoveAll("./Applications.json")
	if err != nil {
		return err
	}
	err = os.RemoveAll("./Hosts")
	if err != nil {
		return err
	}
	return os.RemoveAll("./Applications")
}

func dumpToFile(file string, data interface{}) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, json, 0644)
}

func init() {
	RootCmd.AddCommand(dumpCmd)
}
