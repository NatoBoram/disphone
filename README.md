# Discord-Phone

[![Build Status](https://travis-ci.org/NatoBoram/Discord-Phone.svg?branch=master)](https://travis-ci.org/NatoBoram/Discord-Phone)
[![Go Report Card](https://goreportcard.com/badge/github.com/NatoBoram/Discord-Phone)](https://goreportcard.com/report/github.com/NatoBoram/Discord-Phone)

Simple phone call bot for Discord using [DiscordGo](https://github.com/bwmarrin/discordgo).

You can invite this bot by clicking [here](https://discordapp.com/api/oauth2/authorize?client_id=384692861314007040&permissions=379969&scope=bot). Just be aware it isn't well-polished and you **will** encounter shenanigans.

## Usage

Only server owners can issue commands. Commands are `call ChannelID` and `hang up ChannelID`. Use Discord's [developer mode](https://support.discordapp.com/hc//articles/206346498) to get a channel's ID.

## Installation

```SH
go get -fix -u github.com/NatoBoram/Discord-Phone
```

You need to create a `token.json` file next to `main.go`. Get your token [here](https://discordapp.com/developers/applications/me) and paste it in `token.json` like this :

```JSON
"Mzg3Njk1ODcyOTU3MjE4ODE3.DQiNbw.6Fl3teDG1ieDxcFomfTt8UvnDTY"
```