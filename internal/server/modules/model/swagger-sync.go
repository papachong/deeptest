package model

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"time"
)

type SwaggerSync struct {
	BaseModel
	Switch     consts.SwitchStatus `json:"switch"`
	SyncType   consts.DataSyncType `json:"syncType"`
	CategoryId int                 `json:"categoryId"`
	Url        string              `json:"url"`
	Cron       string              `json:"cron"`
	ProjectId  int                 `json:"projectId" gorm:"unique"`
	ServeId    int                 `json:"ServeId"`
	ExecTime   *time.Time          `json:"execTime"`
}

func (SwaggerSync) TableName() string {
	return "biz_project_serve_swagger_sync"
}
