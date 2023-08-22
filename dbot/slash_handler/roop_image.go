/*
 * @Author: SpenserCai
 * @Date: 2023-08-22 12:58:13
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-22 14:38:09
 * @Description: file content
 */
package slash_handler

import (
	"log"

	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/SpenserCai/sd-webui-discord/utils"

	"github.com/SpenserCai/sd-webui-go/intersvc"
	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) RoopImageOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "roop_image",
		Description: "Image face swap with roop",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "source_url",
				Description: "The url of the source face image",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "target_url",
				Description: "The url of the target image",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "face_restorer",
				Description: "The face restorer model to use",
				Required:    false,
				Choices:     shdl.faceRestorerModelChoice(),
			},
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "restorer_visibility",
				Description: "The visibility of the face restorer 0-1",
				Required:    false,
			},
		},
	}
}

func (shdl SlashHandler) RoopImageSetOptions(dsOpt []*discordgo.ApplicationCommandInteractionDataOption, opt *intersvc.RoopImageRequest) {
	opt.Model = func() *string { v := "inswapper_128.onnx"; return &v }()
	opt.FaceIndex = []int64{0}
	opt.Scale = func() *int64 { var v int64 = 1; return &v }()
	opt.UpscaleVisibility = func() *float64 { var v float64 = 1.0; return &v }()
	opt.RestorerVisibility = func() *float64 { v := 1.0; return &v }()
	for _, v := range dsOpt {
		switch v.Name {
		case "source_url":
			opt.SourceImage, _ = utils.GetImageBase64(v.StringValue())
		case "target_url":
			opt.TargetImage, _ = utils.GetImageBase64(v.StringValue())
		case "face_restorer":
			opt.FaceRestorer = func() *string { v := v.StringValue(); return &v }()
		case "restorer_visibility":
			opt.RestorerVisibility = func() *float64 { v := v.FloatValue(); return &v }()
		}
	}
}

func (shdl SlashHandler) RoopImageAction(s *discordgo.Session, i *discordgo.InteractionCreate, opt *intersvc.RoopImageRequest, node *cluster.ClusterNode) {
	msg, err := shdl.SendStateMessage("Running", s, i)
	if err != nil {
		log.Println(err)
		return
	}
	roop_image := &intersvc.RoopImage{RequestItem: opt}
	roop_image.Action(node.StableClient)
	if roop_image.Error != nil {
		s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
			Content: func() *string { v := roop_image.Error.Error(); return &v }(),
		})
	} else {
		image, err := utils.GetImageReaderByBase64(roop_image.GetResponse().Image)
		if err != nil {
			s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
				Content: func() *string { v := err.Error(); return &v }(),
			})
		} else {
			s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
				Content: func() *string { v := "Success"; return &v }(),
				Files: []*discordgo.File{
					{
						Name:        "roop_image.png",
						ContentType: "image/png",
						Reader:      image,
					},
				},
			})
		}

	}
}

func (shdl SlashHandler) RoopImageCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	option := &intersvc.RoopImageRequest{}
	shdl.ReportCommandInfo(s, i)
	node := global.ClusterManager.GetNodeAuto()
	action := func() (map[string]interface{}, error) {
		shdl.RoopImageSetOptions(i.ApplicationCommandData().Options, option)
		shdl.RoopImageAction(s, i, option, node)
		return nil, nil
	}
	callback := func() {}
	node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
}
