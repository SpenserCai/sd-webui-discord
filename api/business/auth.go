/*
 * @Author: SpenserCai
 * @Date: 2023-09-30 12:53:43
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-30 21:58:49
 * @Description: file content
 */
package business

import (
	ServiceOperations "github.com/SpenserCai/sd-webui-discord/api/gen/restapi/operations/user"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/bwmarrin/discordgo"
	"github.com/go-openapi/runtime/middleware"
)

func (b BusinessBase) checkRandomState(oauthState string) bool {
	for _, state := range b.GenRandomState() {
		if oauthState == state {
			return true
		}
	}
	return false
}

func (b BusinessBase) SetAuthHandler() {
	global.ApiService.UserAuthHandler = ServiceOperations.AuthHandlerFunc(func(params ServiceOperations.AuthParams) middleware.Responder {
		var oauthConfig = b.GetDiscordOauth2Config()
		oauthCode := params.Code
		oauthState := params.State
		// 验证state
		if !b.checkRandomState(oauthState) {
			return ServiceOperations.NewAuthFound().WithLocation("/error?error=auth_state_error")
		}
		// 获取token
		token, err := oauthConfig.Exchange(params.HTTPRequest.Context(), oauthCode)
		if err != nil {
			return ServiceOperations.NewAuthFound().WithLocation("/error?error=auth_error")
		}
		ts, _ := discordgo.New("Bearer " + token.AccessToken)
		user, err := ts.User("@me")
		if err != nil {
			return ServiceOperations.NewAuthFound().WithLocation("/error?error=auth_error")
		}
		// 获取用户信息
		return ServiceOperations.NewAuthFound().WithLocation("/success?id=" + user.ID + "&username=" + user.Username)
	})

}
