// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package format

import (
	"fmt"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"

	"github.com/sapk/go-genesys/api/object"
)

func init() {
	FormaterList["CfgHost"] = Formater{
		func(objType object.ObjectType, obj map[string]interface{}, data map[string][]interface{}) string {
			ret := "# " + obj["name"].(string) + "\n"
			ret += "\n"

			ret += dumpAvailableInformation(obj, data) + "\n"
			ret += formatHostAppList(obj["dbid"].(string), data)
			ret += formatOptions(obj, data)
			ret += formatAnnexes(obj, data)
			ret += dumpBackup(obj)
			return ret
		},
	}
}

func formatHostAppList(dbid string, data map[string][]interface{}) string {
	appList := treemap.NewWithStringComparator()
	for _, _o := range data["CfgApplication"] {
		o := _o.(map[string]interface{})
		hostdbid, ok := o["hostdbid"].(string)
		if ok && hostdbid == dbid {
			appList.Put(o["name"], o)
		}
	}

	appListTxt := ""
	portListTxt := ""
	connListTxt := ""

	for _, id := range appList.Keys() {
		appName := id.(string)
		obj, _ := appList.Get(appName)
		var app object.CfgApplication
		err := mapstructure.Decode(obj, &app)
		if err != nil {
			logrus.Warnf("Fail to convert to CfgApplication")
			continue
		}
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
			if c.Appserverdbid != dbid {
				connections += funcFindByType("CfgApplication")(c.Appserverdbid, data) + "/" + c.ID + "/" + c.Mode + ", "
				//TODO find host of App
			}
		}
		//TODO handle link with backup
		if len(connections) > 2 {
			connListTxt += " - " + appName + " -> (" + connections[:len(connections)-2] + ")\n"
		}
	}

	ret := fmt.Sprintf("## Applications (%d): \n", appList.Size())
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

/*
var result Person
err := mapstructure.Decode(input, &result)
if err != nil {
    panic(err)
}
fmt.Printf("%#v", result)
*/
