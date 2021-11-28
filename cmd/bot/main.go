package main

import (
	"bot/internal/devhubbot"
	"bot/internal/discord"
	"bot/pkg/env"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session, err := discordgo.New("Bot " + env.GetString("DISCORD_BOT_TOKEN", ""))
	if err != nil {
		log.Fatal(err)
	}

	discordSvc := discord.NewDiscordService(session)
	bot := devhubbot.NewBot(discordSvc)

	go func() {
		err := bot.Start(ctx)
		if err != nil && err != context.Canceled {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	fmt.Println("Listening for interupt and term signals")

	<-stop

	cancel()

	bot.Stop()

}
