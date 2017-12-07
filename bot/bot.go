package bot

import (
	"fmt"
	"strings"

	"github.com/NatoBoram/Discord-Phone/config"

	"github.com/bwmarrin/discordgo"
)

// BotID : Numerical ID of the bot
var BotID string
var goBot *discordgo.Session

// Globals
var callPrefix = "call "
var hangUpPrefix = "hang up "

// Start : Starts the bot.
func Start() {

	// Go online!
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Couldn't get online.")
		fmt.Println(err.Error())
		return
	}

	// Get Bot ID
	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println("Couldn't get the BotID.")
		fmt.Println(err.Error())
		return
	}
	BotID = u.ID

	// Hey, listen!
	goBot.AddHandler(messageCreateHandler)

	// Crash on error
	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// It's alive!
	fmt.Println("Discord-Phone is running!")
}

func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	config.Clean(s)

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
	if m.Author.ID == guild.OwnerID {

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
	foward(s, m)
}
