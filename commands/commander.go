package commands

import (
	"regexp"

	dg "github.com/bwmarrin/discordgo"
)

func Commander() func(s *dg.Session, m *dg.MessageCreate) {
	return func(s *dg.Session, m *dg.MessageCreate) {
		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}
		if match, _ := regexp.MatchString("![a-zA-Z]+[ a-zA-Z\"]*[ 0-9/]*", m.Content); match {
			if match, _ := regexp.MatchString("!ping", m.Content); match {
				Ping(s, m)
				return
			}

			if match, _ := regexp.MatchString("(?i)^!stock [a-zA-Z]+$", m.Content); match {
				Stock(s, m)
				return
			}

			if match, _ := regexp.MatchString("(?i)^!er [a-zA-Z]+$", m.Content); match {
				Er(s, m)
				return
			}

			if match, _ := regexp.MatchString("(?i)^!wizdaddy$", m.Content); match {
				Wizdaddy(s, m)
				return
			}

			if match, _ := regexp.MatchString("(?i)^!coin [a-zA-Z]+$", m.Content); match {
				Coin(s, m)
				return
			}

			s.ChannelMessageSend(m.ChannelID, "Invalid Command")

		}
	}
}
