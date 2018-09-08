package commands

import (
	"log"

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
	if url != "" {
		storeurl = url
	} else {
		//Replace with the config from the charts PR
		storeurl = "localhost:9200"
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
	var val string

	val, err := r.Client.Get(date).Result()
	if err != nil {
		log.Println("error when getting %s %v", date, err)
		return err
	}

	//Might err if its nil, if not we need a nil check here
	// if val != nil {
	// 	val = val.append(message)
	// } else {

	// }

	r.Client.Set(date, message)
	return nil
}

func (r *Reminder) Get()
