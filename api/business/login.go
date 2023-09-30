/*
 * @Author: SpenserCai
 * @Date: 2023-09-29 21:25:55
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-30 21:57:18
 * @Description: file content
 */
package business

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	ServiceOperations "github.com/SpenserCai/sd-webui-discord/api/gen/restapi/operations/user"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/go-openapi/runtime/middleware"
)

func (b BusinessBase) GenRandomState() []string {
	// 获取当前时间辍精确到秒和上一秒的时间辍精确到秒
	stateList := []string{
		time.Now().Format("200601021504"),
		time.Now().Add(-1 * time.Minute).Format("200601021504"),
	}
	// 循环stateList，每个拼接上gloabl.Config.Discord.ClientSecret，然后md5后取前6位
	for i, state := range stateList {
		tmpString := state + global.Config.Discord.ClientSecret
		md5Bytes := md5.Sum([]byte(tmpString))
		stateList[i] = hex.EncodeToString(md5Bytes[:16])[0:6]
	}
	return stateList

}

func (b BusinessBase) GetDiscordAuthUrl() string {
	var oauthConfig = b.GetDiscordOauth2Config()
	return oauthConfig.AuthCodeURL(b.GenRandomState()[0])
}

func (b BusinessBase) SetLoginHandler() {
	global.ApiService.UserLoginHandler = ServiceOperations.LoginHandlerFunc(func(params ServiceOperations.LoginParams) middleware.Responder {
		// 转跳到discord授权页面
		return ServiceOperations.NewLoginFound().WithLocation(b.GetDiscordAuthUrl())
	})
}
