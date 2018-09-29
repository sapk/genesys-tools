// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package format

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"

	"github.com/sapk/go-genesys/api/object"
)

type Formater struct {
	Format      func(object.Type, map[string]interface{}, map[string][]interface{}) string
	FormatCSV   func(object.Type, map[string]interface{}, map[string][]interface{}) string
	FormatShort func(object.Type, map[string]interface{}) string
}

func defaultShortFormater(objType object.Type, obj map[string]interface{}) string {
	name := GetFileName(obj)
	//TODO better
	if obj["dbid"] != nil {
		if obj["subtype"] != nil && obj["subtype"] != "" {
			return fmt.Sprintf(" - [%s](./%s/%s \\(%s\\)) (%s)\n", name, objType.Desc, name, obj["dbid"], obj["subtype"])
		}
		return fmt.Sprintf(" - [%s](./%s/%s \\(%s\\))\n", name, objType.Desc, name, obj["dbid"])
	}
	if obj["subtype"] != nil && obj["subtype"] != "" {
		return fmt.Sprintf(" - [%s](./%s/%s) (%s)\n", name, objType.Desc, name, obj["subtype"])
	}
	return fmt.Sprintf(" - [%s](./%s/%s)\n", name, objType.Desc, name)
	//	if objType.IsDumpable {
	//} else {
	//return ""
	//}
}

var FormaterList = map[string]Formater{
	"default": Formater{
		func(objType object.Type, obj map[string]interface{}, data map[string][]interface{}) string {
			name := GetFileName(obj)
			ret := "# " + name + "\n"
			ret += "\n"

			ret += dumpAvailableInformation(obj, data) + "\n"
			ret += formatOptions(obj, data)
			ret += formatAnnexes(obj, data)
			ret += dumpBackup(obj)
			return ret
		},
		nil,
		defaultShortFormater,
	},
	"CfgFormat": Formater{
		func(objType object.Type, obj map[string]interface{}, data map[string][]interface{}) string {
			name := GetFileName(obj)
			ret := "# " + name + "\n"
			ret += "\n"

			ret += dumpAvailableInformation(obj, data) + "\n"
			ret += formatFields(obj, data)
			ret += formatOptions(obj, data)
			ret += formatAnnexes(obj, data)
			ret += dumpBackup(obj)
			return ret
		},
		nil,
		defaultShortFormater,
	},
}

var keyInformations = []struct {
	ID     string
	Name   string
	Format func(interface{}, map[string][]interface{}) string
}{
	//Generic
	{"dbid", "DBID", nil},
	{"name", "Name", nil},
	//In case of Person
	{"username", "Username", nil},
	{"tenantdbid", "Tenant", funcFindByType("CfgTenant")},
	{"hostdbid", "Host", funcFindByType("CfgHost")},
	{"type", "Type", nil},
	{"subtype", "SubType", nil},
	{"componenttype", "Componenttype", nil},
	{"isserver", "Isserver", nil},
	{"version", "Version", nil},
	{"state", "State", nil},
	{"folderid", "Folder path", findFolderPath},
	{"description", "Description", nil},
	{"tenantdbids", "Tenants", findTenants},
	//DN
	{"number", "Number", nil},
	{"switchdbid", "Switch", funcFindByType("CfgSwitch")},
	//Person
	{"firstname", "Firstname", nil},
	{"lastname", "Lastname", nil},
	{"employeeid", "Employeeid", nil},
	{"isagent", "IsAgent", nil},
	{"isexternalauth", "IsExternalAuth", nil},
	{"sitedbid", "Site", funcFindByType("CfgTenant")},
	{"placedbid", "Place", funcFindByType("CfgPlace")},
	{"capacityruledbid", "Capacity Rule", funcFindByType("CfgScript")}, //type=CfgScript&subtype=CFGCapacityRule
	//TODO capacitytabledbid & quotatabledbid for Agent Group
	{"contractdbid", "Contract", funcFindByType("CfgObjectiveTable")}, //type=CfgObjectiveTable&subtype=CFGOTTCostContract
	//Host
	{"ipaddress", "Ipaddress", nil},
	{"scsdbid", "SCS", funcFindByType("CfgApplication")},
	{"lcaport", "Lcaport", nil},
	{"ostype", "Ostype", nil},
	//App
	{"appprototypedbid", "App Template", funcFindByType("CfgAppPrototype")},
	{"startuptype", "Startuptype", nil},
	{"workdirectory", "Workdirectory", nil},
	{"commandline", "Commandline", nil},
	{"commandlinearguments", "Commandlinearguments", nil},
	{"autorestart", "Autorestart", nil},
	{"timeout", "Timeout", nil},
	{"port", "Port principal", nil},
	{"redundancytype", "Redundancytype", nil},
	{"isprimary", "Isprimary", nil},
	{"backupserverdbid", "Backup Server", funcFindByType("CfgApplication")},
	//TODO Add Host key inf and other
}

func formatFields(obj map[string]interface{}, data map[string][]interface{}) string {
	ret := "## Fields: \n"
	var ids struct {
		Id object.CfgDBIDList `json:"id"`
	}
	err := mapstructure.Decode(obj["fielddbids"], &ids)
	if err != nil {
		logrus.Warnf("Fail to convert to CfgDBIDList")
		return err.Error()
	}
	for _, id := range ids.Id {
		o := findObj("CfgField", id.Dbid, data)
		ret += " - " + catchNotString(o["name"]) + " (" + id.Dbid + ") : " + catchNotString(o["description"]) + " \n"
	}
	//
	return ret + " \n"
}

//Find best suitable name
func GetFileName(obj map[string]interface{}) string {

	for _, idPossible := range []string{"name", "username", "number", "logincode"} {
		name, ok := obj[idPossible].(string)
		if ok {
			name = strings.Replace(name, "/", " - ", -1)
			name = strings.Replace(name, "\\", " - ", -1)
			name = strings.Replace(name, "?", "-", -1) //Forbidden char ? on windows fs
			return name
		}
	}

	return ""
}

func findTenants(tenantdbids interface{}, data map[string][]interface{}) string {
	if tenantdbids == nil {
		return "" //Empty
	}
	obj := tenantdbids.(map[string]interface{})
	_, ok := obj["id"]
	if !ok || obj["id"] == nil {
		return "" //Empty
	}
	var ids object.CfgDBIDList
	err := mapstructure.Decode(obj["id"], &ids)
	if err != nil {
		logrus.Warnf("Fail to convert to CfgDBIDList")
		return err.Error()
	}
	ret := " "
	for _, id := range ids {
		ret += findObjName("CfgTenant", id.Dbid, data) + ", "
	}
	return strings.Trim(ret, ", ")
}

func findObj(t string, id string, data map[string][]interface{}) map[string]interface{} {
	if id == "0" {
		return nil
	}
	for _, _o := range data[t] {
		o := _o.(map[string]interface{})
		if o["dbid"].(string) == id {
			return o
		}
	}
	return nil
}

func findObjName(t string, id string, data map[string][]interface{}) string {
	o := findObj(t, id, data)
	if o == nil {
		return id
	}
	name, ok := o["name"].(string)
	if ok {
		return name
	}
	return id
}

func findFolderPath(id interface{}, data map[string][]interface{}) string {
	idFolder := id.(string)
	f := findObj("CfgFolder", idFolder, data) //Chainload to have full path
	if f == nil {
		return idFolder
	}
	name, ok := f["name"].(string)
	if ok {
		parent, ok := f["folderid"].(string)
		if ok {
			return filepath.Join(findFolderPath(parent, data), name)
		}
		return "/" + name
	}
	return idFolder //Chainload to have full path
}

func funcFindByType(t string) func(interface{}, map[string][]interface{}) string {
	return func(id interface{}, data map[string][]interface{}) string {
		return findObjName(t, id.(string), data)
	}
}

func dumpAvailableInformation(obj map[string]interface{}, data map[string][]interface{}) string {
	ret := "## Informations: \n"
	for _, inf := range keyInformations {
		val, ok := obj[inf.ID]
		//.(string)
		if ok {
			if inf.Format != nil {
				val = inf.Format(val, data)
			}
			ret += " " + inf.Name + ": " + val.(string) + "  \n"
		}
	}
	return ret
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

func formatAnnexes(obj map[string]interface{}, data map[string][]interface{}) string {

	var props object.Userproperties
	err := mapstructure.Decode(obj["userproperties"], &props)
	if err != nil {
		logrus.Warnf("Fail to convert to Userproperties")
		return err.Error()
	}

	sectionsAnnex := treeset.NewWithStringComparator()
	annexes := make(map[string]*treemap.Map)
	for _, o := range props.Property {
		sectionsAnnex.Add(o.Section)
		if _, ok := annexes[o.Section]; !ok {
			//Init
			annexes[o.Section] = treemap.NewWithStringComparator()
		}
		annexes[o.Section].Put(o.Key, o.Value)
	}
	annexList := ""
	for _, s := range sectionsAnnex.Values() {
		sec := s.(string)
		annexList += " [" + sec + "]  \n"
		for _, o := range annexes[sec].Keys() {
			opt := o.(string)
			val, _ := annexes[s.(string)].Get(opt)
			annexList += "  " + opt + " = " + val.(string) + "  \n"
		}
		//optList += " - [" + o.Section + "] / " + o.Key + " = " + o.Value + "\n"
	}
	ret := fmt.Sprintf("## Annexes (%d): \n", strings.Count(annexList, "\n")-sectionsAnnex.Size())
	ret += annexList
	ret += "  \n"
	return ret
}

func formatOptions(obj map[string]interface{}, data map[string][]interface{}) string {
	var opts object.Options
	err := mapstructure.Decode(obj["options"], &opts)
	if err != nil {
		logrus.Warnf("Fail to convert to Options")
		return err.Error()
	}

	sections := treeset.NewWithStringComparator()
	options := make(map[string]*treemap.Map)
	for _, o := range opts.Property {
		sections.Add(o.Section)
		if _, ok := options[o.Section]; !ok {
			//Init
			options[o.Section] = treemap.NewWithStringComparator()
		}
		options[o.Section].Put(o.Key, o.Value)
	}
	optList := ""
	for _, s := range sections.Values() {
		sec := s.(string)
		optList += " [" + sec + "]  \n"
		for _, o := range options[sec].Keys() {
			opt := o.(string)
			val, _ := options[s.(string)].Get(opt)
			optList += "  " + opt + " = " + val.(string) + "  \n"
		}
		//optList += " - [" + o.Section + "] / " + o.Key + " = " + o.Value + "\n"
	}

	ret := fmt.Sprintf("## Options (%d): \n", strings.Count(optList, "\n")-sections.Size())
	ret += optList
	ret += "  \n"
	return ret
}
func GetObjectType(obj map[string]interface{}) *object.Type {
	for _, o := range object.TypeList {
		if obj["type"].(string) == o.Name {
			return &o
		}
	}
	return nil
}

//TODO call this from get Filename
func Name(obj map[string]interface{}) string {
	for _, idPossible := range []string{"name", "username", "number", "logincode"} {
		name, ok := obj[idPossible].(string)
		if ok {
			return name
		}
	}
	return ""
}

//Call the good formatter if exist or use the default
func FormatShortObj(obj map[string]interface{}) string {
	objType := GetObjectType(obj)
	if objType == nil {
		return "Object Type Unknown: " + obj["type"].(string)
	}
	if f, ok := FormaterList[objType.Name]; ok {
		return f.FormatShort(*objType, obj)
	}
	return FormaterList["default"].FormatShort(*objType, obj)
}

//Call the good formatter if exist or use the default
func FormatObj(objType object.Type, obj map[string]interface{}, data map[string][]interface{}) string {
	if f, ok := FormaterList[objType.Name]; ok {
		return f.Format(objType, obj, data)
	}
	return FormaterList["default"].Format(objType, obj, data)
}
