// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package format

import (
	"github.com/sapk/go-genesys/api/object"
)

func init() {
	FormaterList["CfgDN"] = Formater{
		func(objType object.ObjectType, obj map[string]interface{}, data map[string][]interface{}) string {
			ret := "# " + obj["number"].(string) + "\n"
			ret += "\n"

			ret += dumpAvailableInformation(obj, data) + "\n"
			//TODO details
			ret += formatOptions(obj, data)
			ret += formatAnnexes(obj, data)
			ret += dumpBackup(obj)
			return ret
		},
	}
}
