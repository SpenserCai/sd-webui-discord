/*
 * @Author: SpenserCai
 * @Date: 2023-10-09 17:31:54
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-09 18:08:41
 * @Description: file content
 */
package business

import (
	"strings"

	"github.com/SpenserCai/sd-webui-discord/api/gen/models"
	ServiceOperations "github.com/SpenserCai/sd-webui-discord/api/gen/restapi/operations/system"
	"github.com/SpenserCai/sd-webui-discord/global"
	DbotUser "github.com/SpenserCai/sd-webui-discord/user"
	"github.com/go-openapi/runtime/middleware"
)

func (b BusinessBase) SetClusterHandler() {
	global.ApiService.SystemClusterHandler = ServiceOperations.ClusterHandlerFunc(func(params ServiceOperations.ClusterParams, principal interface{}) middleware.Responder {
		// 判断是否为管理员principal.(DbotUser.UserInfo).Roles字符串中是否包含"admin"
		if !strings.Contains(principal.(DbotUser.UserInfo).Roles, "admin") {
			return ServiceOperations.NewClusterOK().WithPayload(&models.ClusterInfo{
				Code:    -1,
				Message: "permission denied",
			})
		}
		nodes := []*models.NodeItem{}
		for _, node := range global.ClusterManager.Nodes {
			nodes = append(nodes, &models.NodeItem{
				Host: func() string {
					for _, v := range global.Config.SDWebUi.Servers {
						if v.Name == node.Name {
							return v.Host
						}
					}
					return ""
				}(),
				Name:          node.Name,
				MaxConcurrent: int32(node.ActionQueue.MaxConcurrent),
				Running:       int32(node.ActionQueue.CurrentConcurrent),
				Pending:       int32(len(node.ActionQueue.TaskList)),
			})
		}
		return ServiceOperations.NewClusterOK().WithPayload(&models.ClusterInfo{
			Code:    0,
			Message: "success",
			Data: &models.ClusterInfoData{
				Cluster: nodes,
			},
		})
	})
}
