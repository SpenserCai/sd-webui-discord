/*
 * @Author: SpenserCai
 * @Date: 2023-09-29 21:26:59
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-12 13:48:44
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

func (b BusinessBase) SetUserListHandler() {
	global.ApiService.AdminUserListHandler = ServiceOperations.UserListHandlerFunc(func(params ServiceOperations.UserListParams, principal interface{}) middleware.Responder {
		if !strings.Contains(principal.(DbotUser.UserInfo).Roles, "admin") {
			return ServiceOperations.NewUserListOK().WithPayload(&models.UserList{
				Code:    -1,
				Message: "permission denied",
			})
		}
		query := map[string]interface{}{}
		if params.Body.Query.Username != "" {
			query["name"] = params.Body.Query.Username
		}
		if params.Body.Query.OnlyEnable {
			query["enable"] = params.Body.Query.OnlyEnable
		}
		if params.Body.Query.ID != "" {
			query["id"] = params.Body.Query.ID
		}
		userList, count, err := global.UserCenterSvc.SearchUserInfoList(int(params.Body.PageInfo.Page), int(params.Body.PageInfo.PageSize), query)
		if err != nil {
			return ServiceOperations.NewUserListOK().WithPayload(&models.UserList{
				Code:    -1,
				Message: err.Error(),
			})
		}
		// 获取用户列表中用户生成的图片数量
		userIdList := make([]string, 0)
		for _, item := range userList {
			userIdList = append(userIdList, item.Id)
		}
		userImageCount, err := global.UserCenterSvc.GetUsersImageTotal(userIdList)
		if err != nil {
			return ServiceOperations.NewUserListOK().WithPayload(&models.UserList{
				Code:    -1,
				Message: err.Error(),
			})
		}
		userListRes := make([]*models.UserItem, 0)
		for _, item := range userList {
			userListRes = append(userListRes, &models.UserItem{
				ID:       item.Id,
				Username: item.Name,
				Avatar: func() string {
					if item.Avatar == "" {
						return "https://cdn.discordapp.com/embed/avatars/0.png"
					}
					return item.Avatar
				}(),
				Enable:       item.Enable,
				StableConfig: item.StableConfig,
				Roles:        item.Roles,
				Created:      item.Created,
				ImageCount:   int32(userImageCount[item.Id]),
			})
		}
		return ServiceOperations.NewUserListOK().WithPayload(&models.UserList{
			Code:    0,
			Message: "success",
			Data: &models.UserListData{
				Users: userListRes,
				PageInfo: &models.PageInfoResponse{
					Page:     params.Body.PageInfo.Page,
					PageSize: params.Body.PageInfo.PageSize,
					Total:    int32(count),
				},
			},
		})
	})
}
