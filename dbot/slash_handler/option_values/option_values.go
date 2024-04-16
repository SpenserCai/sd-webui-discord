package option_values

import (
	"fmt"

	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/bwmarrin/discordgo"
)

func ValueOrDefault(value any, defaultvalue string) string {
	valueString := fmt.Sprint(value)
	if valueString == "0" || valueString == "" {
		return defaultvalue
	} else {
		return valueString
	}
}

func NegativePrompt() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "negative_prompt",
		Description: "Negative prompt text. Default: " + ValueOrDefault(global.Config.SDWebUi.DefaultSetting.NegativePrompt, ""),
		Required:    false,
	}
}

func Height() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "height",
		Description: "Height of the generated image. Default: " + ValueOrDefault(global.Config.SDWebUi.DefaultSetting.Height, "512"),
		MinValue:    func() *float64 { v := 64.0; return &v }(),
		MaxValue:    2048.0,
		Required:    false,
	}
}

func Width() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "width",
		Description: "Width of the generated image. Default: " + ValueOrDefault(global.Config.SDWebUi.DefaultSetting.Width, "512"),
		MinValue:    func() *float64 { v := 64.0; return &v }(),
		MaxValue:    2048.0,
		Required:    false,
	}
}

func Sampler() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:         discordgo.ApplicationCommandOptionString,
		Name:         "sampler",
		Description:  "Sampler of the generated image. Default: " + ValueOrDefault(global.Config.SDWebUi.DefaultSetting.Sampler, "Euler"),
		Required:     false,
		Autocomplete: true,
	}
}

func Steps() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "steps",
		Description: "Steps of the generated image. Default: " + ValueOrDefault(global.Config.SDWebUi.DefaultSetting.Steps, "30"),
		Required:    false,
	}
}

func CfgScale() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionNumber,
		Name:        "cfg_scale",
		Description: "Cfg scale of the generated image. Default: " + ValueOrDefault(global.Config.SDWebUi.DefaultSetting.CfgScale, "7"),
		MinValue:    func() *float64 { v := 1.0; return &v }(),
		MaxValue:    30.0,
		Required:    false,
	}
}

func Seed() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "seed",
		Description: "Seed of the generated image. Default: " + ValueOrDefault(global.Config.SDWebUi.DefaultSetting.Seed, "-1"),
		Required:    false,
	}
}

func NIter() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "n_iter",
		Description: "Number of iterations. Default: " + ValueOrDefault(global.Config.SDWebUi.DefaultSetting.NIter, "1"),
		Required:    false,
		MinValue:    func() *float64 { v := 1.0; return &v }(),
		MaxValue:    4.0,
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "1",
				Value: 1,
			},
			{
				Name:  "2",
				Value: 2,
			},
			{
				Name:  "4",
				Value: 4,
			},
		},
	}
}

func Enable(required bool) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "enable",
		Description: "Enable feature",
		Type:        discordgo.ApplicationCommandOptionString,
		Required:    required,
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "True",
				Value: "true",
			},
			{
				Name:  "False",
				Value: "false",
			},
		},
	}
}
