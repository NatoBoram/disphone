package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func addHandlers(s *discordgo.Session) {

	s.AddHandler(messageCreateHandler)

}

func messageCreateHandler(s *discordgo.Session, event *discordgo.MessageCreate) {

	m := event.Message

	// Myself?
	if m.Author.ID == me.ID {
		return
	}

	// Get channel structure
	c, err := stateChannel(s, m.ChannelID)
	if err != nil {
		fmt.Println("Couldn't get the channel structure.")
		fmt.Println("Message : " + m.Content)
		fmt.Println(err.Error())
		return
	}

	// Get guild structure
	g, err := stateGuild(s, c.GuildID)
	if err != nil {
		fmt.Println("Couldn't get the guild structure.")
		fmt.Println("Channel : " + c.Name)
		fmt.Println(err.Error())
		return
	}

	// Guild Owner
	if m.Author.ID == g.OwnerID || m.Author.ID == master.ID {

		// Is it a command?
		if strings.HasPrefix(m.Content, me.Mention()) {
			commands(s, g, c, m)
			return
		}
	}

	// Foward
	// Clean(s)
	// foward(s, m)
}
