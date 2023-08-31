/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 11:05:26
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-31 11:35:34
 * @Description: file content
 */
package global

import (
	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/config"
	"github.com/SpenserCai/sd-webui-discord/user"
	"github.com/bwmarrin/discordgo"
)

var (
	Config         *config.Config
	ClusterManager *cluster.ClusterService
	LongDBotChoice map[string][]*discordgo.ApplicationCommandOptionChoice
	UserCenterSvc  *user.UserCenterService
)
