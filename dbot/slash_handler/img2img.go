package slash_handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/dbot/slash_handler/option_values"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/SpenserCai/sd-webui-discord/utils"

	"github.com/SpenserCai/sd-webui-go/intersvc"
	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) Img2imgOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "img2img",
		Description: "Modify an image.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "init_image",
				Description: "Initial image, a file upload.",
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Required:    true,
			},
			{
				Name:        "resize_by_scale",
				Description: "Resize by scale.",
				Type:        discordgo.ApplicationCommandOptionNumber,
				Required:    false,
				MinValue:    func() *float64 { v := 0.05; return &v }(),
				MaxValue:    4.0,
			},
			{
				Name:        "mask",
				Description: "Mask image, a file upload.",
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "Prompt text",
				Required:    false,
			},
			option_values.NegativePrompt(),
			option_values.Height(),
			option_values.Width(),
			option_values.Sampler(),
			option_values.Steps(),
			option_values.CfgScale(),
			option_values.Seed(),
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "resize_mode",
				Description: "Resize mode of the generated image. Default: Just resize",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Just resize",
						Value: 0,
					},
					{
						Name:  "Crop and resize",
						Value: 1,
					},
					{
						Name:  "Resize and fill",
						Value: 2,
					},
					{
						Name:  "Just resize (latent upscale)",
						Value: 3,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "denoising_strength",
				Description: "Denoising strength of the generated image. Default: 0.7",
				MinValue:    func() *float64 { v := 0.0; return &v }(),
				MaxValue:    1.0,
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "mask_blur", // 蒙板模糊
				Description: "Mask blur of the generated image. Default: 4",
				MinValue:    func() *float64 { v := 0.0; return &v }(),
				MaxValue:    64.0,
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "inpainting_fill", // 蒙版遮住的内容
				Description: "Masked content. Default: original",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Fill",
						Value: 0,
					},
					{
						Name:  "Original",
						Value: 1,
					},
					{
						Name:  "Latent noise",
						Value: 2,
					},
					{
						Name:  "Latent nothing",
						Value: 3,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "mask_mode", // 蒙板模式
				Description: "Mask mode.",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Inpaint masked",
						Value: 0,
					},
					{
						Name:  "Inpaint not masked",
						Value: 1,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "inpaint_mask_only", //  重绘区域, False: whole picture True：only masked
				Description: "Inpaint Area. Default: Whole picture",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "inpaint_full_res_padding",
				Description: "Only masked padding, pixels. Default:32",
				Required:    false,
				MinValue:    func() *float64 { v := 0.0; return &v }(),
				MaxValue:    256.0,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "styles",
				Description: "Style of the generated image,splite with | . Default: None",
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
		},
	}
}

func (shdl SlashHandler) Img2imgSetOptions(cmd discordgo.ApplicationCommandInteractionData, opt *intersvc.SdapiV1Img2imgRequest, i *discordgo.InteractionCreate) {
	opt.Prompt = ""
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
	// img2img default
	opt.DenoisingStrength = func() *float64 { v := 0.7; return &v }()
	opt.MaskBlur = 4
	opt.InpaintingFill = 1
	opt.InpaintingMaskInvert = 0
	opt.InpaintFullRes = func() *bool { v := false; return &v }()
	opt.InpaintFullResPadding = 32
	resizeByScale := 1.0
	isSetSize := false
	isSetCheckpoints := false
	defaultCheckpoints := shdl.GetDefaultSettingFromUser("sd_model_checkpoint", "", i).(string)
	freeu := global.Features.FreeU && strings.ToLower(shdl.GetDefaultSettingFromUser("freeu", "false", i).(string)) == "true"
	sag := global.Features.Sag && strings.ToLower(shdl.GetDefaultSettingFromUser("sag", "false", i).(string)) == "true"
	for _, v := range cmd.Options {
		switch v.Name {
		case "prompt":
			opt.Prompt = v.StringValue()
		case "negative_prompt":
			opt.NegativePrompt = v.StringValue()
		case "resize_by_scale":
			resizeByScale = v.FloatValue()
		case "height":
			opt.Height = func() *int64 { v := v.IntValue(); return &v }()
			isSetSize = true
		case "width":
			opt.Width = func() *int64 { v := v.IntValue(); return &v }()
			isSetSize = true
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
		case "refiner_checkpoint":
			opt.RefinerCheckpoint = v.StringValue()
		case "refiner_switch_at":
			opt.RefinerSwitchAt = v.FloatValue()
		// img2img options
		case "resize_mode":
			opt.ResizeMode = v.IntValue()
		case "init_image":
			initImage, _ := utils.GetImageBase64(cmd.Resolved.Attachments[v.Value.(string)].URL)
			opt.InitImages = append(opt.InitImages, initImage)
		case "mask":
			opt.Mask, _ = utils.GetImageBase64(cmd.Resolved.Attachments[v.Value.(string)].URL)
		case "denoising_strength":
			opt.DenoisingStrength = func() *float64 { v := v.FloatValue(); return &v }()
		case "mask_blur":
			opt.MaskBlur = v.IntValue()
		case "inpainting_fill":
			opt.InpaintingFill = v.IntValue()
		case "mask_mode":
			opt.InpaintingMaskInvert = v.IntValue()
		case "inpaint_mask_only":
			opt.InpaintFullRes = func() *bool { v := v.BoolValue(); return &v }()
		case "inpaint_full_res_padding":
			opt.InpaintFullResPadding = v.IntValue()

		}
	}
	if !isSetSize {
		width, height, err := utils.GetImageSizeFromBase64(opt.InitImages[0].(string))
		if err == nil {
			opt.Width = func() *int64 { v := int64(float64(width) * resizeByScale); return &v }()
			opt.Height = func() *int64 { v := int64(float64(height) * resizeByScale); return &v }()
		}
	}
	if !isSetCheckpoints && defaultCheckpoints != "" {
		tmpOverrideSettings := opt.OverrideSettings.(map[string]interface{})
		tmpOverrideSettings["sd_model_checkpoint"] = defaultCheckpoints
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
}

func (shdl SlashHandler) Img2imgAction(s *discordgo.Session, i *discordgo.InteractionCreate, opt *intersvc.SdapiV1Img2imgRequest, node *cluster.ClusterNode) {

	img2img := &intersvc.SdapiV1Img2img{RequestItem: opt}
	img2img.Action(node.StableClient)
	if img2img.Error != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func() *string { v := img2img.Error.Error(); return &v }(),
		})
	} else {
		files := make([]*discordgo.File, 0)
		outinfo := img2img.GetResponse().Info
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
		for j, v := range img2img.GetResponse().Images {
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
		if len(files) >= 4 {
			files = files[0:4]
		}
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &context,
			Files:   files,
		})
	}
}

func (shdl SlashHandler) Img2imgCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		option := &intersvc.SdapiV1Img2imgRequest{}
		shdl.RespondStateMessage("Running", s, i)
		node := global.ClusterManager.GetNodeAuto()
		action := func() (map[string]interface{}, error) {
			shdl.Img2imgSetOptions(i.ApplicationCommandData(), option, i)
			shdl.Img2imgAction(s, i, option, node)
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
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: repChoices,
			},
		})
	}
}
