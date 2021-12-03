package devhubbot

import (
	"bot/pkg/infra"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// This function will be called (due to AddHandler above) every time a new
// guild is joined.
func (b *Bot) guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "🤖 devhubbot reporting for duty")
			return
		}
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if strings.HasPrefix(m.Content, "!streakcurrent") {

		username := strings.TrimSpace(strings.Replace(m.Content, "!streakcurrent", "", 1))

		// Find the channel that the message came from.
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}

		currentStreak, err := b.githubService.GetCurrentContributionStreakByUsername(ctx, username)
		if err != nil {
			infra.Logger.Error().Err(err).Msg("github service get current contribution streak by username")

			_, _ = s.ChannelMessageSend(c.ID, fmt.Sprintf("something went wrong retrieving current streak for user %s", username))

			return
		}

		_, _ = s.ChannelMessageSend(c.ID, fmt.Sprintf("user %s %s", username, currentStreak.String()))

	} else if strings.HasPrefix(m.Content, "!streaklongest") {

		username := strings.TrimSpace(strings.Replace(m.Content, "!streaklongest", "", 1))

		// Find the channel that the message came from.
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}

		longestStreak, err := b.githubService.GetLongestContributionStreakByUsername(ctx, username)
		if err != nil {
			infra.Logger.Error().Err(err).Msg("github service get longest contribution streak by username")

			_, _ = s.ChannelMessageSend(c.ID, fmt.Sprintf("something went wrong retrieving longest streak for user %s", username))

			return
		}

		_, _ = s.ChannelMessageSend(c.ID, fmt.Sprintf("user %s %s", username, longestStreak.String()))

	}
}
