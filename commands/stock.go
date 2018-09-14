package commands

import (
	"strings"

	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
	iex "github.com/jonwho/go-iex"
)

func Stock(s *dg.Session, m *dg.MessageCreate) {
	slice := strings.Split(m.Content, " ")
	ticker := slice[1]
	iexClient := iex.NewClient()

	quote, err := iexClient.Quote(ticker, true)

	if err != nil {
		rds, iexErr := iexClient.RefDataSymbols()
		if iexErr != nil {
			s.ChannelMessageSend(m.ChannelID, iexErr.Error())
			return
		}

		fuzzySymbols := util.FuzzySearch(ticker, rds.Symbols)

		if len(fuzzySymbols) > 0 {
			fuzzySymbols = fuzzySymbols[:len(fuzzySymbols)%10]
			outputJSON := util.FormatFuzzySymbols(fuzzySymbols)
			s.ChannelMessageSend(m.ChannelID, outputJSON)
			return
		}
	}

	message := util.FormatQuote(quote)

	s.ChannelMessageSend(m.ChannelID, message)
}
