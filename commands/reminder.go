package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

/*
 * Set up redis table in a way where the keys are the dates of when to pull a reminder
 * and the value is an array of events that need to be reminded on that date
 */
type Reminder struct {
	Client *redis.Client
}

func NewReminder(url string) Reminder {
	var (
		storeurl string
		client   *redis.Client
	)
	//For mock testing
	if url != "" {
		storeurl = url
	} else {
		//Replace with the config from the charts PR
		storeurl = "redis:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr:     storeurl,
		Password: "", //no password set
		DB:       0,  //Use default DB
	})
	return Reminder{
		Client: client,
	}
}

// Add adds the new reminder as an entry into the redis table, we assume validation is
// done when the message was recieved and the date should be in the form **/**/**
// Append the reminder to the existing list if it exists, if not create a new list and add it
func (r *Reminder) Add(message string, date string) error {
	var (
		sb  strings.Builder
		val string
	)
	val, err := r.Client.Get(date).Result()
	if err != nil {
		if err.Error() != "redis: nil" {
			return err
		}
	}
	if val != "" {
		sb.WriteString(val)
		sb.WriteString("::")
	}
	sb.WriteString(message)

	untilDate, err := time.Parse("01/02/06", date)
	if err != nil {
		return err
	}

	//Might need to add a buffer for the duration that the redis entry exists so we have time to get the value and output reminder
	timeUntil := time.Until(untilDate) + 1
	err = r.Client.Set(date, sb.String(), timeUntil).Err()
	if err != nil {
		fmt.Println("Add set error")
		return err
	}
	return nil
}

// Get will run on a daily cron job to fetch the raw messages from the redis cache and send to the parser to be formatted before
// being broadcasted.  The expected key is the current date
func (r *Reminder) Get(date string) ([]string, error) {
	messages, err := r.Client.Get(date).Result()
	if err != nil {
		if err.Error() != "redis: nil" {
			return nil, err
		}
		return nil, nil
	}
	output := strings.Split(messages, "::")

	//Delete the entries in the redis cache after it's gotten in case the duration doesn't expire
	r.Client.Del(date)
	return output, nil
}

// TODO:LEFTOVER CODE FROM REMINDER
// update it so it can be used like other commands
// } else if action, _ := regexp.MatchString("(?i)^!remindme$", slice[0]); action {
//   messageArr := strings.Split(m.Content, "\"")
//   if len(messageArr) != 3 {
//     s.ChannelMessageSend(m.ChannelID, "Reminder messages must be surrounded by quotes \"{message}\" ")
//     return
//   }
//   // We store the person who sent the message as well as the channel id into the redis cache so we know where and who to contact later
//   message := m.ChannelID + "~*" + m.Author.Mention() + ": " + messageArr[1]
//   date := slice[len(slice)-1]
//   match, _ := regexp.MatchString("(0?[1-9]|1[012])/(0?[1-9]|[12][0-9]|3[01])/(\\d\\d)", date)
//   if match == false {
//     s.ChannelMessageSend(m.ChannelID, "Invalid date given loser")
//     return
//   }
//
//   // Commenting out the date check for now, weird behavior where you get blocked for
//   // Setting a reminder for the next day
//
//   // dateCheck, _ := time.Parse("01/02/06", date)
//   // if time.Until(dateCheck) < 0 {
//   // 	s.ChannelMessageSend(m.ChannelID, "Date has already passed ya fuck")
//   // 	return
//   // }
//
//   err := reminderClient.Add(message, date)
//   if err != nil {
//     s.ChannelMessageSend(m.ChannelID, err.Error())
//     return
//   }
//   s.ChannelMessageSend(m.ChannelID, "Reminder Set!")
