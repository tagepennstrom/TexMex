package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const FRONTEND_HOST = "localhost:8000"
const SERVER_ADDRESS = "localhost:8080"

type EditDocMessage struct {
	Document string `json:"document"`
}

var document = ""

func getDocument(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(document))
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to write document: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
	}
}

func editDocWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	opts := websocket.AcceptOptions{
		OriginPatterns: []string{FRONTEND_HOST},
	}
	c, err := websocket.Accept(w, r, &opts)
	if err != nil {
		log.Printf("Failed to create websocket connection: %s", err)
		return
	}
	defer c.CloseNow()

	ctx := context.Background()
	var editDocMessage EditDocMessage
	for {
		err = wsjson.Read(ctx, c, &editDocMessage)
		if err != nil {
			log.Printf("Failed to read websocket message: %s", err)
			return
		}

		log.Printf("Received: %v\n", editDocMessage)
		document = editDocMessage.Document
	}
}

func middleware(handlerFunc http.HandlerFunc) http.HandlerFunc {
	frontendUrl := fmt.Sprintf("http://%s", FRONTEND_HOST)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", frontendUrl)
		log.Printf("%s %s\n", r.Method, r.RequestURI)
		handlerFunc(w, r)
	}
}

func main() {
	log.Printf("Server running on %s\n", SERVER_ADDRESS)

	http.HandleFunc("/document", middleware(getDocument))
	http.HandleFunc("/editDocWebsocket", editDocWebsocketHandler)

	err := http.ListenAndServe(SERVER_ADDRESS, nil)
	log.Println(err)
}
