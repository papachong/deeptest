package domain

import (
	"encoding/json"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
)

type DebugInfo struct {
	DebugInterfaceId    uint `json:"debugInterfaceId"`
	EndpointInterfaceId uint `json:"endpointInterfaceId"` // EndpointInterface without DebugInterface init
	CaseInterfaceId     uint `json:"caseInterfaceId"`     // load by endpoint case
	ScenarioProcessorId uint `json:"scenarioProcessorId"` // used to load vars by scenario processor
	DiagnoseInterfaceId uint `json:"diagnoseInterfaceId"` // load by interface diagnose

	UsedBy consts.UsedBy `json:"usedBy"`
}

type SubmitDebugResultRequest struct {
	Request  DebugData     `json:"request"`
	Response DebugResponse `json:"response"`

	PreConditions  []InterfaceExecCondition `json:"preConditions"`
	PostConditions []InterfaceExecCondition `json:"postConditions"`
}

type DebugData struct {
	BaseRequest

	Name string `json:"name"`

	DebugInterfaceId    uint `json:"debugInterfaceId"`
	EndpointInterfaceId uint `json:"endpointInterfaceId"`
	CaseInterfaceId     uint `json:"caseInterfaceId"`
	DiagnoseInterfaceId uint `json:"diagnoseInterfaceId"`

	ScenarioProcessorId uint `json:"scenarioProcessorId"`

	UsedBy consts.UsedBy `json:"usedBy"`

	ServeId   uint `json:"serveId"`
	ServerId  uint `json:"serverId"`
	ProjectId uint `json:"projectId"`

	BaseUrl string `json:"baseUrl"`

	// used for selection and show in right environment tab
	ShareVars    []GlobalVar   `json:"shareVars"`
	EnvVars      []GlobalVar   `json:"envVars"`
	GlobalVars   []GlobalVar   `json:"globalVars"`
	GlobalParams []GlobalParam `json:"globalParams"`

	ProcessorInterfaceSrc consts.ProcessorInterfaceSrc `json:"processorInterfaceSrc"`

	IsScenarioProcessorReset bool `json:"isReset"`
}

type Condition struct {
	Name string `json:"name"`
	Desc string `json:"desc"`

	EntityType string      `json:"entityType"`
	EntityId   uint        `json:"entityId"`
	EntityData interface{} `json:"entityData"`

	Ordr     int  `json:"ordr"`
	Disabled bool `json:"disabled"`
}

type InterfaceExecCondition struct {
	Type consts.ConditionType `json:"debugData"`
	Raw  json.RawMessage      `json:"raw"`
}
