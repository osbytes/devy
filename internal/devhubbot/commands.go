package devhubbot

import (
	"bot/pkg/infra"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	channelFromStateF   = channelFromState
	channelMessageSendF = channelMessageSend
)

func channelFromState(s *discordgo.State, channelID string) (*discordgo.Channel, error) {
	return s.Channel(channelID)
}

func channelMessageSend(s *discordgo.Session, channelID, message string) (*discordgo.Message, error) {
	return s.ChannelMessageSend(channelID, message)
}

type CommandHandler func(session *discordgo.Session, message *discordgo.MessageCreate, channel *discordgo.Channel, bot *Bot)

type Command struct {
	Name        string
	Description string
	Args        []string
	Handler     CommandHandler
}

func (c Command) Usage() string {
	commandUsage := fmt.Sprintf("**%s**", c.Name)

	if len(c.Args) > 0 {
		args := []string{}
		for _, a := range c.Args {
			args = append(args, fmt.Sprintf("{%s}", a))
		}

		commandUsage += fmt.Sprintf(" %s", strings.Join(args, " "))
	}

	return fmt.Sprintf("%s\n\t%s", commandUsage, c.Description)
}

var commandMap = map[string]Command{
	"!streakcurrent": {
		Name:        "!streakcurrent",
		Description: "get the current contribution streak of a github user",
		Args:        []string{"github username"},
		Handler:     streakCurrentCommandHandler,
	},
	"!streaklongest": {
		Name:        "!streaklongest",
		Description: "get the longest contribution streak of a github user",
		Args:        []string{"github username"},
		Handler:     streakLongestCommandHandler,
	},
	"!contributionstotal": {
		Name:        "!contributionstotal",
		Description: "get the all time total contribution of a github user",
		Args:        []string{"github username"},
		Handler:     contributionsTotalCommandHandler,
	},
	"!languages": {
		Name:        "!languages",
		Description: "get a breakdown (in bytes written per language) of all languages used committed to your repositories",
		Args:        []string{"github username"},
		Handler:     languagesCommandHandler,
	},
	"!lastrepo": {
		Name:        "!lastrepo",
		Description: "Get the latest repo the user has updated",
		Args:        []string{"github username"},
		Handler:     lastRepoCommandHandler,
	},
}

func streakCurrentCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate, channel *discordgo.Channel, bot *Bot) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	contentParts := strings.Split(strings.TrimSpace(message.Content), " ")
	if len(contentParts) <= 1 {
		_, _ = channelMessageSendF(session, channel.ID, "missing github username")

		return
	}

	username := contentParts[1]

	currentStreak, err := bot.githubService.GetCurrentContributionStreakByUsername(ctx, username)
	if err != nil {
		infra.Logger.Error().Err(err).Msg("github service get current contribution streak by username")

		_, _ = channelMessageSendF(session, channel.ID, fmt.Sprintf("something went wrong retrieving current streak for github user %s", username))

		return
	}

	_, _ = channelMessageSendF(session, channel.ID, fmt.Sprintf("user %s %s", username, currentStreak.String()))
}

func streakLongestCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate, channel *discordgo.Channel, bot *Bot) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	contentParts := strings.Split(strings.TrimSpace(message.Content), " ")
	if len(contentParts) <= 1 {
		_, _ = channelMessageSendF(session, channel.ID, "missing github username")

		return
	}

	username := contentParts[1]

	longestStreak, err := bot.githubService.GetLongestContributionStreakByUsername(ctx, username)
	if err != nil {
		infra.Logger.Error().Err(err).Msg("github service get longest contribution streak by username")

		_, _ = channelMessageSendF(session, channel.ID, fmt.Sprintf("something went wrong retrieving longest streak for github user %s", username))

		return
	}

	_, _ = channelMessageSendF(session, channel.ID, fmt.Sprintf("user %s %s", username, longestStreak.String()))
}

func contributionsTotalCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate, channel *discordgo.Channel, bot *Bot) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	contentParts := strings.Split(strings.TrimSpace(message.Content), " ")
	if len(contentParts) <= 1 {
		_, _ = channelMessageSendF(session, channel.ID, "missing github username")

		return
	}

	username := contentParts[1]

	totalContributions, err := bot.githubService.GetTotalContributionsByUsername(ctx, username)
	if err != nil {
		infra.Logger.Error().Err(err).Msg("github service get total contributions by username")

		_, _ = channelMessageSendF(session, channel.ID, fmt.Sprintf("something went wrong retrieving total contributions for user %s", username))

		return
	}

	_, _ = channelMessageSendF(session, channel.ID, fmt.Sprintf("user %s %s", username, totalContributions.String()))
}

func languagesCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate, channel *discordgo.Channel, bot *Bot) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	contentParts := strings.Split(strings.TrimSpace(message.Content), " ")
	if len(contentParts) <= 1 {
		_, _ = channelMessageSendF(session, channel.ID, "missing github username")

		return
	}

	username := contentParts[1]

	languages, err := bot.githubService.GetLanguagesByUsername(ctx, username)
	if err != nil {
		infra.Logger.Error().Err(err).Msg("github service get languages by username")

		_, _ = channelMessageSendF(session, channel.ID, fmt.Sprintf("something went wrong retrieving languages for user %s", username))

		return
	}

	_, _ = channelMessageSendF(session, channel.ID, fmt.Sprintf("user %s\n\n%s", username, languages.String()))
}

func lastRepoCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate, channel *discordgo.Channel, bot *Bot) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	contentParts := strings.Split(strings.TrimSpace(message.Content), " ")
	if len(contentParts) <= 1 {
		_, _ = channelMessageSendF(session, channel.ID, "missing github username")

		return
	}

	username := contentParts[1]

	lastRepo, err := bot.githubService.GetLastUpdatedRepoByUsername(ctx, username)
	if err != nil {
		infra.Logger.Error().Err(err).Msg("github service get the last repo updated by username")

		_, _ = channelMessageSendF(session, channel.ID, fmt.Sprintf("something went wrong retrieving last updated repo for github user %s", username))

		return
	}

	_, _ = channelMessageSendF(session, channel.ID, fmt.Sprintf("user %s %s", username, lastRepo.String()))
}
