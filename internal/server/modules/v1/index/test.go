package index

import (
	"github.com/kataras/iris/v12"
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
)

type TestModule struct {
	TestCtrl *controller.TestCtrl `inject:""`
}

func NewTestModule() *TestModule {
	return &TestModule{}
}

// Party 产品
func (m *TestModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())

		index.Get("/", m.TestCtrl.Test).Name = "测试"
	}
	return module.NewModule("/test123", handler)
}
