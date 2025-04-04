package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

const FRONTEND_HOST = "localhost:8000"
const SERVER_ADDRESS = "localhost:8080"

var document = ""

func getDocument(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(document))
	if err != nil {
		http.Error(w, "Failed to write document", http.StatusInternalServerError)
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	opts := websocket.AcceptOptions{
		OriginPatterns: []string{FRONTEND_HOST},
	}
	c, err := websocket.Accept(w, r, &opts)
	if err != nil {
		http.Error(w, "Failed to create websocket connection",
			http.StatusInternalServerError)
		return
	}
	defer c.CloseNow()

	ctx := context.Background()
	for {
		_, message, err := c.Read(ctx)
		if err != nil {
			http.Error(w, "Failed to read websocket message",
				http.StatusInternalServerError)
			return
		}

		message_str := string(message)
		document = message_str
		log.Printf("Received: '%s'\n", message_str)
	}
}

func middleware(handlerFunc http.HandlerFunc) http.HandlerFunc {
	frontend_url := fmt.Sprintf("http://%s", FRONTEND_HOST)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", frontend_url)
		log.Printf("%s %s\n", r.Method, r.RequestURI)
		handlerFunc(w, r)
	}
}

func main() {
	log.Printf("Server running on %s\n", SERVER_ADDRESS)

	http.HandleFunc("/document", middleware(getDocument))
	http.HandleFunc("/websocket", websocketHandler)

	err := http.ListenAndServe(SERVER_ADDRESS, nil)
	log.Println(err)
}
