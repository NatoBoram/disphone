package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func createCall(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Look for ChannelID
	split := strings.SplitAfter(m.Content, callPrefix)

	if len(split) != 2 {
		return
	}

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
		_, err := s.ChannelMessageSend(fromChannel.ID, "Channel `"+to+"` doesn't exist.")
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
		_, err := s.ChannelMessageSend(fromChannel.ID, "<#"+toChannel.ID+"> is not a text channel.")
		if err != nil {
			fmt.Println("Couldn't send a message.")
			fmt.Println("Channel : " + fromChannel.Name)
			fmt.Println(err.Error())
			return
		}
		return
	}

	// First call ever?
	if len(Calls) == 0 {
		Calls = make(map[string][]string)
		Calls[fromChannel.ID] = []string{to}
	} else {

		_, exists := Calls[fromChannel.ID]
		if exists {

			// Existing channel?
			Calls[fromChannel.ID] = rsfa(Calls[fromChannel.ID], to)
			Calls[fromChannel.ID] = append(Calls[fromChannel.ID], to)

			// Feedback
			_, err = s.ChannelMessageSend(fromChannel.ID, "This channel is already calling <#"+toChannel.ID+">.")
			if err != nil {
				fmt.Println("Couldn't send a message.")
				fmt.Println("Channel : " + fromChannel.ID)
				fmt.Println(err.Error())
			}

			// Don't bother with the rest if refreshing a call.
			return
		}

		// New channel?
		Calls[fromChannel.ID] = []string{to}
	}

	// Save
	go WriteCalls()

	// Get the guild structure
	toGuild, err := s.State.Guild(toChannel.GuildID)
	if err != nil {
		fmt.Println("Couldn't get a guild structure.")
		fmt.Println("Channel : " + toChannel.Name)
		fmt.Println(err.Error())
		return
	}

	// Feedback
	_, err = s.ChannelMessageSend(fromChannel.ID, "Call to <#"+toChannel.ID+"> has started.")
	if err != nil {
		fmt.Println("Couldn't send a message.")
		fmt.Println("Channel : " + fromChannel.ID)
		fmt.Println(err.Error())
	}

	// Alert the other
	_, err = s.ChannelMessageSend(toChannel.ID, "<@"+toGuild.OwnerID+"> **"+m.Author.Username+"** is calling from <#"+fromChannel.ID+">. Type `call "+fromChannel.ID+"` to accept.")
	if err != nil {
		fmt.Println("Couldn't send a message.")
		fmt.Println("Channel : " + fromChannel.ID)
		fmt.Println(err.Error())
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
		return
	}
	to := strings.Trim(split[1], " ")

	// Remove call
	_, exists := Calls[m.ChannelID]
	if exists {
		Calls[m.ChannelID] = rsfa(Calls[m.ChannelID], to)
	}

	// Feedback
	_, err := s.ChannelMessageSend(m.ChannelID, "Call to <#"+to+"> was interrupted.")
	if err != nil {
		fmt.Println("Couldn't send a message.")
		fmt.Println("Channel : " + m.ChannelID)
		fmt.Println(err.Error())

	}

	// Save
	WriteCalls()
}

func foward(s *discordgo.Session, m *discordgo.MessageCreate) {

	// No calls are saved?
	if len(Calls) == 0 {
		return
	}

	// Destinations?
	tos, exists := Calls[m.ChannelID]
	if exists && len(tos) > 0 {

		// Get source channel
		fromChannel, err := s.State.Channel(m.ChannelID)
		if err != nil {
			fmt.Println("Couldn't get a channel structure.")
			fmt.Println("Author : " + m.Author.Username)
			fmt.Println("Message : " + m.Content)
			fmt.Println(err.Error())
			return
		}

		// Clean before fowarding
		Clean(s)

		// For each to in tos
		for _, to := range tos {

			// Check if destination exists
			toChannel, err := s.State.Channel(to)
			if err != nil {
				fmt.Println("Found an invalid destination.")
				fmt.Println(err.Error())
				Clean(s)
				return
			}

			// Check if they call back
			tos2, exists2 := Calls[to]
			if exists2 && len(tos2) > 0 {

				// For each in tos2
				for _, to2 := range tos2 {

					// Check if they call the source
					if m.ChannelID == to2 {

						// Actual fowarding
						_, err = s.ChannelMessageSend(toChannel.ID, "**"+m.Author.Username+"** : "+m.Content)
						if err != nil {
							fmt.Println("Couldn't foward the message from " + fromChannel.Name + " to " + toChannel.Name + ".")
							fmt.Println(err.Error())
							return
						}

						// Foward attachments
						for x := 0; x < len(m.Attachments); x++ {
							_, err = s.ChannelMessageSend(toChannel.ID, m.Attachments[x].URL)
							if err != nil {
								fmt.Println("Couldn't foward an attachment from " + fromChannel.Name + " to " + toChannel.Name + ".")
								fmt.Println(err.Error())
							}
						}
					}
				}
			}
		}
	}
}
