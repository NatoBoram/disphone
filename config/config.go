package config

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

	// Calls : List of calls
	Calls map[string][]string
)

// ReadToken : Reads the whole configuration.
func ReadToken() error {

	// Read the config file
	file, err := ioutil.ReadFile("./token.json")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Json -> String
	err = json.Unmarshal(file, &Token)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

// ReadCalls : Reads the list of calls.
func ReadCalls() error {

	// Read the calls
	file, err := ioutil.ReadFile("./calls.json")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Calls = make(map[string][]string)

	// Json -> String
	err = json.Unmarshal(file, &Calls)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if len(Calls) == 0 {
		Calls = make(map[string][]string)
	}

	return nil
}

// WriteCalls : Writes the calls that are saved in the config
func WriteCalls() error {

	// From Calls to JSON
	json, err := json.Marshal(Calls)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// From JSON to file
	err = ioutil.WriteFile("./calls.json", json, os.FileMode(int(0777)))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

// Clean : Removes non-existent channels from the list of calls.
func Clean(s *discordgo.Session) {

	// For each from
	for key, array := range Calls {

		// Check if channel exists
		from, err := s.State.Channel(key)
		if err != nil {
			delete(Calls, key)
			fmt.Println("Removed Call : " + key)
		} else {

			// For each to
			for _, value := range array {

				// Check if channel exists
				_, err := s.State.Channel(value)
				if err != nil {
					Calls[key] = rsfa(array, value)
					fmt.Println("Removed " + key + " from " + from.Name + ".")
				}
			}
		}
	}
}

// rsfa : Remove String From Array.
func rsfa(a []string, s string) []string {
	var n []string
	for i, v := range a {
		if v != s {
			n = append(n, a[i])
		}
	}
	return n
}
