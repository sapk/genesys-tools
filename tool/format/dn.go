// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package format

import (
	//"github.com/emirpasic/gods/maps/treemap"
	//"github.com/mitchellh/mapstructure"
	"fmt"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/mitchellh/mapstructure"
	"github.com/sapk/go-genesys/api/object"
	"github.com/sirupsen/logrus"
	//"github.com/sirupsen/logrus"
)

func init() {
	FormaterList["CfgDN"] = Formater{
		func(objType object.ObjectType, obj map[string]interface{}, data map[string][]interface{}) string {
			ret := "# " + obj["number"].(string) + "\n"
			ret += "\n"

			ret += dumpAvailableInformation(obj, data) + "\n"
			//TODO details : dnaccessnumber
			ret += formatOptions(obj, data)
			ret += formatAnnexes(obj, data)
			ret += dumpBackup(obj)
			return ret
		},
	}
	FormaterList["CfgDNGroup"] = Formater{
		func(objType object.ObjectType, obj map[string]interface{}, data map[string][]interface{}) string {
			ret := "# " + obj["name"].(string) + "\n"
			ret += "\n"

			ret += dumpAvailableInformation(obj, data) + "\n"
			var dnGroup object.CfgDNGroup
			err := mapstructure.Decode(obj, &dnGroup)
			if err != nil {
				logrus.Warnf("Fail to convert to CfgDNGroup")
				return err.Error()
			}
			dns := treemap.NewWithStringComparator()
			for _, d := range dnGroup.DNS.Dninfo {
				dn := findObj("CfgDN", d.Dndbid, data)
				if dn != nil {
					dns.Put(dn["number"].(string), dn)
				} else {
					dns.Put(d.Dndbid, nil)
				}
			}

			dnList := ""
			for _, d := range dns.Keys() {
				number := d.(string)
				val, _ := dns.Get(number)
				if val != nil {
					var dn object.CfgDN
					err := mapstructure.Decode(val, &dn)
					if err != nil {
						logrus.Warnf("Fail to convert to CfgDN")
						dnList += "  dbid:" + number + "\n"
					} else {
						dnList += "  " + dn.Number + " (" + dn.Subtype + "/" + dn.Dbid + ")\n"
					}
				} else {
					dnList += "  dbid:" + number + "\n"
				}
			}

			ret += fmt.Sprintf("## DNs (%d): \n", dns.Size())
			ret += dnList + "\n"
			ret += formatOptions(obj, data)
			ret += formatAnnexes(obj, data)
			ret += dumpBackup(obj)
			return ret
		},
	}
}
