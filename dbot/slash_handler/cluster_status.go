/*
 * @Author: SpenserCai
 * @Date: 2023-09-24 18:25:37
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-11-01 10:21:21
 * @Description: file content
 */
package slash_handler

import (
	"strconv"

	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) ClusterStatusOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "cluster_status",
		Description: "View cluster status",
		Options:     []*discordgo.ApplicationCommandOption{},
	}
}

func (shdl SlashHandler) ClusterStatusAction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embeds := []*discordgo.MessageEmbed{}
	for _, node := range global.ClusterManager.Nodes {
		nodeEmbed := shdl.MessageEmbedTemplate()
		nodeEmbed.Title = node.Name
		nodeEmbed.Fields = []*discordgo.MessageEmbedField{
			{
				Name: "Host",
				Value: func() string {
					for _, v := range global.Config.SDWebUi.Servers {
						if v.Name == node.Name {
							return v.Host
						}
					}
					return ""
				}(),
				Inline: false,
			},
			{
				Name:   "MaxConcurrent",
				Value:  strconv.Itoa(node.ActionQueue.MaxConcurrent),
				Inline: false,
			},
			{
				Name:   "Running",
				Value:  strconv.Itoa(node.ActionQueue.CurrentConcurrent),
				Inline: true,
			},
			{
				Name:   "Pending",
				Value:  strconv.Itoa(len(node.ActionQueue.TaskList)),
				Inline: true,
			},
		}
		embeds = append(embeds, nodeEmbed)
	}
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &embeds,
	})
}

func (shdl SlashHandler) ClusterStatusCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	shdl.RespondStateMessageWithFlag("Running", s, i, discordgo.MessageFlagsEphemeral)
	shdl.ClusterStatusAction(s, i)
}
