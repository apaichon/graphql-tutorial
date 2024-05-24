package subscription

import (
	"fmt"
	"net/http"
	"time"
)

func NewMessageHandler(w http.ResponseWriter, r *http.Request) {
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
        return
    }

    messageChan := make(chan SubscribeMessage)
    subscribers[messageChan] = true

    // Close the channel and remove the subscriber when done
    defer func() {
        delete(subscribers, messageChan)
        close(messageChan)
    }()

    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

    for {
        select {
        case msg := <-messageChan:
            fmt.Fprintf(w, "data: {\"id\":\"%s\", \"content\":\"%s\" , \"timestamp\":\"%v\" }\n\n", msg.Id, msg.Content, time.Now())
            flusher.Flush()
        case <-r.Context().Done():
            return
        }
    }
}

func SendMessage(msg SubscribeMessage) {
    for subscriber := range subscribers {
        subscriber <- msg
    }
}

