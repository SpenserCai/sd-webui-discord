/*
 * @Author: SpenserCai
 * @Date: 2023-10-13 11:33:20
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-11-01 10:03:32
 * @Description: file content
 */
package slash_handler

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/global"

	"github.com/SpenserCai/sd-webui-go/intersvc"
	"github.com/bwmarrin/discordgo"
)

func (shdl SlashHandler) LoraListOptions() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "lora_list",
		Description: "Get Lora list",
		Options:     []*discordgo.ApplicationCommandOption{},
	}
}

func GetLoraListPageSize() int {
	pageSize := 5
	return pageSize
}

func (shdl SlashHandler) BuildLoraListComponent(currentPage int, pageSize int, totalPage int) *[]discordgo.MessageComponent {
	// 根据当前页和总页数生成按钮，总共5个按钮，第一个是到第一页，第二个是到上一页，第三个是当前页（禁用），第四个是下一页，第五个是到最后一页，如果当前页是第一页则第一个和第二个按钮禁用，如果当前页是最后一页则第四个和第五个按钮禁用
	components := []discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					Label:    "<<",
					CustomID: "lora_list|first",
					Style:    discordgo.SecondaryButton,
					Disabled: currentPage == 1,
				},
				&discordgo.Button{
					Label:    "<",
					CustomID: "lora_list|prev",
					Style:    discordgo.SecondaryButton,
					Disabled: currentPage == 1,
				},
				&discordgo.Button{
					Label:    fmt.Sprintf("%d/%d", currentPage, totalPage),
					CustomID: "lora_list|current",
					Style:    discordgo.SecondaryButton,
					Disabled: true,
				},
				&discordgo.Button{
					Label:    ">",
					CustomID: "lora_list|next",
					Style:    discordgo.SecondaryButton,
					Disabled: currentPage == totalPage,
				},
				&discordgo.Button{
					Label:    ">>",
					CustomID: "lora_list|last",
					Style:    discordgo.SecondaryButton,
					Disabled: currentPage == totalPage,
				},
			},
		},
	}
	return &components

}

func (shdl SlashHandler) GetCurrentPageLoraList(currentPage int, pageSize int, loras *[]intersvc.LoraItem) string {
	// 根据当前页和每页数量返回当前页的Lora列表
	loraListString := ""
	start := (currentPage - 1) * pageSize
	end := currentPage * pageSize
	if end > len(*loras) {
		end = len(*loras)
	}
	for _, lora := range (*loras)[start:end] {
		loraListString += "**Name: " + lora.Name + "**\n"
		if lora.Metadata.SsSdModelName != "" {
			loraListString += "**Base Model: " + lora.Metadata.SsSdModelName + "**\n"
		}
		loraListString += "```<lora:" + lora.Name + ":1>```"
		loraListString += "\n"
	}
	return loraListString
}

func (shdl SlashHandler) LoraListAction(s *discordgo.Session, i *discordgo.InteractionCreate, node *cluster.ClusterNode) {
	lora_list := &intersvc.SdapiV1Loras{}
	lora_list.Action(node.StableClient)
	if lora_list.Error != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func() *string { v := lora_list.Error.Error(); return &v }(),
		})
	} else {
		context := ""
		loras := lora_list.GetResponse()
		loraListString := shdl.GetCurrentPageLoraList(1, GetLoraListPageSize(), loras)
		// log.Println(loraListString)
		mainEmbed := shdl.MessageEmbedTemplate()
		mainEmbed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:  "Lora List",
				Value: loraListString,
			},
		}
		allEmbeds := []*discordgo.MessageEmbed{mainEmbed}
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &context,
			Embeds:  &allEmbeds,
			Components: shdl.BuildLoraListComponent(1, 5, func() int {
				if len(*loras)%5 == 0 {
					return len(*loras) / 5
				} else {
					return len(*loras)/5 + 1
				}
			}()),
		})
		if err != nil {
			log.Println(err)
		}
	}

}

func (shdl SlashHandler) LoraListPageChange(changeType string, pageSize int, s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 先把消息改成Running
	currentPage := 1
	currentCommpoents, err := s.ChannelMessage(i.ChannelID, i.Interaction.Message.ID)
	if err != nil {
		s.ChannelMessageSend(i.ChannelID, "Error: "+err.Error())
		return
	}
	// 获取当前页
	for _, v := range currentCommpoents.Components {
		// 判断是否是Button，如果是转成Button处理
		if v.Type() == discordgo.ActionsRowComponent {
			component := v.(*discordgo.ActionsRow)
			for _, v := range component.Components {
				if v.Type() == discordgo.ButtonComponent {
					button := v.(*discordgo.Button)
					if button.CustomID == "lora_list|current" {
						pageLabel := strings.Split(button.Label, "/")
						currentPage, _ = strconv.Atoi(pageLabel[0])
						break
					}
				}
			}
		}
	}
	s.ChannelMessageDelete(i.ChannelID, i.Interaction.Message.ID)
	shdl.RespondStateMessage("Running", s, i)
	lora_list := &intersvc.SdapiV1Loras{}
	lora_list.Action(global.ClusterManager.GetNodeAuto().StableClient)
	if lora_list.Error != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func() *string { v := lora_list.Error.Error(); return &v }(),
		})
	} else {
		context := ""
		loras := lora_list.GetResponse()
		switch changeType {
		case "first":
			currentPage = 1
		case "prev":
			currentPage -= 1
		case "next":
			currentPage += 1
		case "last":
			if len(*loras)%pageSize == 0 {
				currentPage = len(*loras) / pageSize
			} else {
				currentPage = len(*loras)/pageSize + 1
			}
		}
		loraListString := shdl.GetCurrentPageLoraList(currentPage, 5, loras)
		// log.Println(loraListString)
		mainEmbed := shdl.MessageEmbedTemplate()
		mainEmbed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:  "Lora List",
				Value: loraListString,
			},
		}
		allEmbeds := []*discordgo.MessageEmbed{mainEmbed}
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &context,
			Embeds:  &allEmbeds,
			Components: shdl.BuildLoraListComponent(currentPage, 5, func() int {
				if len(*loras)%5 == 0 {
					return len(*loras) / 5
				} else {
					return len(*loras)/5 + 1
				}
			}()),
		})
		if err != nil {
			log.Println(err)
		}
	}
}

func (shdl SlashHandler) LoraListComponentHandler(s *discordgo.Session, i *discordgo.InteractionCreate, pageSize int) {
	customIDList := strings.Split(i.MessageComponentData().CustomID, "|")
	cmd := fmt.Sprintf("%s|%s", customIDList[0], customIDList[1])
	switch cmd {
	case "lora_list|first":
		shdl.LoraListPageChange("first", pageSize, s, i)
	case "lora_list|prev":
		shdl.LoraListPageChange("prev", pageSize, s, i)
	case "lora_list|next":
		shdl.LoraListPageChange("next", pageSize, s, i)
	case "lora_list|last":
		shdl.LoraListPageChange("last", pageSize, s, i)
	}
}

func (shdl SlashHandler) LoraListCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		shdl.RespondStateMessage("Running", s, i)
		node := global.ClusterManager.GetNodeAuto()
		action := func() (map[string]interface{}, error) {
			shdl.LoraListAction(s, i, node)
			return nil, nil
		}
		callback := func() {}
		node.ActionQueue.AddTask(shdl.GenerateTaskID(i), action, callback)
	case discordgo.InteractionMessageComponent:
		shdl.LoraListComponentHandler(s, i, GetLoraListPageSize())
	}
}
