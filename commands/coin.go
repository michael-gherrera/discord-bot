package commands

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/BryanSLam/discord-bot/datasource"
	dg "github.com/bwmarrin/discordgo"
)

func Coin(s *dg.Session, m *dg.MessageCreate) {
	slice := strings.Split(m.Content, " ")
	ticker := strings.ToUpper(slice[1])
	coinURL := "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=" + ticker + "&tsyms=USD"

	resp, err := http.Get(coinURL)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	coin := datasource.Coin{Symbol: ticker}

	if err = json.NewDecoder(resp.Body).Decode(&coin); err != nil || coin.Response == "Error" {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, coin.OutputJSON())
	defer resp.Body.Close()
}
