/*
 * @Author: SpenserCai
 * @Date: 2023-08-17 09:52:25
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-18 01:48:44
 * @Description: file content
 */
package slash_handler

import (
	"fmt"
	"sd-webui-discord/utils"

	"github.com/bwmarrin/discordgo"
)

type SlashHandler struct{}

func (shdl SlashHandler) GetCommandStr(i *discordgo.Interaction) string {
	cmd := ""
	// 把命令的名字和参数拼接起来
	cmd += fmt.Sprintf("Command: %s\n", i.ApplicationCommandData().Name)
	for _, v := range i.ApplicationCommandData().Options {
		cmd += fmt.Sprintf("%v: %v\n", utils.FormatCommand(v.Name), v.Value)
	}
	return cmd
}

func (shdl SlashHandler) ReportCommandInfo(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: shdl.GetCommandStr(i.Interaction),
		},
	})
}

func (shdl SlashHandler) GenerateTaskID(i *discordgo.InteractionCreate) string {
	return fmt.Sprintf("%s_%s_%s", i.Interaction.ID, i.Interaction.Member.User.ID, i.Interaction.Member.User.Username)
}

func (shdl SlashHandler) SendStateMessage(state string, s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.Message, error) {
	msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("%s...", state),
		Files:   []*discordgo.File{},
	})
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong",
		})
		return nil, err
	}
	return msg, nil
}
