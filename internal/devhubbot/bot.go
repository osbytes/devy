package devhubbot

import (
	"bot/internal/github"
	"bot/internal/quotes"
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type Bot struct {
	discord       *discordgo.Session
	quoteService  quotes.QuoteServicer
	githubService github.GithubServicer
}

func NewBot(discord *discordgo.Session, quoteService quotes.QuoteServicer, githubService github.GithubServicer) *Bot {
	return &Bot{
		discord:       discord,
		quoteService:  quoteService,
		githubService: githubService,
	}
}

func (b Bot) Start(ctx context.Context) error {
	err := b.discord.Open()
	if err != nil {
		return errors.Wrap(err, "discord service open")
	}

	b.discord.AddHandler(b.guildCreate)

	b.discord.AddHandler(b.messageCreate)

	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b Bot) Stop() error {
	err := b.discord.Close()
	if err != nil {
		return errors.Wrap(err, "discord service close")
	}

	return nil
}
