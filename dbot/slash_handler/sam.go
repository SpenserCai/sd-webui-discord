/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 22:27:32
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-18 01:25:55
 * @Description: file content
 */
package slash_handler

import (
	"fmt"
	"log"
	"sd-webui-discord/cluster"
	"sd-webui-discord/global"
	"sd-webui-discord/utils"

	"github.com/SpenserCai/sd-webui-go/intersvc"
	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) samModelChoice() []*discordgo.ApplicationCommandOptionChoice {
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	modesvc := &intersvc.SamSamModel{}
	modesvc.Action(global.ClusterManager.GetNodeAuto().StableClient)
	if modesvc.Error != nil {
		log.Println(modesvc.Error)
		return choices
	}
	models := modesvc.GetResponse()
	for _, model := range *models {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  model,
			Value: model,
		})
	}
	return choices
}

func (shdl SlashHandler) dinoModelChoice() []*discordgo.ApplicationCommandOptionChoice {
	return []*discordgo.ApplicationCommandOptionChoice{
		{
			Name:  "GroundingDINO_SwinT_OGC (694MB)",
			Value: "GroundingDINO_SwinT_OGC (694MB)",
		},
		{
			Name:  "GroundingDINO_SwinB (938MB)",
			Value: "GroundingDINO_SwinB (938MB)",
		},
	}
}

func (shdl SlashHandler) SamOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "sam",
		Description: "Segment Anythin with prompt",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "image_url",
				Description: "The url of the image",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "The prompt of the image",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "model",
				Description: "Choice sam model",
				Required:    true,
				Choices:     shdl.samModelChoice(),
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "dion_model",
				Description: "Choice dion model",
				Required:    false,
				Choices:     shdl.dinoModelChoice(),
			},
		},
	}
}

func (shdl SlashHandler) SamSetOptions(dsOpt []*discordgo.ApplicationCommandInteractionDataOption, opt *intersvc.SamSamPredictRequest) {
	opt.DinoEnabled = func() *bool { v := true; return &v }()
	opt.DinoBoxThreshold = func() *float64 { v := 0.3; return &v }()
	opt.DinoPreviewCheckbox = func() *bool { v := false; return &v }()
	opt.DinoPreviewBoxesSelection = []int64{0}
	opt.SamNegativePoints = [][]float64{}
	opt.SamPositivePoints = [][]float64{}
	for _, v := range dsOpt {
		switch v.Name {
		case "image_url":
			image, _ := utils.GetImageBase64(v.StringValue())
			opt.InputImage = &image
		case "prompt":
			opt.DinoTextPrompt = v.StringValue()
		case "model":
			opt.SamModelName = func() *string { v := v.StringValue(); return &v }()
		case "dion_model":
			opt.DinoModelName = func() *string { v := v.StringValue(); return &v }()
		}
	}

	// default value
	if opt.DinoModelName == nil {
		opt.DinoModelName = func() *string { v := shdl.dinoModelChoice()[0].Value.(string); return &v }()
	}

}

func (shdl SlashHandler) SamAction(s *discordgo.Session, i *discordgo.InteractionCreate, opt *intersvc.SamSamPredictRequest, node *cluster.ClusterNode) {
	shdl.ReportCommandInfo(s, i)
	msg, err := shdl.SendRunningMessage(s, i)
	if err != nil {
		log.Println(err)
		return
	}
	sam := &intersvc.SamSamPredict{RequestItem: opt}
	sam.Action(node.StableClient)
	if sam.Error != nil {
		fmt.Println(sam.Error)
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: func() *string { v := sam.Error.Error(); return &v }(),
		})
	} else {
		images := make([]string, 0)
		images = append(images, (sam.GetResponse().MaskedImages)...)
		images = append(images, (sam.GetResponse().Masks)...)
		images = append(images, (sam.GetResponse().BlendedImages)...)

		files := make([]*discordgo.File, 0)
		for j := 0; j < len(images); j++ {
			imageReader, err := utils.GetImageReaderByBase64(images[j])
			if err != nil {
				s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
					Content: func() *string { v := err.Error(); return &v }(),
				})
				return
			}
			files = append(files, &discordgo.File{
				Name:        fmt.Sprintf("sam_%v.png", j),
				Reader:      imageReader,
				ContentType: "image/png",
			})
		}
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: func() *string { v := "Success!"; return &v }(),
			Files:   files[0:3],
		})
	}

}

func (shdl SlashHandler) SamCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	option := &intersvc.SamSamPredictRequest{}
	shdl.SamSetOptions(i.ApplicationCommandData().Options, option)
	node := global.ClusterManager.GetNodeAuto()
	action := func() (map[string]interface{}, error) {
		shdl.SamAction(s, i, option, node)
		return nil, nil
	}
	callback := func() {}
	node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
}
