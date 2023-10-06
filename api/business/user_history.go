/*
 * @Author: SpenserCai
 * @Date: 2023-10-04 20:26:32
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-06 17:21:04
 * @Description: file content
 */
package business

import (
	"encoding/json"
	"strings"

	"github.com/SpenserCai/sd-webui-discord/api/gen/models"
	ServiceOperations "github.com/SpenserCai/sd-webui-discord/api/gen/restapi/operations/user"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/SpenserCai/sd-webui-go/intersvc"
	"github.com/go-openapi/runtime/middleware"
)

func (b BusinessBase) SetUserHistoryHandler() {
	global.ApiService.UserUserHistoryHandler = ServiceOperations.UserHistoryHandlerFunc(func(params ServiceOperations.UserHistoryParams, principal interface{}) middleware.Responder {
		history, total, err := global.UserCenterSvc.GetUserHistoryList(principal.(string), params.Body.Query.Command, int(params.Body.PageInfo.Page), int(params.Body.PageInfo.PageSize))
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
				Created: item.Created,
				UserID:  item.UserID,
				Images:  strings.Split(item.Images, ","),
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
