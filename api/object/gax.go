// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package object

type LoginRequest struct {
	Username            string `json:"username"`
	Password            string `json:"password"`
	IsPasswordEncrypted bool   `json:"isPasswordEncrypted"`
}
type LoginResponse struct {
	Username       string `json:"username"`
	UserType       string `json:"userType"`
	SessionTimeout int    `json:"sessionTimeout"`
	IsDefaultUser  bool   `json:"isDefaultUser"`
	WriteDefault   bool   `json:"writeDefault"`
}

type CfgObject struct {
	Dbid string `json:"dbid"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type CfgHost struct {
	*CfgObject
	Ipaddress string `json:"ipaddress"`
	Scsdbid   string `json:"scsdbid"`
	Subtype   string `json:"subtype"`
	Lcaport   string `json:"lcaport"`
	Ostype    string `json:"ostype"`
	State     string `json:"state"`
	Folderid  string `json:"folderid"`
}

type CfgApplication struct {
	*CfgObject
	Workdirectory        string `json:"workdirectory"`
	Startuptype          string `json:"startuptype"`
	Autorestart          string `json:"autorestart"`
	Isserver             string `json:"isserver"`
	Startuptimeout       string `json:"startuptimeout"`
	Backupserverdbid     string `json:"backupserverdbid"`
	Version              string `json:"version"`
	Isprimary            string `json:"isprimary"`
	Timeout              string `json:"timeout"`
	Commandline          string `json:"commandline"`
	Folderid             string `json:"folderid"`
	Redundancytype       string `json:"redundancytype"`
	Commandlinearguments string `json:"commandlinearguments"`
	Shutdowntimeout      string `json:"shutdowntimeout"`
	Componenttype        string `json:"componenttype"`
	Appprototypedbid     string `json:"appprototypedbid"`
	Subtype              string `json:"subtype"`
	Port                 string `json:"port"`
	State                string `json:"state"`
	Hostdbid             string `json:"hostdbid"`
	Attempts             string `json:"attempts"`
}

type CfgObjectList []CfgObject

func (l CfgObjectList) Len() int      { return len(l) }
func (l CfgObjectList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l CfgObjectList) Less(i, j int) bool {
	li := l[i].Name
	lj := l[j].Name
	//log.Debugf("Comparing %s < %s", ai, aj, ai < aj)
	return li < lj
}

type CfgApplicationList []CfgApplication

func (l CfgApplicationList) Len() int      { return len(l) }
func (l CfgApplicationList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l CfgApplicationList) Less(i, j int) bool {
	li := l[i].Name
	lj := l[j].Name
	//log.Debugf("Comparing %s < %s", ai, aj, ai < aj)
	return li < lj
}

type CfgHostList []CfgHost

func (l CfgHostList) Len() int      { return len(l) }
func (l CfgHostList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l CfgHostList) Less(i, j int) bool {
	li := l[i].Name
	lj := l[j].Name
	//log.Debugf("Comparing %s < %s", ai, aj, ai < aj)
	return li < lj
}
