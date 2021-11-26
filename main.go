package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func main() {
	discord, err := discordgo.New("Bot " + "authentication token")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(discord)
}