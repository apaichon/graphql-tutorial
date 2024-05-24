package subscription

import (
	"time"
)

type Subscription struct {
	Channel chan interface{}
	Done    chan bool
}

// Message is a struct for message
type SubscribeMessage struct {
	Id        string    `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

var subscribers = make(map[chan SubscribeMessage]bool)

func SetSubcribe(messageChan chan SubscribeMessage) {
	subscribers[messageChan] = true
}

func DeleteSubcribe(messageChan chan SubscribeMessage) {
	delete(subscribers, messageChan)
}

