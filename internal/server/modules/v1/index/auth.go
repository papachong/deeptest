package index

import (
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/middleware"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/controller"
	"github.com/kataras/iris/v12"
)

type AuthModule struct {
	AuthCtrl *controller.AuthCtrl `inject:""`
}

// Party 脚本
func (m *AuthModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Use(middleware.InitCheck(), middleware.JwtHandler(), middleware.OperationRecord(), middleware.Casbin())

		index.Post("/oauth2Authorization", m.AuthCtrl.OAuth2Authorization).Name = "生成OAuth认证信息"
		index.Post("/getOAuth2AccessToken", m.AuthCtrl.GetOAuth2AccessToken).Name = "调用认证服务生成访问令牌"
	}
	return module.NewModule("/auth", handler)
}
