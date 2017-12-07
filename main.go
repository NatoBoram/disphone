package main

import (
	"fmt"

	"github.com/NatoBoram/Discord-Phone/bot"
	"github.com/NatoBoram/Discord-Phone/config"
)

func main() {

	// Reads the token
	err := config.ReadToken()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Reads the active calls
	err = config.ReadCalls()
	if err != nil {
		fmt.Println(err.Error())
		config.Calls = make(map[string][]string)
		config.WriteCalls()
	}

	// License
	fmt.Println("")
	fmt.Println("Discord-Phone : Makes phone calls between Discord servers.")
	fmt.Println("Copyright © 2017 Nato Boram")
	fmt.Println("This program is free software : you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY ; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with this program. If not, see http://www.gnu.org/licenses/.")
	fmt.Println("Contact : https://github.com/NatoBoram/Discord-Phone")
	fmt.Println("")

	// Give this bot some life!
	bot.Start()

	// Wait for future input
	<-make(chan struct{})
	return
}
