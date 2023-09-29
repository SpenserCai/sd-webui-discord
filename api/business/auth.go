/*
 * @Author: SpenserCai
 * @Date: 2023-09-29 21:25:55
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-29 21:43:54
 * @Description: file content
 */
package business

import (
	ServiceOperations "github.com/SpenserCai/sd-webui-discord/api/gen/restapi/operations/user"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/go-openapi/runtime/middleware"
)

func (b *BusinessBase) SetAuthHandler() {
	global.ApiService.UserAuthHandler = ServiceOperations.AuthHandlerFunc(func(params ServiceOperations.AuthParams) middleware.Responder {
		return middleware.NotImplemented("operation user.Auth has not yet been implemented")
	})
}
