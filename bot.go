package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// BotID : Numerical ID of the bot
var BotID string

// Globals
const (
	callPrefix   = "call "
	hangUpPrefix = "hang up "
)

// Start : Starts the bot.
func Start() {

	// Go online!
	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Couldn't get online.")
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	// Get Bot ID
	u, err := session.User("@me")
	if err != nil {
		fmt.Println("Couldn't get the BotID.")
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}
	BotID = u.ID

	// Hey, listen!
	session.AddHandler(messageCreateHandler)

	// Crash on error
	err = session.Open()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	// It's alive!
	fmt.Println("Discord-Phone is running!")
}

func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Myself?
	if m.Author.ID == BotID {
		return
	}

	// Get channel structure
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Couldn't get the channel structure.")
		fmt.Println("Message : " + m.Content)
		fmt.Println(err.Error())
		return
	}

	// Get guild structure
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		fmt.Println("Couldn't get the guild structure.")
		fmt.Println("Channel : " + channel.Name)
		fmt.Println(err.Error())
		return
	}

	// Guild Owner
	if m.Author.ID == guild.OwnerID || m.Author.ID == BotMaster {

		// Starting a call?
		if strings.HasPrefix(m.Content, callPrefix) {
			createCall(s, m)
			return
		}

		// Ending a call?
		if strings.HasPrefix(m.Content, hangUpPrefix) {
			hangUp(s, m)
			return
		}
	}

	// Foward
	Clean(s)
	foward(s, m)
}
