/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 22:02:04
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-17 10:04:31
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
		h(s, i)
	}

}

func (dbot *DiscordBot) AddCommand() {
	log.Println("Adding commands...")
	dbot.AddedCommand = make([]*discordgo.ApplicationCommand, len(dbot.AppCommand))
	for i, v := range dbot.AppCommand {
		cmd, err := dbot.Session.ApplicationCommandCreate(dbot.Session.State.User.ID, dbot.ServerID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			continue
		}
		dbot.AddedCommand[i] = cmd
	}
}

func (dbot *DiscordBot) RemoveCommand() {
	log.Println("Removing commands...")
	for _, v := range dbot.AddedCommand {
		err := dbot.Session.ApplicationCommandDelete(dbot.Session.State.User.ID, dbot.ServerID, v.ID)
		if err != nil {
			log.Panicf("Cannot remove '%v' command: %v", v.Name, err)
			continue
		}
	}
}
