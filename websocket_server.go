package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"net/http"

	"github.com/creack/pty"
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

	////////////////////////////////////////////
	// Terminal

	// Create arbitrary command.
	cmd := exec.Command("bash")

	// Start the command with a pty.
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Printf("Error %s when starting pty", err)
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	fmt.Println("started bash")

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
		log.Println(" => starting bash")

		// streaming ptmx -> websocket
		go func() {
			s := bufio.NewScanner(ptmx)
			s.Split(bufio.ScanRunes)
			for s.Scan() {
				log.Printf("Got rune: '%s'", s.Text())
				c.WriteMessage(websocket.TextMessage, []byte(s.Text()))
			}
		}()

		go func() {
			for {
				ptmx.WriteString("a")
				time.Sleep(2 * time.Second)
			}
		}()

		// streaming websocket -> ptmx
		for {
			mt, msg, err := c.ReadMessage()
			if mt != websocket.TextMessage {
				log.Println("received binary message from websocket")
				continue
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			log.Printf("Got message from websocket: '%s'", msg)
			// fmt.Fprint(ptmx, msg)
			ptmx.Write(msg)
		}

		// _, _ = io.Copy(os.Stdout, ptmx)

		// if strings.Trim(string(message), "\n") != "start" {
		// 	err = c.WriteMessage(websocket.TextMessage, []byte("You did not say the magic word!"))
		// 	if err != nil {
		// 		log.Printf("Error %s when sending message to client", err)
		// 		return
		// 	}
		// 	continue
		// }

		// log.Println("start responding to client...")

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
