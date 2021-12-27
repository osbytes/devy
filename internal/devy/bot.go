package devy

import (
	"bot/internal/github"
	"bot/pkg/emoji"
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

var (
	pollPrefix = fmt.Sprintf("%s **POLL**", emoji.Emojis[":memo:"])
)

type Bot struct {
	discord       *discordgo.Session
	githubService github.GithubServicer
}

func NewBot(discord *discordgo.Session, githubService github.GithubServicer) *Bot {
	return &Bot{
		discord:       discord,
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

	b.discord.AddHandler(b.messageReactionAdd)

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
