/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 11:05:26
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-17 11:24:15
 * @Description: file content
 */
package global

import (
	"sd-webui-discord/cluster"
	"sd-webui-discord/config"
)

var (
	Config         *config.Config
	ClusterManager *cluster.ClusterService
)
