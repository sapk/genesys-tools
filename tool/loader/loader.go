// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package loader

import (
	"github.com/sapk/go-genesys/api/client"
)

type Loader struct {
	FormatCreate func(*client.Client, map[string]interface{}) map[string]interface{}
	FormatUpdate func(*client.Client, map[string]interface{}, map[string]interface{}) map[string]interface{}
}

var LoaderList = map[string]Loader{
	"default": Loader{
		FormatCreate: func(c *client.Client, obj map[string]interface{}) map[string]interface{} {
			//cleanObj(obj, "dbid", "hostdbid", "appprototypedbid") //TODO find matching prototype for app //TODO ask for password
			cleanObj(obj, "dbid") //TODO find matching prototype for app //TODO ask for password
			return obj
		},
		FormatUpdate: func(c *client.Client, src, obj map[string]interface{}) map[string]interface{} {
			//TODO find matching prototype for app //TODO ask for password
			obj["dbid"] = src["dbid"]
			return obj
		},
	},
	"CfgAppPrototype": Loader{
		FormatCreate: func(c *client.Client, obj map[string]interface{}) map[string]interface{} {
			//cleanObj(obj, "dbid", "hostdbid", "appprototypedbid") //TODO find matching prototype for app //TODO ask for password
			cleanObj(obj, "dbid") //TODO find matching prototype for app //TODO ask for password
			return obj
		},
		FormatUpdate: func(c *client.Client, src, obj map[string]interface{}) map[string]interface{} {
			//TODO find matching prototype for app //TODO ask for password
			obj["dbid"] = src["dbid"]
			obj["folderid"] = "101" //Force Application Template Folder //TODO find a better folder
			return obj
		},
	},
}

func cleanObj(obj map[string]interface{}, ids ...string) {
	for _, id := range ids {
		_, ok := obj[id]
		if ok {
			delete(obj, id)
		}
	}
}
