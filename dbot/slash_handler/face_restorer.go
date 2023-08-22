/*
 * @Author: SpenserCai
 * @Date: 2023-08-22 13:10:36
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-22 13:12:54
 * @Description: file content
 */
package slash_handler

import (
	"log"

	"github.com/SpenserCai/sd-webui-discord/global"

	"github.com/SpenserCai/sd-webui-go/intersvc"
	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) faceRestorerModelChoice() []*discordgo.ApplicationCommandOptionChoice {
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	modesvc := &intersvc.SdapiV1FaceRestorers{}
	modesvc.Action(global.ClusterManager.GetNodeAuto().StableClient)
	if modesvc.Error != nil {
		log.Println(modesvc.Error)
		return choices
	}
	models := modesvc.GetResponse()
	for _, model := range *models {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  *model.Name,
			Value: *model.Name,
		})
	}
	return choices
}
