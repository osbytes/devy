package main

import (
	"bot/internal/devhubbot"
	"bot/internal/discord"
	"bot/internal/quotes"
	"bot/pkg/colors"
	"bot/pkg/env"
	"bot/pkg/infra"
	"bot/pkg/universalinspirationalquotes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const banner = `
     _            _           _     _           _
  __| | _____   _| |__  _   _| |__ | |__   ___ | |_
 / _  |/ _ \ \ / / '_ \| | | | '_ \| '_ \ / _ \| __|
| (_| |  __/\ V /| | | | |_| | |_) | |_) | (_) | |_ 
 \__,_|\___| \_/ |_| |_|\__,_|_.__/|_.__/ \___/ \__|
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

	discordService := discord.NewDiscordService(session)

	var quoteService quotes.QuoteServicer = &quotes.NOOPQuoteService{}

	rapidAPIKeyUniversalInspirationalQuotes := env.GetString("RAPID_API_KEY_UNIVERSAL_INSPIRATIONAL_QUOTES", "")
	if len(rapidAPIKeyUniversalInspirationalQuotes) > 0 {
		quotesClient := universalinspirationalquotes.NewHTTPClient(&http.Client{}, rapidAPIKeyUniversalInspirationalQuotes)

		quoteService = quotes.NewQuoteService(quotesClient)
	}

	bot := devhubbot.NewBot(devhubbot.BotOpts{
		DiscordService: discordService,
		QuoteService:   quoteService,
	})

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
