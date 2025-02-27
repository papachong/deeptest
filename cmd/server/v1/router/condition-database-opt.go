package router

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/kataras/iris/v12"
)

type DatabaseOptModule struct {
	DatabaseOptCtrl *handler.DatabaseOptCtrl `inject:""`
}

// Party 检查点
func (m *DatabaseOptModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin(), middleware.ProjectPerm())

		index.Get("/{id:uint}", m.DatabaseOptCtrl.Get).Name = "数据库操作后置条件详情"
		index.Put("/", m.DatabaseOptCtrl.Update).Name = "更新数据库操作后置条件"
		index.Delete("/{id:uint}", m.DatabaseOptCtrl.Delete).Name = "删除数据库操作后置条件"
	}

	return module.NewModule("/databaseOpts", handler)
}
