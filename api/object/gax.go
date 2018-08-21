// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package object

//FROM https://docs.genesys.com/Documentation/PSDK/9.0.x/ConfigLayerRef/CfgObjectType
type ObjectType struct {
	ID         int
	Name       string
	Desc       string
	IsDumpable bool
}

var ObjectTypeListShort = []ObjectType{
	ObjectType{3, "CfgPerson", "Person", true},
	ObjectType{9, "CfgApplication", "Application", true},
	ObjectType{10, "CfgHost", "Host", true},
}
var ObjectTypeList = []ObjectType{
	//ObjectType{0, "CfgNoObject", "Unknown Object",true},
	ObjectType{1, "CfgSwitch", "Switch", true},
	ObjectType{2, "CfgDN", "DN", true},
	ObjectType{3, "CfgPerson", "Person", true},
	ObjectType{4, "CfgPlace", "Place", true},
	ObjectType{5, "CfgAgentGroup", "Agent Group", true},
	ObjectType{6, "CfgPlaceGroup", "Place Group", true},
	ObjectType{7, "CfgTenant", "Tenant", true},
	ObjectType{8, "CfgService", "Solution", true},
	ObjectType{9, "CfgApplication", "Application", true},
	ObjectType{10, "CfgHost", "Host", true},
	ObjectType{11, "CfgPhysicalSwitch", "Switching Office", true},
	ObjectType{12, "CfgScript", "Script", true},
	ObjectType{13, "CfgSkill", "Skill", true},
	ObjectType{14, "CfgActionCode", "Action Code", true},
	ObjectType{15, "CfgAgentLogin", "Agent Login", true},
	ObjectType{16, "CfgTransaction", "Transaction", true},
	ObjectType{17, "CfgDNGroup", "DN Group", true},
	ObjectType{18, "CfgStatDay", "Statistical Day", true},
	ObjectType{19, "CfgStatTable", "Statistical Table", true},
	ObjectType{20, "CfgAppPrototype", "Application Template", true},
	ObjectType{21, "CfgAccessGroup", "Access Group", true},
	ObjectType{22, "CfgFolder", "Folder", true},
	ObjectType{23, "CfgField", "Field", true},
	ObjectType{24, "CfgFormat", "Format", true},
	ObjectType{25, "CfgTableAccess", "Table Access", true},
	ObjectType{26, "CfgCallingList", "Calling List", true},
	ObjectType{27, "CfgCampaign", "Campaign", true},
	ObjectType{28, "CfgTreatment", "Treatment", true},
	ObjectType{29, "CfgFilter", "Filter", true},
	ObjectType{30, "CfgTimeZone", "Time Zone", true},
	ObjectType{31, "CfgVoicePrompt", "Voice Prompt", true},
	ObjectType{32, "CfgIVRPort", "IVR Port", true},
	ObjectType{33, "CfgIVR", "IVR", true},
	ObjectType{34, "CfgAlarmCondition", "Alarm Condition", true},
	ObjectType{35, "CfgEnumerator", "Business Attribute", true},
	ObjectType{36, "CfgEnumeratorValue ", "Business Attribute Value", true},
	ObjectType{37, "CfgObjectiveTable", "Objective Table", true},
	ObjectType{38, "CfgCampaignGroup", "Campaign Group", true},
	//ObjectType{39, "CfgGVPReseller", "GVP Reseller",true},
	//ObjectType{40, "CfgGVPCustomer", "GVP Customer",true},
	ObjectType{41, "CfgGVPIVRProfile", "GVP IVRProfile", true},
	//ObjectType{42, "CfgScheduledTask ", "Scheduled Task",true},
	ObjectType{43, "CfgRole", "Role", true},
	//	ObjectType{44, "CfgPersonLastLogin", "PersonLastLogin",true},
	//	ObjectType{45, "CfgMaxObjectType", "Shortcut",true},
}

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

/*
type CfgDN struct {
	*CfgObject
	//TODO
}

type CfgSwitch struct {
	*CfgObject
	//TODO
}

type CfgPlace struct {
	*CfgObject
	//TODO
}
*/
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
type Userproperties struct {
	Property []struct {
		Section string `json:"section"`
		Value   string `json:"value"`
		Key     string `json:"key"`
	} `json:"property"`
}
type Options struct {
	Property []struct {
		Section string `json:"section"`
		Value   string `json:"value"`
		Key     string `json:"key"`
	} `json:"property"`
}

type CfgApplication struct {
	*CfgObject
	Appservers struct {
		Conninfo []struct {
			Mode          string `json:"mode"`
			Appserverdbid string `json:"appserverdbid"`
			Timoutlocal   string `json:"timoutlocal"`
			Longfield1    string `json:"longfield1"`
			Longfield2    string `json:"longfield2"`
			Longfield3    string `json:"longfield3"`
			Longfield4    string `json:"longfield4"`
			Timoutremote  string `json:"timoutremote"`
			ID            string `json:"id"`
		} `json:"conninfo"`
	} `json:"appservers"`
	Autorestart          string         `json:"autorestart"`
	Userproperties       Userproperties `json:"userproperties"`
	Timeout              string         `json:"timeout"`
	Commandline          string         `json:"commandline"`
	Folderid             string         `json:"folderid"`
	Commandlinearguments string         `json:"commandlinearguments"`
	Subtype              string         `json:"subtype"`
	Options              Options        `json:"options"`
	State                string         `json:"state"`
	Hostdbid             string         `json:"hostdbid"`
	Attempts             string         `json:"attempts"`
	Portinfos            struct {
		Portinfo []struct {
			Longfield1 string `json:"longfield1"`
			Longfield2 string `json:"longfield2"`
			Longfield3 string `json:"longfield3"`
			Port       string `json:"port"`
			Longfield4 string `json:"longfield4"`
			ID         string `json:"id"`
		} `json:"portinfo"`
	} `json:"portinfos"`
	Workdirectory string `json:"workdirectory"`
	Startuptype   string `json:"startuptype"`
	Isserver      string `json:"isserver"`
	Resources     struct {
		Resource []interface{} `json:"resource"`
	} `json:"resources"`
	Startuptimeout   string `json:"startuptimeout"`
	Backupserverdbid string `json:"backupserverdbid"`
	Version          string `json:"version"`
	Isprimary        string `json:"isprimary"`
	Redundancytype   string `json:"redundancytype"`
	Shutdowntimeout  string `json:"shutdowntimeout"`
	Componenttype    string `json:"componenttype"`
	Appprototypedbid string `json:"appprototypedbid"`
	Port             string `json:"port"`
}

/*
type CfgObjectList []CfgObject

func (l CfgObjectList) Len() int      { return len(l) }
func (l CfgObjectList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l CfgObjectList) Less(i, j int) bool {
	li := l[i].Name
	lj := l[j].Name
	//log.Debugf("Comparing %s < %s", ai, aj, ai < aj)
	return li < lj
}
*/
/*
type CfgApplicationList []CfgApplication

func (l CfgApplicationList) Len() int      { return len(l) }
func (l CfgApplicationList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l CfgApplicationList) Less(i, j int) bool {
	li := l[i].Name
	lj := l[j].Name
	//log.Debugf("Comparing %s < %s", ai, aj, ai < aj)
	return li < lj
}
*/
/*
type CfgHostList []CfgHost

/*
func (l CfgHostList) Len() int      { return len(l) }
func (l CfgHostList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l CfgHostList) Less(i, j int) bool {
	li := l[i].Name
	lj := l[j].Name
	//log.Debugf("Comparing %s < %s", ai, aj, ai < aj)
	return li < lj
}
*/
