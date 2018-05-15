package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (

	// Token : Client Secret
	Token string

	// BotMaster : ID of the BotMaster
	BotMaster string

	// Calls : List of calls
	Calls map[string][]string
)

// ReadToken : Reads the token.
func ReadToken() error {

	// Read the config file
	file, err := ioutil.ReadFile(TokenPath)
	if err != nil {
		return err
	}

	// Json -> String
	err = json.Unmarshal(file, &Token)
	if err != nil {
		return err
	}

	return nil
}

// ReadBotMaster : Reads the BotMaster ID
func ReadBotMaster() error {

	// Read the config file
	file, err := ioutil.ReadFile(MasterPath)
	if err != nil {
		return err
	}

	// Json -> String
	err = json.Unmarshal(file, &BotMaster)
	if err != nil {
		return err
	}

	return nil
}

// ReadCalls : Reads the list of calls.
func ReadCalls() error {

	Calls = make(map[string][]string)

	// Read the calls
	file, err := ioutil.ReadFile(CallsPath)
	if err != nil {
		return err
	}

	// Json -> String
	err = json.Unmarshal(file, &Calls)
	if err != nil {
		Calls = make(map[string][]string)
		return err
	}

	if len(Calls) == 0 {
		Calls = make(map[string][]string)
	}

	return nil
}

// WriteToken : Writes the token.
func WriteToken() error {

	// From Calls to JSON
	json, err := json.Marshal(Token)
	if err != nil {
		return err
	}

	// From JSON to file
	err = ioutil.WriteFile(TokenPath, json, os.FileMode(int(0777)))
	if err != nil {
		return err
	}

	return nil
}

// WriteBotMaster : Writes the BotMaster's ID.
func WriteBotMaster() error {

	// From Calls to JSON
	json, err := json.Marshal(BotMaster)
	if err != nil {
		return err
	}

	// From JSON to file
	err = ioutil.WriteFile(MasterPath, json, os.FileMode(int(0777)))
	if err != nil {
		return err
	}

	return nil
}

// WriteCalls : Writes the calls that are saved in the config
func WriteCalls() error {

	// From Calls to JSON
	json, err := json.Marshal(Calls)
	if err != nil {
		return err
	}

	// From JSON to file
	err = ioutil.WriteFile(CallsPath, json, os.FileMode(int(0777)))
	if err != nil {
		return err
	}

	return nil
}

// ReadAll reads all the config and fixes mistakes.
func ReadAll() error {
	os.Mkdir(Folder, os.FileMode(int(0777)))

	// Token
	err := ReadToken()
	if err != nil {
		err = WriteToken()
		if err != nil {
			return err
		}
	}

	// BotMaster
	err = ReadBotMaster()
	if err != nil {
		err = WriteBotMaster()
		if err != nil {
			return err
		}
	}

	// Calls
	err = ReadCalls()
	if err != nil {
		err = WriteCalls()
		if err != nil {
			return err
		}
	}

	return err
}

// WriteAll writes all the config.
func WriteAll() {
	WriteBotMaster()
	WriteToken()
	WriteCalls()
}

// Clean : Removes non-existent channels from the list of calls.
func Clean(s *discordgo.Session) {

	clean := false

	// For each from
	for key, array := range Calls {

		// Check if channel exists
		from, err := s.State.Channel(key)
		if err != nil {
			delete(Calls, key)
			fmt.Println("Removed Call : " + key)
			clean = true
		} else {

			// For each to
			for _, value := range array {

				// Check if channel exists
				_, err := s.State.Channel(value)
				if err != nil {
					Calls[key] = rsfa(Calls[key], value)
					fmt.Println("Removed Call : " + from.Name + " / " + key + ".")
					clean = true
				}
			}
		}
	}

	// Save
	if clean {
		WriteCalls()
	}
}
