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

	"github.com/rs/zerolog/log"
	"github.com/sapk/genesys-tools/tool/check"
	"github.com/sapk/genesys-tools/tool/format"
	"github.com/sapk/genesys-tools/tool/loader"
	"github.com/sapk/go-genesys/api/client"
	"github.com/spf13/cobra"
)

var (
	importUsername       string
	importPassword       string
	importDefaultAnswers []string
	importForceYes       bool
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
	"CfgPerson":       true,
	"CfgAgentGroup":   true,
}

//TODO importe template and metadata first
//TODO afficher les connection et lien manquant , host, ...
func init() {
	importCmd.Flags().StringVarP(&importUsername, "user", "u", "default", "GAX user name")
	importCmd.Flags().StringVarP(&importPassword, "pass", "p", "password", "GAX user password")
	importCmd.Flags().BoolVarP(&importForceYes, "force", "f", false, "Implies yes to each questions")
	importCmd.Flags().StringSliceVarP(&importDefaultAnswers, "default", "d", []string{}, "Default value to answer by object type. (Ex: -d 'CfgTenant=102,CfgTenant=224')")
	RootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "[WIP] Connect to a GAX server to import object from dump",
	Long: `[WIP] Use GAX APIs to load  objects from dump of previous configuration.
	Ex:  genesys-tools import hostb:8080 Application/*.md`,
	//TODO list allowedImportTypes
	Args: func(cmd *cobra.Command, args []string) error {
		log.Debug().Msgf("Checking args for import cmd: %s", args)
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

		defaults := make(map[string]string)
		for _, def := range importDefaultAnswers {
			tmp := strings.Split(def, "=")
			if len(tmp) != 2 {
				log.Warn().Interface("default", def).Msg("Invalid default ignored")
				continue
			}
			log.Debug().Interface("type", tmp[0]).Interface("value", tmp[1]).Msg("Registering default value")
			defaults[tmp[0]] = tmp[1]
		}
		//Login
		c := client.NewClient(gax, false)
		user, err := c.Login(importUsername, importPassword)
		if err != nil {
			log.Panicf("Login failed : %v", err)
		}
		log.WithFields(log.Fields{
			"User": user,
		}).Debugf("Logged as: %s", user.Username)

		for _, file := range args[1:] {
			obj := getObj(file)
			log.Info().Msgf("Parsing %s: %s", obj["type"], format.Name(obj))
			log.Debug().Interface("Object", obj).Msg("Parsing object")

			t, ok := obj["type"].(string)
			if !ok {
				log.Fatal().Msgf("Fail to find type of object %s : %v", file, obj)
			}

			if !allowedImportTypes[t] {
				log.Warn().Msgf("Skipping file %s since type %s is not importable yet.", file, t)
				continue
			}

			l := loader.ListObject(c, t)
			log.Debug().Msgf("List response : %v", l)

			if len(l) == 0 { //no same object so we create
				log.Debug().Msgf("Found no object with type : %v", t)
				createObj(c, obj, defaults)
			} else {
				//Try to find if a app is matching
				list := loader.FilterBy(obj, l, loader.MatchIdName)
				if len(list) == 0 {
					log.Debug().Msgf("Found no object with same DBID and Name")
					list = loader.FilterBy(obj, l, loader.MatchName)
					if len(list) == 0 {
						log.Debug().Msgf("Found no object with same Name")
						/* Temporary disable as it doesn't match change in name for exemple (detected on place)
						list = loader.FilterBy(obj, l, loader.MatchId)
						if len(list) == 0 {
							log.Debug().Msgf("Found no object with same DBID")
						}
						*/
					}
				}
				//TODO less ugly
				//TODO manage errors
				var err error
				switch len(list) {
				case 0: //no same object so we create
					err = createObj(c, obj, defaults)
				case 1:
					err = updateObj(c, list[0], obj, defaults)
				default:
					log.Warn().Msgf("Multiple object matching : %s", file)
					for _, src := range list {
						updateObj(c, src, obj, defaults)
					}
				}
				if err != nil {
					log.Error().Interface("object", obj).Msgf("Failed to import object: %v", err)
				} else {
					log.Info().Interface("object", obj).Msgf("Import object success !")
				}
			}
		}
	},
}

func updateObj(c *client.Client, src map[string]interface{}, obj map[string]interface{}, defaults map[string]string) error {
	log.Info().Interface("Source", src).Interface("Object", obj).Msg("Update object")
	eq := reflect.DeepEqual(obj, src)
	if eq {
		log.Info().Interface("Source", src).Interface("Object", obj).Msg("Skipping update of object because of equality")
		return nil
	}
	if f, ok := loader.LoaderList[obj["type"].(string)]; ok {
		obj = f.FormatUpdate(c, src, obj, defaults)
	} else {
		obj = loader.LoaderList["default"].FormatUpdate(c, src, obj, defaults)
	}
	//TODO check eq after cleaning
	eq = reflect.DeepEqual(obj, src)
	if eq {
		log.Info().Interface("Source", src).Interface("Object", obj).Msg("Skipping update of object because of equality after loading format")
		return nil
	}
	log.Info().Interface("Object", obj).Msg("Sending updated object")
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
func createObj(c *client.Client, obj map[string]interface{}, defaults map[string]string) error {
	log.Debug().Interface("Object", obj).Msg("Create object init")
	if f, ok := loader.LoaderList[obj["type"].(string)]; ok {
		obj = f.FormatCreate(c, obj, defaults)
	} else {
		obj = loader.LoaderList["default"].FormatCreate(c, obj, defaults)
	}
	log.Info().Interface("Object", obj).Msg("Create object")

	if importForceYes || check.AskFor(fmt.Sprintf("Create %s", format.FormatShortObj(obj))) { // ask for confirmation
		_, err := c.PostObject(obj) //TODO check up
		return err
	}
	return nil
}

func getObj(file string) map[string]interface{} {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal().Msgf("Read file %s failed : %v", file, err)
	}
	fileStr := string(b)

	pos := strings.LastIndex(fileStr, "[//]: # ({")
	if pos == -1 {
		log.Fatal().Msgf("Fail to found raw dump in file %s : %v", file, err)
	}
	jsonStr := fileStr[pos+9:]

	//TODO regex
	pos = strings.Index(jsonStr, "})\n")
	if pos == -1 {
		log.Fatal().Msgf("Fail to found raw dump in file %s : %v", file, err)
	}
	jsonStr = jsonStr[:pos+1]
	log.Debug().Msgf("Parsing JSON : %s", jsonStr)

	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		log.Fatal().Msgf("Fail failed to parse %s : %v", jsonStr, err)
	}
	return data
}
