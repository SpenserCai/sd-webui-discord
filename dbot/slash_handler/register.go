/*
 * @Author: SpenserCai
 * @Date: 2023-08-31 13:51:37
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-23 17:13:01
 * @Description: file content
 */
package slash_handler

import (
	"log"

	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/global"

	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) RegisterOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "register",
		Description: "Support custom configuration after registration is completed",
		Options:     []*discordgo.ApplicationCommandOption{},
	}
}

func (shdl SlashHandler) RegisterAction(s *discordgo.Session, i *discordgo.InteractionCreate, node *cluster.ClusterNode) {
	reg_msg, err := global.UserCenterSvc.RegisterUser(shdl.ConvertInteractionToUserInfo(i))
	if err != nil {
		log.Println(err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func() *string { v := "REGISTER ERROR"; return &v }(),
		})
	} else {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func() *string { v := reg_msg; return &v }(),
		})
	}
}

func (shdl SlashHandler) RegisterCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	shdl.RespondStateMessage("Registering", s, i)
	node := global.ClusterManager.GetNodeAuto()
	action := func() (map[string]interface{}, error) {
		shdl.RegisterAction(s, i, node)
		return nil, nil
	}
	callback := func() {}
	node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
}
