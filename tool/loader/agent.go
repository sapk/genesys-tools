// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package loader

import (
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
}
