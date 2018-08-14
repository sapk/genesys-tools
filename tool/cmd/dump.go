// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sapk/genesys/api/client"
	"github.com/sapk/genesys/api/object"
	"github.com/sapk/genesys/tool/check"
)

var (
	dumpNoJSON   bool
	dumpOnlyJSON bool
	dumpFromJSON string
)

//TODO add flag for username and pass
//TODO add from dump (from short flag
//TODO flag for no cleanup before export
//TODO ajouter filter pour dump jsute un app ou un host
//TODO interactive loadin bar
//TODO listing port et connection entre host
//TODO voir pour les liens backup et synchor
//TODO add follow folder structure for application
//TODO manage multi-tenant
//TODO replace strings.Count
//TODO manage switch/dn and agent and routing
//TODO add annex to application dumpt

func init() {
	dumpCmd.Flags().BoolVar(&dumpNoJSON, "no-json", false, "Disable global json dump")
	dumpCmd.Flags().BoolVar(&dumpOnlyJSON, "only-json", false, "Dump only global json")
	dumpCmd.Flags().StringVarP(&dumpFromJSON, "from-json", "f", "", "Read data from JSON and not a live GAX (directory containing all json)")
	//TODO from-json
	RootCmd.AddCommand(dumpCmd)
}

// listCmd represents the list command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Connect to a GAX server to dump its state",
	Args: func(cmd *cobra.Command, args []string) error {
		logrus.Debug("Checking args for list cmd: ", args)
		/* TODO multi arg with subfolder
		if len(args) > 1 {
			return fmt.Errorf("requires at least one GAX server")
		}
		*/
		if len(args) != 1 && dumpFromJSON == "" {
			return fmt.Errorf("requires one GAX server")
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
			//Get DATA
			apps, hosts := getData(gax)

			//sort.Sort(hosts) //order data by name
			//TODO fix ordering
			sort.Slice(hosts, func(i, j int) bool {
				return hosts[i].Name > hosts[j].Name
			})
			if !dumpNoJSON && dumpFromJSON == "" {
				err := dumpToFile("Hosts.json", hosts)
				if err != nil {
					logrus.Panicf("Dump failed : %v", err)
				}
			}
			//sort.Sort(apps) //order data by name
			//TODO fix ordering
			sort.Slice(apps, func(i, j int) bool {
				return apps[i].Name > apps[j].Name
			})
			if !dumpNoJSON && dumpFromJSON == "" {
				err := dumpToFile("Applications.json", apps)
				if err != nil {
					logrus.Panicf("Dump failed : %v", err)
				}
			}
			if !dumpOnlyJSON { //Don't analyze data
				err := clean("Hosts", "Applications")
				if err != nil {
					logrus.Panicf("Clean up failed : %v", err)
				}
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
					err = writeToFile(filepath.Join("Hosts", host.Name+".md"), formatHost(host, apps))
					if err != nil {
						logrus.Panicf("File creation failed : %v", err)
					}
				}
				for _, app := range apps {
					logrus.Infof("App: %s (%s)", app.Name, app.Dbid)
					err = writeToFile(filepath.Join("Applications", app.Name+".md"), formatApplication(app, apps, hosts))
					if err != nil {
						logrus.Panicf("File creation failed : %v", err)
					}
				}

			}
		}
	},
}

func getData(gax string) (object.CfgApplicationList, object.CfgHostList) {
	if dumpFromJSON == "" {
		return getGAXData(gax)
	}
	return getJSONData(dumpFromJSON)
}
func getJSONData(dumpFromJSON string) (object.CfgApplicationList, object.CfgHostList) {
	byteHosts, err := ioutil.ReadFile(filepath.Join(dumpFromJSON, "Hosts.json"))
	if err != nil {
		logrus.Panicf("ListHost failed : %v", err)
	}
	var hosts object.CfgHostList
	json.Unmarshal(byteHosts, &hosts)

	byteApps, err := ioutil.ReadFile(filepath.Join(dumpFromJSON, "Applications.json"))
	if err != nil {
		logrus.Panicf("ListApplication failed : %v", err)
	}
	var apps object.CfgApplicationList
	json.Unmarshal(byteApps, &apps)
	return apps, hosts
}
func getGAXData(gax string) (object.CfgApplicationList, object.CfgHostList) {
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
	logrus.WithFields(logrus.Fields{
		"User": user,
	}).Debugf("Logged as: %s", user.Username)

	//Get DATA
	//Hosts
	hosts, err := c.ListHost()
	if err != nil {
		logrus.Panicf("ListHost failed : %v", err)
	}
	//Applications
	apps, err := c.ListApplication()
	if err != nil {
		logrus.Panicf("ListApplication failed : %v", err)
	}
	return apps, hosts
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
	ret += " Commandlinearguments: " + app.Commandlinearguments + "\n"
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

	ports := treemap.NewWithStringComparator() //TODO pass to int comparator ?
	for _, p := range app.Portinfos.Portinfo {
		ports.Put(p.Port, p.ID)
	}
	portList := ""
	for _, id := range ports.Keys() {
		port := id.(string)
		val, _ := ports.Get(port)
		portList += "  " + val.(string) + " / " + port + "\n"
	}
	ret += fmt.Sprintf("## Listening ports (%d): \n", strings.Count(portList, "\n"))
	ret += portList
	ret += "\n"

	connections := treemap.NewWithStringComparator()
	for _, c := range app.Appservers.Conninfo {
		appserv := c.Appserverdbid
		for _, a := range apps {
			if c.Appserverdbid == a.Dbid {
				appserv = a.Name
				break
			}
		}
		connections.Put(appserv, c.ID+" / "+c.Mode)
	}
	connList := ""
	for _, id := range connections.Keys() {
		appName := id.(string)
		val, _ := connections.Get(appName)
		connList += "  " + appName + " / " + val.(string) + "\n"
	}
	ret += fmt.Sprintf("## Connections (%d): \n", strings.Count(connList, "\n"))
	ret += connList
	ret += "\n"

	sections := treeset.NewWithStringComparator()
	options := make(map[string]*treemap.Map)
	for _, o := range app.Options.Property {
		sections.Add(o.Section)
		if _, ok := options[o.Section]; !ok {
			//Init
			options[o.Section] = treemap.NewWithStringComparator()
		}
		options[o.Section].Put(o.Key, o.Value)
	}
	optList := ""
	for _, s := range sections.Values() {
		sec := s.(string)
		optList += " [" + sec + "]\n"
		for _, o := range options[sec].Keys() {
			opt := o.(string)
			val, _ := options[s.(string)].Get(opt)
			optList += "  " + opt + " = " + val.(string) + "\n"
		}
		//optList += " - [" + o.Section + "] / " + o.Key + " = " + o.Value + "\n"
	}

	ret += fmt.Sprintf("## Options (%d): \n", strings.Count(optList, "\n")-sections.Size())
	ret += optList
	ret += "\n"

	sectionsAnnex := treeset.NewWithStringComparator()
	annexes := make(map[string]*treemap.Map)
	for _, o := range app.Userproperties.Property {
		sectionsAnnex.Add(o.Section)
		if _, ok := annexes[o.Section]; !ok {
			//Init
			annexes[o.Section] = treemap.NewWithStringComparator()
		}
		annexes[o.Section].Put(o.Key, o.Value)
	}
	annexList := ""
	for _, s := range sectionsAnnex.Values() {
		sec := s.(string)
		annexList += " [" + sec + "]\n"
		for _, o := range annexes[sec].Keys() {
			opt := o.(string)
			val, _ := annexes[s.(string)].Get(opt)
			annexList += "  " + opt + " = " + val.(string) + "\n"
		}
		//optList += " - [" + o.Section + "] / " + o.Key + " = " + o.Value + "\n"
	}
	ret += fmt.Sprintf("## Annexes (%d): \n", strings.Count(annexList, "\n")-sectionsAnnex.Size())
	ret += annexList
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

	appList := treemap.NewWithStringComparator()
	for _, app := range apps {
		if app.Hostdbid == host.Dbid {
			appList.Put(app.Name, app)
		}
	}
	appListTxt := ""
	portListTxt := ""
	connListTxt := ""
	for _, id := range appList.Keys() {
		appName := id.(string)
		obj, _ := appList.Get(appName)
		app := obj.(object.CfgApplication)
		appListTxt += " - " + appName + "\n"
		ports := ""
		for _, port := range app.Portinfos.Portinfo {
			ports += port.ID + "/" + port.Port + ", "
		}
		if len(ports) > 2 {
			portListTxt += " - " + appName + " (" + ports[:len(ports)-2] + ")\n"
		}
		connections := ""
		for _, c := range app.Appservers.Conninfo {
			if c.Appserverdbid != host.Dbid {
				appserv := c.Appserverdbid
				for _, a := range apps {
					if c.Appserverdbid == a.Dbid {
						appserv = a.Name
						break
					}
				}
				connections += appserv + "/" + c.ID + "/" + c.Mode + ", "
			}
		}
		//TODO handle link with backup
		if len(connections) > 2 {
			connListTxt += " - " + appName + " -> (" + connections[:len(connections)-2] + ")\n"
		}
	}
	ret += fmt.Sprintf("## Applications (%d): \n", appList.Size())
	ret += appListTxt
	ret += "\n"

	ret += "## Listening ports (all applications): \n"
	ret += portListTxt
	ret += "\n"

	ret += "## [WIP] Connection with (client with connection outside): \n"
	ret += connListTxt
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

func clean(pathList ...string) error {
	for _, p := range pathList {
		err := os.RemoveAll(p)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
func cleanAll() error {
	return clean("Hosts.json", "Applications.json", "Hosts", "Applications")
}
*/
func dumpToFile(file string, data interface{}) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, json, 0644)
}
