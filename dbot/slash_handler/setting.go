/*
 * @Author: SpenserCai
 * @Date: 2023-08-31 14:59:27
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-29 12:14:11
 * @Description: file content
 */
package slash_handler

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/SpenserCai/sd-webui-discord/config"

	"reflect"

	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/dbot/slash_handler/option_values"
	"github.com/SpenserCai/sd-webui-discord/global"

	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) SettingOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "setting",
		Description: "Setting your custom configuration",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "sd_model_checkpoint",
				Description:  "Model checkpoint",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     false,
				Autocomplete: true,
			},
			{
				Name:         "sd_vae",
				Description:  "Vae model (If you don't know what this is, please set it to Automatic)",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     false,
				Autocomplete: true,
			},
			option_values.NegativePrompt(),
			option_values.Height(),
			option_values.Width(),
			option_values.Sampler(),
			option_values.Steps(),
			option_values.CfgScale(),
			option_values.Seed(),
			{
				Name:        "clip_skip",
				Description: "Clip skip",
				Type:        discordgo.ApplicationCommandOptionInteger,
				Required:    false,
				MinValue:    func() *float64 { v := 1.0; return &v }(),
				MaxValue:    12.0,
			},
			option_values.NIter(),
		},
	}
}

func (shdl SlashHandler) SettingSetOptions(dsOpt []*discordgo.ApplicationCommandInteractionDataOption, opt *config.StableConfig) {

	for _, v := range dsOpt {
		switch v.Name {
		case "sd_model_checkpoint":
			opt.Model = v.StringValue()
		case "sd_vae":
			opt.Vae = v.StringValue()
		case "height":
			opt.Height = v.IntValue()
		case "width":
			opt.Width = v.IntValue()
		case "steps":
			opt.Steps = v.IntValue()
		case "cfg_scale":
			opt.CfgScale = v.FloatValue()
		case "negative_prompt":
			opt.NegativePrompt = v.StringValue()
		case "sampler":
			opt.Sampler = v.StringValue()
		case "clip_skip":
			opt.ClipSkip = v.IntValue()
		case "n_iter":
			opt.NIter = v.IntValue()
		case "seed":
			opt.Seed = v.IntValue()
		}
	}
}

func (shdl SlashHandler) SettingAction(s *discordgo.Session, i *discordgo.InteractionCreate, opt *config.StableConfig, node *cluster.ClusterNode) error {
	userInfo, err := shdl.GetUserInfoWithInteraction(i)
	isEmptyOpt := true
	if err == nil {
		// 判断userInfo是否为nil，如果为nil则说明用户没有注册
		if userInfo == nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: func() *string { v := "Please register first!"; return &v }(),
			})
			return errors.New("please register first")
		}
		// 通过反射将opt中的非零值赋值给userInfo.StableConfig
		optVal := reflect.ValueOf(opt).Elem()
		userInfoSdConfig := reflect.ValueOf(&userInfo.StableConfig).Elem()
		for i := 0; i < optVal.NumField(); i++ {
			if optVal.Field(i).Interface() != reflect.Zero(optVal.Field(i).Type()).Interface() {
				isEmptyOpt = false
				userInfoSdConfig.FieldByName(optVal.Type().Field(i).Name).Set(optVal.Field(i))
			}
		}
		err = global.UserCenterSvc.UpdateStableConfig(userInfo)
	}
	if err != nil {
		log.Println(err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func() *string { v := "SETTING ERROR!"; return &v }(),
		})
		return err
	} else {
		content := "SETTING SUCCESS!"
		if isEmptyOpt {
			// 把userInfo.StableConfig转换成json字符串，并格式化输出
			stableConfigJson, _ := json.MarshalIndent(userInfo.StableConfig, "", "    ")
			content = "```json\n" + string(stableConfigJson) + "```\n"
		}
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		return nil
	}
}

func (shdl SlashHandler) SettingCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		option := &config.StableConfig{}
		shdl.RespondStateMessage("Running", s, i)
		node := global.ClusterManager.GetNodeAuto()
		action := func() (map[string]interface{}, error) {
			shdl.SettingSetOptions(i.ApplicationCommandData().Options, option)
			shdl.SettingAction(s, i, option, node)
			return nil, nil
		}
		callback := func() {}
		node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
	case discordgo.InteractionApplicationCommandAutocomplete:
		repChoices := []*discordgo.ApplicationCommandOptionChoice{}
		data := i.ApplicationCommandData()

		for _, opt := range data.Options {
			if opt.Name == "sd_model_checkpoint" && opt.Focused {
				repChoices = shdl.FilterChoice(global.LongDBotChoice["sd_model_checkpoint"], opt)
				continue
			}
			if opt.Name == "sampler" && opt.Focused {
				repChoices = shdl.FilterChoice(global.LongDBotChoice["sampler"], opt)
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
