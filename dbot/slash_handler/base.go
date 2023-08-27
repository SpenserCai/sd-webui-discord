/*
 * @Author: SpenserCai
 * @Date: 2023-08-17 09:52:25
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-27 19:30:05
 * @Description: file content
 */
package slash_handler

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/SpenserCai/sd-webui-discord/utils"
	"github.com/SpenserCai/sd-webui-go/intersvc"

	"github.com/bwmarrin/discordgo"
)

type SlashHandler struct{}

func (shdl SlashHandler) GetCommandStr(i *discordgo.Interaction) string {
	cmd := ""
	// 把命令的名字和参数拼接起来
	cmd += fmt.Sprintf("Command: %s\n", i.ApplicationCommandData().Name)
	for _, v := range i.ApplicationCommandData().Options {
		// 判断v.Value是否是字符串，如果是字符串则判断是否是json字符串，如果是json字符串则格式化输出
		if v.Type == discordgo.ApplicationCommandOptionString && utils.IsJsonString(v.Value.(string)) {
			cmd += fmt.Sprintf("%v: ```json\n%v```\n", utils.FormatCommand(v.Name), v.Value)
			continue
		}
		if v.Type == discordgo.ApplicationCommandOptionAttachment {
			cmd += fmt.Sprintf("%v: %v\n", utils.FormatCommand(v.Name), i.ApplicationCommandData().Resolved.Attachments[v.Value.(string)].URL)
			continue
		}
		cmd += fmt.Sprintf("%v: %v\n", utils.FormatCommand(v.Name), v.Value)
	}
	return cmd
}

func (shdl SlashHandler) ReportCommandInfo(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: shdl.GetCommandStr(i.Interaction),
		},
	})
}

func (shdl SlashHandler) GenerateTaskID(i *discordgo.InteractionCreate) string {
	// 判断是群消息还是私聊消息
	if i.GuildID == "" {
		return fmt.Sprintf("%s_%s_%s", i.Interaction.ID, i.Interaction.User.ID, i.Interaction.User.Username)
	} else {
		return fmt.Sprintf("%s_%s_%s", i.Interaction.ID, i.Interaction.Member.User.ID, i.Interaction.Member.User.Username)
	}
}

func (shdl SlashHandler) SendStateMessage(state string, s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.Message, error) {
	msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("%s...", state),
		Files:   []*discordgo.File{},
	})
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong",
		})
		return nil, err
	}
	return msg, nil
}

func (shdl SlashHandler) GetControlNetScript(jsonStr string) (*intersvc.ControlnetPredictScript, error) {
	script := &intersvc.ControlnetPredictScript{}
	// 把jsonStr转成intersvc.ControlnetScriptArgsItem
	arg := &intersvc.ControlnetPredictArgsItem{}
	err := json.Unmarshal([]byte(jsonStr), arg)
	if err != nil {
		return nil, err
	}
	arg.Image, _ = utils.GetImageBase64(arg.Image)
	script.Args = append(script.Args, *arg)
	return script, nil

}

// Only Step 1,will be change to support every user every setting
func (shdl SlashHandler) GetSdDefaultSetting(key string, defaultValue interface{}) interface{} {
	// 把global.Config.SDWebUi.DefaultSetting转成map[string]interface{}
	defaultSettingMap := make(map[string]interface{})
	defaultSettingJson, _ := json.Marshal(global.Config.SDWebUi.DefaultSetting)
	json.Unmarshal(defaultSettingJson, &defaultSettingMap)

	// 判断key是否被赋值如果没有返回defaultValue
	keyValue, ok := defaultSettingMap[key]
	if ok && utils.IsZeroValue(keyValue) {
		return defaultValue
	} else {
		defaultValueType := reflect.TypeOf(defaultValue)
		if defaultValueType.Kind() == reflect.Ptr {
			defaultValueType = defaultValueType.Elem()
		}
		if keyValue != nil && reflect.TypeOf(keyValue) != defaultValueType {
			convertedValue := reflect.ValueOf(keyValue).Convert(defaultValueType).Interface()
			return convertedValue
		}
		return keyValue
	}
}

func (shdl SlashHandler) FilterChoice(choices []*discordgo.ApplicationCommandOptionChoice, option *discordgo.ApplicationCommandInteractionDataOption) []*discordgo.ApplicationCommandOptionChoice {
	if option.StringValue() == "" {
		// 取得choices的前25个
		if len(choices) > 25 {
			return choices[:25]
		} else {
			return choices
		}
	} else {
		// 如果有输入，就过滤choices
		newChoices := []*discordgo.ApplicationCommandOptionChoice{}
		for _, choice := range choices {
			if strings.Contains(choice.Name, option.StringValue()) {
				newChoices = append(newChoices, choice)
			}
		}
		return newChoices
	}
}
