/*
 * @Author: SpenserCai
 * @Date: 2023-08-22 17:13:19
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2024-03-13 15:59:37
 * @Description: file content
 */
package slash_handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/dbot/slash_handler/option_values"
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
			option_values.NegativePrompt(),
			option_values.Height(),
			option_values.Width(),
			option_values.Sampler(),
			option_values.Steps(),
			option_values.CfgScale(),
			option_values.Seed(),
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "styles",
				Description: "Style of the generated image, split with | . Default: None",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "controlnet_args",
				Description: "Controlnet args of the generated image.multi args split with `,` .",
				Required:    false,
			},
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "checkpoint",
				Description:  "Sd model checkpoint. Default: " + global.Config.SDWebUi.DefaultSetting.Model,
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
			option_values.NIter(),
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
	opt.Seed = func() *int64 { v := shdl.GetDefaultSettingFromUser("seed", int64(-1), i).(int64); return &v }()
	opt.NIter = func() *int64 { v := shdl.GetDefaultSettingFromUser("n_iter", int64(1), i).(int64); return &v }()
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
	clipSkip := shdl.GetDefaultSettingFromUser("clip_skip", int64(1), i).(int64)
	freeu := global.Features.FreeU && strings.ToLower(shdl.GetDefaultSettingFromUser("freeu", "false", i).(string)) == "true"
	sag := global.Features.Sag && strings.ToLower(shdl.GetDefaultSettingFromUser("sag", "false", i).(string)) == "true"

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
		case "clip_skip":
			clipSkip = v.IntValue()
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
	if freeu {
		b1 := shdl.GetDefaultSettingFromUser("freeu_b1", 1.01, i).(float64)
		b2 := shdl.GetDefaultSettingFromUser("freeu_b2", 1.02, i).(float64)
		s1 := shdl.GetDefaultSettingFromUser("freeu_s1", 0.99, i).(float64)
		s2 := shdl.GetDefaultSettingFromUser("freeu_s2", 0.95, i).(float64)
		freeuScript := &FreeUScript{}
		freeuScript.Set(b1, b2, s1, s2)
		tmpAScript := opt.AlwaysonScripts.(map[string]interface{})
		tmpAScript["freeu integrated"] = freeuScript
		opt.AlwaysonScripts = tmpAScript
	}
	if sag {
		scale := shdl.GetDefaultSettingFromUser("sag_scale", 0.5, i).(float64)
		blurSigma := shdl.GetDefaultSettingFromUser("sag_blur_sigma", 2.0, i).(float64)
		sagScript := &SagScript{}
		sagScript.Set(scale, blurSigma)
		tmpAScript := opt.AlwaysonScripts.(map[string]interface{})
		tmpAScript["selfattentionguidance integrated"] = sagScript
		opt.AlwaysonScripts = tmpAScript
	}
	if clipSkip != 1 {
		tmpOverrideSettings := opt.OverrideSettings.(map[string]interface{})
		tmpOverrideSettings["CLIP_stop_at_last_layers"] = clipSkip
		opt.OverrideSettings = tmpOverrideSettings
	}

}

func (shdl SlashHandler) BuildTxt2imgComponent(i *discordgo.InteractionCreate, imgCount int64) *[]discordgo.MessageComponent {

	components := []discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					CustomID: "txt2img|retry",
					Label:    "Retry",
					Style:    discordgo.SecondaryButton,
					Emoji:    &discordgo.ComponentEmoji{Name: "🔄"},
					Disabled: func() bool {
						return !global.Config.UserCenter.Enable
					}(),
				},
				&discordgo.Button{
					CustomID: "txt2img|delete|" + shdl.GetDiscordUserId(i),
					Label:    "Delete",
					Style:    discordgo.SecondaryButton,
					Emoji:    &discordgo.ComponentEmoji{Name: "🗑️"},
				},
			},
		},
	}

	// 图片数量大于1时，添加多图片按钮
	multiImageButton := []discordgo.MessageComponent{}
	if imgCount > 1 {
		for j := int64(0); j < imgCount; j++ {
			multiImageButton = append(multiImageButton, &discordgo.Button{
				CustomID: fmt.Sprintf("txt2img|multi_image|%d", j),
				Label:    fmt.Sprintf("%d", j+1),
				Style:    discordgo.SecondaryButton,
				Emoji:    &discordgo.ComponentEmoji{Name: "🖼️"},
				Disabled: func() bool {
					return !global.Config.UserCenter.Enable
				}(),
			})
		}
		components = append(components, &discordgo.ActionsRow{
			Components: multiImageButton,
		})
	}
	return &components
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
		var mergeAdditionalFile *discordgo.File
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
			// 根据opt.NIter数量拼接
			mergeImageBase64, _ := utils.MergeImageFromBase64(txt2img.GetResponse().Images[:*opt.NIter])
			// 如果图片总数大于opt.NIter，则说明有附加图片，把附加图片单独拼接
			if int64(len(txt2img.GetResponse().Images)) > *opt.NIter {
				mergeAdditionalImageBase64, _ := utils.MergeImageFromBase64(txt2img.GetResponse().Images[*opt.NIter:])
				imageReader, _ := utils.GetImageReaderByBase64(mergeAdditionalImageBase64)
				mergeAdditionalFile = &discordgo.File{
					Name:        "merge_additional.png",
					ContentType: "image/png",
					Reader:      imageReader,
				}
			}
			// 把mergeImageBase64放在第一位
			txt2img.GetResponse().Images = []string{mergeImageBase64}
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

		// 生成主要Embed
		mainEmbed := shdl.MessageEmbedTemplate()
		mainEmbed.Image = &discordgo.MessageEmbedImage{
			URL: fmt.Sprintf("attachment://%s", files[0].Name),
		}

		prompt := opt.Prompt
		if len(prompt) > 600 {
			prompt = prompt[:600] + "..."
		}

		mainEmbed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:  "Prompt",
				Value: prompt,
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
		allEmbeds := []*discordgo.MessageEmbed{mainEmbed}
		// 如果合并的附加图片不为空，则添加附加图片的Embed
		if mergeAdditionalFile != nil {
			additionalEmbed := shdl.MessageEmbedTemplate()
			additionalEmbed.Title = "Additional"
			additionalEmbed.Image = &discordgo.MessageEmbedImage{
				URL: fmt.Sprintf("attachment://%s", mergeAdditionalFile.Name),
			}
			allEmbeds = append(allEmbeds, additionalEmbed)
			files = append(files, mergeAdditionalFile)
		}
		msg, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content:    &context,
			Embeds:     &allEmbeds,
			Files:      files,
			Components: shdl.BuildTxt2imgComponent(i, *opt.NIter),
		})
		if err != nil {
			log.Println(err)
		} else {
			opt.Seed = func() *int64 { v, _ := strconv.ParseInt(seed, 10, 64); return &v }()
			shdl.SetHistory("txt2img", msg.ID, i, opt)
			urls := make([]string, 0)
			for _, v := range msg.Embeds {
				urls = append(urls, v.Image.URL)
			}
			shdl.SetHistoryImages(msg.ID, i, urls, shdl.GetBase64ImageListBlurHash(txt2img.GetResponse().Images))
		}

	}

}

func (shdl SlashHandler) MaxLoadGuard(s *discordgo.Session, i *discordgo.InteractionCreate, option *intersvc.SdapiV1Txt2imgRequest) error {

	maxLoad := shdl.GetDefaultSettingFromUser("max_load", int64(1<<16), i).(int64)
	if maxLoad == 0 {
		maxLoad = 1 << 16
	}
	niter := int64(1)
	if option.NIter != nil {
		niter = *option.NIter
	}
	steps := int64(1)
	if option.Steps != nil {
		steps = *option.Steps
	}
	widthScore := 1 + ((*option.Width - 1) / 64)
	heightScore := 1 + ((*option.Height - 1) / 64)
	if widthScore*heightScore*steps*niter > maxLoad {
		var content string
		if niter > 1 {
			content = "Request is too heavy. Consider decresing widht, height, steps or n_iter."
		} else {
			content = "Request is too heavy. Consider decresing widht, height, or steps."
		}
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		return errors.New("request is too heavy")
	}
	return nil
}

func (shdl SlashHandler) Txt2imgAppHandler(s *discordgo.Session, i *discordgo.InteractionCreate, otherOption *intersvc.SdapiV1Txt2imgRequest, useOtherOption bool) {
	option := &intersvc.SdapiV1Txt2imgRequest{}
	shdl.RespondStateMessage("Running", s, i)
	node := global.ClusterManager.GetNodeAuto()
	action := func() (map[string]interface{}, error) {
		if !useOtherOption {
			shdl.Txt2imgSetOptions(i.ApplicationCommandData().Options, option, i)
		} else {
			option = otherOption
		}
		err := shdl.MaxLoadGuard(s, i, option)
		if err != nil {
			return nil, err
		}
		shdl.Txt2imgAction(s, i, option, node)
		return nil, nil
	}
	callback := func() {}
	node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
}

func (shdl SlashHandler) Txt2imgComponentHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 将CustomID分割为数组
	customIDList := strings.Split(i.MessageComponentData().CustomID, "|")
	cmd := fmt.Sprintf("%s|%s", customIDList[0], customIDList[1])
	switch cmd {
	case "txt2img|delete":
		ownerId := shdl.GetDiscordUserId(i)
		if len(customIDList) == 3 {
			ownerId = customIDList[2]
		}
		if shdl.GetDiscordUserId(i) == ownerId || shdl.IsAdmin(i) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
			})
			err := s.ChannelMessageDelete(i.ChannelID, i.Interaction.Message.ID)
			if err == nil {
				shdl.DeleteHistory(i.Interaction.Message.ID, ownerId)
			}

		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You are not the author of this!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	case "txt2img|retry":
		option := &intersvc.SdapiV1Txt2imgRequest{}
		err := shdl.GetHistory("txt2img", i.Interaction.Message.ID, option)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Original data has been cleared",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
		option.Seed = func() *int64 { v := int64(-1); return &v }()
		shdl.Txt2imgAppHandler(s, i, option, true)
	case "txt2img|multi_image":
		option := &intersvc.SdapiV1Txt2imgRequest{}
		err := shdl.GetHistory("txt2img", i.Interaction.Message.ID, option)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Original data has been cleared",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
		if len(customIDList) == 3 {
			index, _ := strconv.ParseInt(customIDList[2], 10, 64)
			option.NIter = func() *int64 { v := int64(1); return &v }()
			option.Seed = func() *int64 { v := *option.Seed + index; return &v }()
			shdl.Txt2imgAppHandler(s, i, option, true)
		}

	}
}

func (shdl SlashHandler) Txt2imgCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		shdl.Txt2imgAppHandler(s, i, nil, false)
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
	case discordgo.InteractionMessageComponent:
		shdl.Txt2imgComponentHandler(s, i)
	}
}
