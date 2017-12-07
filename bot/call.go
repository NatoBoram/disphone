package bot

import (
	"fmt"
	"strings"

	"github.com/NatoBoram/Discord-Phone/config"
	"github.com/bwmarrin/discordgo"
)

func createCall(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Look for ChannelID
	split := strings.SplitAfter(m.Content, callPrefix)

	if len(split) == 2 {

		// Get channel structure
		fromChannel, err := s.State.Channel(m.ChannelID)
		if err != nil {
			fmt.Println("Couldn't get a channel structure.")
			fmt.Println("Message : " + m.Content)
			fmt.Println(err.Error())
			return
		}

		// Get channel structure
		to := strings.Trim(split[1], " ")
		toChannel, err := s.State.Channel(to)
		if err != nil {

			// Just in case
			fmt.Println("Couldn't get a channel structure.")
			fmt.Println("Channel ID : " + to)
			fmt.Println(err.Error())

			// Feedback
			_, err := s.ChannelMessageSend(fromChannel.ID, "Channel "+to+" doesn't exist.")
			if err != nil {
				fmt.Println("Couldn't send a message.")
				fmt.Println("Channel : " + fromChannel.Name)
				fmt.Println(err.Error())
			}

			return
		}

		fmt.Println("Channel name : " + toChannel.Name)

		// Woah there!
		if toChannel.ID == fromChannel.ID {

			// Channel called itself
			_, err := s.ChannelMessageSend(fromChannel.ID, "You can't call yourself!")
			if err != nil {
				fmt.Println("Couldn't send a message.")
				fmt.Println("Channel : " + fromChannel.Name)
				fmt.Println(err.Error())
			}
			return

		} else if toChannel.Type != discordgo.ChannelTypeGuildText {

			// Channel is not a text channel
			_, err := s.ChannelMessageSend(fromChannel.ID, toChannel.Name+" is not a text channel.")
			if err != nil {
				fmt.Println("Couldn't send a message.")
				fmt.Println("Channel : " + fromChannel.Name)
				fmt.Println(err.Error())
				return
			}
			return
		}

		// First call ever?
		if len(config.Calls) == 0 {
			config.Calls = make(map[string][]string)
			config.Calls[fromChannel.ID] = []string{to}
		} else {

			_, exists := config.Calls[fromChannel.ID]
			if exists {

				// Existing channel?
				config.Calls[fromChannel.ID] = rsfa(config.Calls[fromChannel.ID], to)
				config.Calls[fromChannel.ID] = append(config.Calls[fromChannel.ID], to)

			} else {

				// New channel?
				config.Calls[fromChannel.ID] = []string{to}
			}
		}

		// Save
		config.WriteCalls()

		// Get the guild structure
		toGuild, err := s.State.Guild(toChannel.GuildID)
		if err != nil {
			fmt.Println("Couldn't get a guild structure.")
			fmt.Println("Channel : " + toChannel.Name)
			fmt.Println(err.Error())
			return
		}

		// Feedback
		_, err = s.ChannelMessageSend(fromChannel.ID, "Call to <#"+toChannel.ID+"> has started. Make sure <@"+toGuild.OwnerID+"> calls you back by typing `call "+fromChannel.ID+"`.")
		if err != nil {
			fmt.Println("Couldn't send a message.")
			fmt.Println("Channel : " + fromChannel.ID)
			fmt.Println(err.Error())
			return
		}
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
