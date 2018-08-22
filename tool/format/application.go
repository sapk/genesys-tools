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
	FormaterList["CfgApplication"] = Formater{
		func(objType object.ObjectType, obj map[string]interface{}, data map[string][]interface{}) string {
			ret := "# " + obj["name"].(string) + "\n"
			ret += "\n"

			ret += dumpAvailableInformation(obj, data) + "\n"
			//TODO annex and options
			ret += formatApplication(obj, data)
			ret += formatOptions(obj, data)
			ret += formatAnnexes(obj, data)
			ret += dumpBackup(obj)
			return ret
		},
		defaultShortFormater,
	}
}

func formatApplication(obj map[string]interface{}, data map[string][]interface{}) string {
	var app object.CfgApplication
	err := mapstructure.Decode(obj, &app)
	if err != nil {
		logrus.Warnf("Fail to convert to CfgApplication")
		return err.Error()
	}
	ports := treemap.NewWithStringComparator() //TODO pass to int comparator ?
	for _, p := range app.Portinfos.Portinfo {
		ports.Put(p.Port, p.ID)
	}

	portList := ""
	for _, id := range ports.Keys() {
		port := id.(string)
		val, _ := ports.Get(port)
		portList += "  " + val.(string) + " / " + port + "  \n"
	}
	ret := fmt.Sprintf("## Listening ports (%d): \n", ports.Size())
	ret += portList
	ret += "\n"

	connections := treemap.NewWithStringComparator()
	for _, c := range app.Appservers.Conninfo {
		connections.Put(funcFindByType("CfgApplication")(c.Appserverdbid, data), c.ID+" / "+c.Mode)
	}
	connList := ""
	for _, id := range connections.Keys() {
		appName := id.(string)
		val, _ := connections.Get(appName)
		connList += "  " + appName + " / " + val.(string) + "  \n"
	}
	ret += fmt.Sprintf("## Connections (%d): \n", connections.Size())
	ret += connList
	ret += "\n"
	return ret
}
