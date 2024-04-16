/*
 * @Author: ste-art
 * @Date: 2024-04-09 01:50:02
 * @version:
 * @LastEditors: ste-art
 * @LastEditTime: 2024-04-09 01:50:02
 * @Description: Handler for /preset command
 */
package slash_handler

import (
	//  "log"

	"errors"

	"github.com/SpenserCai/sd-webui-discord/config"

	"reflect"

	"github.com/SpenserCai/sd-webui-discord/global"

	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) PresetOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "preset",
		Description: "Apply settings from preset",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "name",
				Description:  "Preset name",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     true,
				Autocomplete: true,
			},
		},
	}
}

func (shdl SlashHandler) PresetChoice() []*discordgo.ApplicationCommandOptionChoice {
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	for _, config := range global.Config.SDWebUi.Presets {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  config.Name,
			Value: config.Name,
		})
	}

	return choices
}

func (shdl SlashHandler) PresetCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		shdl.PresetApplyValues(s, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		repChoices := []*discordgo.ApplicationCommandOptionChoice{}
		data := i.ApplicationCommandData()

		for _, opt := range data.Options {
			if opt.Name == "name" && opt.Focused {
				repChoices = shdl.FilterChoice(global.LongDBotChoice["presets"], opt)
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

func (shdl SlashHandler) PresetApplyValues(s *discordgo.Session, i *discordgo.InteractionCreate) {
	option := &config.StableConfig{}
	shdl.RespondStateMessage("Running", s, i)
	node := global.ClusterManager.GetNodeAuto()
	action := func() (map[string]interface{}, error) {
		name, err := shdl.PresetSetOptions(i, option)
		if err == nil {
			err = shdl.SettingAction(s, i, option, node)
		}
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: func() *string { v := err.Error(); return &v }(),
			})
		} else if name != "" {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: func() *string { v := "Selected preset: " + name; return &v }(),
			})
		}
		return nil, err
	}
	callback := func() {}
	node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
}

func (shdl SlashHandler) PresetSetOptions(i *discordgo.InteractionCreate, opt *config.StableConfig) (string, error) {

	var selectedName string
	for _, v := range i.ApplicationCommandData().Options {
		switch v.Name {
		case "name":
			selectedName = v.StringValue()
		}
	}

	var selectedPreset config.StableConfig
	found := false
	for _, preset := range global.Config.SDWebUi.Presets {
		if preset.Name == selectedName {
			selectedPreset = preset.Settings
			found = true
		}
	}

	if !found {
		return "", errors.New("preset not found")
	}

	optVal := reflect.ValueOf(opt).Elem()
	presetConfig := reflect.ValueOf(&selectedPreset).Elem()
	for i := 0; i < optVal.NumField(); i++ {
		if presetConfig.Field(i).Interface() != reflect.Zero(presetConfig.Field(i).Type()).Interface() {
			optVal.FieldByName(optVal.Type().Field(i).Name).Set(presetConfig.Field(i))
		}
	}

	return selectedName, nil
}
