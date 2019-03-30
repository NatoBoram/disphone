package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/dgraph-io/badger"
)

// Config is the base configuration of this bot.
type Config struct {
	Token     string
	MasterID  string
	Directory string
}

// Permissions
const (
	permPrivateDirectory os.FileMode = 0700
	permPrivateFile      os.FileMode = 0600
)

// State
var (
	db *badger.DB
)

func getDir() (dir string) {

	// Get the directory
	dir = os.Getenv("DISCORD_PHONE_DB")
	if dir == "" {
		current, err := user.Current()
		if err != nil {
			fmt.Println("Couldn't get the current user.")
			log.Fatalln(err.Error())
			return
		}
		dir = path.Join(current.HomeDir, ".config", "DiscordPhone")
	}

	// Create directories if they don't exist.
	err := os.MkdirAll(dir, permPrivateDirectory)
	if err != nil {
		fmt.Println("Couldn't create the directory.")
		log.Fatalln(err.Error())
		return
	}
	return
}

func getConfig(dir string) (config *Config, err error) {

	// Read the JSON file
	file, err := ioutil.ReadFile(path.Join(dir, "config.json"))
	if err != nil {
		return
	}

	// Put the JSON in the object
	err = json.Unmarshal(file, &config)
	if err != nil {
		return
	}

	// Replace the directory in the config
	config.Directory = dir

	return
}
