// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package loader

import (
	"reflect"

	"github.com/sapk/genesys-tools/api/client"
	"github.com/sapk/genesys-tools/api/object"
	"github.com/sirupsen/logrus"
)

func init() {
	LoaderList["CfgPlace"] = Loader{
		FormatCreate: func(c *client.Client, obj map[string]interface{}, defaults map[string]string) map[string]interface{} {
			logrus.WithFields(logrus.Fields{
				"in": obj,
			}).Debugf("CfgPlace.FormatCreate")
			obj = LoaderList["default"].FormatCreate(c, obj, defaults)
			//lost link to contactdbid capacityruledbid dndbids sitedbid

			if contactdbid, exist := obj["contactdbid"]; exist {
				if contactdbid != "0" {
					logrus.WithFields(logrus.Fields{
						"contactdbid": contactdbid,
					}).Warn("Attached contract link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["contactdbid"] = "0"
				}
			}
			if capacityruledbid, exist := obj["capacityruledbid"]; exist {
				if capacityruledbid != "0" {
					logrus.WithFields(logrus.Fields{
						"capacityruledbid": capacityruledbid,
					}).Warn("Attached capacityrule link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["capacityruledbid"] = "0"
				}
			}
			if dndbids, exist := obj["dndbids"]; exist {
				emptyDBIDList := struct {
					Id object.CfgDBIDList `json:"id"`
				}{Id: object.CfgDBIDList{}}
				//{"id":[{"dbid":"143"}]}
				eq := reflect.DeepEqual(dndbids, emptyDBIDList)
				if !eq {
					//if dndbids != "{\"id\":[]}" {
					logrus.WithFields(logrus.Fields{
						"dndbids": dndbids,
					}).Warn("Attached DNs link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["dndbids"] = emptyDBIDList
				}
			}
			if sitedbid, exist := obj["sitedbid"]; exist {
				if sitedbid != "0" {
					logrus.WithFields(logrus.Fields{
						"sitedbid": sitedbid,
					}).Warn("Attached site link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["sitedbid"] = "0"
				}
			}
			return obj
		},
		FormatUpdate: func(c *client.Client, src, obj map[string]interface{}, defaults map[string]string) map[string]interface{} {
			logrus.WithFields(logrus.Fields{
				"src": src,
				"obj": obj,
			}).Debugf("CfgPlace.FormatUpdate")
			//TODO reuse by default value of src
			obj = LoaderList["default"].FormatUpdate(c, src, obj, defaults)
			//lost link to contactdbid capacityruledbid dndbids sitedbid

			if contactdbid, exist := obj["contactdbid"]; exist {
				if contactdbid != "0" {
					logrus.WithFields(logrus.Fields{
						"contactdbid": contactdbid,
					}).Warn("Attached contract link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["contactdbid"] = "0"
				}
			}
			if capacityruledbid, exist := obj["capacityruledbid"]; exist {
				if capacityruledbid != "0" {
					logrus.WithFields(logrus.Fields{
						"capacityruledbid": capacityruledbid,
					}).Warn("Attached capacityrule link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["capacityruledbid"] = "0"
				}
			}
			if dndbids, exist := obj["dndbids"]; exist {
				emptyDBIDList := struct {
					Id object.CfgDBIDList `json:"id"`
				}{Id: object.CfgDBIDList{}}
				//{"id":[{"dbid":"143"}]}
				eq := reflect.DeepEqual(dndbids, emptyDBIDList)
				if !eq {
					//if dndbids != "{\"id\":[]}" {
					logrus.WithFields(logrus.Fields{
						"dndbids": dndbids,
					}).Warn("Attached DNs link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["dndbids"] = emptyDBIDList
				}
			}
			if sitedbid, exist := obj["sitedbid"]; exist {
				if sitedbid != "0" {
					logrus.WithFields(logrus.Fields{
						"sitedbid": sitedbid,
					}).Warn("Attached site link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["sitedbid"] = "0"
				}
			}
			return obj
		},
	}
}
