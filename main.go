package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/badger"
)

func main() {

	// License
	fmt.Println("")
	fmt.Println("Discord-Phone : Makes phone calls between Discord servers.")
	fmt.Println("Copyright © 2019 Nato Boram")
	fmt.Println("This program is free software : you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY ; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with this program. If not, see http://www.gnu.org/licenses/.")
	fmt.Println("Contact : https://github.com/NatoBoram/Discord-Phone")
	fmt.Println("")

	config, err := getConfig(getDir())
	if err != nil {
		fmt.Println("Couldn't get the bot's configuration.")
		log.Fatalln(err.Error())
	}

	db = initDB(config)
	router := initRouter()
	initDiscord(config, router)

	// Wait for future input
	<-make(chan struct{})
}

func initDB(config *Config) (db *badger.DB) {

	// Open the Badger database located in the ~/.config/discord_phone directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions
	opts.Dir = config.Directory
	opts.ValueDir = config.Directory
	db, err := badger.Open(opts)
	if err != nil {
		fmt.Println("Couldn't open the database.")
		log.Fatal(err)
	}

	// Close the DB when the program is interrupted.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println(sig.String())
			db.Close()
			os.Exit(0)
		}
	}()

	return db
}

func initDiscord(config *Config, router *exrouter.Route) {

	// Go online!
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Couldn't get online.")
		log.Fatalln(err.Error())
		return
	}

	// Get Bot ID
	_, err = session.User("@me")
	if err != nil {
		fmt.Println("Couldn't get the BotID.")
		log.Fatalln(err.Error())
		return
	}

	// Hey, listen!
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		router.FindAndExecute(s, "", s.State.User.ID, m.Message)
	})

	// Crash on error
	err = session.Open()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	// It's alive!
	fmt.Println("Discord-Phone is running!")
}

func initRouter() *exrouter.Route {

	router := exrouter.New()

	// Call
	router.On("call", func(ctx *exrouter.Context) {

		// Get the calling guild
		guildFrom, err := ctx.Guild(ctx.Msg.GuildID)
		if err != nil {
			logctx(ctx, "Couldn't get the calling guild.", err)
			return
		}

		// Check for guild owner
		if guildFrom.OwnerID != ctx.Msg.Author.ID {
			ctx.Reply("Only the server owner can issue this command.")
			return
		}

		// Get current channel
		channelFrom, err := ctx.Channel(ctx.Msg.ChannelID)
		if err != nil {
			logctx(ctx, "Couldn't get the calling channel.", err)
			return
		}

		// For each arguments
		if len(ctx.Args) > 1 {
			for i := 1; i < len(ctx.Args); i++ {
				var duplicate bool

				// Check if this channel exists
				channelTo, err := ctx.Channel(ctx.Args[i])
				if err != nil {
					ctx.Reply("Channel ", ctx.Args[i], " doesn't exist.")
					continue
				}

				// Get the called guild
				guildTo, err := ctx.Guild(channelTo.GuildID)
				if err != nil {
					logctx(ctx, "Couldn't get the called guild.", err)
					return
				}

				// Update the database
				err = db.Update(func(txn *badger.Txn) (err error) {

					bfrom, err := encodeString(channelFrom.ID)
					if err != nil {
						logctx(ctx, "Couldn't encode a Channel ID.", nil)
						return
					}

					item, err := txn.Get(bfrom)
					if err != nil && err != badger.ErrKeyNotFound {
						logctx(ctx, "Couldn't get the called channels.", nil)
						return
					} else if err == badger.ErrKeyNotFound {
						return insertCall(ctx, channelFrom.ID, []string{channelTo.ID}, txn)
					}

					fmt.Println("Adding this channel to calls")

					duplicate, err = appendCall(ctx, channelFrom.ID, channelTo.ID, txn, item)
					if err != nil {
						logctx(ctx, "Couldn't append a call to this channel's calls.", nil)
						return
					}

					return txn.Commit()
				})
				if err != nil {
					logctx(ctx, "Couldn't update the database.", err)
					return
				}

				// Craft a mention from memberFrom and channelTo.
				mentionMemberFrom := "**" + ctx.Msg.Author.Username + "**"
				mentionChannelTo := "**" + channelTo.Name + "**"
				memberFrom, err := ctx.Ses.GuildMember(channelTo.GuildID, ctx.Msg.Author.ID)
				if err == nil {
					mentionMemberFrom = memberFrom.Mention()
					mentionChannelTo = channelTo.Mention()
				} else {
					err = nil
				}

				// Craft a mention from channelFrom.
				mentionChannelFrom := "**" + channelFrom.Name + "**"
				memberTo, err := ctx.Ses.GuildMember(channelFrom.GuildID, guildTo.OwnerID)
				if err == nil {
					mentionChannelFrom = channelFrom.Mention()
				} else {
					err = nil
				}

				if duplicate {
					ctx.Reply("There's already a call towards ", mentionChannelTo, ".")
					return
				}

				// Inform the other channel that I'm calling.
				ctx.Ses.ChannelMessageSend(channelTo.ID, memberTo.Mention()+" *Ring, ring!* "+mentionMemberFrom+" is calling from "+mentionChannelFrom+".")
				ctx.Reply("Calling ", mentionChannelTo, "...")
			}
		} else {
			ctx.Reply("You must enter one or more Channel ID to initiate a call.")
		}
	}).Desc("Issue a call from this channel to another channel.")

	// End Call
	router.On("end", func(ctx *exrouter.Context) {

		// Create the call
		ctx.Reply("Let me check if this channel exists.")
		fmt.Println("A new call should be created.")

	}).Desc("End a call from this channel to another channel.")

	// Help
	router.Default = router.On("help", func(ctx *exrouter.Context) {
		var f func(depth int, r *exrouter.Route) string
		f = func(depth int, r *exrouter.Route) string {
			text := ""
			for _, v := range r.Routes {
				text += strings.Repeat("  ", depth) + v.Name + " : " + v.Description + "\n"
				text += f(depth+1, &exrouter.Route{Route: v})
			}
			return text
		}
		ctx.Reply("```" + f(0, router) + "```")
	}).Desc("prints this help menu")

	return router
}

func logctx(ctx *exrouter.Context, msg string, err error) {
	var serr string
	if err != nil {
		serr = "\n```go\n" + err.Error() + "\n```"
	}
	ctx.Reply("**Error :** ", msg, serr)
	fmt.Println(msg, "\n", err.Error())
}
