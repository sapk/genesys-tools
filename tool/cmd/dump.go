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
	"github.com/sapk/go-genesys/tool/format"
	"github.com/sapk/go-genesys/tool/fs"
)

var (
	dumpFull     bool
	dumpZip      bool
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
//TODO add timeout to connection in app format
//TODO export AgentGroup script
//TODO dump log on as of application
//TODO find a solution for applciation  like UCS that use port from config (annex [ports])
//TODO dump object at end of file to re-import them back
//TODO add exemple liek flowtester (and fix typo lister in flowtester)

func init() {
	dumpCmd.Flags().BoolVarP(&dumpFull, "extended", "e", false, "[WIP] Get also switch, dn, person, place, ...")
	dumpCmd.Flags().BoolVarP(&dumpZip, "zip", "z", false, "zip the output folder")
	dumpCmd.Flags().BoolVar(&dumpNoJSON, "no-json", false, "Disable global json dump")
	dumpCmd.Flags().BoolVar(&dumpOnlyJSON, "only-json", false, "Dump only global json")
	dumpCmd.Flags().StringVarP(&dumpFromJSON, "from-json", "f", "", "Read data from JSON and not a live GAX (directory containing all json)")
	dumpCmd.Flags().StringVarP(&dumpUsername, "user", "u", "default", "GAX user name")
	dumpCmd.Flags().StringVarP(&dumpPassword, "pass", "p", "password", "GAX user password")
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
					if !objType.IsDumpable {
						continue //Skip
					}
					//TODO skip if empty array ?
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
							err = fs.WriteToFile(filepath.Join(outFolder, obj["dbid"].(string)+"-"+name+".md"), formatObj(objType, obj, data))
							if err != nil {
								logrus.Panicf("File creation failed : %v", err)
							}
						} else {
							//Second try with username (default user)
							name, ok := obj["username"].(string)
							if ok {
								err = fs.WriteToFile(filepath.Join(outFolder, obj["dbid"].(string)+"-"+name+".md"), formatObj(objType, obj, data))
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
			if dumpZip {
				logrus.Infof("Compacting folder: %s", gax)
				err := fs.RecursiveZip(gax, gax+".zip")
				if err != nil {
					logrus.Panicf("Failed to zip folder : %v", err)
				}
				err = fs.Clean(gax)
				if err != nil {
					logrus.Panicf("Clean up failed : %v", err)
				}
			}
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
		bytes, err := ioutil.ReadFile(filepath.Join(dumpFromJSON, gax, objType.Desc+".json"))
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

//Call the good formatter if exist or use the default
func formatObj(objType object.ObjectType, obj map[string]interface{}, data map[string][]interface{}) string {
	if f, ok := format.FormaterList[objType.Name]; ok {
		return f.Format(objType, obj, data)
	}
	return format.FormaterList["default"].Format(objType, obj, data)
}
