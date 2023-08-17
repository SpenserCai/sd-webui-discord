/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 22:10:00
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-17 15:17:16
 * @Description: file content
 */
package dbot

import (
	"log"
	"reflect"
	"sd-webui-discord/dbot/slash_handler"
	"sd-webui-discord/utils"

	"github.com/bwmarrin/discordgo"
)

func (dbot *DiscordBot) GenerateSlashMap() error {
	// 遍历AppCommands，取name
	for _, v := range dbot.AppCommand {
		commandName := v.Name
		// 如果name中有_则用下划线分割后每个首字母专大写，如果没有_则直接首字母转大写
		commandName = utils.FormatCommand(commandName) + "CommandHandler"
		// 通过反射找到对应的方法赋值给map
		pkgValue := reflect.ValueOf(slash_handler.SlashHandler{})
		methodValue := pkgValue.MethodByName(commandName)

		if !methodValue.IsValid() {
			log.Println("Function not found:", commandName)
		}
		dbot.SlashHandlerMap[v.Name] = methodValue.Interface().(func(s *discordgo.Session, i *discordgo.InteractionCreate))
	}
	return nil
}

func (dbot *DiscordBot) GenerateCommandList() {
	dbot.AppCommand = append(dbot.AppCommand, slash_handler.SlashHandler{}.DeoldifyOptions())
	dbot.AppCommand = append(dbot.AppCommand, slash_handler.SlashHandler{}.SamOptions())
}
