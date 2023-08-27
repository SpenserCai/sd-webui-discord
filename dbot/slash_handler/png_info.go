/*
 * @Author: SpenserCai
 * @Date: 2023-08-19 18:27:34
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-27 23:43:10
 * @Description: file content
 */
package slash_handler

import (
	"encoding/json"
	"log"

	"github.com/SpenserCai/sd-webui-discord/utils"

	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/global"

	"github.com/SpenserCai/sd-webui-go/intersvc"
	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) PngInfoOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "png_info",
		Description: "Get image info",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "image",
				Description: "The image",
				Required:    true,
			},
		},
	}
}

func (shdl SlashHandler) PngInfoSetOptions(cmd discordgo.ApplicationCommandInteractionData, opt *intersvc.SdapiV1PngInfoRequest) {
	for _, v := range cmd.Options {
		switch v.Name {
		case "image":
			opt.Image = func() *string {
				v, _ := utils.GetImageBase64(cmd.Resolved.Attachments[v.Value.(string)].URL)
				return &v
			}()
		}
	}
}

func (shdl SlashHandler) PngInfoAction(s *discordgo.Session, i *discordgo.InteractionCreate, opt *intersvc.SdapiV1PngInfoRequest, node *cluster.ClusterNode) {
	msg, err := shdl.SendStateMessage("Running", s, i)
	if err != nil {
		log.Println(err)
		return
	}
	png_info := &intersvc.SdapiV1PngInfo{RequestItem: opt}
	png_info.Action(node.StableClient)
	if png_info.Error != nil {
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: func() *string { v := png_info.Error.Error(); return &v }(),
		})
	} else {
		items, _ := json.MarshalIndent(png_info.GetResponse().Items, "", "    ")
		outString := "items:\n```json\n" + string(items) + "\n```"
		outString += "info:\n" + *png_info.GetResponse().Info
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: &outString,
		})
	}
}

func (shdl SlashHandler) PngInfoCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	option := &intersvc.SdapiV1PngInfoRequest{}
	shdl.ReportCommandInfo(s, i)
	node := global.ClusterManager.GetNodeAuto()
	action := func() (map[string]interface{}, error) {
		shdl.PngInfoSetOptions(i.ApplicationCommandData(), option)
		shdl.PngInfoAction(s, i, option, node)
		return nil, nil
	}
	callback := func() {}
	node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
}
