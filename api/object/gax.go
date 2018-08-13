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
	DBID int64  `json:"dbid"`
}

type CfgHost struct {
	*CfgObject
}

type CfgApplication struct {
	*CfgObject
	WorkDirectory string `json:"workdirectory"`
}
