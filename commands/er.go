package commands

import (
	"strings"

	"github.com/BryanSLam/discord-bot/config"
	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
	iex "github.com/jonwho/go-iex"
)

func Er(s *dg.Session, m *dg.MessageCreate) {
	slice := strings.Split(m.Content, " ")
	ticker := slice[1]
	iexClient := iex.NewClient()
	earnings, err := iexClient.Earnings(ticker)
	logger := util.Logger{Session: s, ChannelID: config.GetConfig().BotLogChannelID}

	logger.Info("Fetching earnings report info for " + ticker)
	if err != nil {
		logger.Trace("IEX request failed. Message: " + err.Error())
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	message := util.FormatEarnings(earnings)

	s.ChannelMessageSend(m.ChannelID, message)
}
