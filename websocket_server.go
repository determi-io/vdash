package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"net/http"

	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c, err := wsh.upgrader.Upgrade(w, r, nil)

	if err != nil {

		log.Printf("error %s when upgrading connection to websocket", err)

		return

	}

	log.Printf("Upgraded connection to %s", c.RemoteAddr().String())

	defer c.Close()

	for {

		mt, message, err := c.ReadMessage()

		if err != nil {
			log.Printf("Error %s when reading message from client", err)
			return
		}

		if mt == websocket.BinaryMessage {
			err = c.WriteMessage(websocket.TextMessage, []byte("server doesn't support binary messages"))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)

			}
			return
		}

		log.Printf("Receive message %s", string(message))

		if strings.Trim(string(message), "\n") != "start" {
			err = c.WriteMessage(websocket.TextMessage, []byte("You did not say the magic word!"))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
				return
			}
			continue
		}

		log.Println("start responding to client...")

		i := 1
		for {
			response := fmt.Sprintf("Notification %d", i)
			err = c.WriteMessage(websocket.TextMessage, []byte(response))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
				return
			}

			i = i + 1
			time.Sleep(2 * time.Second)
		}
	}
}

func runWebsocketServer() {

	webSocketHandler := webSocketHandler{

		upgrader: websocket.Upgrader{},
	}
	webSocketHandler.upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	http.Handle("/", webSocketHandler)

	log.Print("Starting server...")

	log.Fatal(http.ListenAndServe("localhost:8080", nil))

}
