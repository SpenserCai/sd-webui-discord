/*
 * @Author: SpenserCai
 * @Date: 2023-09-21 16:27:24
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-29 11:58:06
 * @Description: file content
 */
package slash_handler

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/config"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/SpenserCai/sd-webui-discord/user"

	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) SettingUiOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "setting_ui",
		Description: "Setting with Ui",
		Options:     []*discordgo.ApplicationCommandOption{},
	}
}

func (shdl SlashHandler) SettingUiSetOptions(dsOpt []*discordgo.ApplicationCommandInteractionDataOption, opt *config.StableConfig) {

	for _, v := range dsOpt {
		switch v.Name {

		}
	}
}

func (shdl SlashHandler) BuildSettingUiComponent(opt *config.StableConfig, i *discordgo.InteractionCreate) *[]discordgo.MessageComponent {
	adminRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				CustomID: shdl.GetDiscordUserCustomId("setting_ui", "ban", i),
				Label:    "Ban",
				Style:    discordgo.DangerButton,
				Emoji:    &discordgo.ComponentEmoji{Name: "🔒"},
			},
			discordgo.Button{
				CustomID: shdl.GetDiscordUserCustomId("setting_ui", "unban", i),
				Label:    "Unban",
				Style:    discordgo.SuccessButton,
				Emoji:    &discordgo.ComponentEmoji{Name: "🔓"},
			},
		},
	}
	component := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    shdl.GetDiscordUserCustomId("setting_ui", "sd_model_checkpoint", i),
					Placeholder: "Choose a model checkpoint",
					Options:     shdl.ConvertCommandOptionChoiceToMenuOption(global.LongDBotChoice["sd_model_checkpoint"], opt.Model),
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    shdl.GetDiscordUserCustomId("setting_ui", "sd_vae", i),
					Placeholder: "Choose a vae model",
					Options:     shdl.ConvertCommandOptionChoiceToMenuOption(global.LongDBotChoice["sd_vae"], opt.Vae),
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    shdl.GetDiscordUserCustomId("setting_ui", "sampler", i),
					Placeholder: "Choose a sampler",
					Options:     shdl.ConvertCommandOptionChoiceToMenuOption(global.LongDBotChoice["sampler"], opt.Sampler),
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: shdl.GetDiscordUserCustomId("setting_ui", "set_size", i),
					Label:    "Set Image size",
					Style:    discordgo.PrimaryButton,
					Emoji:    &discordgo.ComponentEmoji{Name: "📐"},
				},
				discordgo.Button{
					CustomID: shdl.GetDiscordUserCustomId("setting_ui", "set_steps", i),
					Label:    "Set Steps",
					Style:    discordgo.PrimaryButton,
					Emoji:    &discordgo.ComponentEmoji{Name: "🔢"},
				},
				discordgo.Button{
					CustomID: shdl.GetDiscordUserCustomId("setting_ui", "set_cfg_scale", i),
					Label:    "Set Cfg Scale",
					Style:    discordgo.PrimaryButton,
					Emoji:    &discordgo.ComponentEmoji{Name: "📏"},
				},
				discordgo.Button{
					CustomID: shdl.GetDiscordUserCustomId("setting_ui", "set_negative_prompt", i),
					Label:    "Set Negative Prompt",
					Style:    discordgo.PrimaryButton,
					Emoji:    &discordgo.ComponentEmoji{Name: "🚫"},
				},
				discordgo.Button{
					CustomID: shdl.GetDiscordUserCustomId("setting_ui", "set_clip_skip", i),
					Label:    "Clip Skip",
					Style:    discordgo.PrimaryButton,
					Emoji:    &discordgo.ComponentEmoji{Name: "📌"},
				},
			},
		},
	}
	if i.ApplicationCommandData().Name == "user_info" {
		component = append(component, adminRow)
	}
	return &component
}

func (shdl SlashHandler) SettingUiAction(s *discordgo.Session, i *discordgo.InteractionCreate, node *cluster.ClusterNode) {
	userInfo, err := shdl.GetUserInfoWithInteraction(i)
	if err == nil {
		// 判断userInfo是否为nil，如果为nil则说明用户没有注册
		if userInfo == nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: func() *string { v := "Please register first!"; return &v }(),
			})
			return
		}
		component := shdl.BuildSettingUiComponent(&userInfo.StableConfig, i)
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func() *string {
				v := fmt.Sprintf("**Setting GUI - [%s]**\nIf you need more options, please use the `/setting` command", userInfo.Name)
				return &v
			}(),
			Components: component,
		})
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Println(err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func() *string { v := "SETTING ERROR!"; return &v }(),
		})
	}

}

func (shdl SlashHandler) SettingUiComponentHandler(s *discordgo.Session, i *discordgo.InteractionCreate, userInfo *user.UserInfo) (isFinish bool) {
	isSenderAdmin, _ := global.UserCenterSvc.IsAdmin(shdl.GetDiscordUserId(i))
	customIDList := strings.Split(i.MessageComponentData().CustomID, "|")
	cmd := fmt.Sprintf("%s|%s", customIDList[0], customIDList[1])
	if len(customIDList) == 3 {
		tmpUserInfo, err := global.UserCenterSvc.GetUserInfo(customIDList[2])
		if err == nil {
			*userInfo = *tmpUserInfo
		}
	}
	switch cmd {
	case "setting_ui|sd_model_checkpoint":
		userInfo.StableConfig.Model = i.MessageComponentData().Values[0]
		return false
	case "setting_ui|sd_vae":
		userInfo.StableConfig.Vae = i.MessageComponentData().Values[0]
		return false
	case "setting_ui|sampler":
		userInfo.StableConfig.Sampler = i.MessageComponentData().Values[0]
		return false
	case "setting_ui|ban":
		if isSenderAdmin {
			global.UserCenterSvc.BanUser(userInfo.Id)
			shdl.SendTextInteractionRespondWithFlag(fmt.Sprintf("Ban **`%s`** Success", userInfo.Name), s, i, discordgo.MessageFlagsEphemeral)
		} else {
			shdl.SendTextInteractionRespondWithFlag("You are not admin!", s, i, discordgo.MessageFlagsEphemeral)
		}
		return true
	case "setting_ui|unban":
		if isSenderAdmin {
			global.UserCenterSvc.UnBanUser(userInfo.Id)
			shdl.SendTextInteractionRespondWithFlag(fmt.Sprintf("UnBan **`%s`** Success", userInfo.Name), s, i, discordgo.MessageFlagsEphemeral)
		} else {
			shdl.SendTextInteractionRespondWithFlag("You are not admin!", s, i, discordgo.MessageFlagsEphemeral)
		}
		return true
	// 显示图片大小设置窗口
	case "setting_ui|set_size":
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: shdl.GetDiscordUserCustomIdWithUserId("setting_ui", "set_size_modal", userInfo.Id),
				Title:    "Set Image Size",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "setting_ui|height",
								Label:       "Height",
								Style:       discordgo.TextInputShort,
								Placeholder: "Set image height",
								Value:       fmt.Sprintf("%d", userInfo.StableConfig.Height),
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "setting_ui|width",
								Label:       "Width",
								Style:       discordgo.TextInputShort,
								Placeholder: "Set image width",
								Value:       fmt.Sprintf("%d", userInfo.StableConfig.Width),
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Println(err)
		}
	case "setting_ui|set_steps":
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: shdl.GetDiscordUserCustomIdWithUserId("setting_ui", "set_steps_modal", userInfo.Id),
				Title:    "Set Steps",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "setting_ui|steps",
								Label:       "Steps",
								Style:       discordgo.TextInputShort,
								Placeholder: "Set steps",
								Value:       fmt.Sprintf("%d", userInfo.StableConfig.Steps),
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Println(err)
		}
	case "setting_ui|set_cfg_scale":
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: shdl.GetDiscordUserCustomIdWithUserId("setting_ui", "set_cfg_scale_modal", userInfo.Id),
				Title:    "Set Cfg Scale",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "setting_ui|cfg_scale",
								Label:       "Cfg Scale",
								Style:       discordgo.TextInputShort,
								Placeholder: "Set cfg scale",
								// 小数点后保留两位
								Value: fmt.Sprintf("%.2f", userInfo.StableConfig.CfgScale),
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Println(err)
		}
	case "setting_ui|set_negative_prompt":
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: shdl.GetDiscordUserCustomIdWithUserId("setting_ui", "set_negative_prompt_modal", userInfo.Id),
				Title:    "Set Negative Prompt",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "setting_ui|negative_prompt",
								Label:       "Negative Prompt",
								Style:       discordgo.TextInputParagraph,
								Placeholder: "Set negative prompt",
								MinLength:   0,
								MaxLength:   200,
								Value:       userInfo.StableConfig.NegativePrompt,
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Println(err)
		}
	case "setting_ui|set_clip_skip":
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: shdl.GetDiscordUserCustomIdWithUserId("setting_ui", "clip_skip_modal", userInfo.Id),
				Title:    "Set Clip Skip",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "setting_ui|clip_skip",
								Label:       "Clip Skip",
								Style:       discordgo.TextInputShort,
								Placeholder: "Set clip skip",
								Value:       fmt.Sprintf("%d", userInfo.StableConfig.ClipSkip),
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Println(err)
		}
	}
	return true
}

func (shdl SlashHandler) SettingUiModalSubmitHander(s *discordgo.Session, i *discordgo.InteractionCreate, userInfo *user.UserInfo) (isSuccess bool) {
	customIDList := strings.Split(i.ModalSubmitData().CustomID, "|")
	cmd := fmt.Sprintf("%s|%s", customIDList[0], customIDList[1])
	if len(customIDList) == 3 {
		tmpUserInfo, err := global.UserCenterSvc.GetUserInfo(customIDList[2])
		if err == nil {
			*userInfo = *tmpUserInfo
		}
	}
	switch cmd {
	case "setting_ui|set_size_modal":
		modal_data := i.ModalSubmitData()
		tmpHeight, formathErr := strconv.ParseInt(modal_data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value, 10, 64)
		tmpWidth, formatwErr := strconv.ParseInt(modal_data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value, 10, 64)
		if formathErr != nil || formatwErr != nil {
			shdl.SendTextInteractionRespondWithFlag("Format Error", s, i, discordgo.MessageFlagsEphemeral)
			return false
		}
		// 判断Height和Width是否都>=64,<=2048
		if tmpHeight < 64 || tmpHeight > 2048 || tmpWidth < 64 || tmpWidth > 2048 {
			shdl.SendTextInteractionRespondWithFlag("Height and Width must be >=64 and <=2048", s, i, discordgo.MessageFlagsEphemeral)
			return false
		}
		userInfo.StableConfig.Height = tmpHeight
		userInfo.StableConfig.Width = tmpWidth
	case "setting_ui|set_steps_modal":
		modal_data := i.ModalSubmitData()
		tmpSteps, formatErr := strconv.ParseInt(modal_data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value, 10, 64)
		if formatErr != nil {
			shdl.SendTextInteractionRespondWithFlag("Format Error", s, i, discordgo.MessageFlagsEphemeral)
			return false
		}
		// 判断Steps是否>=15,<=100
		if tmpSteps < 15 || tmpSteps > 100 {
			shdl.SendTextInteractionRespondWithFlag("Steps must be >=15 and <=100", s, i, discordgo.MessageFlagsEphemeral)
			return false
		}
		userInfo.StableConfig.Steps = tmpSteps
	case "setting_ui|set_cfg_scale_modal":
		modal_data := i.ModalSubmitData()
		tmpCfgScale, formatErr := strconv.ParseFloat(modal_data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value, 64)
		if formatErr != nil {
			shdl.SendTextInteractionRespondWithFlag("Format Error", s, i, discordgo.MessageFlagsEphemeral)
			return false
		}
		// 判断CfgScale是否>=0.1,<=1.0
		if tmpCfgScale < 1.0 || tmpCfgScale > 30.0 {
			shdl.SendTextInteractionRespondWithFlag("CfgScale must be >=1.0 and <=30.0", s, i, discordgo.MessageFlagsEphemeral)
			return false
		}
		userInfo.StableConfig.CfgScale = tmpCfgScale
	case "setting_ui|set_negative_prompt_modal":
		modal_data := i.ModalSubmitData()
		userInfo.StableConfig.NegativePrompt = modal_data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	case "setting_ui|clip_skip_modal":
		modal_data := i.ModalSubmitData()
		tmpClipSkip, formatErr := strconv.ParseInt(modal_data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value, 10, 64)
		if formatErr != nil {
			shdl.SendTextInteractionRespondWithFlag("Format Error", s, i, discordgo.MessageFlagsEphemeral)
			return false
		}
		// 范围1-12
		if tmpClipSkip < 1 || tmpClipSkip > 12 {
			shdl.SendTextInteractionRespondWithFlag("ClipSkip must be >=1 and <=12", s, i, discordgo.MessageFlagsEphemeral)
			return false
		}
		userInfo.StableConfig.ClipSkip = tmpClipSkip
	}
	return true
}

func (shdl SlashHandler) SettingUiCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	node := global.ClusterManager.GetNodeAuto()
	userInfo, err := shdl.GetUserInfoWithInteraction(i)
	if err != nil {
		log.Println(err)
		return
	}
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		shdl.RespondStateMessageWithFlag("Running", s, i, discordgo.MessageFlagsEphemeral)
		action := func() (map[string]interface{}, error) {
			// shdl.SettingUiSetOptions(i.ApplicationCommandData().Options, option)
			shdl.SettingUiAction(s, i, node)
			return nil, nil
		}
		callback := func() {}
		node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
		return
	case discordgo.InteractionMessageComponent:
		isFinish := shdl.SettingUiComponentHandler(s, i, userInfo)
		if isFinish {
			return
		}
	case discordgo.InteractionModalSubmit:
		isSuccess := shdl.SettingUiModalSubmitHander(s, i, userInfo)
		if !isSuccess {
			return
		}
	}
	err = global.UserCenterSvc.UpdateStableConfig(userInfo)
	if err == nil {
		shdl.SendTextInteractionRespondWithFlag("", s, i, discordgo.MessageFlagsEphemeral)
	} else {
		sendErr := shdl.SendTextInteractionRespondWithFlag("Setting Error: "+err.Error(), s, i, discordgo.MessageFlagsEphemeral)
		if sendErr != nil {
			log.Println(sendErr)
		}
	}
}
