package domain

import "github.com/aaronchen2k/deeptest/internal/pkg/consts"

type GlobalVar struct {
	VarId       uint   `gorm:"-" json:"varId"`
	Name        string `json:"name"`
	RightValue  string `gorm:"type:text" json:"rightValue"`
	LocalValue  string `gorm:"type:text" json:"localValue"`
	RemoteValue string `gorm:"type:text" json:"remoteValue"`
}
type GlobalParam struct {
	Name         string           `json:"name"`
	Type         consts.ParamType `json:"type"`
	In           consts.ParamIn   `json:"in"`
	Required     bool             `json:"Required"`
	DefaultValue string           `gorm:"type:text" json:"defaultValue"`
}

type InterfaceToEnvMap map[uint]uint        // interfaceId -> envId
type EnvToVariables map[uint][]GlobalVar    // envId -> vars
type Datapools map[string][]VarKeyValuePair // datapoolName -> array of map<colName, colValue>

type VarKeyValuePair map[string]interface{}
