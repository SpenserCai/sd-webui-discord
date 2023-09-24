/*
 * @Author: SpenserCai
 * @Date: 2023-08-22 12:58:13
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-24 21:32:49
 * @Description: file content
 */
package slash_handler

import (
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
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "source_image",
				Description: "The source face image",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "target_image",
				Description: "The target image",
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

func (shdl SlashHandler) RoopImageSetOptions(cmd discordgo.ApplicationCommandInteractionData, opt *intersvc.RoopImageRequest) {
	opt.Model = func() *string { v := "inswapper_128.onnx"; return &v }()
	opt.FaceIndex = []int64{0}
	opt.Scale = func() *int64 { var v int64 = 1; return &v }()
	opt.UpscaleVisibility = func() *float64 { var v float64 = 1.0; return &v }()
	opt.RestorerVisibility = func() *float64 { v := 1.0; return &v }()
	for _, v := range cmd.Options {
		switch v.Name {
		case "source_image":
			// opt.SourceImage, _ = utils.GetImageBase64(cmd.Resolved.Attachments[v.Value.(string)].URL)
			opt.SourceImage = cmd.Resolved.Attachments[v.Value.(string)].URL
		case "target_image":
			// opt.TargetImage, _ = utils.GetImageBase64(cmd.Resolved.Attachments[v.Value.(string)].URL)
			opt.TargetImage = cmd.Resolved.Attachments[v.Value.(string)].URL
		case "face_restorer":
			opt.FaceRestorer = func() *string { v := v.StringValue(); return &v }()
		case "restorer_visibility":
			opt.RestorerVisibility = func() *float64 { v := v.FloatValue(); return &v }()
		}
	}
}

func (shdl SlashHandler) RoopImageAction(s *discordgo.Session, i *discordgo.InteractionCreate, opt *intersvc.RoopImageRequest, node *cluster.ClusterNode) {
	roop_image := &intersvc.RoopImage{RequestItem: opt}
	roop_image.Action(node.StableClient)
	if roop_image.Error != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func() *string { v := roop_image.Error.Error(); return &v }(),
		})
	} else {
		image, err := utils.GetImageReaderByBase64(roop_image.GetResponse().Image)
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: func() *string { v := err.Error(); return &v }(),
			})
		} else {
			msg, _ := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: func() *string { v := "Success"; return &v }(),
				Files: []*discordgo.File{
					{
						Name:        "roop_image.png",
						ContentType: "image/png",
						Reader:      image,
					},
				},
			})
			shdl.SetHistory("roop_image", msg.ID, i, opt)
		}

	}
}

func (shdl SlashHandler) RoopImageCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	option := &intersvc.RoopImageRequest{}
	shdl.RespondStateMessage("Running", s, i)
	node := global.ClusterManager.GetNodeAuto()
	action := func() (map[string]interface{}, error) {
		shdl.RoopImageSetOptions(i.ApplicationCommandData(), option)
		shdl.RoopImageAction(s, i, option, node)
		return nil, nil
	}
	callback := func() {}
	node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
}
