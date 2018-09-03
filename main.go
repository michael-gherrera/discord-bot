package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/BryanSLam/discord-bot/datasource"
	"github.com/BryanSLam/discord-bot/util"
	iex "github.com/jonwho/go-iex"

	"github.com/bwmarrin/discordgo"
	"github.com/tkanos/gonfig"
)

type botConfig struct {
	CoinAPIURL            string
	InvalidCommandMessage string
}

// Variables to initialize
var (
	token     string
	config    botConfig
	iexClient *iex.Client
)

func init() {
	// Run the program with `go run main.go -t <token>`
	// flag.Parse() will assign to token var
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()

	// If no value was provided from flag look for env var BOT_TOKEN
	if token == "" {
		token = os.Getenv("BOT_TOKEN")
	}

	// Initialize iexClient with new client
	iexClient = iex.NewClient()

	chart, errr := iexClient.Chart("ADBE", "7d")
	if errr != nil {
		log.Print(errr)
		return
	}

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
			quote, err := iexClient.Quote(ticker, true)

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}

			outputJSON := util.FormatQuote(quote)

			s.ChannelMessageSend(m.ChannelID, outputJSON)
		} else if action, _ := regexp.MatchString("(?i)^!er$", slice[0]); action {
			ticker := slice[1]
			earnings, err := iexClient.Earnings(ticker)

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}

			outputJSON := util.FormatEarnings(earnings)

			s.ChannelMessageSend(m.ChannelID, outputJSON)
		} else if action, _ := regexp.MatchString("(?i)^!coin$", slice[0]); action {
			ticker := strings.ToUpper(slice[1])
			coinURL := config.CoinAPIURL + ticker + "&tsyms=USD"

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
		} else {
			s.ChannelMessageSend(m.ChannelID, config.InvalidCommandMessage)
		}
	}
}
