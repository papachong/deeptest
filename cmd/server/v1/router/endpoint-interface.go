package router

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/kataras/iris/v12"
)

type EndpointInterfaceModule struct {
	EndpointInterfaceCtrl *handler.EndpointInterfaceCtrl `inject:""`
}

func NewEndpointInterfaceModule() *EndpointInterfaceModule {
	return &EndpointInterfaceModule{}
}

// Party 注册模块
func (m *EndpointInterfaceModule) Party() module.WebModule {
	handler := func(public iris.Party) {
		public.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())

		public.Post("/listForSelection", m.EndpointInterfaceCtrl.ListForSelection).Name = "接口列表"
	}
	return module.NewModule("/endpoints/interfaces", handler)
}
