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
	FormatCreate func(*client.Client, map[string]interface{}, map[string]string) map[string]interface{}
	FormatUpdate func(*client.Client, map[string]interface{}, map[string]interface{}, map[string]string) map[string]interface{}
}

var LoaderList = map[string]Loader{
	"default": Loader{
		FormatCreate: func(c *client.Client, obj map[string]interface{}, defaults map[string]string) map[string]interface{} {
			logrus.WithFields(logrus.Fields{
				"in": obj,
			}).Debugf("default.FormatCreate")
			//cleanObj(obj, "dbid", "hostdbid", "appprototypedbid") //TODO find matching prototype for app //TODO ask for password
			cleanObj(obj, "dbid")
			if tenant, exist := obj["tenantdbid"]; exist {
				obj["tenantdbid"] = searchFor(c, "CfgTenant", tenant.(string), defaults)
			}
			if folder, exist := obj["folderid"]; exist {
				obj["folderid"] = searchFor(c, "CfgFolder", folder.(string), defaults)
			}
			logrus.WithFields(logrus.Fields{
				"out": obj,
			}).Debugf("default.FormatCreate")
			return obj
		},
		FormatUpdate: func(c *client.Client, src, obj map[string]interface{}, defaults map[string]string) map[string]interface{} {
			//TODO find matching prototype for app //TODO ask for password
			obj["dbid"] = src["dbid"]
			//TODO use default of src
			if tenant, exist := src["tenantdbid"]; exist {
				obj["tenantdbid"] = searchFor(c, "CfgTenant", tenant.(string), defaults)
			}
			if folder, exist := src["folderid"]; exist {
				obj["folderid"] = searchFor(c, "CfgFolder", folder.(string), defaults)
			}
			return obj
		},
	},
}

func searchFor(c *client.Client, t string, id string, defaults map[string]string) string {
	//TODO search for seam id
	//TODO history

	//Use default value by default
	if def, ok := defaults[t]; ok {
		return def
	}

	list := ListObject(c, t)
	logrus.WithField("list", list).WithField("type", t).Debugf("Fetched list")
	logrus.Infof("Please choose a %s :", t)
	val := prompt.Input("> ", func(d prompt.Document) []prompt.Suggest {
		//log.Print(d.Text)
		//		if d.Text == "" {
		//			d.Text = "TEST"
		//		}
		//logrus.WithField("list", list).Info("Fetched list")
		//TODO put id corresponding obj if any first
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
	logrus.WithFields(logrus.Fields{
		"cmp":    "dbid",
		"src":    src["dbid"],
		"dst":    dst["dbid"],
		"result": src["dbid"] == dst["dbid"],
	}).Debug("Matching dbid")
	return src["dbid"] == dst["dbid"]
}
func MatchName(src, dst map[string]interface{}) bool {
	if src["name"] != nil && dst["name"] != nil {
		logrus.WithFields(logrus.Fields{
			"cmp":    "name",
			"src":    src["name"],
			"dst":    dst["name"],
			"result": src["name"] == dst["name"],
		}).Debug("Matching name")
		return src["name"] == dst["name"]
	}
	if src["username"] != nil && dst["username"] != nil {
		logrus.WithFields(logrus.Fields{
			"cmp":    "username",
			"src":    src["username"],
			"dst":    dst["username"],
			"result": src["username"] == dst["username"],
		}).Debug("Matching username")
		return src["username"] == dst["username"]
	}
	if src["number"] != nil && dst["number"] != nil {
		logrus.WithFields(logrus.Fields{
			"cmp":    "number",
			"src":    src["number"],
			"dst":    dst["number"],
			"result": src["number"] == dst["number"],
		}).Debug("Matching number")
		return src["number"] == dst["number"]
	}
	if src["logincode"] != nil && dst["logincode"] != nil {
		logrus.WithFields(logrus.Fields{
			"cmp":    "logincode",
			"src":    src["logincode"],
			"dst":    dst["logincode"],
			"result": src["logincode"] == dst["logincode"],
		}).Debug("Matching logincode")
		return src["logincode"] == dst["logincode"]
	}
	return false
}
func MatchIdName(src, dst map[string]interface{}) bool { //TODO Manage Person (username)
	return MatchName(src, dst) && MatchId(src, dst)
}
