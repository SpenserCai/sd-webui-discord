/*
 * @Author: SpenserCai
 * @Date: 2023-10-09 20:01:58
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-14 17:18:22
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

func (b BusinessBase) SetCommunityHistoryHandler() {
	global.ApiService.UserCommunityHistoryHandler = ServiceOperations.CommunityHistoryHandlerFunc(func(params ServiceOperations.CommunityHistoryParams, principal interface{}) middleware.Responder {
		// 判断command是否为txt2img，如果不是则返回错误：暂不支持此类型
		if params.Body.Query.Command != "txt2img" {
			return ServiceOperations.NewCommunityHistoryOK().WithPayload(&models.HistoryList{
				Code:    -1,
				Message: "Not supported",
			})
		}
		history, total, err := global.UserCenterSvc.GetUserHistoryList("", params.Body.Query.Command, int(params.Body.PageInfo.Page), int(params.Body.PageInfo.PageSize))
		if err != nil {
			return ServiceOperations.NewCommunityHistoryOK().WithPayload(&models.HistoryList{
				Code:    -1,
				Message: err.Error(),
			})
		}
		userIds := make([]string, 0)
		for _, item := range history {
			userIds = append(userIds, item.UserID)
		}
		userInfos, err := global.UserCenterSvc.GetUserInfoList(userIds)
		if err != nil {
			return ServiceOperations.NewCommunityHistoryOK().WithPayload(&models.HistoryList{
				Code:    -1,
				Message: err.Error(),
			})
		}
		historyRes := make([]*models.HistoryItem, 0)
		for _, item := range history {
			userInfo := &DbotUser.UserInfo{}
			for _, v := range userInfos {
				if v.Id == item.UserID {
					userInfo = v
					break
				}
			}
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
				ImagesBlurhash: func() []string {
					blurs := make([]string, 0)
					json.Unmarshal([]byte(item.ImageBlurHashs), &blurs)
					return blurs
				}(),
			})
		}
		return ServiceOperations.NewCommunityHistoryOK().WithPayload(&models.HistoryList{
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
