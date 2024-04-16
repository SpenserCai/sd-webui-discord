/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 22:10:00
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-13 11:59:33
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
	for _, v := range dbot.AppCommands {
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
	dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.DeoldifyOptions())
	dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.SamOptions())
	dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.RembgOptions())
	dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.ExtraSingleOptions())
	dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.PngInfoOptions())
	dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.ControlnetDetectOptions())
	dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.RoopImageOptions())
	dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.Txt2imgOptions())
	dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.Img2imgOptions())
	dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.LoraListOptions())
	dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.PresetOptions())
	if global.Config.UserCenter.Enable {
		dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.RegisterOptions())
		dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.SettingOptions())
		dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.SettingUiOptions())
		dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.ClusterStatusOptions())
		dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.UserInfoOptions())
		if global.Features.Sag {
			dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.SagOptions())
		}
		if global.Features.FreeU {
			dbot.AppCommands = append(dbot.AppCommands, slash_handler.SlashHandler{}.FreeUOptions())
		}
	}
}

func (dbot *DiscordBot) SetLongChoice() {
	global.LongDBotChoice = make(map[string][]*discordgo.ApplicationCommandOptionChoice)
	global.LongDBotChoice["control_net_module"] = slash_handler.SlashHandler{}.ControlnetModuleChoice()
	global.LongDBotChoice["control_net_model"] = slash_handler.SlashHandler{}.ControlnetModelChoice()
	global.LongDBotChoice["sd_model_checkpoint"] = slash_handler.SlashHandler{}.SdModelChoice()
	global.LongDBotChoice["sampler"] = slash_handler.SlashHandler{}.SamplerChoice()
	global.LongDBotChoice["sd_vae"] = slash_handler.SlashHandler{}.SdVaeChoice()
	global.LongDBotChoice["presets"] = slash_handler.SlashHandler{}.PresetChoice()
}
