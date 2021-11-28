package main

import (
	"bot/internal/devhubbot"
	"bot/internal/discord"
	"bot/pkg/colors"
	"bot/pkg/env"
	"bot/pkg/infra"
	"context"
	"fmt"
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

	env.Env = env.GetString("ENV", "local")

	infra.InitLogging(env.GetString("LOG_LEVEL", "info"))

	if env.IsLocal() {
		fmt.Println(colors.Purple, banner, colors.Reset)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session, err := discordgo.New("Bot " + env.GetString("DISCORD_BOT_TOKEN", ""))
	if err != nil {
		infra.Logger.Fatal().Err(err).Msg("discordgo new session")
	}

	discordSvc := discord.NewDiscordService(session)
	bot := devhubbot.NewBot(discordSvc)

	go func() {
		infra.Logger.Info().Msg("starting bot")

		err := bot.Start(ctx)
		if err != nil && err != context.Canceled {
			infra.Logger.Fatal().Err(err).Msg("bot start")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	infra.Logger.Info().Msg("Listening for interupt and term signals")

	<-stop

	cancel()

	infra.Logger.Info().Msg("stopping bot")
	err = bot.Stop()
	if err != nil {
		infra.Logger.Fatal().Err(err).Msg("bot stop")
	}

}
