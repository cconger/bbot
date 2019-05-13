package finance

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cconger/bbot/internal/iex"
	"github.com/cconger/bbot/pkg/discord"
)

type Quote struct {
	IEXClient *iex.Client
}

// Verify that we implement the interface at compile time.
var _ discord.Command = &Quote{}

func (q *Quote) Match(m *discordgo.Message) (bool, int) {
	args := strings.Split(m.Content, " ")
	if len(args) != 2 || strings.ToLower(args[0]) != "!quote" {
		return false, 0
	}

	return true, 2
}

func (q *Quote) Run(s *discordgo.Session, m *discordgo.Message) error {
	args := strings.Split(m.Content, " ")
	symbol := args[1]
	quote, err := q.IEXClient.GetQuote(symbol)
	if err != nil {
		return err
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: $%f as of %s", quote.CompanyName, quote.LatestPrice, quote.LatestTime))
	return nil
}

func (q *Quote) String() string {
	return "Quote"
}
