package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var token string

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	token = os.Getenv("DISCORD_BOT_TOKEN")
}

func main() {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println(err)
	}

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}