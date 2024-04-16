/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 22:02:04
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-11-01 12:25:22
 * @Description: file content
 */
package dbot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (dbot *DiscordBot) Ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
}

func (dbot *DiscordBot) InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand, discordgo.InteractionApplicationCommandAutocomplete:
		if h, ok := dbot.SlashHandlerMap[i.ApplicationCommandData().Name]; ok {
			if dbot.CheckPermission(i.ApplicationCommandData().Name, s, i) {
				h(s, i)
			}
		}
	case discordgo.InteractionMessageComponent:
		component_command := strings.Split(i.MessageComponentData().CustomID, "|")
		if h, ok := dbot.SlashHandlerMap[component_command[0]]; ok {
			if dbot.CheckPermission(component_command[0], s, i) {
				h(s, i)
			}
		}
	case discordgo.InteractionModalSubmit:
		modal_command := strings.Split(i.ModalSubmitData().CustomID, "|")
		if h, ok := dbot.SlashHandlerMap[modal_command[0]]; ok {
			if dbot.CheckPermission(modal_command[0], s, i) {
				h(s, i)
			}
		}
	}

}

func (dbot *DiscordBot) SyncCommands() {
	commands, err := dbot.Session.ApplicationCommands(dbot.Session.State.User.ID, dbot.ServerID)
	dbot.RegisteredCommands = commands
	if err != nil {
		log.Panicf("Cannot get commands: %v", err)
	}
	if len(dbot.RegisteredCommands) > 0 {
		// log.Println("Clearing other bots commands...")
		for _, v := range dbot.RegisteredCommands {

			// if command is not in the command list, remove it
			if dbot.CheckCommandInList(v.Name) {
				// log.Printf("'%v' command is in the command list, skip...", v.Name)
				continue
			}
			log.Printf("Removing '%v' command...", v.Name)
			err := dbot.Session.ApplicationCommandDelete(dbot.Session.State.User.ID, dbot.ServerID, v.ID)
			if err != nil {
				// log.Panicf("Cannot remove '%v' command: %v", v.Name, err)
				continue
			}
		}
	}

	log.Println("Adding/Updating commands...")
	for _, v := range dbot.AppCommands {

		// check if command needs update
		if !dbot.CommandNeedsUpdate(v) {
			// log.Printf("'%v' command options are unchanged, skipping...", v.Name)
			continue
		}
		// delete old version of command
		for _, r := range dbot.RegisteredCommands {
			if r.Name == v.Name {
				err := dbot.Session.ApplicationCommandDelete(dbot.Session.State.User.ID, dbot.ServerID, r.ID)
				if err != nil {
					log.Panicf("Cannot remove '%v' command: %v", v.Name, err)
					continue
				}
			}
		}

		// add new version of command
		// log.Printf("Adding '%v' command...", v.Name)
		_, err := dbot.Session.ApplicationCommandCreate(dbot.Session.State.User.ID, dbot.ServerID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			continue
		}
	}
	updatedCommands, err := dbot.Session.ApplicationCommands(dbot.Session.State.User.ID, dbot.ServerID)
	if err != nil {
		log.Panicf("Cannot get commands: %v", err)
	}
	dbot.RegisteredCommands = updatedCommands

}

func (dbot *DiscordBot) SetBaseCommandHelperInfo() {
	commandInfo := "\n"
	needAddCmd := []string{
		"setting",
		"setting_ui",
		"txt2img",
	}
	for _, v := range dbot.RegisteredCommands {
		for _, vv := range needAddCmd {
			if v.Name == vv {
				commandInfo += fmt.Sprintf("</%s:%s> ", v.Name, v.ID)
			}
		}
	}
	fmt.Println("Please update this to App description:", commandInfo)
}

func (dbot *DiscordBot) CheckCommandInList(name string) bool {
	for _, v := range dbot.AppCommands {
		if v.Name == name {
			return true
		}
	}
	return false
}

func (dbot *DiscordBot) CommandNeedsUpdate(command *discordgo.ApplicationCommand) bool {
	for _, registeredCommand := range dbot.RegisteredCommands {
		if registeredCommand.Name == command.Name {
			// log.Println("Processing command", command.Name)
			// new description
			if registeredCommand.Description != command.Description {
				// log.Println("New description for command", command.Name, "old: '"+registeredCommand.Description+"', new: '"+command.Description+"', not equal:", registeredCommand.Description != command.Description)
				return true
			}

			// new location
			if registeredCommand.DescriptionLocalizations != nil {
				// 先判断有没有新的location
				if len(*registeredCommand.DescriptionLocalizations) != len(*command.DescriptionLocalizations) {
					// log.Println("New localization for command", command.Name)
					return true
				}
				// 再判断有没有新的location的description
				for k, v := range *command.DescriptionLocalizations {
					if (*registeredCommand.DescriptionLocalizations)[k] != v {
						// log.Println("Different localization for command", command.Name)
						return true
					}
				}
			}

			// new options
			if len(registeredCommand.Options) != len(command.Options) {
				// log.Println("New options for command", command.Name)
				return true
			}

			for i, option := range command.Options {
				// new option description
				if option.Description != registeredCommand.Options[i].Description {
					// log.Println("Registered description '", registeredCommand.Options[i].Description, "' is different from command description '", option.Description, "' for command", command.Name)
					return true
				}

				// if option.Autocomplete {
				// 	return true
				// }

				// new choices
				for k, choice := range option.Choices {

					if len(registeredCommand.Options[i].Choices) != len(option.Choices) {
						// log.Println("Length of choices is different for command", command.Name)
						// log.Println("Registered command:", registeredCommand.Options[i].Choices)
						// // print all the choices names and their description
						// for _, v := range command.Options[i].Choices {
						// 	log.Println("Command", v.Name, v.Value)
						// }
						// for _, v := range registeredCommand.Options[i].Choices {
						// 	log.Println("RegisteredCommand", v.Name, v.Value)
						// }
						return true
					}
					if choice.Name == registeredCommand.Options[i].Choices[k].Name {
						if choice.Value != registeredCommand.Options[i].Choices[k].Value {
							continue
						}
					} else {
						return false
					}

				}

				// if min/max values are different
				newIsNull := option.MinValue == nil
				oldIsNull := registeredCommand.Options[i].MinValue == nil
				if newIsNull != oldIsNull || (!newIsNull && !oldIsNull && *option.MinValue != *registeredCommand.Options[i].MinValue) || option.MaxValue != registeredCommand.Options[i].MaxValue {
					//log.Println("New min/max values for command", command.Name, "old min:", registeredCommand.Options[i].MinValue, "new min:", option.MinValue, "old max:", registeredCommand.Options[i].MaxValue, "new max:", option.MaxValue)
					return true
				}
				// no new choices
			}
			// no new options or no options
			return false
		}
	}
	return true
}
