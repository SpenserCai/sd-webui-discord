/*
 * @Author: ste-art
 * @Date: 2024-04-09 01:50:02
 * @version:
 * @LastEditors: ste-art
 * @LastEditTime: 2024-04-09 01:50:02
 * @Description: Handler for /freeu command
 */
package slash_handler

import (
	"github.com/SpenserCai/sd-webui-discord/config"
	"github.com/SpenserCai/sd-webui-discord/dbot/slash_handler/option_values"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/bwmarrin/discordgo"
)

type FreeUArgsItem struct {
	Enabled bool    `json:"Enabled"`
	B1      float64 `json:"B1"`
	B2      float64 `json:"B2"`
	S1      float64 `json:"S1"`
	S2      float64 `json:"S2"`
}

type FreeUScript struct {
	Args []FreeUArgsItem `json:"args"`
}

func (script *FreeUScript) Set(b1 float64, b2 float64, s1 float64, s2 float64) {
	arg := FreeUArgsItem{Enabled: true, B1: b1, B2: b2, S1: s1, S2: s2}
	args := []FreeUArgsItem{arg}
	script.Args = args
}

func (shdl SlashHandler) FreeUOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "freeu",
		Description: "FreeU",
		Options: []*discordgo.ApplicationCommandOption{
			option_values.Enable(true),
			{
				Name:        "b1",
				Description: "B1",
				Type:        discordgo.ApplicationCommandOptionNumber,
				Required:    false,
				MinValue:    func() *float64 { v := 0.0; return &v }(),
				MaxValue:    2.0,
			},
			{
				Name:        "b2",
				Description: "B2",
				Type:        discordgo.ApplicationCommandOptionNumber,
				Required:    false,
				MinValue:    func() *float64 { v := 0.0; return &v }(),
				MaxValue:    2.0,
			},
			{
				Name:        "s1",
				Description: "S1",
				Type:        discordgo.ApplicationCommandOptionNumber,
				Required:    false,
				MinValue:    func() *float64 { v := 0.0; return &v }(),
				MaxValue:    4.0,
			},
			{
				Name:        "s2",
				Description: "S2",
				Type:        discordgo.ApplicationCommandOptionNumber,
				Required:    false,
				MinValue:    func() *float64 { v := 0.0; return &v }(),
				MaxValue:    4.0,
			},
		},
	}
}

func (shdl SlashHandler) FreeUSetOptions(dsOpt []*discordgo.ApplicationCommandInteractionDataOption, opt *config.StableConfig) {

	for _, v := range dsOpt {
		switch v.Name {
		case "enable":
			opt.FreeU = v.StringValue()
		case "b1":
			opt.FreeUB1 = v.FloatValue()
		case "b2":
			opt.FreeUB2 = v.FloatValue()
		case "s1":
			opt.FreeUS1 = v.FloatValue()
		case "s2":
			opt.FreeUS1 = v.FloatValue()
		}
	}
}

func (shdl SlashHandler) FreeuCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		option := &config.StableConfig{}
		shdl.RespondStateMessage("Running", s, i)
		node := global.ClusterManager.GetNodeAuto()
		action := func() (map[string]interface{}, error) {
			shdl.FreeUSetOptions(i.ApplicationCommandData().Options, option)
			err := shdl.SettingAction(s, i, option, node)
			if err != nil {
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: func() *string { v := err.Error(); return &v }(),
				})
			} else {
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: func() *string { v := "FreeU setting updated"; return &v }(),
				})
			}
			return nil, err
		}
		callback := func() {}
		node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
	}
}
