package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/cconger/bbot/internal/iex"
	"github.com/cconger/bbot/pkg/discord"
	"github.com/cconger/bbot/pkg/discord/finance"
)

func main() {
	discordToken := os.GetEnv("DISCORD_TOKEN")
	iexToken := os.GetEnv("IEX_TOKEN")
	if discordToken == "" || iexToken == "" {
		// TODO: More helpful error... move loading to the actual commands so they can fail indepedently.
		log.Fatal("Missing required environment variables")
	}

	session, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatal("Unable to create a discord session")
	}

	iexClient := &iex.Client{
		HTTPClient: &http.Client{},
		Token:      iexToken,
		BaseURL:    "https://cloud.iexapis.com/beta",
	}

	bot := discord.Bot{
		DiscordClient: session,
		Commands: []discord.Command{
			&finance.Quote{IEXClient: iexClient},
		},
	}

	closeFn, err := bot.Listen()
	if err != nil {
		log.Fatal(err)
	}
	defer closeFn()

	log.Println("Bot started")
	closeChan := make(chan os.Signal)
	signal.Notify(closeChan, os.Interrupt, os.Kill)
	<-closeChan
}
