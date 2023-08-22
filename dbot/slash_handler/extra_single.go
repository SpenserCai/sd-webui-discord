/*
 * @Author: SpenserCai
 * @Date: 2023-08-19 16:21:45
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-22 14:35:37
 * @Description: file content
 */
package slash_handler

import (
	"log"

	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/SpenserCai/sd-webui-discord/utils"

	"github.com/SpenserCai/sd-webui-go/intersvc"
	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) upscalerModelChoise() []*discordgo.ApplicationCommandOptionChoice {
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	modesvc := &intersvc.SdapiV1Upscalers{}
	modesvc.Action(global.ClusterManager.GetNodeAuto().StableClient)
	if modesvc.Error != nil {
		log.Println(modesvc.Error)
		return choices
	}
	models := modesvc.GetResponse()
	for _, model := range *models {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  *model.Name,
			Value: *model.Name,
		})
	}
	return choices
}

func (shdl SlashHandler) ExtraSingleOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "extra_single",
		Description: "Upscaler and face restorer for single image",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "image_url",
				Description: "The url of the image",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "resize_mode",
				Description: "The resize mode",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Scaling",
						Value: 0,
					},
					{
						Name:  "Specify Size",
						Value: 1,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "upscaler_1",
				Description: "The upscaler1 to use",
				Required:    true,
				Choices:     shdl.upscalerModelChoise(),
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "upscaler_2",
				Description: "The upscaler2 to use",
				Required:    false,
				Choices:     shdl.upscalerModelChoise(),
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "upscaler_2_visibility",
				Description: "The upscaler2 visibility",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "gfpgan",
				Description: "The gfpgan visibility",
				MinValue:    func() *float64 { v := 0.0; return &v }(),
				MaxValue:    1.0,
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "upscaling_resize_w",
				Description: "The resize width",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "upscaling_resize_h",
				Description: "The resize height",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "upscaling_resize",
				Description: "The resize scale",
				MinValue:    func() *float64 { v := 0.0; return &v }(),
				MaxValue:    8.0,
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "upscaling_crop",
				Description: "The resize crop",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "upscaling_first",
				Description: "The upscaling first",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "codeformer_visibility",
				Description: "The codeformer visibility",
				MinValue:    func() *float64 { v := 0.0; return &v }(),
				MaxValue:    1.0,
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "codeformer_weight",
				Description: "The codeformer weight",
				MinValue:    func() *float64 { v := 0.0; return &v }(),
				MaxValue:    1.0,
				Required:    false,
			},
		},
	}
}

func (shdl SlashHandler) ExtraSingleSetOptions(dsOpt []*discordgo.ApplicationCommandInteractionDataOption, opt *intersvc.SdapiV1ExtraSingleImageRequest) {
	// default value
	opt.UpscalingResize = 2.0
	opt.ShowExtrasResults = func() *bool { v := true; return &v }()
	opt.CodeformerVisibility = func() *float64 { v := 0.0; return &v }()
	opt.CodeformerWeight = func() *float64 { v := 0.0; return &v }()
	opt.UpscaleFirst = func() *bool { v := false; return &v }()
	opt.UpscalingCrop = func() *bool { v := true; return &v }()
	opt.UpscalingResizew = 512
	opt.UpscalingResizeh = 512
	opt.Upscaler2 = func() *string { v := "None"; return &v }()
	for _, v := range dsOpt {
		switch v.Name {
		case "image_url":
			opt.Image, _ = utils.GetImageBase64(v.StringValue())
		case "upscaler_1":
			opt.Upscaler1 = func() *string { v := v.StringValue(); return &v }()
		case "upscaler_2":
			opt.Upscaler2 = func() *string { v := v.StringValue(); return &v }()
		case "upscaler_2_visibility":
			opt.ExtrasUpscaler2Visibility = func() *float64 { v := v.FloatValue(); return &v }()
		case "gfpgan":
			opt.GfpganVisibility = func() *float64 { v := v.FloatValue(); return &v }()
		case "resize_mode":
			opt.ResizeMode = v.IntValue()
		case "upscaling_resize_w":
			opt.UpscalingResizew = v.IntValue()
		case "upscaling_resize_h":
			opt.UpscalingResizeh = v.IntValue()
		case "upscaling_resize":
			opt.UpscalingResize = v.FloatValue()
		case "upscaling_crop":
			opt.UpscalingCrop = func() *bool { v := v.BoolValue(); return &v }()
		case "upscaling_first":
			opt.UpscaleFirst = func() *bool { v := v.BoolValue(); return &v }()
		case "codeformer_visibility":
			opt.CodeformerVisibility = func() *float64 { v := v.FloatValue(); return &v }()
		case "codeformer_weight":
			opt.CodeformerWeight = func() *float64 { v := v.FloatValue(); return &v }()
		}
	}
}

func (shdl SlashHandler) ExtraSingleAction(s *discordgo.Session, i *discordgo.InteractionCreate, opt *intersvc.SdapiV1ExtraSingleImageRequest, node *cluster.ClusterNode) {
	msg, err := shdl.SendStateMessage("Running", s, i)
	if err != nil {
		log.Println(err)
		return
	}
	extra_single := &intersvc.SdapiV1ExtraSingleImage{RequestItem: opt}
	extra_single.Action(node.StableClient)
	if extra_single.Error != nil {
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: func() *string { v := extra_single.Error.Error(); return &v }(),
		})
	} else {
		image, err := utils.GetImageReaderByBase64(extra_single.GetResponse().Image)
		if err != nil {
			s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
				Content: func() *string { v := err.Error(); return &v }(),
			})
		} else {
			s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
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

func (shdl SlashHandler) ExtraSingleCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	option := &intersvc.SdapiV1ExtraSingleImageRequest{}
	shdl.ReportCommandInfo(s, i)
	node := global.ClusterManager.GetNodeAuto()
	action := func() (map[string]interface{}, error) {
		shdl.ExtraSingleSetOptions(i.ApplicationCommandData().Options, option)
		shdl.ExtraSingleAction(s, i, option, node)
		return nil, nil
	}
	callback := func() {}
	node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
}
