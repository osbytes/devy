package devy

import (
	"bot/pkg/infra"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// This function will be called (due to AddHandler above) every time a new
// guild is joined.
func (b *Bot) guildCreate(session *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = session.ChannelMessageSend(channel.ID, "🤖 devhubbot reporting for duty")
			return
		}
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func (b *Bot) messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if message.Author.ID == session.State.User.ID {
		return
	}

	// Find the channel that the message came from.
	channel, err := channelFromStateF(session.State, message.ChannelID)
	if err != nil {
		// Could not find channel.
		return
	}

	commandName := strings.Split(message.Content, " ")[0]

	if commandName == "!help" {
		commandUsages := []string{}

		for _, c := range commandMap {
			commandUsages = append(commandUsages, c.Usage())
		}

		usage := strings.Join(commandUsages, "\n\n")

		_, _ = channelMessageSend(session, channel.ID, usage)

		return
	}

	command, exists := commandMap[commandName]

	if !exists {
		return
	}

	command.Handler(session, message, channel, b)

}

func (b *Bot) messageReactionAdd(session *discordgo.Session, message *discordgo.MessageReactionAdd) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if message.UserID == session.State.User.ID {
		return
	}

	msg, err := channelMessageF(session, message.ChannelID, message.MessageID)
	if err != nil {
		infra.Logger.Error().Err(err).Msg("get channel message")

		return
	}

	if !strings.HasPrefix(msg.Content, pollPrefix) {
		return
	}

	// remove emoji if not one of poll emojis
	// TODO: need a better way to determine if the reaction emoji is one of the poll emojis. This code does not consider that the question itself may contain an emoji and that emoji would be allowed to be used in the reactions
	if !strings.Contains(msg.Content, message.Emoji.Name) {
		_ = messageReactionRemoveF(session, message.ChannelID, message.MessageID, message.Emoji.Name, message.UserID)

		return
	}

	for _, reaction := range msg.Reactions {
		if reaction.Emoji.Name == message.MessageReaction.Emoji.Name {
			continue
		}

		_ = messageReactionRemoveF(session, message.ChannelID, message.MessageID, reaction.Emoji.Name, message.UserID)
	}

}
