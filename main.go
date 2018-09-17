package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/BryanSLam/discord-bot/commands"
	"github.com/robfig/cron"

	"github.com/bwmarrin/discordgo"
)

// Variables to initialize
var (
	token          string
	reminderClient commands.Reminder
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

	// Initalize new reminder goroutine
	reminderClient = commands.NewReminder("")
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register handlers for discordgo
	dg.AddHandler(commands.Commander())

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// 5 AM everyday Monday - Friday
	go func() {
		c := cron.New()
		c.AddFunc("0 5 * * 1-5", func() {
			fmt.Println("test")
			err := reminderRoutine(dg)
			if err != nil {
				fmt.Println(err)
			}
		})
		c.Start()

	}()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// Function run during the daily reminder check
func reminderRoutine(s *discordgo.Session) error {
	output, err := reminderClient.Get(time.Now().Format("01/02/06"))
	if err != nil {
		return err
	}
	for _, reminder := range output {
		cacheEntry := strings.Split(reminder, "~*")
		channel := cacheEntry[0]
		s.ChannelMessageSend(channel, cacheEntry[1])
	}
	return nil
}
