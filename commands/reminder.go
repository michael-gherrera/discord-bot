package commands

import (
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
		storeurl = "localhost:6380"
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

//Add adds the new reminder as an entry into the redis table, we assume validation is
//done when the message was recieved and the date should be in the form **/**/**
//Append the reminder to the existing list if it exists, if not create a new list and add it
func (r *Reminder) Add(message string, date string) error {

	// val, err := r.Client.Get(date).Result()
	// if err != nil {
	// 	return err
	// }

	//Might err if its nil, if not we need a nil check here
	// if val != nil {
	// 	val = val.append(message)
	// } else {

	// }

	untilDate, err := time.Parse("01/02/06", date)
	if err != nil {
		return err
	}

	//Might need to add a buffer for the duration that the redis entry exists so we have time to get the value and output reminder
	timeUntil := time.Until(untilDate) + 1
	err = r.Client.Set(date, message, timeUntil).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get will run on a daily cron job to fetch the raw message from the redish cache and send to the parser to be formatted before
// being broadcasted.  The expected key is the current date
func (r *Reminder) Get(date string) (string, error) {
	output, err := r.Client.Get(date).Result()
	if err != nil {
		return "", err
	}
	return output, nil
}
