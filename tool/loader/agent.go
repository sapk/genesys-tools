// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package loader

import (
	"github.com/sapk/go-genesys/api/client"
	"github.com/sirupsen/logrus"
)

func init() {
	LoaderList["CfgAgentLogin"] = Loader{
		FormatCreate: func(c *client.Client, obj map[string]interface{}) map[string]interface{} {
			logrus.WithFields(logrus.Fields{
				"in": obj,
			}).Debugf("CfgAgentLogin.FormatCreate")
			obj = LoaderList["default"].FormatCreate(c, obj)
			if sw, exist := obj["switchdbid"]; exist {
				obj["switchdbid"] = searchFor(c, "CfgSwitch", sw.(string))
			}
			return obj
		},
		FormatUpdate: func(c *client.Client, src, obj map[string]interface{}) map[string]interface{} {
			logrus.WithFields(logrus.Fields{
				"in": obj,
			}).Debugf("CfgAgentLogin.FormatUpdate")
			obj = LoaderList["default"].FormatUpdate(c, src, obj)
			if sw, exist := obj["switchdbid"]; exist {
				obj["switchdbid"] = searchFor(c, "CfgSwitch", sw.(string))
			}
			return obj
		},
	}
}
