/*
 * @Author: SpenserCai
 * @Date: 2023-09-30 12:53:43
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-02 22:32:09
 * @Description: file content
 */
package business

import (
	"log"
	"net/http"

	ServiceOperations "github.com/SpenserCai/sd-webui-discord/api/gen/restapi/operations/user"
	apiMiddleware "github.com/SpenserCai/sd-webui-discord/api/middleware"
	"github.com/SpenserCai/sd-webui-discord/global"
	DbotUser "github.com/SpenserCai/sd-webui-discord/user"
	"github.com/bwmarrin/discordgo"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

type AuthResponse struct {
	responder middleware.Responder
	token     string
}

func (ap AuthResponse) WriteResponse(rw http.ResponseWriter, p runtime.Producer) {
	cookie := http.Cookie{Name: "token", Value: ap.token, Path: "/"}
	http.SetCookie(rw, &cookie)
	ap.responder.WriteResponse(rw, p)
}

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
		rsgMsg, err := global.UserCenterSvc.RegisterUser(&DbotUser.UserInfo{
			Id:   user.ID,
			Name: user.Username,
		})
		if err != nil {
			log.Println("RegisterUser error:", err)
			return ServiceOperations.NewAuthFound().WithLocation("/error?error=login_error")
		}
		if rsgMsg != "REGISTERED USER INFO" && rsgMsg != "UPDATED USER INFO" {
			return ServiceOperations.NewAuthFound().WithLocation("/error?error=login_error|" + rsgMsg)
		}
		userInfo, err := global.UserCenterSvc.GetUserInfo(user.ID)
		if err != nil {
			log.Println("GetUserInfo error:", err)
			return ServiceOperations.NewAuthFound().WithLocation("/error?error=login_error")
		}
		// 通过用户信息构建jwt，包含用户id，用户名，角色列表，过期时间
		jwt, err := apiMiddleware.BuildJwt(*userInfo, map[string]string{
			"avatar": user.AvatarURL(""),
		})
		if err != nil {
			log.Println("BuildJwt error:", err)
			return ServiceOperations.NewAuthFound().WithLocation("/error?error=login_error")
		}
		response := ServiceOperations.NewAuthFound().WithLocation("/api/user_info")
		authResponse := AuthResponse{
			responder: response,
			token:     jwt,
		}
		return authResponse
	})

}
