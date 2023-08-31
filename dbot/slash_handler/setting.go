/*
 * @Author: SpenserCai
 * @Date: 2023-08-31 14:59:27
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-31 16:49:13
 * @Description: file content
 */
package slash_handler

import (
	"encoding/json"
	"log"

	"github.com/SpenserCai/sd-webui-discord/user"

	"reflect"

	"github.com/SpenserCai/sd-webui-discord/cluster"
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
				Name:        "height",
				Description: "Height of the output image",
				Type:        discordgo.ApplicationCommandOptionInteger,
				Required:    false,
			},
			{
				Name:        "width",
				Description: "Width of the output image",
				Type:        discordgo.ApplicationCommandOptionInteger,
				Required:    false,
			},
			{
				Name:        "steps",
				Description: "Number of steps to run",
				Type:        discordgo.ApplicationCommandOptionInteger,
				Required:    false,
			},
			{
				Name:        "cfg_scale",
				Description: "Scale of the config",
				Type:        discordgo.ApplicationCommandOptionNumber,
				Required:    false,
			},
			{
				Name:        "negative_prompt",
				Description: "Negative prompt",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
	}
}

func (shdl SlashHandler) SettingSetOptions(dsOpt []*discordgo.ApplicationCommandInteractionDataOption, opt *user.StableConfig) {

	for _, v := range dsOpt {
		switch v.Name {
		case "sd_model_checkpoint":
			opt.Model = v.StringValue()
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
		}
	}
}

func (shdl SlashHandler) SettingAction(s *discordgo.Session, i *discordgo.InteractionCreate, opt *user.StableConfig, node *cluster.ClusterNode) {
	msg, err := shdl.SendStateMessage("Setting", s, i)
	if err != nil {
		log.Println(err)
		return
	}
	userInfo, err := shdl.GetUserInfoWithInteraction(i)
	isEmptyOpt := true
	if err == nil {
		// 判断userInfo是否为nil，如果为nil则说明用户没有注册
		if userInfo == nil {
			s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
				Content: func() *string { v := "Please register first!"; return &v }(),
			})
			return
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
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: func() *string { v := "SETTING ERROR!"; return &v }(),
		})
	} else {
		content := "SETTING SUCCESS!"
		if isEmptyOpt {
			// 把userInfo.StableConfig转换成json字符串，并格式化输出
			stableConfigJson, _ := json.MarshalIndent(userInfo.StableConfig, "", "    ")
			content = "```json\n" + string(stableConfigJson) + "```\n"
		}
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: &content,
		})
	}
}

func (shdl SlashHandler) SettingCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		option := &user.StableConfig{}
		shdl.ReportCommandInfo(s, i)
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
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: repChoices,
			},
		})
	}
}
