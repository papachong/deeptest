package controller

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type AuthCtrl struct {
	AuthService *service.AuthService `inject:""`
	BaseCtrl
}

// OAuth2Authorization
func (c *AuthCtrl) OAuth2Authorization(ctx iris.Context) {
	projectId, err := ctx.URLParamInt("currProjectId")
	if projectId == 0 {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: "projectId"})
		return
	}

	req := model.InterfaceOAuth20{}
	err = ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.ParamErr.Code, Msg: err.Error()})
		return
	}

	data, err := c.AuthService.OAuth2Authorization(req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: data})
}

// GetOAuth2AccessToken
func (c *AuthCtrl) GetOAuth2AccessToken(ctx iris.Context) {
	accessTokenURL := ctx.URLParam("accessTokenURL")
	clientId := ctx.URLParam("clientId")
	clientSecret := ctx.URLParam("clientSecret")
	code := ctx.URLParam("code")

	data, err := c.AuthService.GenOAuth2AccessToken(accessTokenURL, clientId, clientSecret, code)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: data})
}
