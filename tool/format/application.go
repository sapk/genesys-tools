// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package format

import "github.com/sapk/go-genesys/api/object"

func init() {
	FormaterList["CfgApplication"] = Formater{
		func(objType object.ObjectType, obj map[string]interface{}, data map[string][]interface{}) string {
			ret := "# " + obj["name"].(string) + "\n"
			ret += "\n"

			ret += "## Informations: \n"
			ret += " Dbid: " + obj["dbid"].(string) + "\n"
			ret += dumpBackup(obj)
			return ret
		},
	}
}
