package discord

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestNewDiscordService(t *testing.T) {
	assert := assert.New(t)

	session := &discordgo.Session{}

	resp := NewDiscordService(session)

	assert.Equal(session, resp.session)
}

type DiscordServiceMock struct {

}

func(d *DiscordServiceMock) Open() error {
	return nil
}


func TestNewDiscordService_Open(t *testing.T) {
	assert := assert.New(t)

	test := &DiscordServiceMock{}

	test.Open()

	session := &discordgo.Session{}

	resp := NewDiscordService(session)

	assert.Equal(session, resp.session)
}