/*
 * @Author: SpenserCai
 * @Date: 2023-09-29 21:25:55
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-30 00:40:48
 * @Description: file content
 */
package business

import (
	ServiceOperations "github.com/SpenserCai/sd-webui-discord/api/gen/restapi/operations/user"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/go-openapi/runtime/middleware"
	"golang.org/x/oauth2"
)

func (b BusinessBase) GetDiscordAuthUrl() string {
	var oauthConfig = oauth2.Config{
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
		Scopes:       []string{"identify", "email", "guilds.members.read", "guilds.join", "role_connections.write"},
		ClientID:     global.Config.Discord.AppId,
		ClientSecret: global.Config.Discord.ClientSecret,
		RedirectURL:  global.Config.Discord.OAuth2RedirectUrl + "/api/auth",
	}
	return oauthConfig.AuthCodeURL("random-state")
}

func (b BusinessBase) SetLoginHandler() {
	global.ApiService.UserLoginHandler = ServiceOperations.LoginHandlerFunc(func(params ServiceOperations.LoginParams) middleware.Responder {
		// 转跳到discord授权页面
		return ServiceOperations.NewLoginFound().WithLocation(b.GetDiscordAuthUrl())
	})
}
