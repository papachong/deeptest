package router

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/module"
	"github.com/kataras/iris/v12"
)

type ScenarioExecModule struct {
	ScenarioExecCtrl *handler.ScenarioExecCtrl `inject:""`
}

func NewScenarioExecModule() *ScenarioExecModule {
	return &ScenarioExecModule{}
}

// Party 场景
func (m *ScenarioExecModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		// index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())

		index.Get("/loadExecScenario", m.ScenarioExecCtrl.LoadExecData).Name = "加载执行场景"
		index.Get("/loadExecResult", m.ScenarioExecCtrl.LoadExecResult).Name = "加载执行结果"

		index.Post("/submitResult/{id:uint}", m.ScenarioExecCtrl.SubmitResult).Name = "提交测试结果"
	}

	return module.NewModule("/scenarios/exec/", handler)
}
