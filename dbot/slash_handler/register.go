/*
 * @Author: SpenserCai
 * @Date: 2023-08-31 13:51:37
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-31 14:06:50
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
	msg, err := shdl.SendStateMessage("Registering", s, i)
	if err != nil {
		log.Println(err)
		return
	}
	reg_msg, err := global.UserCenterSvc.RegisterUser(shdl.ConvertInteractionToUserInfo(i))
	if err != nil {
		log.Println(err)
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: func() *string { v := "REGISTER ERROR"; return &v }(),
		})
	} else {
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: func() *string { v := reg_msg; return &v }(),
		})
	}
}

func (shdl SlashHandler) RegisterCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	shdl.ReportCommandInfo(s, i)
	node := global.ClusterManager.GetNodeAuto()
	action := func() (map[string]interface{}, error) {
		shdl.RegisterAction(s, i, node)
		return nil, nil
	}
	callback := func() {}
	node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
}
