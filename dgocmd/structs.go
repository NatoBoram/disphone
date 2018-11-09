package dgocmd

import "github.com/bwmarrin/discordgo"

// Command is a command
type Command struct {

	// Name of the command.
	Name string

	// Description of a command and what it does.
	Description string

	// Shim is the string used to call a command.
	Shim string

	// Alias are alternative shims for a command.
	Alias []string

	// Subcommands that can be ran using this command.
	Commands []*Command

	// Function to call when this command is called. If empty, an auto-generated help message will take its place.
	Function func(s *discordgo.Session, m *discordgo.MessageCreate, args string) (err error)
}
