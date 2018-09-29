// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package format

import (
	"fmt"
	"strings"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/mitchellh/mapstructure"
	"github.com/sapk/go-genesys/api/object"
	"github.com/sirupsen/logrus"
)

func init() {
	FormaterList["CfgAgentGroup"] = Formater{
		func(objType object.Type, obj map[string]interface{}, data map[string][]interface{}) string {
			ret := "# " + obj["name"].(string) + "\n"
			ret += "\n"

			ret += dumpAvailableInformation(obj, data) + "\n"
			ret += formatAgentGroupDetails(obj, data)
			ret += formatOptions(obj, data)
			ret += formatAnnexes(obj, data)
			ret += dumpBackup(obj)
			return ret
		},
		nil,
		defaultShortFormater,
	}

	FormaterList["CfgPerson"] = Formater{
		func(objType object.Type, obj map[string]interface{}, data map[string][]interface{}) string {
			ret := "# " + obj["username"].(string) + "\n"
			ret += "\n"
			ret += dumpAvailableInformation(obj, data) + "\n"
			ret += formatPersonDetails(obj, data)
			ret += formatOptions(obj, data)
			ret += formatAnnexes(obj, data)
			ret += dumpBackup(obj)
			return ret
		},
		nil,
		func(objType object.Type, obj map[string]interface{}) string {
			name := GetFileName(obj)
			displayname := strings.TrimSpace(catchNotString(obj["firstname"]) + " " + catchNotString(obj["lastname"]))
			if displayname == "" {
				displayname = obj["username"].(string)
			}
			//	if objType.IsDumpable {
			return fmt.Sprintf(" - [%s](./%s/%s \\(%s\\)) (%s/%s)\n", displayname, objType.Desc, name, obj["dbid"], obj["username"], obj["employeeid"])
		},
	}
	FormaterList["CfgAccessGroup"] = Formater{
		func(objType object.Type, obj map[string]interface{}, data map[string][]interface{}) string {
			ret := "# " + obj["name"].(string) + "\n"
			ret += "\n"

			ret += dumpAvailableInformation(obj, data) + "\n"
			ret += formatAccessGroupDetails(obj, data)
			ret += formatOptions(obj, data)
			ret += formatAnnexes(obj, data)
			ret += dumpBackup(obj)
			return ret
		},
		nil,
		defaultShortFormater,
	}
}

func catchNotString(obj interface{}) string {
	/*
		if obj == nil {
			return ""
		}
	*/
	str, ok := obj.(string)
	if ok {
		return str
	}
	return ""
}
func formatPersonDetails(obj map[string]interface{}, data map[string][]interface{}) string {
	//TODO appranks
	//TODO skilllevels
	//TODO agentlogins
	ret := ""
	var person object.CfgPerson
	err := mapstructure.Decode(obj, &person)
	if err != nil {
		logrus.Warnf("Fail to convert to CfgPerson")
		return err.Error()
	}
	return ret
}
func formatPersonList(list object.CfgDBIDList, data map[string][]interface{}) (string, int) {
	users := treemap.NewWithStringComparator()
	for _, u := range list {
		user := findObj("CfgPerson", u.Dbid, data)
		if user != nil {
			users.Put(user["username"].(string), user)
		} else {
			users.Put(u.Dbid, nil)
		}
	}
	userList := ""
	for _, u := range users.Keys() {
		username := u.(string)
		val, _ := users.Get(username)
		if val != nil {
			var user object.CfgPerson
			err := mapstructure.Decode(val, &user)
			if err != nil {
				logrus.Warnf("Fail to convert to CfgPerson")
				userList += "  dbid:" + username + "  \n"
			} else {
				userList += "  " + user.Username + " / " + user.Firstname + " " + user.Lastname + " (" + user.Employeeid + "/" + user.Dbid + ")  \n"
			}
		} else {
			userList += "  dbid:" + username + "  \n"
		}
	}
	return userList, users.Size()
}

func formatAccessGroupDetails(obj map[string]interface{}, data map[string][]interface{}) string {
	var ag object.CfgAccessGroup
	err := mapstructure.Decode(obj, &ag)
	if err != nil {
		logrus.Warnf("Fail to convert to CfgAccessGroup")
		return err.Error()
	}
	uL, uS := formatPersonList(ag.Memberids.Idtype, data) //TODO catch when memberids is not CFGPerson
	ret := fmt.Sprintf("## Members (%d): \n", uS)
	ret += uL + "\n"

	return ret
}
func formatAgentGroupDetails(obj map[string]interface{}, data map[string][]interface{}) string {
	var ag object.CfgAgentGroup
	err := mapstructure.Decode(obj, &ag)
	if err != nil {
		logrus.Warnf("Fail to convert to CfgAgentGroup")
		return err.Error()
	}
	uL, uS := formatPersonList(ag.Agentdbids.ID, data)
	ret := fmt.Sprintf("## Agents (%d): \n", uS)
	ret += uL + "\n"

	mL, mS := formatPersonList(ag.Managerdbids.ID, data)
	ret += fmt.Sprintf("## Managers (%d): \n", mS)
	ret += mL + "\n"

	return ret
}
