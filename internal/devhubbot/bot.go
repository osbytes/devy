package devhubbot

import (
	"bot/internal/discord"
	"bot/internal/github"
	"bot/internal/quotes"
	"context"

	"github.com/pkg/errors"
)

type Bot struct {
	discordService discord.DiscordServicer
	quoteService   quotes.QuoteServicer
	githubService  github.GithubServicer
}

func NewBot(discordService discord.DiscordServicer, quoteService quotes.QuoteServicer, githubService github.GithubServicer) *Bot {
	return &Bot{
		discordService: discordService,
		quoteService:   quoteService,
		githubService:  githubService,
	}
}

func (b Bot) Start(ctx context.Context) error {
	err := b.discordService.Open()
	if err != nil {
		return errors.Wrap(err, "discord service open")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b Bot) Stop() error {
	err := b.discordService.Close()
	if err != nil {
		return errors.Wrap(err, "discord service close")
	}

	return nil
}
