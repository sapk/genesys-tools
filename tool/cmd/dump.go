// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	//"path/filepath"
	//"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sapk/go-genesys/api/client"
	"github.com/sapk/go-genesys/api/object"
	"github.com/sapk/go-genesys/tool/check"
	"github.com/sapk/go-genesys/tool/fs"
)

var (
	dumpFull     bool
	dumpNoJSON   bool
	dumpOnlyJSON bool
	dumpFromJSON string
	dumpUsername string
	dumpPassword string
)

//TODO add scheme (http/https) flag or assume default http
//TODO ajouter filter pour dump jsute un app ou un host
//TODO interactive loading bar
//TODO voir pour les liens backup et synchor (pour connection entre host et app)
//TODO add follow folder structure for application
//TODO manage multi-tenant
//TODO manage switch/dn and agent and routing
//TODO add --zip to package all output to zip
//TODO add timeout to connection in app format
//TODO export AgentGroup script
//TODO dump log on as of application
//TODO find a solution for applciation  like UCS that use port from config (annex [ports])
//TODO dump object at end of file to re-import them back
//TODO add exemple liek flowtester (and fix typo lister in flowtester)

func init() {
	dumpCmd.Flags().BoolVarP(&dumpFull, "extended", "e", false, "[WIP] Get also switch, dn, person, place, ...")
	dumpCmd.Flags().BoolVar(&dumpNoJSON, "no-json", false, "Disable global json dump")
	dumpCmd.Flags().BoolVar(&dumpOnlyJSON, "only-json", false, "Dump only global json")
	dumpCmd.Flags().StringVarP(&dumpFromJSON, "from-json", "f", "", "Read data from JSON and not a live GAX (directory containing all json)")
	dumpCmd.Flags().StringVarP(&dumpUsername, "user", "u", "default", "GAX user name")
	dumpCmd.Flags().StringVarP(&dumpPassword, "pass", "p", "password", "GAX user password")
	//TODO from-json
	RootCmd.AddCommand(dumpCmd)
}

// listCmd represents the list command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Connect to a GAX server to dump its state",
	Args: func(cmd *cobra.Command, args []string) error {
		logrus.Debug("Checking args for list cmd: ", args)
		if len(args) > 1 {
			return fmt.Errorf("requires at least one GAX server")
		}
		/*
			if len(args) != 1 && dumpFromJSON == "" {
				return fmt.Errorf("requires one GAX server")
			}
		*/
		for _, arg := range args {
			if !check.IsValidClientArg(arg) {
				return fmt.Errorf("invalid argument specified (ex: gax_host:8080): %s", arg)
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		for _, gax := range args {

			//tmp := strings.Split(gax, ":")
			//host := tmp[0]
			logrus.Infof("Get info from GAX: %s", gax)
			if dumpFromJSON == "" { //Create folder if not from restore
				if _, err := os.Stat(gax); err == nil {
					logrus.Warnf("Overwriting old export %s", gax)
					err = fs.Clean(gax)
					if err != nil {
						logrus.Panicf("Clean up failed : %v", err)
					}
				}
				err := os.Mkdir(gax, 0755)
				if err != nil {
					logrus.Panicf("Folder creation failed : %v", err)
				}
			}

			list := object.ObjectTypeListShort
			if dumpFull {
				list = object.ObjectTypeList
			}
			//Get DATA
			data := getData(gax, list)
			//TODO fix ordering

			if !dumpNoJSON && dumpFromJSON == "" {
				for _, objType := range list {
					err := fs.DumpToFile(filepath.Join(gax, objType.Desc+".json"), data[objType.Name])
					if err != nil {
						logrus.Panicf("Dump failed : %v", err)
					}
				}
			}
			if !dumpOnlyJSON {
				for _, objType := range list {
					outFolder := filepath.Join(gax, objType.Desc)
					if _, err := os.Stat(outFolder); err == nil {
						logrus.Warnf("Overwriting old export %s", outFolder)
						err = fs.Clean(outFolder)
						if err != nil {
							logrus.Panicf("Clean up failed : %v", err)
						}
					}
					err := os.Mkdir(outFolder, 0755)
					if err != nil {
						logrus.Panicf("Folder creation failed : %v", err)
					}
					for _, o := range data[objType.Name] {
						obj := o.(map[string]interface{})
						logrus.Infof("%s: %s (%s)", objType.Name, obj["name"], obj["dbid"])
						name, ok := obj["name"].(string)
						if ok {
							err = fs.WriteToFile(filepath.Join(outFolder, name+".md"), formatObj(objType, obj, data))
							if err != nil {
								logrus.Panicf("File creation failed : %v", err)
							}
						} else {
							//Second try with username (default user)
							name, ok := obj["username"].(string)
							if ok {
								err = fs.WriteToFile(filepath.Join(outFolder, name+".md"), formatObj(objType, obj, data))
								if err != nil {
									logrus.Panicf("File creation failed : %v", err)
								}
							} else {
								logrus.Warnf("Ignoring invalid object / %s: %s (%s)", objType.Name, obj["name"], obj["dbid"])
							}
						}
					}
				}
			}

			/*
				if !dumpOnlyJSON { //Don't analyze data
					err := clean("Hosts", "Applications", "Switchs", "DNs", "Places")
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
					if dumpFull {
						//TODO inprove by refactor repetive code
						err = os.Mkdir("Switchs", 0755)
						if err != nil {
							logrus.Panicf("Folder creation failed : %v", err)
						}
						err = os.Mkdir("DNs", 0755)
						if err != nil {
							logrus.Panicf("Folder creation failed : %v", err)
						}
						err = os.Mkdir("Places", 0755)
						if err != nil {
							logrus.Panicf("Folder creation failed : %v", err)
						}
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
					if dumpFull {
						//TODO inprove by refactor repetive code
						for _, s := range switchs {
							logrus.Infof("Switch: %s (%s)", s.Name, s.Dbid)
							err = writeToFile(filepath.Join("Switchs", s.Name+".md"), formatSwitch(s))
							if err != nil {
								logrus.Panicf("File creation failed : %v", err)
							}
						}
						for _, d := range dns {
							logrus.Infof("DN: %s (%s)", d.Name, d.Dbid)
							err = writeToFile(filepath.Join("DNs", d.Name+".md"), formatDN(d))
							if err != nil {
								logrus.Panicf("File creation failed : %v", err)
							}
						}
						for _, p := range places {
							logrus.Infof("Place: %s (%s)", p.Name, p.Dbid)
							err = writeToFile(filepath.Join("Places", p.Name+".md"), formatPlace(p))
							if err != nil {
								logrus.Panicf("File creation failed : %v", err)
							}
						}
					}
				}
			*/
		}
	},
}

func getData(gax string, list []object.ObjectType) map[string][]interface{} {
	if dumpFromJSON == "" {
		return getGAXData(gax, list)
	}
	return getJSONData(dumpFromJSON, gax, list)
}

func getJSONData(dumpFromJSON string, gax string, list []object.ObjectType) map[string][]interface{} {
	var res = make(map[string][]interface{})
	for _, objType := range list {
		//Get objects
		var data []interface{}
		bytes, err := ioutil.ReadFile(filepath.Join(dumpFromJSON, gax, objType.Name+".json"))
		if err != nil {
			logrus.Warnf("List %s failed : %v", objType.Name, err)
		}
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			logrus.Warnf("List %s failed : %v", objType.Name, err)
		} else {
			res[objType.Name] = data
		}
	}
	return res
}
func getGAXData(gax string, list []object.ObjectType) map[string][]interface{} {
	if !strings.Contains(gax, ":") {
		//By default use port 8080
		gax += "8080"
	}
	//Login
	c := client.NewClient(gax)
	user, err := c.Login(dumpUsername, dumpPassword)
	if err != nil {
		logrus.Panicf("Login failed : %v", err)
	}
	logrus.WithFields(logrus.Fields{
		"User": user,
	}).Debugf("Logged as: %s", user.Username)
	var res = make(map[string][]interface{})
	for _, objType := range list {
		//Get objects
		var data []interface{}
		_, err := c.ListObject(objType.Name, &data)
		if err != nil {
			logrus.Warnf("List %s failed : %v", objType.Name, err)
		} else {
			res[objType.Name] = data
		}
	}
	return res
}
func formatObj(objType object.ObjectType, obj map[string]interface{}, data map[string][]interface{}) string {
	return "TODO"
}

/*

//TODO order applications conn and port
func formatApplication(app object.CfgApplication, apps []object.CfgApplication, hosts []object.CfgHost) string {
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
	ret += fmt.Sprintf("## Listening ports (%d): \n", ports.Size())
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
	ret += fmt.Sprintf("## Connections (%d): \n", connections.Size())
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

func formatHost(host object.CfgHost, apps []object.CfgApplication) string {
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

func formatSwitch(s object.CfgSwitch) string {
	ret := "# " + s.Name + "\n"
	ret += "\n"

	return ret
}
func formatDN(d object.CfgDN) string {
	ret := "# " + d.Name + "\n"
	ret += "\n"

	return ret
}
func formatPlace(p object.CfgPlace) string {
	ret := "# " + p.Name + "\n"
	ret += "\n"

	return ret
}
*/
