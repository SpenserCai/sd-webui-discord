/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 22:27:15
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-23 17:18:24
 * @Description: file content
 */
package slash_handler

import (
	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/SpenserCai/sd-webui-discord/utils"

	"github.com/SpenserCai/sd-webui-go/intersvc"
	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) DeoldifyOptions() *discordgo.ApplicationCommand {
	renderFactorMin := 1.0
	renderFactorMax := 50.0
	return &discordgo.ApplicationCommand{
		Name:        "deoldify",
		Description: "Deoldify a image",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "image",
				Description: "The image",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "render_factor",
				Description: "The render factor of the image",
				Required:    false,
				MinValue:    &renderFactorMin,
				MaxValue:    float64(renderFactorMax),
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "artistic",
				Description: "Whether to use artistic mode",
				Required:    false,
			},
		},
	}
}

func (shdl SlashHandler) DeoldifySetOptions(cmd discordgo.ApplicationCommandInteractionData, opt *intersvc.DeoldifyImageRequest) {
	for _, v := range cmd.Options {
		switch v.Name {
		case "image":
			opt.InputImage = cmd.Resolved.Attachments[v.Value.(string)].URL
		case "render_factor":
			opt.RenderFactor = func() *int64 { v := v.IntValue(); return &v }()
		case "artistic":
			opt.Artistic = func() *bool { v := v.BoolValue(); return &v }()
		}
	}
	// default value
	if opt.RenderFactor == nil {
		v := int64(35)
		opt.RenderFactor = &v
	}
}

func (shdl SlashHandler) DeoldifyAction(s *discordgo.Session, i *discordgo.InteractionCreate, opt *intersvc.DeoldifyImageRequest, node *cluster.ClusterNode) {
	deoldify := &intersvc.DeoldifyImage{RequestItem: opt}
	deoldify.Action(node.StableClient)
	if deoldify.Error != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func() *string { v := deoldify.Error.Error(); return &v }(),
		})
	} else {
		image, err := utils.GetImageReaderByBase64(deoldify.GetResponse().Image)
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: func() *string { v := err.Error(); return &v }(),
			})
		} else {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: func() *string { v := "Success"; return &v }(),
				Files: []*discordgo.File{
					{
						Name:        "deoldify.png",
						ContentType: "image/png",
						Reader:      image,
					},
				},
			})
		}
	}
}

func (shdl SlashHandler) DeoldifyCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	option := &intersvc.DeoldifyImageRequest{}
	shdl.RespondStateMessage("Running", s, i)
	node := global.ClusterManager.GetNodeAuto()
	action := func() (map[string]interface{}, error) {
		shdl.DeoldifySetOptions(i.ApplicationCommandData(), option)
		shdl.DeoldifyAction(s, i, option, node)
		return nil, nil
	}
	callback := func() {}
	node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
}
