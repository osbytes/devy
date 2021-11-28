package main

import (
	"bot/internal/devhubbot"
	"bot/internal/discord"
	"bot/pkg/colors"
	"bot/pkg/env"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const banner = `
       __          __          __          __
  ____/ /__ _   __/ /_  __  __/ /_  ____  / /_
 / __  / _ \ | / / __ \/ / / / __ \/ __ \/ __/
/ /_/ /  __/ |/ / / / / /_/ / /_/ / /_/ / /_
\__,_/\___/|___/_/ /_/\__,_/_.___/\____/\__/
`

func main() {

	fmt.Println(colors.Purple, banner, colors.Reset)

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
