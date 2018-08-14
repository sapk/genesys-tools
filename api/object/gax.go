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
	Name string `json:"name"`
	DBID string `json:"dbid"`
}

type CfgHost struct {
	*CfgObject
}

type CfgApplication struct {
	*CfgObject
	WorkDirectory string `json:"workdirectory"`
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
