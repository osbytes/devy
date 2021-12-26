package main

import (
	"bot/internal/devy"
	"bot/internal/github"
	"bot/pkg/colors"
	"bot/pkg/env"
	"bot/pkg/infra"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const banner = `
 ______   _______  __   __  __   __
|      | |       ||  | |  ||  | |  |
|  _    ||    ___||  |_|  ||  |_|  |
| | |   ||   |___ |       ||       |
| |_|   ||    ___||       ||_     _|
|       ||   |___  |     |   |   |
|______| |_______|  |___|    |___|
`

func main() {

	env.Env = env.GetString("ENV", "local")

	infra.InitLogging(env.GetString("LOG_LEVEL", "info"))

	if env.IsLocal() {
		fmt.Println(colors.Purple, banner, colors.Reset)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	discord, err := discordgo.New("Bot " + env.GetString("DISCORD_BOT_TOKEN", ""))
	if err != nil {
		infra.Logger.Fatal().Err(err).Msg("discordgo new")
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: env.GetString("GITHUB_TOKEN", "")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	githubService := github.NewGithubService(client)

	bot := devy.NewBot(discord, githubService)

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
