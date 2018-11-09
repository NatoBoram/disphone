package main

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
)

func createTableCalls() (res sql.Result, err error) {
	return db.Exec("create table if not exists `calls` (`from` varchar(32) not null, `to` varchar(32) not null, constraint `pk_calls` primary key (`from`, `to`)) default charset=utf8mb4 collate=utf8mb4_unicode_ci;")
}

func createTables() (res sql.Result, err error) {

	// Declare tables to create
	functs := [...]func() (res sql.Result, err error){
		createTableCalls,
	}

	// Create the tables
	for _, funct := range functs {
		res, err = funct()
		if err != nil {
			return
		}
	}

	return
}

func insertCall(from *discordgo.Channel, to *discordgo.Channel) (res sql.Result, err error) {
	return db.Exec("insert into `calls`(`from`, `to`) values(?, ?);", from.ID, to.ID)
}

func selectCall(from *discordgo.Channel, to *discordgo.Channel) (call PhoneCall, err error) {
	err = db.QueryRow("select `from`, `to` from `calls` where `from` = ? and `to` = ?;", from.ID, to.ID).Scan(&call.From, &call.To)
	return
}

func selectCalls(from *discordgo.Channel) (rows *sql.Rows, err error) {
	return db.Query("select `from`, `to` from `calls` where `from` = ?;", from.ID)
}

func deleteCall(from string, to string) (res sql.Result, err error) {
	return db.Exec("delete from `calls` where `from` = ? and `to` = ?;", from, to)
}

func deleteCalls(from string) (res sql.Result, err error) {
	return db.Exec("delete from `calls` where `from` = ? or `to` = ?;", from, from)
}
