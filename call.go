package main

import (
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

/*
func createCall(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Look for ChannelID
	split := strings.SplitAfter(m.Content, "call")

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

			// Currently calling this channel?
			if csia(Calls[fromChannel.ID], to) {

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

			// Existing channel?
			Calls[fromChannel.ID] = rsfa(Calls[fromChannel.ID], to)
			Calls[fromChannel.ID] = append(Calls[fromChannel.ID], to)
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

	if m.Author.ID == me.ID {
		return
	}

	command := m.Content

	// Look for ChannelID
	split := strings.SplitAfter(command, "hang up")

	if len(split) != 2 {
		return
	}
	to := strings.Trim(split[1], " ")

	// Remove call
	_, exists := Calls[m.ChannelID]
	if exists {

		// If TO is inside FROM
		if csia(Calls[m.ChannelID], to) {

			Calls[m.ChannelID] = rsfa(Calls[m.ChannelID], to)

			// Feedback
			_, err := s.ChannelMessageSend(m.ChannelID, "Call to <#"+to+"> was interrupted.")
			if err != nil {
				fmt.Println("Couldn't send a message.")
				fmt.Println("Channel : " + m.ChannelID)
				fmt.Println(err.Error())
			}

			// Save
			WriteCalls()
			return
		}
	}

	// Nope
	_, err := s.ChannelMessageSend(m.ChannelID, "This channel wasn't calling <#"+to+">.")
	if err != nil {
		fmt.Println("Couldn't send a message.")
		fmt.Println("Channel : " + m.ChannelID)
		fmt.Println(err.Error())
	}
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
*/

func getCalls(c *discordgo.Channel) (calls []PhoneCall, err error) {

	// Select calls
	rows, err := selectCalls(c)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {

		// For each call
		var call PhoneCall
		err = rows.Scan(&call.From, &call.To)
		if err != nil {
			fmt.Println("Couldn't select a call.")
			fmt.Println(err.Error())
			continue
		}

		// Append
		calls = append(calls, call)
	}

	err = rows.Err()
	return
}

func getChannels(s *discordgo.Session, calls []PhoneCall) (channels []*discordgo.Channel) {

	// For each calls
	for _, call := range calls {

		// Check if the channel it's calling exists
		channel, err := stateChannel(s, call.To)
		if err != nil {

			fmt.Println("Channel doesn't exist.")
			fmt.Println(err.Error())

			// Channel doesn't exist, should be removed.
			_, err = deleteCalls(call.To)
			if err != nil {
				fmt.Println("Couldn't delete all references to a channel.")
				fmt.Println(err.Error())
			}

			continue
		}

		// Append
		channels = append(channels, channel)
	}

	return
}

func getValidCalls(s *discordgo.Session, c *discordgo.Channel) (channels []*discordgo.Channel, err error) {

	// Get this channel's calls
	calls, err := getCalls(c)
	if err != nil {
		return
	}

	// Get channels from calls
	called := getChannels(s, calls)
	if err != nil {
		return
	}

	// For each channel
	for _, channel := range called {

		// Check if the channel calls back.
		_, err = selectCall(channel, c)
		if err == sql.ErrNoRows {

			// This channel doesn't call back.
			continue

		} else if err != nil {
			fmt.Println("Couldn't select a call.")
			fmt.Println(err.Error())
			continue
		}

		// Append
		channels = append(channels, channel)
	}

	return
}

func createMessageEmbed(s *discordgo.Session, g *discordgo.Guild, c *discordgo.Channel, m *discordgo.Message, b *discordgo.Member) (embed *discordgo.MessageEmbed) {

	// Embed
	embed = &discordgo.MessageEmbed{
		Color: s.State.UserColor(m.Author.ID, m.ID),
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://canary.discordapp.com/channels/" + g.ID + "/" + m.ChannelID + "/" + m.ID + "/",
			Name:    m.Author.Username,
			IconURL: m.Author.AvatarURL(""),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    g.Name,
			IconURL: discordgo.EndpointGuildIcon(g.ID, g.Icon),
		},
		Timestamp: string(m.Timestamp),
	}

	// Nick
	if b.Nick != "" {
		embed.Author.Name = b.Nick
	}

	// Description
	if m.Content != "" {
		embed.Description = m.Content
	}

	return
}

func forward(s *discordgo.Session, g *discordgo.Guild, from *discordgo.Channel, m *discordgo.Message, b *discordgo.Member, to *discordgo.Channel) {

	embed := createMessageEmbed(s, g, from, m, b)

	s.ChannelMessageSendEmbed(to.ID, embed)
}
