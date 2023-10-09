/*
 * @Author: SpenserCai
 * @Date: 2023-10-09 11:01:47
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-09 11:03:58
 * @Description: file content
 */
package business

import (
	"github.com/SpenserCai/sd-webui-discord/api/gen/models"
	ServiceOperations "github.com/SpenserCai/sd-webui-discord/api/gen/restapi/operations/system"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/go-openapi/runtime/middleware"
)

func (b BusinessBase) SetDiscordServerHandler() {
	global.ApiService.SystemDiscordServerHandler = ServiceOperations.DiscordServerHandlerFunc(func(params ServiceOperations.DiscordServerParams) middleware.Responder {
		return ServiceOperations.NewDiscordServerOK().WithPayload(&models.DiscordServer{
			Code:    0,
			Message: "success",
			Data: &models.DiscordServerData{
				URL: global.Config.Discord.DiscordServerUrl,
			},
		})
	})
}
