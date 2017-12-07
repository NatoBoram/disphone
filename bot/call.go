package bot

import (
	"fmt"
	"strings"

	"github.com/NatoBoram/Discord-Phone/config"
	"github.com/bwmarrin/discordgo"
)

func createCall(s *discordgo.Session, m *discordgo.MessageCreate) {

	command := m.Content

	// Look for ChannelID
	split := strings.SplitAfter(command, callPrefix)

	if len(split) == 2 {
		to := strings.Trim(split[1], " ")
		fmt.Println("Begin call to " + to + ".")

		// Get channel structure
		channel, err := s.State.Channel(to)
		if err != nil {

			// Just in case
			fmt.Println(err.Error())

			// Feedback
			_, err := s.ChannelMessageSend(m.ChannelID, "Channel "+to+" doesn't exist.")
			if err != nil {
				fmt.Println("Couldn't say that " + to + " doesn't exist.")
				fmt.Println(err.Error())
				return
			}

			return
		}
		fmt.Println("Channel name : " + channel.Name)

		// Woah there!
		if channel.ID == m.ChannelID {

			// Channel called itself
			_, err := s.ChannelMessageSend(m.ChannelID, "You can't call yourself!")
			if err != nil {
				fmt.Println("Couldn't tell " + m.Author.Username + " that it can't make a channel call itself.")
				fmt.Println(err.Error())
				return
			}
			return

		} else if channel.Type != discordgo.ChannelTypeGuildText {

			// Channel is not a text channel
			_, err := s.ChannelMessageSend(m.ChannelID, channel.Name+" is not a text channel.")
			if err != nil {
				fmt.Println("Couldn't tell " + m.Author.Username + " that " + channel.Name + " is not a text channel.")
				fmt.Println(err.Error())
				return
			}
			return
		} else {

			// Get the guild structure
			guild, err := s.State.Guild(channel.GuildID)
			if err != nil {
				fmt.Println("Couldn't get a guild structure.")
				fmt.Println("Channel : " + channel.Name)
				fmt.Println(err.Error())
				return
			}

			// Feedback
			_, err = s.ChannelMessageSend(m.ChannelID, "Call to <#"+channel.ID+"> has started. Make sure <@"+guild.OwnerID+"> calls you back by typing `call "+m.ChannelID+"`.")
			if err != nil {
				fmt.Println("Couldn't send a message.")
				fmt.Println("Channel : " + m.ChannelID)
				fmt.Println(err.Error())
				return
			}
		}

		// First call ever?
		if len(config.Calls) == 0 {
			config.Calls = make(map[string][]string)
			config.Calls[m.ChannelID] = []string{to}
		} else {

			_, exists := config.Calls[m.ChannelID]
			if exists {

				// Existing channel?
				config.Calls[m.ChannelID] = rsfa(config.Calls[m.ChannelID], to)
				config.Calls[m.ChannelID] = append(config.Calls[m.ChannelID], to)

			} else {

				// New channel?
				config.Calls[m.ChannelID] = []string{to}
			}
		}

		// Save
		config.WriteCalls()
	}
}

func hangUp(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotID {
		return
	}

	command := m.Content

	// Look for ChannelID
	split := strings.SplitAfter(command, hangUpPrefix)

	if len(split) == 2 {
		to := strings.Trim(split[1], " ")

		// Remove call
		_, exists := config.Calls[m.ChannelID]
		if exists {
			config.Calls[m.ChannelID] = rsfa(config.Calls[m.ChannelID], to)
		}

		// Save
		config.WriteCalls()
	}
}

// rsfa : Remove String From Array.
func rsfa(a []string, s string) []string {
	var n []string
	for i, v := range a {
		if v != s {
			n = append(n, a[i])
		}
	}
	return n
}

func foward(s *discordgo.Session, m *discordgo.MessageCreate) {

	// No calls are saved?
	if len(config.Calls) == 0 {
		return
	}

	// Destinations?
	tos, exists := config.Calls[m.ChannelID]
	if exists && len(tos) > 0 {

		// For each to in tos
		for _, to := range tos {

			// Check if they call back
			tos2, exists2 := config.Calls[to]
			if exists2 && len(tos2) > 0 {

				// For each in tos2
				for _, to2 := range tos2 {

					// Check if they call the source
					if m.ChannelID == to2 {
						_, err := s.ChannelMessageSend(to, "<@"+m.Author.ID+"> : "+m.Content)
						if err != nil {
							fmt.Println("Couldn't foward the message from " + m.ChannelID + " to " + to)
							fmt.Println(err.Error())
							return
						}
					}
				}
			}
		}
	}
}
