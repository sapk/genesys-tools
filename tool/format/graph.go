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
					*links += fmt.Sprintf("    %s(%s) ==>|%s| %s(%s)\n", app.Dbid, app.Name, port, remote.Dbid, remote.Name)
				} else {
					*links += fmt.Sprintf("    %s(%s) -.->|%s| %s(%s)\n", app.Dbid, app.Name, port, remote.Dbid, remote.Name)
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
func formatDotHost(name, dbid string, links *string, linksCount *int, data map[string][]interface{}) string {

	ret := "  subgraph cluster_" + dbid + " {\n"
	ret += "    label = \"" + name + "\";\n"
	//ret += "    style=filled;\n"
	//ret += "    color=grey95;\n"
	//ret += "    color=grey;\n"
	ret += "    style=rounded;"
	ret += "    fontsize=32;"
	ret += "    bgcolor=\"grey95\";"
	ret += "    node [fontsize=16,style=filled,shape=record];\n"

	appTypes := make(map[string][]*object.CfgApplication)
	for _, a := range data["CfgApplication"] {
		var app object.CfgApplication
		err := mapstructure.Decode(a, &app)
		if err != nil {
			logrus.Warnf("Fail to convert to CfgApplication")
			continue
		}
		if dbid == app.Hostdbid {
			appTypes[app.Subtype] = append(appTypes[app.Subtype], &app)

			if app.State != "CFGEnabled" {
				ret += fmt.Sprintf("    \"%s\" [color=gray42,label = \"<name> %s", app.Dbid, app.Name)
			} else if app.Isserver != "CFGTrue" {
				ret += fmt.Sprintf("    \"%s\" [color=cyan3,label = \"<name> %s", app.Dbid, app.Name)
			} else if app.Isprimary == "CFGTrue" {
				ret += fmt.Sprintf("    \"%s\" [color=palegreen,label = \"<name> %s", app.Dbid, app.Name)
			} else {
				ret += fmt.Sprintf("    \"%s\" [color=grey,label = \"<name> %s", app.Dbid, app.Name)
			}

			for _, p := range app.Portinfos.Portinfo {
				ret += fmt.Sprintf("| <%s> %s", p.ID, p.Port)
			}
			ret += "\"];\n"

			for _, c := range app.Appservers.Conninfo {
				r := findObj("CfgApplication", c.Appserverdbid, data)
				var remote object.CfgApplication
				err := mapstructure.Decode(r, &remote)
				if err != nil {
					logrus.Warnf("Fail to convert to CfgApplication")
					continue
				}

				//TODO line backup primaire
				if c.Mode == "CFGTMBoth" {
					if remote.Hostdbid != app.Hostdbid {
						*links += fmt.Sprintf("\"%s\":name -> \"%s\":%s [color=green3,id = %d];\n", app.Dbid, remote.Dbid, c.ID, *linksCount)
					} else {
						*links += fmt.Sprintf("\"%s\":name -> \"%s\":%s [style=dashed,color=green3,id = %d];\n", app.Dbid, remote.Dbid, c.ID, *linksCount)
					}
				} else {
					if remote.Hostdbid != app.Hostdbid {
						*links += fmt.Sprintf("\"%s\":name -> \"%s\":%s [id = %d];\n", app.Dbid, remote.Dbid, c.ID, *linksCount)
					} else {
						*links += fmt.Sprintf("\"%s\":name -> \"%s\":%s [style=dotted,id = %d];\n", app.Dbid, remote.Dbid, c.ID, *linksCount)
					}
				}
				*linksCount++
			}
		}
	}
	//{rank = same; B; D; Y;}
	for _, appT := range appTypes {
		ret += "{rank = same; "
		for _, app := range appT {
			ret += app.Dbid + "; "
		}
		ret += "}\n"
	}

	return ret + "  }\n"
}

func GenerateDotGraphByApp(data map[string][]interface{}) string {
	links := "\n"
	linksCount := 0
	ret := "digraph g {\n"
	ret += " graph [rankdir = \"LR\"];\n"
	for _, h := range data["CfgHost"] {
		host := h.(map[string]interface{})
		ret += formatDotHost(host["name"].(string), host["dbid"].(string), &links, &linksCount, data)
	}
	//Repass apps for apps without host
	ret += formatDotHost("Non définit", "0", &links, &linksCount, data)
	ret += formatDotHost("Non applicable", "", &links, &linksCount, data)

	//ORder outsite host to override blocks
	appTypes := make(map[string][]*object.CfgApplication)
	for _, a := range data["CfgApplication"] {
		var app object.CfgApplication
		err := mapstructure.Decode(a, &app)
		if err != nil {
			logrus.Warnf("Fail to convert to CfgApplication")
			continue
		}
		appTypes[app.Subtype] = append(appTypes[app.Subtype], &app)
	}
	//{rank = same; B; D; Y;}
	for _, appT := range appTypes {
		ret += "{rank = same; "
		for _, app := range appT {
			ret += app.Dbid + "; "
		}
		ret += "}\n"
	}

	return ret + links + "}\n"
}
func GenerateDotGraphByHost(data map[string][]interface{}) string {
	links := "\n"
	linksCount := 0
	ret := "digraph g {\n"
	ret += " graph [rankdir = \"LR\"];\n"
	for _, h := range data["CfgHost"] {
		host := h.(map[string]interface{})
		ret += formatDotHost(host["name"].(string), host["dbid"].(string), &links, &linksCount, data)
	}
	//Repass apps for apps without host
	ret += formatDotHost("Non définit", "0", &links, &linksCount, data)
	ret += formatDotHost("Non applicable", "", &links, &linksCount, data)
	//TODO order
	/*
		appTypes := make(map[string][]*object.CfgApplication)
		for _, a := range data["CfgApplication"] {
			var app object.CfgApplication
			err := mapstructure.Decode(a, &app)
			if err != nil {
				logrus.Warnf("Fail to convert to CfgApplication")
				continue
			}
			appTypes[app.Subtype] = append(appTypes[app.Subtype], &app)
		}
		//{rank = same; B; D; Y;}
		for _, appT := range appTypes {
			ret += "{rank = same; "
			for _, app := range appT {
				ret += app.Dbid + "; "
			}
			ret += "}\n"
		}
	*/
	return ret + links + "}\n"
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
