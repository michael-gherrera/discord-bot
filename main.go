package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/bwmarrin/discordgo"
	"github.com/tkanos/gonfig"
)

type stock struct {
	Quote struct {
		Symbol        string  `json:"symbol"`
		CompanyName   string  `json:"companyName"`
		Current       float32 `json:"latestPrice"`
		High          float32 `json:"high"`
		Low           float32 `json:"low"`
		Open          float32 `json:"open"`
		Close         float32 `json:"close"`
		PercentChange float32 `json:"changePercent"`
		Volume        int32   `json:"latestVolume"`
	} `json:"quote"`
}

type coin struct {
	Symbol        string
	Current       string
	Open          string
	High          string
	Low           string
	PercentChange string
	Response      string
}

type botConfig struct {
	StockAPIURL           string
	CoinAPIURL            string
	InvalidCommandMessage string
	InvalidSymbolMessage  string
}

// Variables to initialize
var (
	token  string
	config botConfig
)

func init() {
	// Run the program with `go run main.go -t <token>`
	// flag.Parse() will assign to token var
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()

	// Use gonfig to fetch the config variables from config.json
	err := gonfig.GetConf("config.json", &config)
	if err != nil {
		fmt.Println("error fetching config values", err)
		return
	}
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "!ping" reply with "pong!"
	if match, _ := regexp.MatchString("!ping", m.Content); match {
		s.ChannelMessageSend(m.ChannelID, "pong!")
	}

	if match, _ := regexp.MatchString("![a-zA-Z]+ [a-zA-Z]+", m.Content); match {
		slice := strings.Split(m.Content, " ")

		if action, _ := regexp.MatchString("(?i)^!stock$", slice[0]); action {
			ticker := slice[1]
			tickerURL := config.StockAPIURL + ticker + "/batch?types=quote"

			fmt.Println("The stock url is: ", tickerURL)

			resp, err := http.Get(tickerURL)

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}

			if resp.StatusCode != 200 {
				s.ChannelMessageSend(m.ChannelID, config.InvalidSymbolMessage)
				return
			}

			defer resp.Body.Close()
			stock := stock{}

			if err = json.NewDecoder(resp.Body).Decode(&stock); err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}

			output := map[string]string{
				"Symbol":           stock.Quote.Symbol,
				"Company Name":     stock.Quote.CompanyName,
				"Current":          fmt.Sprintf("%#v", stock.Quote.Current),
				"High":             fmt.Sprintf("%#v", stock.Quote.High),
				"Low":              fmt.Sprintf("%#v", stock.Quote.Low),
				"Open":             fmt.Sprintf("%#v", stock.Quote.Open),
				"Close":            fmt.Sprintf("%#v", stock.Quote.Close),
				"Change % (1 day)": fmt.Sprintf("%#v", stock.Quote.PercentChange*100) + " %",
				"Volume":           fmt.Sprintf("%#v", stock.Quote.Volume),
			}

			outputJSON := formatStockJSON(output)

			s.ChannelMessageSend(m.ChannelID, outputJSON)
		} else if action, _ := regexp.MatchString("(?i)^!coin$", slice[0]); action {
			ticker := strings.ToUpper(slice[1])
			coinURL := config.CoinAPIURL + ticker + "&tsyms=USD"

			fmt.Println("The coin url is: ", coinURL)

			resp, err := http.Get(coinURL)

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}

			coin := coin{Symbol: ticker}

			if err = json.NewDecoder(resp.Body).Decode(&coin); err != nil || coin.Response == "Error" {
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}

			output := map[string]string{
				"Symbol":              ticker,
				"Current Price (USD)": coin.Current,
				"Open (24 Hours)":     coin.Open,
				"High (24 Hours)":     coin.High,
				"Low (24 Hours)":      coin.Low,
				"Change % (24 Hours)": coin.PercentChange,
			}

			outputJSON := formatCoinJSON(output)

			s.ChannelMessageSend(m.ChannelID, outputJSON)
			defer resp.Body.Close()
		} else {
			s.ChannelMessageSend(m.ChannelID, config.InvalidCommandMessage)
		}
	}
}

func formatStockJSON(output map[string]string) string {
	stringOrder := []string{
		"Symbol",
		"Company Name",
		"Current",
		"High",
		"Low",
		"Open",
		"Close",
		"Change % (1 day)",
		"Volume",
	}

	printer := message.NewPrinter(language.English)
	fmtStr := "```\n"

	for _, k := range stringOrder {
		if k == "Volume" {
			i, _ := strconv.Atoi(output[k])
			fmtStr += printer.Sprintf("%-17s %d\n", k, i)
		} else {
			fmtStr += printer.Sprintf("%-17s %-20s\n", k, output[k])
		}
	}

	fmtStr += "```\n"

	return fmtStr
}

func formatCoinJSON(output map[string]string) string {
	stringOrder := []string{
		"Symbol",
		"Current Price (USD)",
		"Open (24 Hours)",
		"High (24 Hours)",
		"Low (24 Hours)",
		"Change % (24 Hours)",
	}

	fmtStr := "```\n"

	for _, k := range stringOrder {
		if k == "Change % (24 Hours)" {
			f, _ := strconv.ParseFloat(output[k], 64)
			fmtStr += fmt.Sprintf("%-20s %.2f %%\n", k, f)
		} else {
			fmtStr += fmt.Sprintf("%-20s %-20s\n", k, output[k])
		}
	}

	fmtStr += "```\n"

	return fmtStr
}

func (c *coin) UnmarshalJSON(data []byte) error {
	// auxiliary struct to help map json
	var aux struct {
		Display  map[string]interface{}
		Response string
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&aux); err != nil {
		return fmt.Errorf("decode coin: %v", err)
	}

	if aux.Response == "Error" {
		return fmt.Errorf("could not find coin: %v", c.Symbol)
	}

	c.Current = aux.Display[c.Symbol].(map[string]interface{})["USD"].(map[string]interface{})["PRICE"].(string)
	c.Open = aux.Display[c.Symbol].(map[string]interface{})["USD"].(map[string]interface{})["OPEN24HOUR"].(string)
	c.High = aux.Display[c.Symbol].(map[string]interface{})["USD"].(map[string]interface{})["HIGH24HOUR"].(string)
	c.Low = aux.Display[c.Symbol].(map[string]interface{})["USD"].(map[string]interface{})["LOW24HOUR"].(string)
	c.PercentChange = aux.Display[c.Symbol].(map[string]interface{})["USD"].(map[string]interface{})["CHANGEPCT24HOUR"].(string)

	return nil
}
