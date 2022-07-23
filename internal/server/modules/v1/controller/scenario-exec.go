package controller

import (
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	"github.com/kataras/iris/v12"
)

type ScenarioExecCtrl struct {
	ScenarioExecService *service.ScenarioExecService `inject:""`
	BaseCtrl
}

// loadExecData
func (c *ScenarioExecCtrl) LoadExecData(ctx iris.Context) {
	scenarioId, err := ctx.URLParamInt("scenarioId")

	data, err := c.ScenarioExecService.Load(scenarioId)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: data, Msg: _domain.NoErr.Msg})
}
