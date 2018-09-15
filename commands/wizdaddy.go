package commands

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BryanSLam/discord-bot/datasource"
	dg "github.com/bwmarrin/discordgo"
)

func Wizdaddy(s *dg.Session, m *dg.MessageCreate) {
	wizdaddyURL := "http://dev.wizdaddy.io/api/giveItToMeDaddy"

	resp, err := http.Get(wizdaddyURL)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Daddy is down")
		return
	}

	var daddyResponse datasource.WizdaddyResponse
	if err = json.NewDecoder(resp.Body).Decode(&daddyResponse); err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID,
		fmt.Sprintf("%s %s %s %s", daddyResponse.Symbol,
			daddyResponse.StrikePrice, daddyResponse.ExpirationDate, daddyResponse.Type))
}
