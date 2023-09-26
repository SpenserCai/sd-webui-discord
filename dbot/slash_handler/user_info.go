/*
 * @Author: SpenserCai
 * @Date: 2023-09-26 23:24:02
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-27 00:58:39
 * @Description: file content
 */
package slash_handler

import (
	"log"

	"github.com/SpenserCai/sd-webui-discord/cluster"

	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) UserInfoOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "user_info",
		Description: "User Manager",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "user",
				Description:  "Select a user",
				Required:     true,
				Autocomplete: true,
			},
		},
	}
}

func (shdl SlashHandler) UserInfoSetOptions(dsOpt []*discordgo.ApplicationCommandInteractionDataOption, s *discordgo.Session, uId *string) {

	for _, v := range dsOpt {
		switch v.Name {
		case "user":
			*uId = v.StringValue()
		}
	}
}

// TODO: 由于所有的组件和modal的id都是固定了命令的，所以需要实现id中的命令是动态的，同时由于参数中的用户id无法传递到按钮上所以需要在最开始设置按钮的时候就把id带到按钮的组件id中
func (shdl SlashHandler) UserInfoAction(s *discordgo.Session, i *discordgo.InteractionCreate, uId *string, node *cluster.ClusterNode) {
	shdl.SetDiscordUserId(i, *uId)
	shdl.SettingUiCommandHandler(s, i)

}

func (shdl SlashHandler) UserInfoCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand, discordgo.InteractionMessageComponent, discordgo.InteractionModalSubmit:
		option := func() *string { v := ""; return &v }()
		shdl.UserInfoSetOptions(i.ApplicationCommandData().Options, s, option)
		log.Println(*option)
		shdl.UserInfoAction(s, i, option, nil)
	case discordgo.InteractionApplicationCommandAutocomplete:
		repChoices := []*discordgo.ApplicationCommandOptionChoice{}
		data := i.ApplicationCommandData()
		for _, opt := range data.Options {
			if opt.Name == "user" && opt.Focused {
				repChoices = shdl.FilterChoice(shdl.GetUserCommandOptionChoice(i), opt)
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
