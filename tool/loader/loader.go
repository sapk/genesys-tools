// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package loader

import (
	"github.com/c-bata/go-prompt"
	"github.com/sapk/go-genesys/api/client"
	"github.com/sapk/go-genesys/tool/format"
	"github.com/sirupsen/logrus"
)

var data map[string][]map[string]interface{}

func init() {
	data = make(map[string][]map[string]interface{})
}

func ListObject(c *client.Client, t string) []map[string]interface{} {
	logrus.WithField("type", t).Debugf("ListObject")
	if l, ok := data[t]; ok {
		logrus.WithField("list", l).Debugf("ListObject in cache")
		return l //don't get list allready done
	}
	var list []map[string]interface{}
	c.ListObject(t, &list)
	logrus.WithField("list", list).Debugf("ListObject fetched")
	if t == "CfgTenant" {
		list = append(list, map[string]interface{}{"name": "Environment", "dbid": "1"})
	}
	data[t] = list
	return list
}

type Loader struct {
	FormatCreate func(*client.Client, map[string]interface{}) map[string]interface{}
	FormatUpdate func(*client.Client, map[string]interface{}, map[string]interface{}) map[string]interface{}
}

var LoaderList = map[string]Loader{
	"default": Loader{
		FormatCreate: func(c *client.Client, obj map[string]interface{}) map[string]interface{} {
			logrus.WithFields(logrus.Fields{
				"in": obj,
			}).Debugf("default.FormatCreate")
			//cleanObj(obj, "dbid", "hostdbid", "appprototypedbid") //TODO find matching prototype for app //TODO ask for password
			cleanObj(obj, "dbid")
			if tenant, exist := obj["tenantdbid"]; exist {
				obj["tenantdbid"] = searchFor(c, "CfgTenant", tenant.(string))
			}
			logrus.WithFields(logrus.Fields{
				"out": obj,
			}).Debugf("default.FormatCreate")
			return obj
		},
		FormatUpdate: func(c *client.Client, src, obj map[string]interface{}) map[string]interface{} {
			//TODO find matching prototype for app //TODO ask for password
			obj["dbid"] = src["dbid"]
			cleanObj(obj, "tenantdbid")
			return obj
		},
	},
	/*
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
	*/
}

func searchFor(c *client.Client, t string, id string) string {
	//TODO search for seam id

	list := ListObject(c, t)
	logrus.WithField("list", list).WithField("type", t).Info("Fetched list")
	logrus.Infof("Please choose a %s :", t)
	val := prompt.Input("> ", func(d prompt.Document) []prompt.Suggest {
		//logrus.WithField("list", list).Info("Fetched list")
		s := make([]prompt.Suggest, len(list))
		for i, o := range list {
			//logrus.WithField("obj", o).Info("Add to list")
			s[i] = prompt.Suggest{o["dbid"].(string), format.GetFileName(o)}
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	})

	logrus.Infof("You selected " + val)
	if val == "" {
		switch t {
		case "CfgTenant":
			return "1"
		}
	}
	return val
}

func cleanObj(obj map[string]interface{}, ids ...string) {
	for _, id := range ids {
		_, ok := obj[id]
		if ok {
			delete(obj, id)
		}
	}
}

//TODO match for each format like format

func FilterBy(obj map[string]interface{}, data []map[string]interface{}, cmp func(map[string]interface{}, map[string]interface{}) bool) []map[string]interface{} {
	ret := make([]map[string]interface{}, 0)
	for _, o := range data {
		if cmp(obj, o) {
			logrus.WithField("obj", o).Info("Matched object")
			ret = append(ret, obj) //TODO best allocate
		}
	}
	return ret
}

func MatchId(src, dst map[string]interface{}) bool {
	return src["dbid"] == dst["dbid"]
}
func MatchName(src, dst map[string]interface{}) bool {
	//TODO check not nil
	return src["name"] == dst["name"] || src["username"] == dst["username"] || src["number"] == dst["number"] || src["logincode"] == dst["logincode"]
}
func MatchIdName(src, dst map[string]interface{}) bool { //TODO Manage Person (username)
	return MatchName(src, dst) && MatchId(src, dst)
}
