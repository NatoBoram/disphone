package bot

import (
	"fmt"

	"github.com/NatoBoram/Discord-Phone/config"

	"github.com/bwmarrin/discordgo"
)

// BotID : Numerical ID of the bot
var BotID string
var goBot *discordgo.Session

// Start : Starts the bot.
func Start() {

	// Go online!
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Get Bot ID
	st, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}
	BotID = st.ID

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

	// Myself?
	if m.Author.ID == BotID {
		return
	}

	// Get channel structure
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Get guild structure
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Guild Owner
	if m.Author.ID == guild.OwnerID {

	}
}
