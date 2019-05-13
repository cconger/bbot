package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	DiscordClient *discordgo.Session
	Commands      []Command
}

// Command is the interface you should implement to create a command that can be directly invoked by users.
type Command interface {
	Match(m *discordgo.Message) (bool, int)
	Run(s *discordgo.Session, m *discordgo.Message) error
	String() string
}

func (b *Bot) Listen() (func() error, error) {
	b.DiscordClient.AddHandler(b.handleMessage)

	err := b.DiscordClient.Open()
	if err != nil {
		return nil, err
	}

	return b.DiscordClient.Close, nil
}

func (b *Bot) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		// We wrote this msg
		return
	}

	var bestDepth int
	var bestMatch Command
	for _, cmd := range b.Commands {
		match, depth := cmd.Match(m.Message)
		if !match || depth < bestDepth {
			// No match, or worse match
			continue
		}

		if depth == bestDepth && bestMatch != nil {
			log.Printf("Error: two commands matched the message equally well: String: %s, CmdA: %s, CmdB: %s", cmd, bestMatch)
		}

		bestDepth = depth
		bestMatch = cmd
	}

	if bestMatch == nil {
		// We don't need to handle every message...
		return
	}

	err := bestMatch.Run(s, m.Message)
	if err != nil {
		log.Printf("Error executing cmd %s: %s", bestMatch, err)
		s.ChannelMessageSend(m.ChannelID, "Encountered an error fulfilling your request.")
	}
}
