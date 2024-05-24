package subscription

import ("github.com/graphql-go/graphql")

func MessageSubscribe(params graphql.ResolveParams) (interface{}, error) {
	messageChan := make(chan SubscribeMessage)
    subscribers[messageChan] = true

    // Close the channel and remove the subscriber when done
    defer func() {
        delete(subscribers, messageChan)
        close(messageChan)
    }()

    for {
        select {
        case msg := <-messageChan:
           return msg, nil
        case <- messageChan:
            return nil, nil
        }
    }
	
}
