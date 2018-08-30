// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package loader

import (
	"reflect"

	"github.com/sapk/go-genesys/api/client"
	"github.com/sirupsen/logrus"
)

func init() {
	LoaderList["CfgAgentLogin"] = Loader{
		FormatCreate: func(c *client.Client, obj map[string]interface{}, defaults map[string]string) map[string]interface{} {
			logrus.WithFields(logrus.Fields{
				"in": obj,
			}).Debugf("CfgAgentLogin.FormatCreate")
			obj = LoaderList["default"].FormatCreate(c, obj, defaults)
			if sw, exist := obj["switchdbid"]; exist {
				obj["switchdbid"] = searchFor(c, "CfgSwitch", sw.(string), defaults)
			}
			return obj
		},
		FormatUpdate: func(c *client.Client, src, obj map[string]interface{}, defaults map[string]string) map[string]interface{} {
			logrus.WithFields(logrus.Fields{
				"src": src,
				"obj": obj,
			}).Debugf("CfgAgentLogin.FormatUpdate")
			//TODO reuse by default value of src
			obj = LoaderList["default"].FormatUpdate(c, src, obj, defaults)
			if sw, exist := obj["switchdbid"]; exist {
				obj["switchdbid"] = searchFor(c, "CfgSwitch", sw.(string), defaults)
			}
			return obj
		},
	}
	LoaderList["CfgPerson"] = Loader{
		FormatCreate: func(c *client.Client, obj map[string]interface{}, defaults map[string]string) map[string]interface{} {
			logrus.WithFields(logrus.Fields{
				"in": obj,
			}).Debugf("CfgPerson.FormatCreate")
			obj = LoaderList["default"].FormatCreate(c, obj, defaults)
			//agentlogins":{"agentlogininfo":[{"agentlogindbid":"151","wrapuptime":"0"}]},"appranks":{"apprank":[]},"
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

			if agentlogins, exist := obj["agentlogins"]; exist {
				emptyDBIDList := struct {
					Agentlogininfo []struct {
						Agentlogindbid string `json:"agentlogindbid"`
						Wrapuptime     string `json:"wrapuptime"`
					} `json:"agentlogininfo"`
				}{Agentlogininfo: []struct {
					Agentlogindbid string `json:"agentlogindbid"`
					Wrapuptime     string `json:"wrapuptime"`
				}{}}
				eq := reflect.DeepEqual(agentlogins, emptyDBIDList)
				if !eq {
					logrus.WithFields(logrus.Fields{
						"agentlogins": agentlogins,
					}).Warn("Attached Agent Login link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["agentlogins"] = emptyDBIDList
				}
			}
			if appranks, exist := obj["appranks"]; exist {
				//TODO Apprank struct
				emptyDBIDList := struct {
					apprank []struct {
						/*
							Agentlogindbid string `json:"agentlogindbid"`
							Wrapuptime     string `json:"wrapuptime"`
						*/
					} `json:"apprank"`
				}{apprank: []struct{}{}}
				eq := reflect.DeepEqual(appranks, emptyDBIDList)
				if !eq {
					logrus.WithFields(logrus.Fields{
						"appranks": appranks,
					}).Warn("Attached appranks link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["appranks"] = emptyDBIDList
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
			if placedbid, exist := obj["placedbid"]; exist {
				if placedbid != "0" {
					logrus.WithFields(logrus.Fields{
						"placedbid": placedbid,
					}).Warn("Attached place link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["placedbid"] = "0"
				}
			}
			if skilllevels, exist := obj["skilllevels"]; exist {
				//TODO skilllevels struct
				//"skilllevels":{"skilllevel":[{"level":"5","skilldbid":"105"}]}
				emptyDBIDList := struct {
					skilllevel []struct {
						Level     string `json:"level"`
						Skilldbid string `json:"skilldbid"`
					} `json:"skilllevel"`
				}{skilllevel: []struct {
					Level     string `json:"level"`
					Skilldbid string `json:"skilldbid"`
				}{}}
				eq := reflect.DeepEqual(skilllevels, emptyDBIDList)
				if !eq {
					logrus.WithFields(logrus.Fields{
						"skilllevels": skilllevels,
					}).Warn("Attached skilllevels link will be lost")
					//TODO search
					obj["skilllevels"] = emptyDBIDList
				}
			}
			delete(obj, "password") //Clear Password
			logrus.Warn("Possible attached Agent password will be lost")
			//TODO password
			return obj
		},
		FormatUpdate: func(c *client.Client, src, obj map[string]interface{}, defaults map[string]string) map[string]interface{} {
			logrus.WithFields(logrus.Fields{
				"src": src,
				"obj": obj,
			}).Debugf("CfgPerson.FormatUpdate")
			obj = LoaderList["default"].FormatCreate(c, obj, defaults)
			//TODO reuse by default value of src
			//agentlogins":{"agentlogininfo":[{"agentlogindbid":"151","wrapuptime":"0"}]},"appranks":{"apprank":[]},"
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

			if agentlogins, exist := obj["agentlogins"]; exist {
				emptyDBIDList := struct {
					Agentlogininfo []struct {
						Agentlogindbid string `json:"agentlogindbid"`
						Wrapuptime     string `json:"wrapuptime"`
					} `json:"agentlogininfo"`
				}{Agentlogininfo: []struct {
					Agentlogindbid string `json:"agentlogindbid"`
					Wrapuptime     string `json:"wrapuptime"`
				}{}}
				eq := reflect.DeepEqual(agentlogins, emptyDBIDList)
				if !eq {
					logrus.WithFields(logrus.Fields{
						"agentlogins": agentlogins,
					}).Warn("Attached Agent Login link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["agentlogins"] = emptyDBIDList
				}
			}
			if appranks, exist := obj["appranks"]; exist {
				//TODO Apprank struct
				emptyDBIDList := struct {
					apprank []struct {
						/*
							Agentlogindbid string `json:"agentlogindbid"`
							Wrapuptime     string `json:"wrapuptime"`
						*/
					} `json:"apprank"`
				}{apprank: []struct{}{}}
				eq := reflect.DeepEqual(appranks, emptyDBIDList)
				if !eq {
					logrus.WithFields(logrus.Fields{
						"appranks": appranks,
					}).Warn("Attached appranks link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["appranks"] = emptyDBIDList
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
			if placedbid, exist := obj["placedbid"]; exist {
				if placedbid != "0" {
					logrus.WithFields(logrus.Fields{
						"placedbid": placedbid,
					}).Warn("Attached place link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["placedbid"] = "0"
				}
			}
			if skilllevels, exist := obj["skilllevels"]; exist {
				//TODO skilllevels struct
				//"skilllevels":{"skilllevel":[{"level":"5","skilldbid":"105"}]}
				emptyDBIDList := struct {
					skilllevel []struct {
						Level     string `json:"level"`
						Skilldbid string `json:"skilldbid"`
					} `json:"skilllevel"`
				}{skilllevel: []struct {
					Level     string `json:"level"`
					Skilldbid string `json:"skilldbid"`
				}{}}
				eq := reflect.DeepEqual(skilllevels, emptyDBIDList)
				if !eq {
					logrus.WithFields(logrus.Fields{
						"skilllevels": skilllevels,
					}).Warn("Attached skilllevels link will be lost")
					//TODO search
					obj["skilllevels"] = emptyDBIDList
				}
			}

			delete(obj, "password") //Clear Password
			logrus.Warn("Possible attached Agent password will be lost")
			//TODO password
			return obj

		},
	}
}
