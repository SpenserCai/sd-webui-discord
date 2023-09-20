/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 22:10:00
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-20 12:12:28
 * @Description: file content
 */
package dbot

import (
	"log"
	"reflect"

	"github.com/SpenserCai/sd-webui-discord/dbot/slash_handler"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/SpenserCai/sd-webui-discord/utils"

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
	dbot.AppCommand = append(dbot.AppCommand, slash_handler.SlashHandler{}.RembgOptions())
	dbot.AppCommand = append(dbot.AppCommand, slash_handler.SlashHandler{}.ExtraSingleOptions())
	dbot.AppCommand = append(dbot.AppCommand, slash_handler.SlashHandler{}.PngInfoOptions())
	dbot.AppCommand = append(dbot.AppCommand, slash_handler.SlashHandler{}.ControlnetDetectOptions())
	dbot.AppCommand = append(dbot.AppCommand, slash_handler.SlashHandler{}.RoopImageOptions())
	dbot.AppCommand = append(dbot.AppCommand, slash_handler.SlashHandler{}.Txt2imgOptions())
	dbot.AppCommand = append(dbot.AppCommand, slash_handler.SlashHandler{}.Img2imgOptions())
	if global.Config.UserCenter.Enable {
		dbot.AppCommand = append(dbot.AppCommand, slash_handler.SlashHandler{}.RegisterOptions())
		dbot.AppCommand = append(dbot.AppCommand, slash_handler.SlashHandler{}.SettingOptions())
	}
}

func (dbot *DiscordBot) SetLongChoice() {
	global.LongDBotChoice = make(map[string][]*discordgo.ApplicationCommandOptionChoice)
	global.LongDBotChoice["control_net_module"] = slash_handler.SlashHandler{}.ControlnetModuleChoice()
	global.LongDBotChoice["control_net_model"] = slash_handler.SlashHandler{}.ControlnetModelChoice()
	global.LongDBotChoice["sd_model_checkpoint"] = slash_handler.SlashHandler{}.SdModelChoice()
	global.LongDBotChoice["sampler"] = slash_handler.SlashHandler{}.SamplerChoice()
	global.LongDBotChoice["sd_vae"] = slash_handler.SlashHandler{}.SdVaeChoice()
}
