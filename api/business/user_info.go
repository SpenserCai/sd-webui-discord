/*
 * @Author: SpenserCai
 * @Date: 2023-09-29 21:26:43
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-30 00:35:00
 * @Description: file content
 */
package business

import (
	ServiceOperations "github.com/SpenserCai/sd-webui-discord/api/gen/restapi/operations/user"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/go-openapi/runtime/middleware"
)

func (b BusinessBase) SetUserInfoHandler() {
	global.ApiService.UserUserInfoHandler = ServiceOperations.UserInfoHandlerFunc(func(params ServiceOperations.UserInfoParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation user.UserInfo has not yet been implemented")
	})
}
