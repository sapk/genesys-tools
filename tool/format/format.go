// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package format

import (
	"encoding/json"
	"fmt"

	"github.com/sapk/go-genesys/api/object"
)

type Formater struct {
	Format func(object.ObjectType, map[string]interface{}, map[string][]interface{}) string
}

var FormaterList = map[string]Formater{
	"default": Formater{
		func(objType object.ObjectType, obj map[string]interface{}, data map[string][]interface{}) string {
			ret := "# " + obj["name"].(string) + "\n"
			ret += "\n"

			ret += "## Informations: \n"
			ret += " Dbid: " + obj["dbid"].(string) + "\n"
			t, ok := obj["type"].(string)
			if ok {
				ret += " Type: " + t + "\n"
			}
			st, ok := obj["subtype"].(string)
			if ok {
				ret += " Subtype: " + st + "\n"
			}
			ret += dumpBackup(obj)
			return ret
		},
	},
	"CfgPerson": Formater{
		func(objType object.ObjectType, obj map[string]interface{}, data map[string][]interface{}) string {
			ret := "# " + obj["username"].(string) + "\n"
			ret += "\n"

			ret += "## Informations: \n"
			ret += " Dbid: " + obj["dbid"].(string) + "\n"
			ret += dumpBackup(obj)
			return ret
		},
	},
}

func dumpBackup(obj map[string]interface{}) string {
	json, err := json.Marshal(obj)
	if err != nil {
		//return fmt.Sprintf("\n## Dump :\n<!-- %s -->\n", err)
		return fmt.Sprintf("\n\n[//]: # (%s)\n", err)
	}
	//return fmt.Sprintf("\n## Dump :\n<!-- %s -->\n", json)
	return fmt.Sprintf("\n\n[//]: # (%s)\n", json)
}
