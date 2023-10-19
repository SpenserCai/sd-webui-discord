/*
 * @Author: SpenserCai
 * @Date: 2023-10-19 14:11:00
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-19 14:14:32
 * @Description: file content
 */
package business

import (
	"strings"

	"github.com/SpenserCai/sd-webui-discord/api/gen/models"
	ServiceOperations "github.com/SpenserCai/sd-webui-discord/api/gen/restapi/operations/admin"
	"github.com/SpenserCai/sd-webui-discord/global"
	DbotUser "github.com/SpenserCai/sd-webui-discord/user"
	"github.com/go-openapi/runtime/middleware"
)

func (b BusinessBase) SetSetUserPrivateHandler() {
	global.ApiService.AdminSetUserPrivateHandler = ServiceOperations.SetUserPrivateHandlerFunc(func(params ServiceOperations.SetUserPrivateParams, principal interface{}) middleware.Responder {
		if !strings.Contains(principal.(DbotUser.UserInfo).Roles, "admin") {
			return ServiceOperations.NewSetUserPrivateOK().WithPayload(&models.BaseResponse{
				Code:    -1,
				Message: "permission denied",
			})
		}
		err := global.UserCenterSvc.SetUserPrivate(params.Body.UserID, params.Body.IsPrivate)
		if err != nil {
			return ServiceOperations.NewSetUserPrivateOK().WithPayload(&models.BaseResponse{
				Code:    -1,
				Message: err.Error(),
			})
		}
		return ServiceOperations.NewSetUserPrivateOK().WithPayload(&models.BaseResponse{
			Code:    0,
			Message: "success",
		})
	})
}
