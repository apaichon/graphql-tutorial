package main

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"sync"
	"time"
	"net/http"
)

// Message represents the type that will be published and subscribed to
type Message struct {
	Content string `json:"content"`
}

// PubSub represents the pub-sub system
type PubSub struct {
	subscribers map[chan Message]struct{}
	mutex       sync.RWMutex
}

// NewPubSub creates a new PubSub instance
func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[chan Message]struct{}),
	}
}

// Subscribe adds a new subscriber to the pub-sub system
func (ps *PubSub) Subscribe() chan Message {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	sub := make(chan Message)
	ps.subscribers[sub] = struct{}{}
	return sub
}

// Unsubscribe removes a subscriber from the pub-sub system
func (ps *PubSub) Unsubscribe(sub chan Message) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	delete(ps.subscribers, sub)
	close(sub)
}

// Publish sends a message to all subscribers
func (ps *PubSub) Publish(msg Message) {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()
	for sub := range ps.subscribers {
		sub <- msg
	}
}

func main() {
	ps := NewPubSub()

	// GraphQL schema definition
	fields := graphql.Fields{
		"message": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Dummy resolver, not used in subscription
				return "Hello, World!", nil
			},
		},
	}

	subscriptionFields := graphql.Fields{
		"messageReceived": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Subscribe to messages
				sub := ps.Subscribe()
				defer ps.Unsubscribe(sub)
				for {
					select {
					case message := <-sub:
						return message.Content, nil
					}
				}
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	rootSubscription := graphql.ObjectConfig{Name: "RootSubscription", Fields: subscriptionFields}

	schemaConfig := graphql.SchemaConfig{
		Query:        graphql.NewObject(rootQuery),
		Subscription: graphql.NewObject(rootSubscription),
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		panic(err)
	}

	// GraphQL handler
	handler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Start GraphQL server
	go func() {
		fmt.Println("GraphQL server running on http://localhost:8080/graphql")
		http.Handle("/graphql", handler)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic(err)
		}
	}()

	// Publish some messages
	go func() {
		for i := 1; i <= 5; i++ {
			ps.Publish(Message{Content: fmt.Sprintf("Message %d", i)})
			time.Sleep(2 * time.Second)
		}
	}()

	// Wait forever
	select {}
}
