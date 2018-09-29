// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package format

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"

	"github.com/sapk/genesys-tools/api/object"
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
		func(objType object.ObjectType, obj map[string]interface{}, data map[string][]interface{}) string {
			var app object.CfgApplication
			err := mapstructure.Decode(obj, &app)
			if err != nil {
				logrus.Warnf("Fail to convert to CfgApplication")
				return err.Error()
			}
			_, portList := getApplicationPorts(app)
			portList = strings.Replace(portList, " ", "", -1)
			portList = strings.Replace(portList, "\n", ",", -1)
			portList = strings.Trim(portList, ",")
			buf := new(bytes.Buffer)
			wr := csv.NewWriter(buf)
			val, ok := obj["hostdbid"] //TODO use app.Hostdbid
			if ok {
				//TODO add folder
				wr.Write([]string{app.Dbid, app.Name, trimCFGString(app.Subtype), app.Version, funcFindByType("CfgHost")(val, data), portList, trimCFGString(app.State), app.Backupserverdbid})
			} else {
				wr.Write([]string{app.Dbid, app.Name, trimCFGString(app.Subtype), app.Version, "", portList, trimCFGString(app.State), app.Backupserverdbid}) //empty host
			}
			wr.Flush()
			return buf.String()
		},
		func(objType object.ObjectType, obj map[string]interface{}) string {
			name := GetFileName(obj)
			/* Short should not need acces to data
			if obj["hostdbid"] != nil && obj["hostdbid"] != "" {
				host := findObjName("CfgHost", obj["hostdbid"].(string), data)
				return fmt.Sprintf(" - [%s](./%s/%s \\(%s\\)) (%s) @ %s\n", name, objType.Desc, name, obj["dbid"], obj["subtype"], host)
			}
			*/
			return fmt.Sprintf(" - [%s](./%s/%s \\(%s\\)) (%s)\n", name, objType.Desc, name, obj["dbid"], obj["subtype"])
		},
	}
}

func getApplicationPorts(app object.CfgApplication) (int, string) {
	ports := treemap.NewWithStringComparator()
	for _, p := range app.Portinfos.Portinfo {
		ports.Put(p.Port, p.ID)
	}
	portList := ""
	for _, id := range ports.Keys() {
		port := id.(string)
		val, _ := ports.Get(port)
		portList += "  " + val.(string) + " / " + port + "  \n"
	}
	return ports.Size(), portList
}

func trimCFGString(srt string) string {
	return strings.TrimPrefix(srt, "CFG")
}

func formatApplication(obj map[string]interface{}, data map[string][]interface{}) string {
	var app object.CfgApplication
	err := mapstructure.Decode(obj, &app)
	if err != nil {
		logrus.Warnf("Fail to convert to CfgApplication")
		return err.Error()
	}

	portSize, portList := getApplicationPorts(app)
	ret := fmt.Sprintf("## Listening ports (%d): \n", portSize)
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
