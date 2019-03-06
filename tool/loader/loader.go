// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package loader

import (
	"github.com/c-bata/go-prompt"
	"github.com/rs/zerolog/log"

	"github.com/sapk/genesys-tools/tool/format"
	"github.com/sapk/go-genesys/api/client"
)

var data map[string][]map[string]interface{}

func init() {
	data = make(map[string][]map[string]interface{})
}

func ListObject(c *client.Client, t string) []map[string]interface{} {
	log.Debug().Interface("type", t).Msgf("ListObject")
	if l, ok := data[t]; ok {
		log.Debug().Interface("list", l).Msgf("ListObject in cache")
		return l //don't get list allready done
	}
	var list []map[string]interface{}
	c.ListObject(t, &list)
	log.Debug().Interface("list", list).Msgf("ListObject fetched")
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
			log.Debug().Interface("in", obj).Msgf("default.FormatCreate")
			//cleanObj(obj, "dbid", "hostdbid", "appprototypedbid") //TODO find matching prototype for app //TODO ask for password
			cleanObj(obj, "dbid")
			if tenant, exist := obj["tenantdbid"]; exist {
				obj["tenantdbid"] = searchFor(c, "CfgTenant", tenant.(string), defaults)
			}
			if folder, exist := obj["folderid"]; exist {
				obj["folderid"] = searchFor(c, "CfgFolder", folder.(string), defaults)
			}
			/*
				if userproperties, exist := obj["userproperties"]; exist {
						obj["userproperties"] = cleanEmptyAnnexes(userproperties.(map[string]interface{}))
				}
			*/
			log.Debug().Interface("out", obj).Msgf("default.FormatCreate")
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
			/*
				if userproperties, exist := obj["userproperties"]; exist {
					obj["userproperties"] = cleanEmptyAnnexes(userproperties.(map[string]interface{}))
				}
			*/
			return obj
		},
	},
}

/*Some key failed but this is not the solution you are looking for
func cleanEmptyAnnexes(o map[string]interface{}) interface{} {
	var obj object.Userproperties
	err := mapstructure.Decode(o, &obj)
	if err != nil {
		log.Warnf("Fail to convert to Userproperties -> Skipping cleaning")
		return o
	}
	for id, val := range obj.Property {
		if val.Key == "" {
			obj.Property = append(obj.Property[:id], obj.Property[id+1:]...)
		}
	}
	return obj
}
*/
func searchFor(c *client.Client, t string, id string, defaults map[string]string) string {
	//TODO search for seam id
	//TODO history

	//Use default value by default
	if def, ok := defaults[t]; ok {
		log.Debug().Interface("def", def).Interface("type", t).Msgf("Using default value")
		return def
	}

	list := ListObject(c, t)
	log.Debug().Interface("list", list).Interface("type", t).Msgf("Fetched list")
	log.Info().Msgf("Please choose a %s :", t)
	val := prompt.Input("> ", func(d prompt.Document) []prompt.Suggest {
		//log.Print(d.Text)
		//		if d.Text == "" {
		//			d.Text = "TEST"
		//		}
		//log.Interface("list", list).Info("Fetched list")
		//TODO put id corresponding obj if any first
		s := make([]prompt.Suggest, len(list))
		for i, o := range list {
			//log.Interface("obj", o).Info("Add to list")
			s[i] = prompt.Suggest{o["dbid"].(string), format.GetFileName(o)}
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	})

	log.Info().Msgf("You selected " + val)
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
			log.Info().Interface("obj", o).Msg("Matched object")
			ret = append(ret, o) //TODO best allocate
		}
	}
	return ret
}

func MatchId(src, dst map[string]interface{}) bool {
	log.Debug().Str("cmp", "dbid").Interface("src", src["dbid"]).Interface("dst", dst["dbid"]).Interface("result", src["dbid"] == dst["dbid"]).Msg("Matching dbid")
	return src["dbid"] == dst["dbid"]
}

func MatchByEl(src, dst map[string]interface{}, el string) bool {
	log.Debug().Str("cmp", el).Interface("src", src[el]).Interface("dst", dst[el]).Interface("result", src[el] == dst[el]).Msg("Matching name")
	return src[el] == dst[el]
}

func MatchName(src, dst map[string]interface{}) bool {
	if src["name"] != nil && dst["name"] != nil {
		return MatchByEl(src, dst, "name")
	}
	if src["username"] != nil && dst["username"] != nil {
		return MatchByEl(src, dst, "username")
	}
	if src["number"] != nil && dst["number"] != nil {
		return MatchByEl(src, dst, "number")
	}
	if src["logincode"] != nil && dst["logincode"] != nil {
		return MatchByEl(src, dst, "logincode")
	}
	return false
}

func MatchIdName(src, dst map[string]interface{}) bool { //TODO Manage Person (username)
	return MatchName(src, dst) && MatchId(src, dst)
}
