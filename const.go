package main

import "os"

// Paths
const (
	rootFolder   = "./DiscordPhone"
	discordPath  = rootFolder + "/discord.json"
	databasePath = rootFolder + "/database.json"
)

// Permissions
const (
	permPrivateDirectory os.FileMode = 0700
	permPrivateFile      os.FileMode = 0600
)
