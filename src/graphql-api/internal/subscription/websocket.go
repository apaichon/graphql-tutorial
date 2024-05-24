package subscription

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Check the origin of the request here if needed
		return true
	},
}

func SubscribeWsHandler(schema graphql.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer ws.Close()

		messageChan := make(chan SubscribeMessage)
		SetSubcribe(messageChan)

		// Close the channel and remove the subscriber when done
		defer func() {
			DeleteSubcribe(messageChan)
			close(messageChan)
		}()

		for {
			select {
			case msg := <-messageChan:
				jsonBytes, err := json.Marshal(msg)
				if err != nil {
					fmt.Println(err)
				}
				// Write message back to browser
				if err := ws.WriteMessage(1, jsonBytes); err != nil {
					fmt.Println(err)
					return
				}
			case <-r.Context().Done():
				return
			}
		}
	}
}