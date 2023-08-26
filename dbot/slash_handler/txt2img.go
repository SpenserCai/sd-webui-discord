/*
 * @Author: SpenserCai
 * @Date: 2023-08-22 17:13:19
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-27 01:36:35
 * @Description: file content
 */
package slash_handler

import (
	"fmt"
	"log"

	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/SpenserCai/sd-webui-discord/utils"

	"github.com/SpenserCai/sd-webui-go/intersvc"
	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) samplerChoice() []*discordgo.ApplicationCommandOptionChoice {
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	modesvc := &intersvc.SdapiV1Samplers{}
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

func (shdl SlashHandler) Txt2imgOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "txt2img",
		Description: "Text generate image",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "Prompt text",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "negative_prompt",
				Description: "Negative prompt text",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "height",
				Description: "Height of the generated image. Default: 512",
				MinValue:    func() *float64 { v := 64.0; return &v }(),
				MaxValue:    2048.0,
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "width",
				Description: "Width of the generated image. Default: 512",
				MinValue:    func() *float64 { v := 64.0; return &v }(),
				MaxValue:    2048.0,
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "sampler",
				Description: "Sampler of the generated image. Default: Euler",
				Required:    false,
				Choices:     shdl.samplerChoice(),
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "steps",
				Description: "Steps of the generated image. Default: 20",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "cfg_scale",
				Description: "Cfg scale of the generated image. Default:7",
				MinValue:    func() *float64 { v := 1.0; return &v }(),
				MaxValue:    30.0,
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "seed",
				Description: "Seed of the generated image. Default: -1",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "controlnet_args",
				Description: "Controlnet args of the generated image. Default: {}",
				Required:    false,
			},
		},
	}
}

func (shdl SlashHandler) Txt2imgSetOptions(dsOpt []*discordgo.ApplicationCommandInteractionDataOption, opt *intersvc.SdapiV1Txt2imgRequest) {
	opt.NegativePrompt = shdl.GetSdDefaultSetting("negative_prompt", "").(string)
	opt.Height = func() *int64 { v := shdl.GetSdDefaultSetting("height", int64(512)).(int64); return &v }()
	opt.Width = func() *int64 { v := shdl.GetSdDefaultSetting("width", int64(512)).(int64); return &v }()
	opt.SamplerIndex = func() *string { v := "Euler"; return &v }()
	opt.Steps = func() *int64 { v := shdl.GetSdDefaultSetting("steps", int64(20)).(int64); return &v }()
	opt.CfgScale = func() *float64 { v := shdl.GetSdDefaultSetting("cfg_scale", 7.0).(float64); return &v }()
	opt.Seed = func() *int64 { v := int64(-1); return &v }()
	opt.NIter = func() *int64 { v := int64(1); return &v }()
	opt.Styles = []string{}
	opt.ScriptArgs = []interface{}{}
	opt.AlwaysonScripts = map[string]interface{}{}
	opt.OverrideSettings = map[string]interface{}{}

	for _, v := range dsOpt {
		switch v.Name {
		case "prompt":
			opt.Prompt = v.StringValue()
		case "negative_prompt":
			opt.NegativePrompt = v.StringValue()
		case "height":
			opt.Height = func() *int64 { v := v.IntValue(); return &v }()
		case "width":
			opt.Width = func() *int64 { v := v.IntValue(); return &v }()
		case "sampler":
			opt.SamplerIndex = func() *string { v := v.StringValue(); return &v }()
		case "steps":
			opt.Steps = func() *int64 { v := v.IntValue(); return &v }()
		case "cfg_scale":
			opt.CfgScale = func() *float64 { v := v.FloatValue(); return &v }()
		case "seed":
			opt.Seed = func() *int64 { v := v.IntValue(); return &v }()
		case "controlnet_args":
			script, err := shdl.GetControlNetScript(v.StringValue())
			if err == nil {
				tmpAScript := opt.AlwaysonScripts.(map[string]interface{})
				tmpAScript["controlnet"] = script
				opt.AlwaysonScripts = tmpAScript
			}
		}
	}

	// optJson, _ := json.Marshal(opt)
	// log.Println(string(optJson))
}

func (shdl SlashHandler) Txt2imgAction(s *discordgo.Session, i *discordgo.InteractionCreate, opt *intersvc.SdapiV1Txt2imgRequest, node *cluster.ClusterNode) {
	msg, err := shdl.SendStateMessage("Running", s, i)
	if err != nil {
		log.Println(err)
		return
	}
	txt2img := &intersvc.SdapiV1Txt2img{RequestItem: opt}
	txt2img.Action(node.StableClient)
	if txt2img.Error != nil {
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: func() *string { v := txt2img.Error.Error(); return &v }(),
		})
	} else {
		files := make([]*discordgo.File, 0)
		outinfo := txt2img.GetResponse().Info
		for j, v := range txt2img.GetResponse().Images {
			imageReader, err := utils.GetImageReaderByBase64(v)
			if err != nil {
				s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
					Content: func() *string { v := err.Error(); return &v }(),
				})
				return
			}
			files = append(files, &discordgo.File{
				Name:        fmt.Sprintf("image_%d.png", j),
				Reader:      imageReader,
				ContentType: "image/png",
			})
		}
		if len(files) >= 4 {
			files = files[0:4]
		}
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: func() *string { v := fmt.Sprintf("```\n%v```\n", *outinfo); return &v }(),
			Files:   files,
		})
	}

}

func (shdl SlashHandler) Txt2imgCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	option := &intersvc.SdapiV1Txt2imgRequest{}
	shdl.ReportCommandInfo(s, i)
	node := global.ClusterManager.GetNodeAuto()
	action := func() (map[string]interface{}, error) {
		shdl.Txt2imgSetOptions(i.ApplicationCommandData().Options, option)
		shdl.Txt2imgAction(s, i, option, node)
		return nil, nil
	}
	callback := func() {}
	node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
}
