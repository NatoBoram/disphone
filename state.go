package main

import (
	"github.com/bwmarrin/discordgo"
)

// Guild
func stateGuild(s *discordgo.Session, guildID string) (guild *discordgo.Guild, err error) {

	// State
	guild, err = s.State.Guild(guildID)
	if err == discordgo.ErrStateNotFound {

		// Session
		guild, err = s.Guild(guildID)
		if err != nil {
			return
		}
	}
	return
}

// Channel
func stateChannel(s *discordgo.Session, channelID string) (channel *discordgo.Channel, err error) {

	// State
	channel, err = s.State.Channel(channelID)
	if err == discordgo.ErrStateNotFound {

		// Session
		channel, err = s.Channel(channelID)
		if err != nil {
			return
		}
	}
	return
}

// Member
func stateMember(s *discordgo.Session, guildID, userID string) (member *discordgo.Member, err error) {

	// State
	member, err = s.State.Member(guildID, userID)
	if err == discordgo.ErrStateNotFound {

		// Session
		member, err = s.GuildMember(guildID, userID)
		if err != nil {
			return
		}
	}
	return
}

// Message
func stateMessage(s *discordgo.Session, channelID, messageID string) (message *discordgo.Message, err error) {

	// State
	message, err = s.State.Message(channelID, messageID)
	if err == discordgo.ErrStateNotFound {

		// Session
		message, err = s.ChannelMessage(channelID, messageID)
		if err != nil {
			return
		}
	}
	return
}
