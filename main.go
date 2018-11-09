package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"gitlab.com/NatoBoram/Discord-Phone/dgocmd"
)

func main() {

	// License
	fmt.Println("")
	fmt.Println("Discord-Phone : Makes phone calls between Discord servers.")
	fmt.Println("Copyright © 2018 Nato Boram")
	fmt.Println("This program is free software : you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY ; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with this program. If not, see http://www.gnu.org/licenses/.")
	fmt.Println("Contact : https://gitlab.com/NatoBoram/Discord-Phone")
	fmt.Println("")

	// Database
	err := initDatabase()
	if err != nil {
		return
	}
	defer db.Close()

	// Discord
	session, err = initDiscord()
	if err != nil {
		return
	}
	defer session.Close()

	// Commands
	err = initCommands()
	if err != nil {
		return
	}

	// It's alive!
	fmt.Println("Hi, " + master.Username + ". I am " + me.Username + ", and everything's all right!")

	// Wait for future input
	<-make(chan struct{})
}

func initDatabase() (err error) {

	// Read the database config
	var database Database
	err = readDatabase(&database)
	if err != nil {
		fmt.Println("Could not load the database configuration.")
		fmt.Println(err.Error())
		writeTemplateDatabase()
		return
	}

	// Check for empty JSON
	if database.isEmpty() {
		err = errors.New("Configuration is missing inside " + databasePath)
		fmt.Println(err.Error())
		return
	}

	// Connect to MariaDB
	db, err = sql.Open("mysql", database.toConnectionString())
	if err != nil {
		fmt.Println("Could not connect to the database.")
		fmt.Println(err.Error())
		return
	}

	// Test the connection to MariaDB
	err = db.Ping()
	if err != nil {
		fmt.Println("Could not ping the database.")
		fmt.Println(err.Error())
		return
	}

	// Create the tables if they don't exist
	_, err = createTables()
	if err != nil {
		fmt.Println("Could not create a table in the database.")
		fmt.Println(err.Error())
		return
	}

	return
}

func initDiscord() (s *discordgo.Session, err error) {

	// Read the Discord config
	var discord Discord
	err = readDiscord(&discord)
	if err != nil {
		fmt.Println("Could not load the Discord configuration.")
		fmt.Println(err.Error())
		writeTemplateDiscord()
		return
	}

	// Check for empty JSON
	if discord.isEmpty() {
		err = errors.New("Configuration is missing inside " + discordPath)
		fmt.Println(err.Error())
		return
	}

	// Create a Discord session
	s, err = discordgo.New("Bot " + discord.Token)
	if err != nil {
		fmt.Println("Could not create a Discord session.")
		fmt.Println(err.Error())
		return
	}

	// Connect to Discord
	err = s.Open()
	if err != nil {
		fmt.Println("Could not connect to Discord.")
		fmt.Println(err.Error())
		return
	}

	// Myself
	me, err = s.User("@me")
	if err != nil {
		fmt.Println("Couldn't get myself.")
		fmt.Println(err.Error())
		return
	}

	// Master
	master, err = s.User(discord.MasterID)
	if err != nil {
		fmt.Println("Couldn't recognize my master.")
		fmt.Println(err.Error())
		return
	}

	// Handlers
	addHandlers(s)

	return
}

func initCommands() (err error) {

	tree = &dgocmd.Command{
		Name:        "Mention",
		Description: "Mention me to start a command!",
		Shim:        me.Mention(),
		Commands: []*dgocmd.Command{
			&dgocmd.Command{
				Name:        "Start Call",
				Description: "Start a call with another channel.",
				Shim:        "start call",
				Alias:       []string{"call"},
				Function:    commandStartCall,
			},
			&dgocmd.Command{
				Name:        "Stop Call",
				Description: "Stop a call with another channel.",
				Shim:        "stop call",
				Alias:       []string{"hang up"},
				Function:    commandStopCall,
			},
			&dgocmd.Command{
				Name:        "Set",
				Description: "Set a configuration within the bot.",
				Shim:        "set",
				Commands: []*dgocmd.Command{
					&dgocmd.Command{
						Name:        "Embeds",
						Description: "Toggle embeds on / off",
						Shim:        "embeds",
						Alias:       []string{"rich embeds", "rich embed", "embed"},
						Commands: []*dgocmd.Command{
							&dgocmd.Command{
								Name:        "On",
								Description: "Turn embeds on",
								Shim:        "on",
								Alias:       []string{"on", "yes", "enable"},
								//Function:    commandStopCall,
							},
							&dgocmd.Command{
								Name:        "Embed",
								Description: "Turn embeds off",
								Shim:        "embed",
								Alias:       []string{"off", "no", "disable"},
								//Function:    commandStopCall,
							},
						},
					},
				},
			},
		},
	}

	return
}
