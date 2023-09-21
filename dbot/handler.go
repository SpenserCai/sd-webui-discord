/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 22:02:04
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-20 13:58:02
 * @Description: file content
 */
package dbot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (dbot *DiscordBot) Ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
}

func (dbot *DiscordBot) InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := dbot.SlashHandlerMap[i.ApplicationCommandData().Name]; ok {
		if dbot.CheckPermission(i.ApplicationCommandData().Name, s, i) {
			h(s, i)
		}
	}

}

func (dbot *DiscordBot) AddCommand() {
	log.Println("Adding commands...")
	dbot.AddedCommand = make([]*discordgo.ApplicationCommand, len(dbot.AppCommand))
	for i, v := range dbot.AppCommand {
		// check if command options are the same
		if dbot.OptionsUnchanged(v) {
			log.Printf("'%v' command options are unchanged, skipping...", v.Name)
			continue
		}
		for _, r := range dbot.RegisteredCommands {
			if r.Name == v.Name {
				err := dbot.Session.ApplicationCommandDelete(dbot.Session.State.User.ID, dbot.ServerID, r.ID)
				if err != nil {
					log.Panicf("Cannot remove '%v' command: %v", v.Name, err)
					continue
				}
			}
		}


		log.Printf("Adding '%v' command...", v.Name)
		cmd, err := dbot.Session.ApplicationCommandCreate(dbot.Session.State.User.ID, dbot.ServerID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			continue
		}
		dbot.AddedCommand[i] = cmd
	}
}

func (dbot *DiscordBot) ClearCommand() {
	commands, err := dbot.Session.ApplicationCommands(dbot.Session.State.User.ID, dbot.ServerID)
	dbot.RegisteredCommands = commands
	if err != nil {
		log.Panicf("Cannot get commands: %v", err)
	}
	if len(dbot.RegisteredCommands) > 0 {
		log.Println("Clearing commands...")
		for _, v := range dbot.RegisteredCommands {
			
			// if command is not in the command list, remove it
			if dbot.CheckCommandInList(v.Name) {
				log.Printf("'%v' command is in the command list, skip...", v.Name)
				continue
			}
			log.Printf("Removing '%v' command...", v.Name)
			err := dbot.Session.ApplicationCommandDelete(dbot.Session.State.User.ID, dbot.ServerID, v.ID)
			if err != nil {
				log.Panicf("Cannot remove '%v' command: %v", v.Name, err)
				continue
			}
		}
	}

}

func (dbot *DiscordBot) CheckCommandInList(name string) bool {
	for _, v := range dbot.AppCommand {
		if v.Name == name {
			return true
		}
	}
	return false
}

func (dbot *DiscordBot) OptionsUnchanged(command *discordgo.ApplicationCommand) bool {
	for _, registeredCommand := range dbot.RegisteredCommands {
		if registeredCommand.Name == command.Name {
			if registeredCommand.Description != command.Description { return false }
			if len(registeredCommand.Options) != len(command.Options) {
				return false
			}
			for i, option := range command.Options {
				if option.Description != registeredCommand.Options[i].Description {
					log.Println("Registered description '", registeredCommand.Options[i].Description, "' is different from command description '", option.Description, "' for command", command.Name)
					return false
				}
				for k, choice := range option.Choices {

					if len(registeredCommand.Options[i].Choices) != len(option.Choices) {
						log.Println("Length of choices is different for command", command.Name)
						log.Println("Registered command:", registeredCommand.Options[i].Choices)
						// print all the choices names and their description
						for _, v := range command.Options[i].Choices {
							log.Println("Command", v.Name, v.Value)
						}
						for _, v := range registeredCommand.Options[i].Choices {
							log.Println("RegisteredCommand", v.Name, v.Value)
						}
						return false
					}
					if choice.Name == registeredCommand.Options[i].Choices[k].Name {
						if choice.Value != registeredCommand.Options[i].Choices[k].Value {
							continue

						}
					} else {
						return false
					}

				}
			}
		}
	}
	return true
}
