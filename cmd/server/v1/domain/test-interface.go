package serverDomain

import (
	serverConsts "github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/kataras/iris/v12"
)

type DiagnoseInterface struct {
	Id int64 `json:"id"`

	Title  string                             `json:"title"`
	Desc   string                             `json:"desc"`
	Type   serverConsts.DiagnoseInterfaceType `json:"type"`
	IsLeaf bool                               `json:"isLeaf"`

	DebugInterfaceId uint  `json:"debugInterfaceId"`
	ParentId         int64 `json:"parentId"`
	ProjectId        uint  `json:"projectId"`
	ServeId          uint  `json:"serveId"`
	UseID            uint  `json:"useId"`

	Ordr     int                  `json:"ordr"`
	Children []*DiagnoseInterface `json:"children"`
	Slots    iris.Map             `json:"slots"`
}

type DiagnoseInterfaceLoadReq struct {
	ServeId   int `json:"serveId"`
	ProjectId int `json:"projectId"`
}

type DiagnoseInterfaceSaveReq struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Mode      string `json:"mode"`
	ParentId  uint   `json:"parentId"`
	ServeId   uint   `json:"serveId"`
	ProjectId uint   `json:"projectId"`

	Type serverConsts.DiagnoseInterfaceType `json:"type"`
}

type DiagnoseInterfaceReq struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Parent uint
}

type DiagnoseInterfaceMoveReq struct {
	DragKey int                  `json:"dragKey"`
	DropKey int                  `json:"dropKey"`
	DropPos serverConsts.DropPos `json:"dropPos"`
}

type DiagnoseInterfaceImportReq struct {
	InterfaceIds []int `json:"interfaceIds"`
	TargetId     uint  `json:"targetId"`
	CreateBy     uint  `json:"createBy"`
}
