/*
 * @Author: SpenserCai
 * @Date: 2023-10-04 20:26:32
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-09 15:26:53
 * @Description: file content
 */
package business

import (
	"encoding/json"
	"strings"

	"github.com/SpenserCai/sd-webui-discord/api/gen/models"
	ServiceOperations "github.com/SpenserCai/sd-webui-discord/api/gen/restapi/operations/user"
	"github.com/SpenserCai/sd-webui-discord/global"
	DbotUser "github.com/SpenserCai/sd-webui-discord/user"
	"github.com/SpenserCai/sd-webui-go/intersvc"
	"github.com/go-openapi/runtime/middleware"
)

func (b BusinessBase) SetUserHistoryHandler() {
	global.ApiService.UserUserHistoryHandler = ServiceOperations.UserHistoryHandlerFunc(func(params ServiceOperations.UserHistoryParams, principal interface{}) middleware.Responder {
		history, total, err := global.UserCenterSvc.GetUserHistoryList(principal.(DbotUser.UserInfo).Id, params.Body.Query.Command, int(params.Body.PageInfo.Page), int(params.Body.PageInfo.PageSize))
		if err != nil {
			return ServiceOperations.NewUserHistoryOK().WithPayload(&models.HistoryList{
				Code:    -1,
				Message: err.Error(),
			})
		}
		userInfo, err := global.UserCenterSvc.GetUserInfo(principal.(DbotUser.UserInfo).Id)
		if err != nil {
			return ServiceOperations.NewUserHistoryOK().WithPayload(&models.HistoryList{
				Code:    -1,
				Message: err.Error(),
			})
		}
		historyRes := make([]*models.HistoryItem, 0)
		for _, item := range history {
			historyRes = append(historyRes, &models.HistoryItem{
				ID:      item.MessageID,
				Command: item.CommandName,
				Options: func() interface{} {
					option := &intersvc.SdapiV1Txt2imgRequest{}
					json.Unmarshal([]byte(item.OptionJson), &option)
					return option
				}(),
				Created:  item.Created,
				UserID:   item.UserID,
				UserName: userInfo.Name,
				UserAvatar: func() string {
					if userInfo.Avatar == "" {
						return "https://cdn.discordapp.com/embed/avatars/0.png"
					}
					return userInfo.Avatar
				}(),
				Images: strings.Split(item.Images, ","),
			})
		}
		return ServiceOperations.NewUserHistoryOK().WithPayload(&models.HistoryList{
			Code:    0,
			Message: "success",
			Data: &models.HistoryListData{
				History: historyRes,
				PageInfo: &models.PageInfoResponse{
					Page:     params.Body.PageInfo.Page,
					PageSize: params.Body.PageInfo.PageSize,
					Total:    int32(total),
				},
			},
		})
	})
}
