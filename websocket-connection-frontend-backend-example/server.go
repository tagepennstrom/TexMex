package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const FRONTEND_URL = "http://localhost:8000"

var document = ""

func getDocument(w http.ResponseWriter, r *http.Request) {
	fmt.Println(document)
	w.Write([]byte(document))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == FRONTEND_URL
	},
}
func websocketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to create websocket connection", http.StatusBadRequest)
	}

	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return
		}

		fmt.Printf("Received: %s\n", string(message))
		document = string(message)
	}
}

func main() {
	fmt.Println("Server running")

	r := mux.NewRouter()
	r.HandleFunc("/document", getDocument)
	r.HandleFunc("/websocket", websocketHandler)

	allowed_origins := handlers.AllowedOrigins([]string{FRONTEND_URL})
	err := http.ListenAndServe("localhost:8080", handlers.CORS(allowed_origins)(r))
	fmt.Println(err)
}
