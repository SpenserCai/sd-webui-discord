/*
 * @Author: SpenserCai
 * @Date: 2023-08-20 12:45:58
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-22 14:36:15
 * @Description: file content
 */

package slash_handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/SpenserCai/sd-webui-discord/utils"

	"github.com/SpenserCai/sd-webui-go/intersvc"
	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) controlnetModuleChoice() []*discordgo.ApplicationCommandOptionChoice {
	exclued := []string{
		"clip_vision",
		"t2ia_color_grid",
		"pidinet",
		"pidinet_safe",
		"t2ia_sketch_pidi",
		"scribble_pidinet",
		"scribble_xdog",
		"scribble_hed",
		"normal_bae",
		"lineart_realistic",
		"lineart_coarse",
		"lineart_anime",
		"pidinet",
		"pidinet_safe",
		"pidinet_sketch",
		"pidinet_scribble",
		"inpaint_global_harmonious",
		"inpaint_only",
		"inpaint_only+lama",
		"normal_map",
		"invert",
		"shuffle",
		"tile_colorfix",
		"tile_colorfix+sharp",
		"reference_adain+attn",
		"mediapipe_face",
	}
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	modulesvc := &intersvc.ControlnetModuleList{}
	modulesvc.Action(global.ClusterManager.GetNodeAuto().StableClient)
	if modulesvc.Error != nil {
		log.Println(modulesvc.Error)
		return choices
	}
	model_list := modulesvc.GetResponse().ModuleList
	for _, model := range model_list {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  model,
			Value: model,
		})
	}
	newChoices := []*discordgo.ApplicationCommandOptionChoice{}
	for _, choice := range choices {
		exclu := false
		for _, ex := range exclued {
			if strings.Contains(choice.Name, ex) {
				exclu = true
				break
			}
		}
		if !exclu {
			newChoices = append(newChoices, choice)
		}
	}
	return newChoices
}

func (shdl SlashHandler) ControlnetDetectOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "controlnet_detect",
		Description: "ControlNet detect",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "image_url",
				Description: "The url of the images,split by ','",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "module",
				Description: "The module to use",
				Required:    true,
				Choices:     shdl.controlnetModuleChoice(),
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "processor_res",
				Description: "The resolution of the processor",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "threshold_a",
				Description: "The threshold of the processor",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "threshold_b",
				Description: "The threshold of the processor",
				Required:    false,
			},
		},
	}
}

func (shdl SlashHandler) ControlnetDetectSetOptions(dsOpt []*discordgo.ApplicationCommandInteractionDataOption, opt *intersvc.ControlnetDetectRequest) {
	opt.ControlnetProcessorRes = func() *int64 { v := int64(512); return &v }()
	opt.ControlnetThresholda = func() *float64 { v := float64(64.0); return &v }()
	opt.ControlnetThresholdb = func() *float64 { v := float64(64.0); return &v }()
	for _, v := range dsOpt {
		switch v.Name {
		case "image_url":
			imgUrls := strings.Split(v.StringValue(), ",")
			imgs := []string{}
			for _, imgUrl := range imgUrls {
				img, _ := utils.GetImageBase64(imgUrl)
				imgs = append(imgs, img)
			}
			opt.ControlnetInputImages = imgs
		case "module":
			opt.ControlnetModule = func() *string { v := v.StringValue(); return &v }()
		case "processor_res":
			opt.ControlnetProcessorRes = func() *int64 { v := v.IntValue(); return &v }()
		case "threshold_a":
			opt.ControlnetThresholda = func() *float64 { v := v.FloatValue(); return &v }()
		case "threshold_b":
			opt.ControlnetThresholdb = func() *float64 { v := v.FloatValue(); return &v }()
		}
	}
}

func (shdl SlashHandler) ControlnetDetectAction(s *discordgo.Session, i *discordgo.InteractionCreate, opt *intersvc.ControlnetDetectRequest, node *cluster.ClusterNode) {
	msg, err := shdl.SendStateMessage("Running", s, i)
	if err != nil {
		log.Println(err)
		return
	}
	if len(opt.ControlnetInputImages) > 4 {
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: func() *string { v := "Too many images, please input less than 4 images"; return &v }(),
		})
		return
	}
	controlnet_detect := &intersvc.ControlnetDetect{RequestItem: opt}
	controlnet_detect.Action(node.StableClient)
	if controlnet_detect.Error != nil {
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: func() *string { v := controlnet_detect.Error.Error(); return &v }(),
		})
	} else {
		files := make([]*discordgo.File, 0)
		for n, img := range controlnet_detect.GetResponse().Images {
			imageReader, err := utils.GetImageReaderByBase64(img)
			if err != nil {
				s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
					Content: func() *string { v := err.Error(); return &v }(),
				})
				return
			}
			files = append(files, &discordgo.File{
				Name:        fmt.Sprintf("result_%d.png", n),
				Reader:      imageReader,
				ContentType: "image/png",
			})
		}
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: func() *string { v := "Success"; return &v }(),
			Files:   files,
		})
	}
}

func (shdl SlashHandler) ControlnetDetectCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	option := &intersvc.ControlnetDetectRequest{}
	shdl.ReportCommandInfo(s, i)
	node := global.ClusterManager.GetNodeAuto()
	action := func() (map[string]interface{}, error) {
		shdl.ControlnetDetectSetOptions(i.ApplicationCommandData().Options, option)
		shdl.ControlnetDetectAction(s, i, option, node)
		return nil, nil
	}
	callback := func() {}
	node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
}
