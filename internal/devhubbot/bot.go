package devhubbot

import (
	"bot/internal/discord"
	"bot/pkg/infra"
	"context"

	"github.com/pkg/errors"
)

type Bot struct {
	discordSvc discord.DiscordServicer
}

func NewBot(discordSvc discord.DiscordServicer) *Bot {
	return &Bot{
		discordSvc: discordSvc,
	}
}

func (b Bot) Start(ctx context.Context) error {
	infra.Logger.Info().Msg("starting bot")

	err := b.discordSvc.Open()
	if err != nil {
		return errors.Wrap(err, "discord service open")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b Bot) Stop() error {
	infra.Logger.Info().Msg("stopping bot")

	err := b.discordSvc.Close()
	if err != nil {
		return errors.Wrap(err, "discord service close")
	}

	return nil
}
