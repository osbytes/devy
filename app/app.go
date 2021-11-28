package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"bot/envars"

	"github.com/bwmarrin/discordgo"
)

var (
	Discord *discordgo.Session
	err error
)

func Run() {
	Discord, err = discordgo.New("Bot " + envars.DISCORD_BOT_TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	err = Discord.Open()
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
