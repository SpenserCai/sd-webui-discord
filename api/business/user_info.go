/*
 * @Author: SpenserCai
 * @Date: 2023-09-29 21:26:43
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-09 15:21:47
 * @Description: file content
 */
package business

import (
	"github.com/SpenserCai/sd-webui-discord/api/gen/models"
	ServiceOperations "github.com/SpenserCai/sd-webui-discord/api/gen/restapi/operations/user"
	"github.com/SpenserCai/sd-webui-discord/global"
	DbotUser "github.com/SpenserCai/sd-webui-discord/user"
	"github.com/go-openapi/runtime/middleware"
)

func (b BusinessBase) SetUserInfoHandler() {
	global.ApiService.UserUserInfoHandler = ServiceOperations.UserInfoHandlerFunc(func(params ServiceOperations.UserInfoParams, principal interface{}) middleware.Responder {
		userInfo, err := global.UserCenterSvc.GetUserInfo(principal.(DbotUser.UserInfo).Id)
		if err != nil {
			return ServiceOperations.NewUserInfoOK().WithPayload(&models.UserInfo{
				Code:    -1,
				Message: err.Error(),
			})
		}
		return ServiceOperations.NewUserInfoOK().WithPayload(&models.UserInfo{
			Code:    0,
			Message: "success",
			Data: &models.UserInfoData{
				User: &models.UserItem{
					ID:           userInfo.Id,
					Username:     userInfo.Name,
					Avatar:       userInfo.Avatar,
					Enable:       userInfo.Enable,
					StableConfig: userInfo.StableConfig,
					Roles:        userInfo.Roles,
					Created:      userInfo.Created,
				},
			},
		})
	})
}
