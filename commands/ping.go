package commands

import (
	dg "github.com/bwmarrin/discordgo"
)

func Ping(s *dg.Session, m *dg.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "pong!")
}
