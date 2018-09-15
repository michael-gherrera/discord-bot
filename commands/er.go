package commands

import (
	"strings"

	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
	iex "github.com/jonwho/go-iex"
)

func Er(s *dg.Session, m *dg.MessageCreate) {
	slice := strings.Split(m.Content, " ")
	ticker := slice[1]
	iexClient := iex.NewClient()
	earnings, err := iexClient.Earnings(ticker)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	message := util.FormatEarnings(earnings)

	s.ChannelMessageSend(m.ChannelID, message)
}
