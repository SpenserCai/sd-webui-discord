/*
 * @Author: SpenserCai
 * @Date: 2023-08-22 17:13:19
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-23 19:06:08
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

func (shdl SlashHandler) SamplerChoice() []*discordgo.ApplicationCommandOptionChoice {
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

func (shdl SlashHandler) SdModelChoice() []*discordgo.ApplicationCommandOptionChoice {
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	modesvc := &intersvc.SdapiV1SdModels{}
	modesvc.Action(global.ClusterManager.GetNodeAuto().StableClient)
	if modesvc.Error != nil {
		log.Println(modesvc.Error)
		return choices
	}
	models := modesvc.GetResponse()
	for _, model := range *models {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  *model.ModelName,
			Value: *model.Title,
		})
	}
	return choices
}

func (shdl SlashHandler) SdVaeChoice() []*discordgo.ApplicationCommandOptionChoice {
	choice := []*discordgo.ApplicationCommandOptionChoice{}
	// add Automatic
	choice = append(choice, &discordgo.ApplicationCommandOptionChoice{
		Name:  "Automatic",
		Value: "Automatic",
	})
	vaesvc := &intersvc.SdapiV1SdVae{}
	vaesvc.Action(global.ClusterManager.GetNodeAuto().StableClient)
	if vaesvc.Error != nil {
		log.Println(vaesvc.Error)
		return choice
	}
	vaes := vaesvc.GetResponse()
	for _, vae := range *vaes {
		choice = append(choice, &discordgo.ApplicationCommandOptionChoice{
			Name:  *vae.ModelName,
			Value: *vae.ModelName,
		})
	}
	return choice
}

func (shdl SlashHandler) Txt2imgOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "txt2img",
		Description: "Generate an img from text.",
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
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "sampler",
				Description:  "Sampler of the generated image. Default: Euler",
				Required:     false,
				Autocomplete: true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "steps",
				Description: "Steps of the generated image. Default: 30",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "cfg_scale",
				Description: "Cfg scale of the generated image. Default: 7",
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
				Name:        "styles",
				Description: "Style of the generated image, split with | . Default: None",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "controlnet_args",
				Description: "Controlnet args of the generated image.",
				Required:    false,
			},
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "checkpoint",
				Description:  "Sd model checkpoint. Default: SDXL 1.0",
				Required:     false,
				Autocomplete: true,
			},
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "sd_vae",
				Description:  "Sd vae. Default: Automatic",
				Required:     false,
				Autocomplete: true,
			},
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "refiner_checkpoint",
				Description:  "Refiner checkpoint. Default: None",
				Required:     false,
				Autocomplete: true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "refiner_switch_at",
				Description: "Refiner switch at. Default: 0.0",
				Required:    false,
				MinValue:    func() *float64 { v := 0.0; return &v }(),
				MaxValue:    1.0,
			},
			// {
			// 	Type:        discordgo.ApplicationCommandOptionInteger,
			// 	Name:        "n_iter",
			// 	Description: "Number of iterations. Default: 1",
			// 	Required:    false,
			// 	MinValue:    func() *float64 { v := 1.0; return &v }(),
			// 	MaxValue:    4.0,
			// 	Choices: []*discordgo.ApplicationCommandOptionChoice{
			// 		{
			// 			Name:  "1",
			// 			Value: 1,
			// 		},
			// 		{
			// 			Name:  "2",
			// 			Value: 2,
			// 		},
			// 		{
			// 			Name:  "4",
			// 			Value: 4,
			// 		},
			// 	},
			// },
		},
	}
}

func (shdl SlashHandler) Txt2imgSetOptions(dsOpt []*discordgo.ApplicationCommandInteractionDataOption, opt *intersvc.SdapiV1Txt2imgRequest, i *discordgo.InteractionCreate) {
	opt.NegativePrompt = shdl.GetDefaultSettingFromUser("negative_prompt", "", i).(string)
	opt.Height = func() *int64 { v := shdl.GetDefaultSettingFromUser("height", int64(512), i).(int64); return &v }()
	opt.Width = func() *int64 { v := shdl.GetDefaultSettingFromUser("width", int64(512), i).(int64); return &v }()
	opt.SamplerIndex = func() *string { v := shdl.GetDefaultSettingFromUser("sampler", "Euler", i).(string); return &v }()
	opt.Steps = func() *int64 { v := shdl.GetDefaultSettingFromUser("steps", int64(20), i).(int64); return &v }()
	opt.CfgScale = func() *float64 { v := shdl.GetDefaultSettingFromUser("cfg_scale", 7.0, i).(float64); return &v }()
	opt.Seed = func() *int64 { v := int64(-1); return &v }()
	opt.NIter = func() *int64 { v := int64(1); return &v }()
	opt.Styles = []string{}
	opt.RefinerCheckpoint = ""
	opt.RefinerSwitchAt = float64(0.0)
	opt.ScriptArgs = []interface{}{}
	opt.AlwaysonScripts = map[string]interface{}{}
	opt.OverrideSettings = map[string]interface{}{}
	isSetCheckpoints := false
	isSetVae := false
	defaultCheckpoints := shdl.GetDefaultSettingFromUser("sd_model_checkpoint", "", i).(string)
	defaultVae := shdl.GetDefaultSettingFromUser("sd_vae", "", i).(string)

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
		case "styles":
			styleList := strings.Split(v.StringValue(), "|")
			outStyleList := []string{}
			for _, style := range styleList {
				outStyleList = append(outStyleList, strings.TrimSpace(style))
			}
			opt.Styles = outStyleList
		case "controlnet_args":
			script, err := shdl.GetControlNetScript(v.StringValue())
			if err == nil {
				tmpAScript := opt.AlwaysonScripts.(map[string]interface{})
				tmpAScript["controlnet"] = script
				opt.AlwaysonScripts = tmpAScript
			}
		case "checkpoint":
			tmpOverrideSettings := opt.OverrideSettings.(map[string]interface{})
			tmpOverrideSettings["sd_model_checkpoint"] = v.StringValue()
			opt.OverrideSettings = tmpOverrideSettings
			isSetCheckpoints = true
		case "sd_vae":
			tmpOverrideSettings := opt.OverrideSettings.(map[string]interface{})
			tmpOverrideSettings["sd_vae"] = v.StringValue()
			opt.OverrideSettings = tmpOverrideSettings
			isSetVae = true
		case "refiner_checkpoint":
			opt.RefinerCheckpoint = v.StringValue()
		case "refiner_switch_at":
			opt.RefinerSwitchAt = v.FloatValue()
		case "n_iter":
			opt.NIter = func() *int64 { v := v.IntValue(); return &v }()
		}
	}
	if !isSetCheckpoints && defaultCheckpoints != "" {
		tmpOverrideSettings := opt.OverrideSettings.(map[string]interface{})
		tmpOverrideSettings["sd_model_checkpoint"] = defaultCheckpoints
		opt.OverrideSettings = tmpOverrideSettings
	}
	if !isSetVae && defaultVae != "" && defaultVae != "Automatic" {
		tmpOverrideSettings := opt.OverrideSettings.(map[string]interface{})
		tmpOverrideSettings["sd_vae"] = defaultVae
		opt.OverrideSettings = tmpOverrideSettings
	}

	// optJson, _ := json.Marshal(opt)
	// log.Println(string(optJson))
}

func (shdl SlashHandler) Txt2imgAction(s *discordgo.Session, i *discordgo.InteractionCreate, opt *intersvc.SdapiV1Txt2imgRequest, node *cluster.ClusterNode) {
	txt2img := &intersvc.SdapiV1Txt2img{RequestItem: opt}
	txt2img.Action(node.StableClient)
	if txt2img.Error != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func() *string { v := txt2img.Error.Error(); return &v }(),
		})
	} else {
		files := make([]*discordgo.File, 0)
		outinfo := txt2img.GetResponse().Info
		// parse outinfo from json
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(*outinfo), &data); err != nil {
			fmt.Println("Error:", err)
			return
		}

		context := ""
		if !global.Config.DisableReturnGenInfo {
			// 如果outinfo长度大于2000则context为：Success！，并创建info.json文件
			if len(*outinfo) > 1800 {
				context = "Success!"
				infoJson, _ := utils.GetJsonReaderByJsonString(*outinfo)
				files = append(files, &discordgo.File{
					Name:        "info.json",
					ContentType: "application/json",
					Reader:      infoJson,
				})
			} else {
				var fOutput bytes.Buffer
				json.Indent(&fOutput, []byte(*outinfo), "", "  ")
				context = fmt.Sprintf("```json\n%v```\n", fOutput.String())
			}
		}
		seed := fmt.Sprintf("%.0f", data["seed"])

		if len(txt2img.GetResponse().Images) > 1 {
			mergeImageBase64, _ := utils.MergeImageFromBase64(txt2img.GetResponse().Images)
			// 把mergeImageBase64放在第一位
			txt2img.GetResponse().Images = append([]string{mergeImageBase64}, txt2img.GetResponse().Images...)
		}

		for j, v := range txt2img.GetResponse().Images {
			imageReader, err := utils.GetImageReaderByBase64(v)
			if err != nil {
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
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
		if len(files) >= 5 {
			files = files[0:5]
		}

		// 生成主要Embed
		mainEmbed := shdl.MessageEmbedTemplate()
		mainEmbed.Image = &discordgo.MessageEmbedImage{
			URL: fmt.Sprintf("attachment://%s", files[0].Name),
		}
		mainEmbed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:  "Prompt",
				Value: opt.Prompt,
			},
			{
				Name:  "Model",
				Value: data["sd_model_name"].(string),
			},
			{
				Name: "VAE",
				Value: func() string {
					vae, ok := data["sd_vae_name"]
					if ok && vae != nil {
						return vae.(string)
					} else {
						return "Automatic"
					}
				}(),
			},
			{
				Name:  "Sampler",
				Value: data["sampler_name"].(string),
			},
			{
				Name:   "Size",
				Value:  fmt.Sprintf("%dx%d", *opt.Height, *opt.Width),
				Inline: true,
			},
			{
				Name:   "Steps",
				Value:  fmt.Sprintf("%v", data["steps"]),
				Inline: true,
			},
			{
				Name:   "Cfg Scale",
				Value:  fmt.Sprintf("%v", data["cfg_scale"]),
				Inline: true,
			},
			{
				Name:   "Seed",
				Value:  seed,
				Inline: true,
			},
			{
				Name:   "User",
				Value:  fmt.Sprintf("<@%s>", shdl.GetDiscordUserId(i)),
				Inline: true,
			},
		}

		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &context,
			Embeds: &[]*discordgo.MessageEmbed{
				mainEmbed,
			},
			Files: files,
		})
		if err != nil {
			log.Println(err)
		}
	}

}

func (shdl SlashHandler) Txt2imgCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		option := &intersvc.SdapiV1Txt2imgRequest{}
		shdl.RespondStateMessage("Running", s, i)
		node := global.ClusterManager.GetNodeAuto()
		action := func() (map[string]interface{}, error) {
			shdl.Txt2imgSetOptions(i.ApplicationCommandData().Options, option, i)
			shdl.Txt2imgAction(s, i, option, node)
			return nil, nil
		}
		callback := func() {}
		node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
	case discordgo.InteractionApplicationCommandAutocomplete:
		repChoices := []*discordgo.ApplicationCommandOptionChoice{}
		data := i.ApplicationCommandData()

		for _, opt := range data.Options {
			if opt.Name == "checkpoint" && opt.Focused {
				repChoices = shdl.FilterChoice(global.LongDBotChoice["sd_model_checkpoint"], opt)
				continue
			}
			if opt.Name == "sampler" && opt.Focused {
				repChoices = shdl.FilterChoice(shdl.SamplerChoice(), opt)
				continue
			}
			if opt.Name == "refiner_checkpoint" && opt.Focused {
				repChoices = shdl.FilterChoice(global.LongDBotChoice["sd_model_checkpoint"], opt)
				continue
			}
			if opt.Name == "sd_vae" && opt.Focused {
				repChoices = shdl.FilterChoice(global.LongDBotChoice["sd_vae"], opt)
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
