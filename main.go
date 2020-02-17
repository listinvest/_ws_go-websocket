package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func homePage(w http.ResponseWriter, r *http.Request) {
	var (
		homeTemplate *template.Template
		err          error
	)

	homeTemplate, err = template.ParseFiles("index.gohtml")
	if err != nil {
		log.Panicln(err)
	}

	w.Header().Set("Content-Type", "text/html")
	if err := homeTemplate.Execute(w, nil); err != nil {
		log.Println(err)
	}

}

// This is a listener on server side to listen to in coming client
func reader(conn *websocket.Conn) {
	// Loop to keep connection open and listen indefinitely
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Panicln(err)
			return
		}

		// print out message received from client
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func wsEndPoint(w http.ResponseWriter, r *http.Request) {
	// upgrade current HTTP connection to websocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	// write back to Client
	err = ws.WriteMessage(1, []byte("hi Client!"))
	if err != nil {
		log.Println(err)
	}

	// keep the connection open
	reader(ws)
}

func setupRoute() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndPoint)
}

func main() {
	fmt.Println("hello, starting server on :8000")
	setupRoute()
	http.ListenAndServe(":8000", nil)
}
