// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package format

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/sapk/go-genesys/api/object"
	"github.com/sirupsen/logrus"
)

//TODO grpah macro between host ?

func formatHost(name, dbid string, links *string, linksCount *int, data map[string][]interface{}) string {
	//TODO create bject for each port as it will aggregate some links and be more readable
	ret := "    subgraph " + name + "\n"
	for _, a := range data["CfgApplication"] {
		var app object.CfgApplication
		err := mapstructure.Decode(a, &app)
		if err != nil {
			logrus.Warnf("Fail to convert to CfgApplication")
			continue
		}
		if dbid == app.Hostdbid {
			ret += fmt.Sprintf("    %s(%s)\n", app.Dbid, app.Name)
			//TODO app.Isserver
			//TODO "state":"CFGEnabled",
			//TODO "redundancytype":"CFGHTWarmStanby"
			//TODO "isprimary":"CFGTrue","isserver":"CFGTrue"
			//TODO "backupserverdbid":"105"
			if app.State != "CFGEnabled" {
				ret += fmt.Sprintf("    class %s disabled\n", app.Dbid)
			} else if app.Isserver != "CFGTrue" {
				ret += fmt.Sprintf("    class %s client\n", app.Dbid)
			} else if app.Isprimary == "CFGTrue" {
				ret += fmt.Sprintf("    class %s primary\n", app.Dbid)
			}

			for _, p := range app.Portinfos.Portinfo {
				ret += fmt.Sprintf("    %s-%s[%s]\n", app.Dbid, p.Port, p.Port)
				*links += fmt.Sprintf("    %s(%s) --- %s-%s[%s]\n", app.Dbid, app.Name, app.Dbid, p.Port, p.Port)
			}
			for _, c := range app.Appservers.Conninfo {
				r := findObj("CfgApplication", c.Appserverdbid, data)
				var remote object.CfgApplication
				err := mapstructure.Decode(r, &remote)
				if err != nil {
					logrus.Warnf("Fail to convert to CfgApplication")
					continue
				}
				//TODO Color when c.Mode = CFGTMBoth
				port := c.ID
				for _, p := range remote.Portinfos.Portinfo {
					if p.ID == c.ID {
						port = p.Port
						break
					}
				}

				//TODO line backup primaire
				if remote.Hostdbid != app.Hostdbid {
					*links += fmt.Sprintf("    %s(%s) ==>  %s-%s[%s]\n", app.Dbid, app.Name, remote.Dbid, port, port)
				} else {
					*links += fmt.Sprintf("    %s(%s) -.->  %s-%s[%s]\n", app.Dbid, app.Name, remote.Dbid, port, port)
				}
				if c.Mode == "CFGTMBoth" {
					*links += fmt.Sprintf("    linkStyle %d stroke:red,stroke-width:4px;\n", *linksCount)
				}
				*linksCount++
			}
		}
	}
	ret += "    end\n"
	return ret
}
func GenerateMermaidGraph(data map[string][]interface{}) string {
	//TODO add link to https://mermaidjs.github.io/mermaid-live-editor preview
	//TODO better performance
	ret := "# Graphique Mermaid\n\n"
	links := "\n"
	linksCount := 0
	ret += "graph TB\n"
	ret += "classDef primary fill:#9f6,stroke:#333,stroke-width:2px;\n"
	ret += "classDef disabled fill:#CCC,stroke:#333,stroke-width:2px;\n"
	ret += "classDef client fill:#E6E6E6,stroke:#42992C,stroke-width:2px,stroke-dasharray:12px;\n"
	ret += "\n"
	//TODO display app witouh host
	for _, h := range data["CfgHost"] {
		host := h.(map[string]interface{})
		ret += formatHost(host["name"].(string), host["dbid"].(string), &links, &linksCount, data)
	}
	//Repass apps for apps without host
	ret += formatHost("Non définit", "0", &links, &linksCount, data)
	ret += formatHost("Non applicable", "", &links, &linksCount, data)
	return ret + links
}
