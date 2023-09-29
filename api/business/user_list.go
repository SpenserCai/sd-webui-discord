/*
 * @Author: SpenserCai
 * @Date: 2023-09-29 21:26:59
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-29 22:02:14
 * @Description: file content
 */
package business

import (
	ServiceOperations "github.com/SpenserCai/sd-webui-discord/api/gen/restapi/operations/admin"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/go-openapi/runtime/middleware"
)

func (b *BusinessBase) SetUserListHandler() {
	global.ApiService.AdminUserListHandler = ServiceOperations.UserListHandlerFunc(func(params ServiceOperations.UserListParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation user.UserInfo has not yet been implemented")
	})
}
