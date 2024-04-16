/*
 * @Author: ste-art
 * @Date: 2024-04-09 01:50:02
 * @version:
 * @LastEditors: ste-art
 * @LastEditTime: 2024-04-09 01:50:02
 * @Description: Handler for /sag command (Self Attention Guidance)
 */
package slash_handler

import (
	"github.com/SpenserCai/sd-webui-discord/config"
	"github.com/SpenserCai/sd-webui-discord/dbot/slash_handler/option_values"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/bwmarrin/discordgo"
)

type SagArgsItem struct {
	Enabled   bool    `json:"Enabled"`
	Scale     float64 `json:"Scale"`
	BlurSigma float64 `json:"Blur Sigma"`
}

type SagScript struct {
	Args []SagArgsItem `json:"args"`
}

func (script *SagScript) Set(scale float64, blurSigma float64) {
	arg := SagArgsItem{Scale: scale, BlurSigma: blurSigma, Enabled: true}
	args := []SagArgsItem{arg}
	script.Args = args
}

func (shdl SlashHandler) SagOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "sag",
		Description: "Self Attention Guidance",
		Options: []*discordgo.ApplicationCommandOption{
			option_values.Enable(true),
			{
				Name:        "scale",
				Description: "Scale",
				Type:        discordgo.ApplicationCommandOptionNumber,
				Required:    false,
				MinValue:    func() *float64 { v := -2.0; return &v }(),
				MaxValue:    5.0,
			},
			{
				Name:        "blur_sigma",
				Description: "Blur Sigma",
				Type:        discordgo.ApplicationCommandOptionNumber,
				Required:    false,
				MinValue:    func() *float64 { v := 0.0; return &v }(),
				MaxValue:    10.0,
			},
		},
	}
}

func (shdl SlashHandler) SagSetOptions(dsOpt []*discordgo.ApplicationCommandInteractionDataOption, opt *config.StableConfig) {

	for _, v := range dsOpt {
		switch v.Name {
		case "enable":
			opt.Sag = v.StringValue()
		case "scale":
			opt.SagScale = v.FloatValue()
		case "blur_sigma":
			opt.SagBlurSigma = v.FloatValue()
		}
	}
}

func (shdl SlashHandler) SagCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		option := &config.StableConfig{}
		shdl.RespondStateMessage("Running", s, i)
		node := global.ClusterManager.GetNodeAuto()
		action := func() (map[string]interface{}, error) {
			shdl.SagSetOptions(i.ApplicationCommandData().Options, option)
			err := shdl.SettingAction(s, i, option, node)
			if err != nil {
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: func() *string { v := err.Error(); return &v }(),
				})
			} else {
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: func() *string { v := "Self Attention Guidance setting updated"; return &v }(),
				})
			}
			return nil, err
		}
		callback := func() {}
		node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
	}
}
