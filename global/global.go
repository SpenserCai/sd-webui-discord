/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 11:05:26
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-18 13:02:54
 * @Description: file content
 */
package global

import (
	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/config"
)

var (
	Config         *config.Config
	ClusterManager *cluster.ClusterService
)
