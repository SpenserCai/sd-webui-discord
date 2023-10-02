/*
 * @Author: SpenserCai
 * @Date: 2023-09-29 19:24:52
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-02 21:08:46
 * @Description: file content
 */
package business

import (
	apiMiddleware "github.com/SpenserCai/sd-webui-discord/api/middleware"
	"github.com/SpenserCai/sd-webui-discord/global"
	DbotUser "github.com/SpenserCai/sd-webui-discord/user"
	"golang.org/x/oauth2"
)

type BusinessBase struct{}

func (b BusinessBase) GetDiscordOauth2Config() oauth2.Config {
	return oauth2.Config{
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
		Scopes:       []string{"identify", "email", "guilds.members.read", "guilds.join", "role_connections.write"},
		ClientID:     global.Config.Discord.AppId,
		ClientSecret: global.Config.Discord.ClientSecret,
		RedirectURL:  global.Config.Discord.OAuth2RedirectUrl + "/api/auth",
	}
}

func (b BusinessBase) GetUserInfo(token string) (*DbotUser.UserInfo, error) {
	jwt, err := apiMiddleware.DecodeJwt(token)
	if err != nil {
		return &DbotUser.UserInfo{}, err
	}
	return global.UserCenterSvc.GetUserInfo(jwt["id"].(string))

}
