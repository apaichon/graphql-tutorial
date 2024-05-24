package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	"github.com/gorilla/websocket"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Connect to the server
	u := "ws://localhost:8080/ws"
	fmt.Printf("connecting to %s\n", u)

	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	// Start a goroutine to listen for messages from the server
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				return
			}
			fmt.Printf("recv: %s\n", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				fmt.Println("write:", err)
				return
			}
		case <-interrupt:
			fmt.Println("interrupt")

			// Cleanly close the connection by sending a close message
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
