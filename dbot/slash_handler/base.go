/*
 * @Author: SpenserCai
 * @Date: 2023-08-17 09:52:25
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-14 13:45:26
 * @Description: file content
 */
package slash_handler

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/SpenserCai/sd-webui-discord/user"
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

func (shdl SlashHandler) RespondStateMessage(state string, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s...", state),
		},
	})
}

func (shdl SlashHandler) RespondStateMessageWithFlag(state string, s *discordgo.Session, i *discordgo.InteractionCreate, flags discordgo.MessageFlags) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s...", state),
			Flags:   flags,
		},
	})
}

func (shdl SlashHandler) SendStateMessageWithFlag(state string, s *discordgo.Session, i *discordgo.InteractionCreate, flags discordgo.MessageFlags) (*discordgo.Message, error) {
	msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("%s...", state),
		Flags:   flags,
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

func (shdl SlashHandler) SendTextInteractionRespondWithFlag(msg string, s *discordgo.Session, i *discordgo.InteractionCreate, flags discordgo.MessageFlags) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   flags,
		},
	})
}

func (shdl SlashHandler) GetControlNetScript(jsonStr string) (*intersvc.ControlnetPredictScript, error) {
	jsonStr = fmt.Sprintf(`{"args": [%s]}`, jsonStr)
	script := &intersvc.ControlnetPredictScript{}
	// 把jsonStr转成intersvc.ControlnetScriptArgsItem
	// arg := &intersvc.ControlnetPredictArgsItem{}
	err := json.Unmarshal([]byte(jsonStr), script)
	if err != nil {
		return nil, err
	}
	// arg.Image, _ = utils.GetImageBase64(arg.Image)
	// script.Args = append(script.Args, *arg)
	return script, nil

}

func (shdl SlashHandler) GetDefaultSettingFromUser(key string, defaultValue interface{}, i *discordgo.InteractionCreate) interface{} {
	if global.Config.UserCenter.Enable {
		userInfo, err := shdl.GetUserInfoWithInteraction(i)
		if userInfo != nil && err == nil {
			value, err := global.UserCenterSvc.GetUserStableConfigItem(userInfo.Id, key, shdl.GetSdDefaultSetting(key, defaultValue))
			if err == nil {
				return value
			}
		}
	}
	return shdl.GetSdDefaultSetting(key, defaultValue)
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
			if strings.Contains(strings.ToLower(choice.Name), strings.ToLower(option.StringValue())) {
				newChoices = append(newChoices, choice)
			}
		}
		if len(newChoices) > 25 {
			return newChoices[:25]
		} else {
			return newChoices
		}
	}
}

func (shdl SlashHandler) ConvertCommandOptionChoiceToMenuOption(choices []*discordgo.ApplicationCommandOptionChoice, default_v string) []discordgo.SelectMenuOption {
	menuOption := []discordgo.SelectMenuOption{}
	for _, choice := range choices {
		selectMenueOption := discordgo.SelectMenuOption{
			Label: choice.Name,
			Value: choice.Value.(string),
		}
		if choice.Value.(string) == default_v {
			selectMenueOption.Default = true
			// 把default的值放到第一个
			menuOption = append([]discordgo.SelectMenuOption{selectMenueOption}, menuOption...)
		} else {
			menuOption = append(menuOption, selectMenueOption)
		}
	}
	// 如果超过25个，就只取前25个
	if len(menuOption) > 25 {
		menuOption = menuOption[:25]
	}
	return menuOption
}

func (shdl SlashHandler) GetUserInfoWithInteraction(i *discordgo.InteractionCreate) (*user.UserInfo, error) {
	// 判断是群消息还是私聊消息
	return global.UserCenterSvc.GetUserInfo(shdl.GetDiscordUserId(i))
}

func (shdl SlashHandler) GetDiscordUserId(i *discordgo.InteractionCreate) string {
	// 判断是群消息还是私聊消息
	if i.GuildID == "" {
		return i.Interaction.User.ID
	} else {
		return i.Interaction.Member.User.ID
	}
}

// 判断是否是管理员
func (shdl SlashHandler) IsAdmin(i *discordgo.InteractionCreate) bool {
	userInfo, err := global.UserCenterSvc.GetUserInfo(shdl.GetDiscordUserId(i))
	if err != nil {
		return false
	}
	if strings.Contains(userInfo.Roles, "admin") {
		return true
	}
	return false
}

func (shdl SlashHandler) SetDiscordUserId(i *discordgo.InteractionCreate, userId string) {
	// 判断是群消息还是私聊消息
	if i.GuildID == "" {
		i.Interaction.User.ID = userId
	} else {
		i.Interaction.Member.User.ID = userId
	}
}

func (shdl SlashHandler) ConvertInteractionToUserInfo(i *discordgo.InteractionCreate) *user.UserInfo {
	// 判断是群消息还是私聊消息
	if i.GuildID == "" {
		return &user.UserInfo{
			Id:   i.Interaction.User.ID,
			Name: i.Interaction.User.Username,
			Avatar: func() string {
				if i.Interaction.User.Avatar == "" {
					return ""
				} else {
					return "https://cdn.discordapp.com/avatars/" + i.Interaction.User.ID + "/" + i.Interaction.User.Avatar + ".png"
				}
			}(),
		}
	} else {
		return &user.UserInfo{
			Id:   i.Interaction.Member.User.ID,
			Name: i.Interaction.Member.User.Username,
			Avatar: func() string {
				if i.Interaction.Member.User.Avatar == "" {
					return ""
				} else {
					return "https://cdn.discordapp.com/avatars/" + i.Interaction.Member.User.ID + "/" + i.Interaction.Member.User.Avatar + ".png"
				}
			}(),
		}
	}
}

// MessageEmbed模板
func (shdl SlashHandler) MessageEmbedTemplate() *discordgo.MessageEmbed {
	bName := func() string {
		if global.Config.Discord.BotName == "" {
			return "SD-WEBUI-BOT"
		} else {
			return global.Config.Discord.BotName
		}
	}()
	bAvatar := func() string {
		if global.Config.Discord.BotAvatar == "" {
			return "https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/logo.png"
		} else {
			return global.Config.Discord.BotAvatar
		}
	}()
	bUrl := func() string {
		if global.Config.Discord.BotUrl == "" {
			return "https://github.com/SpenserCai/sd-webui-discord"
		} else {
			return global.Config.Discord.BotUrl
		}
	}()
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    bName,
			IconURL: bAvatar,
			URL:     bUrl,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Powered by sd-webui-discord",
			IconURL: "https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/logo.png",
		},
	}
	return embed
}

func (shdl SlashHandler) SetHistory(command string, messageId string, i *discordgo.InteractionCreate, opt any) {
	if global.Config.UserCenter.Enable {
		optJson, _ := json.Marshal(opt)
		userId := shdl.GetDiscordUserId(i)
		global.UserCenterSvc.WriteUserHistory(messageId, userId, command, string(optJson))
	}
}

func (shdl SlashHandler) SetHistoryImages(messageId string, i *discordgo.InteractionCreate, images []string, imagesBlurHash []string) {
	if global.Config.UserCenter.Enable {
		userId := shdl.GetDiscordUserId(i)
		// 把imagesBlurHash转成json
		imagesBlurHashJson, _ := json.Marshal(imagesBlurHash)
		global.UserCenterSvc.WriteUserHistoryImages(messageId, userId, strings.Join(images, ","), string(imagesBlurHashJson))
	}
}

func (shdl SlashHandler) DeleteHistory(messageId string, userId string) {
	if global.Config.UserCenter.Enable {
		global.UserCenterSvc.DeleteUserHistory(messageId, userId)
	}
}

func (shdl SlashHandler) GetHistory(command string, messageId string, opt any) error {
	if global.Config.UserCenter.Enable {
		history, err := global.UserCenterSvc.GetUserHistoryOptWithMessageId(messageId, command)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(history), opt)
		if err != nil {
			return err
		}
	}
	return nil
}

// 获取用户Choice
func (shdl SlashHandler) GetUserCommandOptionChoice(i *discordgo.InteractionCreate) []*discordgo.ApplicationCommandOptionChoice {
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	userList, err := global.UserCenterSvc.GetUserList(shdl.GetDiscordUserId(i))
	if err != nil {
		return choices
	}
	for _, v := range userList {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  v.Name,
			Value: v.Id,
		})
	}
	return choices
}

// 获取带有cmmd和用户ID的CustomID
func (shdl SlashHandler) GetDiscordUserCustomId(command string, customId string, i *discordgo.InteractionCreate) string {
	return shdl.GetDiscordUserCustomIdWithUserId(command, customId, shdl.GetDiscordUserId(i))
}

func (shdl SlashHandler) GetDiscordUserCustomIdWithUserId(command string, customId string, userId string) string {
	return fmt.Sprintf("%s|%s|%s", command, customId, userId)
}

// 获取图片列表的BlurHash
func (shdl SlashHandler) GetBase64ImageListBlurHash(imageList []string) []string {
	blurHashList := []string{}
	for _, image := range imageList {
		blurHash, err := utils.GetImageBlurHashFromBase64(image)
		if err != nil {
			blurHashList = append(blurHashList, "")
		} else {
			blurHashList = append(blurHashList, blurHash)
		}
	}
	return blurHashList
}

func (shdl SlashHandler) GetSctipts() *[]string {

	scriptsSvc := &intersvc.SdapiV1Scripts{}
	scriptsSvc.Action(global.ClusterManager.GetNodeAuto().StableClient)
	if scriptsSvc.Error != nil {
		return &[]string{}
	}
	scripts := scriptsSvc.GetResponse()
	scriptsResult := make([]string, len(scripts.Txt2img))
	for i, script := range scripts.Txt2img {
		scriptsResult[i] = script.(string)
	}
	return &scriptsResult
}
