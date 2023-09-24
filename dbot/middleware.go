/*
 * @Author: SpenserCai
 * @Date: 2023-09-11 13:43:11
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-24 19:11:19
 * @Description: file content
 */
package dbot

import (
	shdl "github.com/SpenserCai/sd-webui-discord/dbot/slash_handler"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/bwmarrin/discordgo"
)

// 权限验证
func (dbot *DiscordBot) CheckPermission(cmd string, s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	if !global.Config.UserCenter.Enable {
		return true
	}
	userId := shdl.SlashHandler{}.GetDiscordUserId(i)

	if global.UserCenterSvc.CheckUserPermission(userId, cmd) {
		return true
	} else {
		// 发送无权访问的消息
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You don't have permission to access this command",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return false
	}
}
