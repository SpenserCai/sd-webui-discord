/*
 * @Author: SpenserCai
 * @Date: 2023-08-20 12:45:58
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-06 23:59:27
 * @Description: file content
 */

package slash_handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/SpenserCai/sd-webui-discord/utils"

	"github.com/SpenserCai/sd-webui-go/intersvc"
	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) controlnetControlModeChoice() []*discordgo.ApplicationCommandOptionChoice {
	return []*discordgo.ApplicationCommandOptionChoice{
		{
			Name:  "Balanced",
			Value: 0,
		},
		{
			Name:  "My prompt is more important",
			Value: 1,
		},
		{
			Name:  "ControlNet is more important",
			Value: 2,
		},
	}
}

func (shdl SlashHandler) controlnetZoomModeChoice() []*discordgo.ApplicationCommandOptionChoice {
	return []*discordgo.ApplicationCommandOptionChoice{
		{
			Name:  "Resize Only (Stretch)",
			Value: 0,
		},
		{
			Name:  "Crop and Resize",
			Value: 1,
		},
		{
			Name:  "Resize and Fill",
			Value: 2,
		},
	}
}

func (shdl SlashHandler) ControlnetModelChoice() []*discordgo.ApplicationCommandOptionChoice {
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	modesvc := &intersvc.ControlnetModelList{}
	modesvc.Action(global.ClusterManager.GetNodeAuto().StableClient)
	if modesvc.Error != nil {
		log.Println(modesvc.Error)
		return choices
	}
	models := modesvc.GetResponse().ModelList
	for _, model := range models {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  model,
			Value: model,
		})
	}
	return choices
}

func (shdl SlashHandler) ControlnetModuleChoice() []*discordgo.ApplicationCommandOptionChoice {
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
	return choices
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
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "module",
				Description:  "The module to use",
				Required:     true,
				Autocomplete: true,
			},
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "model",
				Description:  "The model to use",
				Required:     true,
				Autocomplete: true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "weight",
				Description: "The weight of the module. Default: 1.0",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "control_mode",
				Description: "Control mode. Default: Balanced",
				Required:    false,
				Choices:     shdl.controlnetControlModeChoice(),
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "zoom_mode",
				Description: "Zoom mode. Default: Crop and Resize",
				Required:    false,
				Choices:     shdl.controlnetZoomModeChoice(),
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "processor_res",
				Description: "The resolution of the processor. Default: 512",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "threshold_a",
				Description: "The threshold of the processor. Default: 64.0",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "threshold_b",
				Description: "The threshold of the processor. Default: 64.0",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "pixel_perfect",
				Description: "Whether to use pixel perfect. Default: false",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "guidance_start",
				Description: "The guidance start. Default: 0.0",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "guidance_end",
				Description: "The guidance end. Default: 1.0",
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

func (shdl SlashHandler) ControlnetArgJsonGen(dsOpt []*discordgo.ApplicationCommandInteractionDataOption) string {
	cnArg := &intersvc.ControlnetPredictArgsItem{}
	cnArg.Enabled = true
	cnArg.Weight = 1.0
	cnArg.ControlMode = 0
	cnArg.ResizeMode = 1
	cnArg.ProcessorRes = 512
	cnArg.ThresholdA = 64.0
	cnArg.ThresholdB = 64.0
	cnArg.GuidanceStart = 0.0
	cnArg.GuidanceEnd = 1.0
	cnArg.PixelPerFect = false
	for _, v := range dsOpt {
		switch v.Name {
		case "control_mode":
			cnArg.ControlMode = v.IntValue()
		case "zoom_mode":
			cnArg.ResizeMode = v.IntValue()
		case "processor_res":
			cnArg.ProcessorRes = v.IntValue()
		case "threshold_a":
			cnArg.ThresholdA = v.FloatValue()
		case "threshold_b":
			cnArg.ThresholdB = v.FloatValue()
		case "pixel_perfect":
			cnArg.PixelPerFect = v.BoolValue()
		case "module":
			cnArg.Module = v.StringValue()
		case "model":
			cnArg.Model = v.StringValue()
		case "image_url":
			cnArg.Image = v.StringValue()
		case "guidance_start":
			cnArg.GuidanceStart = v.FloatValue()
		case "guidance_end":
			cnArg.GuidanceEnd = v.FloatValue()
		case "weight":
			cnArg.Weight = v.FloatValue()
		}
	}
	jsonStr, _ := json.Marshal(cnArg)
	return string(jsonStr)
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
		outinfo := shdl.ControlnetArgJsonGen(i.ApplicationCommandData().Options)
		context := ""
		// 如果outinfo长度大于2000则context为：Success！，并创建info.json文件
		if len(outinfo) > 2000 {
			context = "Success!"
			infoJson, _ := utils.GetJsonReaderByJsonString(outinfo)
			files = append(files, &discordgo.File{
				Name:        "args.json",
				ContentType: "application/json",
				Reader:      infoJson,
			})
		} else {
			context = fmt.Sprintf("```json\n%v```\n", outinfo)
		}
		for n, img := range controlnet_detect.GetResponse().Images {
			var imageReader *bytes.Reader
			var err error
			if *opt.ControlnetModule == "none" {
				imageReader, err = utils.GetImageReaderByBase64(opt.ControlnetInputImages[0])
			} else {
				imageReader, err = utils.GetImageReaderByBase64(img)
			}
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
			Content: &context,
			Files:   files,
		})
	}
}

func (shdl SlashHandler) ControlnetDetectCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
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
	case discordgo.InteractionApplicationCommandAutocomplete:
		repChoices := []*discordgo.ApplicationCommandOptionChoice{}
		data := i.ApplicationCommandData()
		for _, opt := range data.Options {
			if opt.Name == "module" && opt.Focused {
				repChoices = shdl.FilterChoice(global.LongDBotChoice["control_net_module"], opt)
				continue
			}
			if opt.Name == "model" && opt.Focused {
				repChoices = shdl.FilterChoice(global.LongDBotChoice["control_net_model"], opt)
				continue
			}

		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: repChoices,
			},
		})
	}
}
