package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type DiscordServicer interface {
	Open() error
	Close() error
}

var _ DiscordServicer = (*DiscordService)(nil)

type DiscordService struct {
	session *discordgo.Session
}

func NewDiscordService(session *discordgo.Session) *DiscordService {
	return &DiscordService{
		session: session,
	}
}

func (d *DiscordService) Open() error {
	err := d.session.Open()
	if err != nil {
		return errors.Wrap(err, "discordgo session open")
	}

	return nil
}

func (d *DiscordService) Close() error {
	err := d.session.Close()
	if err != nil {
		return errors.Wrap(err, "discordgo session close")
	}

	return nil
}
