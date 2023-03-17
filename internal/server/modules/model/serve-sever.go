package model

type ServeServer struct {
	BaseModel
	EnvironmentId uint   `json:"environmentId"`
	ServeId       uint   `json:"serveId"`
	Url           string `json:"url"`
	Description   string `json:"description"`
	ServeName     string `gorm:"-" json:"serveName"`
}

func (ServeServer) TableName() string {
	return "biz_project_serve_server"
}
