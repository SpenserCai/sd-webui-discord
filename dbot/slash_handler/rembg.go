/*
 * @Author: SpenserCai
 * @Date: 2023-08-18 11:11:48
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-23 17:11:39
 * @Description: file content
 */
package slash_handler

import (
	"github.com/SpenserCai/sd-webui-discord/utils"

	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/global"

	"github.com/SpenserCai/sd-webui-go/intersvc"
	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) rembgModelChoice() []*discordgo.ApplicationCommandOptionChoice {
	model_list := []string{
		"u2net",
		"u2netp",
		"u2net_human_seg",
		"u2net_cloth_seg",
		"silueta",
		"isnet-general-use",
		"isnet-anime",
	}
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	for _, model := range model_list {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  model,
			Value: model,
		})
	}
	return choices
}

func (shdl SlashHandler) RembgOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "rembg",
		Description: "Remove background from image",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "image",
				Description: "The image",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "model",
				Description: "The model to use",
				Required:    true,
				Choices:     shdl.rembgModelChoice(),
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "return_mask",
				Description: "Whether to return mask",
				Required:    false,
			},
		},
	}
}

func (shdl SlashHandler) RembgSetOptions(cmd discordgo.ApplicationCommandInteractionData, opt *intersvc.RembgRequest) {
	opt.AlphaMatting = func() *bool { v := false; return &v }()
	for _, v := range cmd.Options {
		switch v.Name {
		case "image":
			opt.InputImage, _ = utils.GetImageBase64(cmd.Resolved.Attachments[v.Value.(string)].URL)
		case "model":
			opt.Model = func() *string { v := v.StringValue(); return &v }()
		case "return_mask":
			opt.ReturnMask = func() *bool { v := v.BoolValue(); return &v }()
		}
	}
}

func (shdl SlashHandler) RembgAction(s *discordgo.Session, i *discordgo.InteractionCreate, opt *intersvc.RembgRequest, node *cluster.ClusterNode) {

	rembg := &intersvc.Rembg{RequestItem: opt}
	rembg.Action(node.StableClient)
	if rembg.Error != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func() *string { v := rembg.Error.Error(); return &v }(),
		})
	} else {
		image, err := utils.GetImageReaderByBase64(rembg.GetResponse().Image)
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: func() *string { v := err.Error(); return &v }(),
			})
		} else {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: func() *string { v := "Success"; return &v }(),
				Files: []*discordgo.File{
					{
						Name:        "rembg.png",
						ContentType: "image/png",
						Reader:      image,
					},
				},
			})
		}
	}
}

func (shdl SlashHandler) RembgCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	option := &intersvc.RembgRequest{}
	shdl.RespondStateMessage("Running", s, i)
	node := global.ClusterManager.GetNodeAuto()
	action := func() (map[string]interface{}, error) {
		shdl.RembgSetOptions(i.ApplicationCommandData(), option)
		shdl.RembgAction(s, i, option, node)
		return nil, nil
	}
	callback := func() {}
	node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
}
