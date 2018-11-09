package dgocmd

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Add adds a command to another command.
// If there's no function in the new command, a help dialog is automatically added.
func (cmd *Command) Add(toadd *Command) {
	if toadd.Function == nil {
		toadd.Function = toadd.Help
	}
	cmd.Commands = append(cmd.Commands, toadd)
}

// Help dialog for a command.
func (cmd *Command) Help(s *discordgo.Session, m *discordgo.MessageCreate, args string) (err error) {

	var help string
	help += "**Name :** " + cmd.Name + "\n"
	help += "**Description :** " + cmd.Description + "\n"
	help += "**Shim :** `" + cmd.Shim + "`\n"

	// Alias
	if len(cmd.Alias) > 0 {
		help += "\n**Alias :** "
	}

	// Alias loop
	for index, alias := range cmd.Alias {
		help += "`" + alias + "`"

		// Comma
		if index < len(cmd.Alias) {
			help += ", "
		}
	}

	// Commands
	if len(cmd.Commands) > 0 {
		help += "\n**Commands**\n"
	}

	// Commands loop
	for _, command := range cmd.Commands {
		help += "**" + command.Name + "** `" + command.Shim + "`\n" + command.Description + "\n"
	}

	_, err = s.ChannelMessageSend(m.ChannelID, help)
	return
}

func (cmd *Command) follow(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
	if strings.HasPrefix(args, cmd.Name) {
		cmd.execute(s, m, args)
	}
}

func (cmd *Command) next(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
	for _, sub := range cmd.Commands {

		// Check this command
		if strings.HasPrefix(args, sub.Name) {

			// Remove the command's name from the args
			splittedArgs := strings.SplitAfter(args, sub.Name)

			// strings.TrimPrefix
			if len(splittedArgs) > 1 {

				cmd.next(s, m, strings.Join(splittedArgs, sub.Name))
			}

		} else {

			// End of the recursion
			cmd.execute(s, m, args)
		}
	}
}

func (cmd *Command) execute(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
	if cmd.Function == nil {
		cmd.Help(s, m, args)
	} else {
		cmd.Function(s, m, args)
	}
}
