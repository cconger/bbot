package finance

import (
	"strings"
	"time"

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

	quoteTime := time.Unix(quote.LatestUpdate/1000, 0)
	color := 0xababab
	if quote.Change.IsNegative() {
		color = 0xff0000
	}
	if quote.Change.IsPositive() {
		color = 0x00ff00
	}

	//s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: $%f as of %s", quote.CompanyName, quote.LatestPrice, quote.LatestTime))
	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:     quote.CompanyName,
			Timestamp: quoteTime.Format(time.RFC3339),
			Color:     color,
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Price", Value: quote.LatestPrice.String()},
				{Name: "Change", Value: quote.Change.String()},
				{Name: "High", Value: quote.High.String()},
				{Name: "Low", Value: quote.Low.String()},
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (q *Quote) String() string {
	return "Quote"
}
