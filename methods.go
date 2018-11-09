package main

import (
	"strconv"
)

func (discord Discord) isEmpty() bool {
	return discord.MasterID == "" || discord.Token == ""
}

func (database Database) isEmpty() bool {
	return database.Address == "" ||
		database.Database == "" ||
		database.Password == "" ||
		database.Port == 0 ||
		database.User == ""
}

func (database Database) toConnectionString() string {
	return database.User + ":" + database.Password + "@tcp(" + database.Address + ":" + strconv.Itoa(database.Port) + ")/" + database.Database
}
