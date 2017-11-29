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
			fmt.Println("Channel" + to + " doesn't exist.")
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Channel name : " + channel.Name)

		// Woah there!
		if channel.ID == m.ChannelID {
			fmt.Println(m.Author.Username + " wanted " + m.ChannelID + " to call itself.")
			return
		}

		// First call ever?
		if len(config.Calls) == 0 {
			config.Calls = make(map[string][]string)
			config.Calls[m.ChannelID] = []string{to}
		} else {

			_, exists := config.Calls[m.ChannelID]
			if exists {

				// Existing channel?
				fmt.Println("Trying to add channel " + to + " to the list!")
				config.Calls[m.ChannelID] = append(config.Calls[m.ChannelID], to)
				fmt.Println("Success!")
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

		// Get channel structure
		_, err := s.State.Channel(to)
		if err != nil {
			fmt.Println("This channel doesn't exist.")
			fmt.Println(err.Error())
			return
		}

		_, exists := config.Calls[m.ChannelID]
		if exists {
			config.Calls[m.ChannelID] = rsfa(config.Calls[m.ChannelID], to)
		}

		// Save
		config.WriteCalls()
	}
}

// rsfa : Remove String From Array. https://stackoverflow.com/a/34070691/5083247
func rsfa(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
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
				for _, to2 := range tos2 {

					if m.ChannelID == to2 {
						_, err := s.ChannelMessageSend(to, "<@"+m.Author.ID+"> : "+m.Content)
						if err != nil {
							fmt.Println("Couldn't foward the message!")
							fmt.Println(err.Error())
							return
						}
					}
				}
			}
		}
	}
}
