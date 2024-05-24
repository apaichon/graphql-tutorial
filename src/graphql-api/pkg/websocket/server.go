package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

// Upgrade the HTTP server connection to the WebSocket protocol
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Check the origin of the request here if needed
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	// Infinite loop to continuously listen for messages from the WebSocket
	for {
		// Read message from browser
		messageType, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the message to the console
		fmt.Printf("Received: %s\n", msg)

		// Write message back to browser
		if err := ws.WriteMessage(messageType, msg); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("WebSocket server starting at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
