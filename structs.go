package main

// Database hosts the bot's database configuration.
type Database struct {
	User     string
	Password string
	Address  string
	Port     int
	Database string
}

// Discord hosts the bot's Discord configuration.
type Discord struct {
	Token    string
	MasterID string
}

// PhoneCall is a call from a channel to another channel.
type PhoneCall struct {
	From string
	To   string
}
