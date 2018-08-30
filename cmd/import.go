// +build linux darwin windows
// +build amd64 386

// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/sapk/go-genesys/api/client"
	"github.com/sapk/go-genesys/tool/check"
	"github.com/sapk/go-genesys/tool/format"
	"github.com/sapk/go-genesys/tool/loader"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	importUsername string
	importPassword string
	importForceYes bool
)

//TODO add help message for what is not imported
var allowedImportTypes = map[string]bool{
	//"CfgApplication":  true,
	"CfgPlace":        true, //lost link to contactdbid capacityruledbid dndbids sitedbid
	"CfgDN":           true,
	"CfgAppPrototype": true,
	"CfgField":        true,
	"CfgScript":       true,
	"CfgAgentLogin":   true,
}

//TODO importe template and metadata first
//TODO afficher les connection et lien manquant , host, ...
func init() {
	importCmd.Flags().StringVarP(&importUsername, "user", "u", "default", "GAX user name")
	importCmd.Flags().StringVarP(&importPassword, "pass", "p", "password", "GAX user password")
	importCmd.Flags().BoolVarP(&importForceYes, "force", "f", false, "Implies yes to each questions")
	RootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "[WIP] Connect to a GAX server to import object from dump",
	Long: `[WIP] Use GAX APIs to load  objects from dump of previous configuration.
	Ex:  genesys-tools import hostb:8080 Application/*.md`,
	//TODO list allowedImportTypes
	Args: func(cmd *cobra.Command, args []string) error {
		logrus.Debug("Checking args for import cmd: ", args)
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
		if !strings.Contains(gax, ":") {
			//By default use port 8080
			gax += ":8080"
		}
		//Login
		c := client.NewClient(gax)
		user, err := c.Login(importUsername, importPassword)
		if err != nil {
			logrus.Panicf("Login failed : %v", err)
		}
		logrus.WithFields(logrus.Fields{
			"User": user,
		}).Debugf("Logged as: %s", user.Username)

		for _, file := range args[1:] {
			obj := getObj(file)
			logrus.Infof("Parsing %s: %s", obj["type"], format.Name(obj))
			logrus.WithFields(logrus.Fields{
				"Object": obj,
			}).Debug("Parsing object")

			t, ok := obj["type"].(string)
			if !ok {
				logrus.Fatalf("Fail to find type of object %s : %v", file, obj)
			}

			if !allowedImportTypes[t] {
				logrus.Warnf("Skipping file %s since type %s is not importable yet.", file, t)
				continue
			}

			l := loader.ListObject(c, t)
			logrus.Debugf("List response : %v", l)

			if len(l) == 0 { //no same object so we create
				logrus.Debugf("Found no object with type : %v", t)
				createObj(c, obj)
			} else {
				//Try to find if a app is matching
				list := loader.FilterBy(obj, l, loader.MatchIdName)
				if len(list) == 0 {
					logrus.Debugf("Found no object with same DBID and Name")
					list = loader.FilterBy(obj, l, loader.MatchName)
					if len(list) == 0 {
						logrus.Debugf("Found no object with same Name")
						/* Temporary disable as it doesn't match change in name for exemple (detected on place)
						list = loader.FilterBy(obj, l, loader.MatchId)
						if len(list) == 0 {
							logrus.Debugf("Found no object with same DBID")
						}
						*/
					}
				}
				//TODO less ugly
				//TODO manage errors
				var err error
				switch len(list) {
				case 0: //no same object so we create
					err = createObj(c, obj)
				case 1:
					err = updateObj(c, list[0], obj)
				default:
					logrus.Warnf("Multiple object matching : %s", file)
					for _, src := range list {
						updateObj(c, src, obj)
					}
				}
				if err != nil {
					logrus.WithField("object", obj).Errorf("Failed to import object: %v", err)
				} else {
					logrus.WithField("object", obj).Infof("Import object success !")
				}
			}
		}
	},
}

func updateObj(c *client.Client, src map[string]interface{}, obj map[string]interface{}) error {
	logrus.WithFields(logrus.Fields{
		"Source": src,
		"Object": obj,
	}).Info("Update object")
	eq := reflect.DeepEqual(obj, src)
	if eq {
		logrus.WithFields(logrus.Fields{
			"Source": src,
			"Object": obj,
		}).Info("Skipping update of object because of equality")
		return nil
	}
	if f, ok := loader.LoaderList[obj["type"].(string)]; ok {
		obj = f.FormatUpdate(c, src, obj)
	} else {
		obj = loader.LoaderList["default"].FormatUpdate(c, src, obj)
	}
	//TODO check eq after cleaning
	eq = reflect.DeepEqual(obj, src)
	if eq {
		logrus.WithFields(logrus.Fields{
			"Source": src,
			"Object": obj,
		}).Info("Skipping update of object because of equality after loading format")
		return nil
	}
	logrus.WithFields(logrus.Fields{
		"Object": obj,
	}).Debugf("Sending updated object")
	//TODO ask for ovveride
	//TODO get dbid for older one ?
	//TODO check possible deps
	//TODO check if no change
	if importForceYes || check.AskFor(fmt.Sprintf("Update %s", format.FormatShortObj(obj))) { // ask for confirmation
		_, err := c.UpdateObject(src["type"].(string), src["dbid"].(string), obj) //TODO check up
		return err
	}
	return nil
}
func createObj(c *client.Client, obj map[string]interface{}) error {
	logrus.WithFields(logrus.Fields{
		"Object": obj,
	}).Info("Create object")
	if f, ok := loader.LoaderList[obj["type"].(string)]; ok {
		obj = f.FormatCreate(c, obj)
	} else {
		obj = loader.LoaderList["default"].FormatCreate(c, obj)
	}
	logrus.WithFields(logrus.Fields{
		"Object": obj,
	}).Debugf("Sending new object")

	if importForceYes || check.AskFor(fmt.Sprintf("Create %s", format.FormatShortObj(obj))) { // ask for confirmation
		_, err := c.PostObject(obj) //TODO check up
		return err
	}
	return nil
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
