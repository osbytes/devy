package devhubbot

import (
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
