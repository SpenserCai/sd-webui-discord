/*
 * @Author: SpenserCai
 * @Date: 2023-10-13 11:33:20
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-13 12:16:23
 * @Description: file content
 */
package slash_handler

import (
	"log"

	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/global"

	"github.com/SpenserCai/sd-webui-go/intersvc"
	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) LoraListOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "lora_list",
		Description: "Get Lora list",
		Options:     []*discordgo.ApplicationCommandOption{},
	}
}

func (shdl SlashHandler) LoraListAction(s *discordgo.Session, i *discordgo.InteractionCreate, node *cluster.ClusterNode) {
	lora_list := &intersvc.SdapiV1Loras{}
	lora_list.Action(node.StableClient)
	if lora_list.Error != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func() *string { v := lora_list.Error.Error(); return &v }(),
		})
	} else {
		context := ""
		loras := lora_list.GetResponse()
		loraListString := ""
		for _, lora := range *loras {
			if len(loraListString) > 920 {
				loraListString += "`More...`"
				break
			}
			// loraListString += "Name: " + lora.Name + "\n"
			// if lora.Metadata.SsSdModelName != "" {
			// 	loraListString += "Base Model: " + lora.Metadata.SsSdModelName + "\n"
			// }
			loraListString += "```<lora:" + lora.Name + ":1>```"
			// loraListString += "\n"
		}
		// log.Println(loraListString)
		mainEmbed := shdl.MessageEmbedTemplate()
		mainEmbed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:  "Lora List",
				Value: loraListString,
			},
		}
		allEmbeds := []*discordgo.MessageEmbed{mainEmbed}
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &context,
			Embeds:  &allEmbeds,
		})
		if err != nil {
			log.Println(err)
		}
	}

}

func (shdl SlashHandler) LoraListCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	shdl.RespondStateMessage("Running", s, i)
	node := global.ClusterManager.GetNodeAuto()
	action := func() (map[string]interface{}, error) {
		shdl.LoraListAction(s, i, node)
		return nil, nil
	}
	callback := func() {}
	node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
}
