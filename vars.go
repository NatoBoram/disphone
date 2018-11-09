package main

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/NatoBoram/Discord-Phone/dgocmd"
)

var (
	db      *sql.DB
	me      *discordgo.User
	master  *discordgo.User
	session *discordgo.Session
	tree    *dgocmd.Command
)
