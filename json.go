package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func readDatabase(object *Database) error {

	// Read the JSON file
	file, err := ioutil.ReadFile(databasePath)
	if err != nil {
		return err
	}

	// Put the JSON in the object
	err = json.Unmarshal(file, &object)
	if err != nil {
		return err
	}

	return nil
}

func readDiscord(object *Discord) error {

	// Read the JSON file
	file, err := ioutil.ReadFile(discordPath)
	if err != nil {
		return err
	}

	// Put the JSON in the object
	err = json.Unmarshal(file, &object)
	if err != nil {
		return err
	}

	return nil
}

func writeDatabase(object Database) error {

	// Create required directories
	os.MkdirAll(rootFolder, permPrivateDirectory)

	// From object to JSON
	json, err := json.Marshal(object)
	if err != nil {
		return err
	}

	// From JSON to File
	err = ioutil.WriteFile(databasePath, json, permPrivateFile)
	if err != nil {
		return err
	}

	return nil
}

func writeDiscord(object Discord) error {

	// Create required directories
	os.MkdirAll(rootFolder, permPrivateDirectory)

	// From object to JSON
	json, err := json.Marshal(object)
	if err != nil {
		return err
	}

	// From JSON to File
	err = ioutil.WriteFile(discordPath, json, permPrivateFile)
	if err != nil {
		return err
	}

	return nil
}

func writeTemplateDatabase() error {
	fmt.Println("Writing a new database configuration template...")
	return writeDatabase(Database{
		User:     "DiscordPhone",
		Address:  "localhost",
		Port:     3306,
		Database: "DiscordPhone",
	})
}

func writeTemplateDiscord() error {
	fmt.Println("Writing a new Discord configuration template...")
	return writeDiscord(Discord{})
}
