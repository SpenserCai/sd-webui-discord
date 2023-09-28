/*
 * @Author: SpenserCai
 * @Date: 2023-09-11 13:43:11
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-28 14:02:47
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
	// 判断是否开启必须注册
	if global.Config.UserCenter.MustRegister && !global.UserCenterSvc.IsRegistered(userId) && cmd != "register" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "**USER NOT REGISTERED**\nPlease register with `/register` first",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return false
	}
	if isAdmin, _ := global.UserCenterSvc.IsAdmin(userId); isAdmin {
		return true
	} else {
		// 检查用户是否被禁用
		if isDisabled, _ := global.UserCenterSvc.IsBaned(userId); isDisabled {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "**USER BANED**",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return false
		}
	}
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
