// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	//"path/filepath"
	//"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sapk/genesys-tools/api/client"
	"github.com/sapk/genesys-tools/api/object"
	"github.com/sapk/genesys-tools/tool/check"
	"github.com/sapk/genesys-tools/tool/format"
	"github.com/sapk/genesys-tools/tool/fs"
)

var (
	dumpFull     bool
	dumpZip      bool
	dumpNoJSON   bool
	dumpOnlyJSON bool
	dumpCSV      bool
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
//TODO find a solution for applciation  like UCS that use port from config (annex [ports])
//TODO list agent in agent Group
//TODO list skils of Agent
//TODO add a timer of 6 month of end of life

func init() {
	dumpCmd.Flags().BoolVarP(&dumpFull, "extended", "e", false, "[WIP] Get also switch, dn, person, place, ...")
	dumpCmd.Flags().BoolVarP(&dumpZip, "zip", "z", false, "zip the output folder")
	dumpCmd.Flags().BoolVar(&dumpCSV, "csv", false, "output some csv table for some type (ex: Application)")
	dumpCmd.Flags().BoolVar(&dumpNoJSON, "no-json", false, "Disable global json dump")
	dumpCmd.Flags().BoolVar(&dumpOnlyJSON, "only-json", false, "Dump only global json")
	dumpCmd.Flags().StringVarP(&dumpFromJSON, "from-json", "f", "", "Read data from JSON and not a live GAX (directory containing all json)")
	dumpCmd.Flags().StringVarP(&dumpUsername, "user", "u", "default", "GAX user name")
	dumpCmd.Flags().StringVarP(&dumpPassword, "pass", "p", "password", "GAX user password")
	RootCmd.AddCommand(dumpCmd)
}

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Connect to a GAX server to dump his state",
	Long: `Use GAX APIs to get all objects from config server.
This command can dump multiple gax at a time. One folder for each GAX is created.
	Ex:  genesys-tools dump 172.18.0.5:8080 hosta hostb:4242`,
	Args: func(cmd *cobra.Command, args []string) error {
		logrus.Debug("Checking args for dump cmd: ", args)
		if len(args) < 1 {
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
			gaxFolder := strings.Replace(gax, ":", "-", -1)
			//tmp := strings.Split(gax, ":")
			//host := tmp[0]
			logrus.Infof("Get info from GAX: %s", gax)
			if dumpFromJSON == "" { //Create folder if not from restore
				if _, err := os.Stat(gaxFolder); err == nil {
					logrus.Warnf("Overwriting old export %s", gaxFolder)
					err = fs.Clean(gaxFolder)
					if err != nil {
						logrus.Panicf("Clean up failed : %v", err)
					}
				}
				err := os.Mkdir(gaxFolder, 0755)
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
					err := fs.DumpToFile(filepath.Join(gaxFolder, objType.Desc+".json"), data[objType.Name])
					if err != nil {
						logrus.Panicf("Dump failed : %v", err)
					}
				}
			}
			if !dumpOnlyJSON {
				sig := fmt.Sprintf("\n[//]: # (generated @ %s by genesys-tools-%s/%s developed by Antoine GIRARD)\n", time.Now().Format(time.RFC3339), Version, Commit)
				resume := "# " + gax + "\n\n"
				for _, objType := range list {
					if !objType.IsDumpable {
						continue //Skip
					}

					csvData := ""
					var csvFormater func(object.ObjectType, map[string]interface{}, map[string][]interface{}) string
					if dumpCSV {
						if f, ok := format.FormaterList[objType.Name]; ok && f.FormatCSV != nil {
							csvData = "dbid,name,type,version,host,ports,status,backup\n"
							csvFormater = f.FormatCSV
						}
					}
					//TODO skip if empty array ?
					outFolder := filepath.Join(gaxFolder, objType.Desc)
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
					resume += fmt.Sprintf("## %s (%d) :\n", objType.Desc, len(data[objType.Name]))
					//TODO order objects
					for _, o := range data[objType.Name] {
						obj := o.(map[string]interface{})
						name := format.GetFileName(obj)

						resume += format.FormatShortObj(obj)
						logrus.Infof("%s: %s (%s)", objType.Name, name, obj["dbid"])

						if name != "" {
							if csvFormater != nil {
								csvData += csvFormater(objType, obj, data)
							}
							err = fs.WriteToFile(filepath.Join(outFolder, name+" ("+obj["dbid"].(string)+").md"), format.FormatObj(objType, obj, data), sig)
							if err != nil {
								logrus.Warnf("File creation failed : %v", err) //Dont't panic and keep continue even in case of error
							}
						} else {
							logrus.Warnf("Ignoring invalid object / %s: %s (%s)", objType.Name, name, obj["dbid"])
						}
					}
					//Dump CSV at end
					if csvFormater != nil {
						err := fs.WriteToFile(filepath.Join(gaxFolder, objType.Desc+".csv"), csvData, "")
						if err != nil {
							logrus.Panicf("CSV failed : %v", err)
						}
					}
					resume += "\n"
				}
				resume += format.GenerateMermaidGraph(data)
				err := fs.WriteToFile(filepath.Join(gaxFolder, "index.md"), resume, sig)
				if err != nil {
					logrus.Panicf("File creation failed : %v", err)
				}
				err = fs.WriteToFile(filepath.Join(gaxFolder, "graph-by-app.dot"), format.GenerateDotGraphByApp(data), "")
				if err != nil {
					logrus.Panicf("File creation failed : %v", err)
				}
				err = fs.WriteToFile(filepath.Join(gaxFolder, "graph-by-host.dot"), format.GenerateDotGraphByHost(data), "")
				if err != nil {
					logrus.Panicf("File creation failed : %v", err)
				}
			}
			if dumpZip {
				time.Sleep(250 * time.Millisecond) //Sleep to let the prog release access on file
				logrus.Infof("Compacting folder: %s", gaxFolder)
				err := fs.RecursiveZip(gaxFolder, gaxFolder+".zip")
				if err != nil {
					logrus.Panicf("Failed to zip folder : %v", err)
				}
				time.Sleep(250 * time.Millisecond) //Sleep to let the prog release access on file
				err = fs.Clean(gaxFolder)
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
	gaxFolder := strings.Replace(gax, ":", "-", -1)
	for _, objType := range list {
		//Get objects
		var data []interface{}
		bytes, err := ioutil.ReadFile(filepath.Join(dumpFromJSON, gaxFolder, objType.Desc+".json"))
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
		gax += ":8080"
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
