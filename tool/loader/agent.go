// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package loader

import (
	"reflect"

	"github.com/rs/zerolog/log"

	"github.com/sapk/go-genesys/api/client"
	"github.com/sapk/go-genesys/api/object"
)

func init() {
	LoaderList["CfgAgentGroup"] = Loader{
		FormatCreate: func(c *client.Client, obj map[string]interface{}, defaults map[string]string) map[string]interface{} {
			log.Debug().Interface("in", obj).Msg("CfgAgentGroup.FormatCreate")
			obj = LoaderList["default"].FormatCreate(c, obj, defaults)

			for _, val := range []string{"capacityruledbid", "capacitytabledbid", "contractdbid", "quotatabledbid", "sitedbid"} {
				if dbid, exist := obj[val]; exist {
					if dbid != "0" {
						log.Warn().Interface("val", dbid).Msgf("Attached %s link will be lost", val)
						//TODO search
						obj[val] = "0"
					}
				}
			}

			for _, val := range []string{"agentdbids", "managerdbids"} {
				if dbids, exist := obj[val]; exist {
					emptyDBIDList := struct {
						Id object.CfgDBIDList `json:"id"`
					}{Id: object.CfgDBIDList{}}
					eq := reflect.DeepEqual(dbids, emptyDBIDList)
					if !eq {
						log.Warn().Interface("val", dbids).Msgf("Attached %s link will be lost", val)
						obj[val] = emptyDBIDList
					}
				}
			}
			return obj
		},
		FormatUpdate: func(c *client.Client, src, obj map[string]interface{}, defaults map[string]string) map[string]interface{} {
			log.Debug().Interface("src", src).Interface("obj", obj).Msg("CfgAgentGroup.FormatUpdate")
			//TODO reuse by default value of src
			obj = LoaderList["default"].FormatUpdate(c, src, obj, defaults)
			for _, val := range []string{"capacityruledbid", "capacitytabledbid", "contractdbid", "quotatabledbid", "sitedbid"} {
				if dbid, exist := obj[val]; exist {
					if dbid != "0" {
						log.Warn().Interface("val", dbid).Msgf("Attached %s link will be lost", val)
						//TODO search
						obj[val] = "0"
					}
				}
			}

			for _, val := range []string{"agentdbids", "managerdbids"} {
				if dbids, exist := obj[val]; exist {
					emptyDBIDList := struct {
						Id object.CfgDBIDList `json:"id"`
					}{Id: object.CfgDBIDList{}}
					eq := reflect.DeepEqual(dbids, emptyDBIDList)
					if !eq {
						log.Warn().Interface("val", dbids).Msgf("Attached %s link will be lost", val)
						obj[val] = emptyDBIDList
					}
				}
			}
			return obj
		},
	}
	LoaderList["CfgAgentLogin"] = Loader{
		FormatCreate: func(c *client.Client, obj map[string]interface{}, defaults map[string]string) map[string]interface{} {
			log.Debug().Interface("in", obj).Msg("CfgAgentLogin.FormatCreate")
			obj = LoaderList["default"].FormatCreate(c, obj, defaults)
			if sw, exist := obj["switchdbid"]; exist {
				obj["switchdbid"] = searchFor(c, "CfgSwitch", sw.(string), defaults)
			}
			return obj
		},
		FormatUpdate: func(c *client.Client, src, obj map[string]interface{}, defaults map[string]string) map[string]interface{} {
			log.Debug().Interface("src", src).Interface("obj", obj).Msg("CfgAgentLogin.FormatUpdate")
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
			log.Debug().Interface("in", obj).Msg("CfgPerson.FormatCreate")
			obj = LoaderList["default"].FormatCreate(c, obj, defaults)
			//agentlogins":{"agentlogininfo":[{"agentlogindbid":"151","wrapuptime":"0"}]},"appranks":{"apprank":[]},"
			if contactdbid, exist := obj["contactdbid"]; exist {
				if contactdbid != "0" {
					log.Warn().Interface("contactdbid", contactdbid).Msg("Attached contract link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["contactdbid"] = "0"
				}
			}
			if capacityruledbid, exist := obj["capacityruledbid"]; exist {
				if capacityruledbid != "0" {
					log.Warn().Interface("capacityruledbid", capacityruledbid).Msg("Attached capacityrule link will be lost")
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
					log.Warn().Interface("agentlogins", agentlogins).Msg("Attached Agent Login link will be lost")
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
					log.Warn().Interface("appranks", appranks).Msg("Attached appranks link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["appranks"] = emptyDBIDList
				}
			}
			if sitedbid, exist := obj["sitedbid"]; exist {
				if sitedbid != "0" {
					log.Warn().Interface("sitedbid", sitedbid).Msg("Attached site link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["sitedbid"] = "0"
				}
			}

			if sitedbid, exist := obj["sitedbid"]; exist {
				if sitedbid != "0" {
					log.Warn().Interface("sitedbid", sitedbid).Msg("Attached site link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["sitedbid"] = "0"
				}
			}
			if placedbid, exist := obj["placedbid"]; exist {
				if placedbid != "0" {
					log.Warn().Interface("placedbid", placedbid).Msg("Attached place link will be lost")
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
					log.Warn().Interface("skilllevels", skilllevels).Msg("Attached skilllevels link will be lost")
					//TODO search
					obj["skilllevels"] = emptyDBIDList
				}
			}
			delete(obj, "password") //Clear Password
			log.Warn().Msg("Possible attached Agent password will be lost")
			//TODO password
			return obj
		},
		FormatUpdate: func(c *client.Client, src, obj map[string]interface{}, defaults map[string]string) map[string]interface{} {
			log.Debug().Interface("src", src).Interface("obj", obj).Msg("CfgPerson.FormatUpdate")
			obj = LoaderList["default"].FormatCreate(c, obj, defaults)
			//TODO reuse by default value of src
			//agentlogins":{"agentlogininfo":[{"agentlogindbid":"151","wrapuptime":"0"}]},"appranks":{"apprank":[]},"
			if contactdbid, exist := obj["contactdbid"]; exist {
				if contactdbid != "0" {
					log.Warn().Interface("contactdbid", contactdbid).Msg("Attached contract link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["contactdbid"] = "0"
				}
			}
			if capacityruledbid, exist := obj["capacityruledbid"]; exist {
				if capacityruledbid != "0" {
					log.Warn().Interface("capacityruledbid", capacityruledbid).Msg("Attached capacityrule link will be lost")
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
					log.Warn().Interface("agentlogins", agentlogins).Msg("Attached Agent Login link will be lost")
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
					log.Warn().Interface("appranks", appranks).Msg("Attached appranks link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["appranks"] = emptyDBIDList
				}
			}
			if sitedbid, exist := obj["sitedbid"]; exist {
				if sitedbid != "0" {
					log.Warn().Interface("sitedbid", sitedbid).Msg("Attached site link will be lost")
					//TODO search
					//obj["contactdbid"] = searchFor(c, "CfgScript", contract.(string))
					obj["sitedbid"] = "0"
				}
			}
			if placedbid, exist := obj["placedbid"]; exist {
				if placedbid != "0" {
					log.Warn().Interface("placedbid", placedbid).Msg("Attached place link will be lost")
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
					log.Warn().Interface("skilllevels", skilllevels).Msg("Attached skilllevels link will be lost")
					//TODO search
					obj["skilllevels"] = emptyDBIDList
				}
			}

			delete(obj, "password") //Clear Password
			log.Warn().Msg("Possible attached Agent password will be lost")
			//TODO password
			return obj

		},
	}
}
