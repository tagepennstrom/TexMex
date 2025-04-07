package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const frontendHost = "localhost:8000"
const serverAddress = "localhost:8080"

type EditDocMessage struct {
	Document string `json:"document"`
}

var document = ""
var channels []chan EditDocMessage

func getDocument(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(document))
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to write document: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
	}
}

func listenChannel(ctx context.Context, c *websocket.Conn, channel chan EditDocMessage) {
	for {
		new_doc := <-channel
		log.Printf("Channel received: %v\n", new_doc)
		err := wsjson.Write(ctx, c, new_doc)
		if err != nil {
			log.Printf("Failed to write websocket message: %s", err)
		}
	}
}

func listenWebsocket(ctx context.Context, c *websocket.Conn) {
	var editDocMessage EditDocMessage
	for {
		err := wsjson.Read(ctx, c, &editDocMessage)
		if err != nil {
			log.Printf("Failed to read websocket message: %s", err)
			return
		}

		log.Printf("Received: %v\n", editDocMessage)
		document = editDocMessage.Document
		for _, channel := range channels {
			channel <- editDocMessage
		}
	}
}

func editDocWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	opts := websocket.AcceptOptions{
		OriginPatterns: []string{frontendHost},
	}
	c, err := websocket.Accept(w, r, &opts)
	if err != nil {
		log.Printf("Failed to create websocket connection: %s", err)
		return
	}
	channel := make(chan EditDocMessage)
	channels = append(channels, channel)
	ctx := context.Background()
	go listenChannel(ctx, c, channel)
	go listenWebsocket(ctx, c)
}

func compileDocument(w http.ResponseWriter, r *http.Request) {
	document, err := io.ReadAll(r.Body)
	if err != nil {
		errorMessage := fmt.Sprintf("Error reading request body: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}
	r.Body.Close()

	const filename = "document"
	filenameLatex := fmt.Sprintf("%s.tex", filename)
	filenamePdf := fmt.Sprintf("%s.pdf", filename)

	const writeReadPermission = os.FileMode(0666)
	err = os.WriteFile(filenameLatex, document, writeReadPermission)
	if err != nil {
		errorMessage := fmt.Sprintf("Error creating LaTeX file: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", filenameLatex)
	err = cmd.Run()
	if err != nil {
		errorMessage := fmt.Sprintf("Error compiling LaTeX file: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, filenamePdf)
}

func middleware(handlerFunc http.HandlerFunc) http.HandlerFunc {
	frontendUrl := fmt.Sprintf("http://%s", frontendHost)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", frontendUrl)
		log.Printf("%s %s\n", r.Method, r.RequestURI)
		handlerFunc(w, r)
	}
}

func main() {
	log.Printf("Server running on %s\n", serverAddress)

	http.HandleFunc("/document", middleware(getDocument))
	http.HandleFunc("/editDocWebsocket", editDocWebsocketHandler)
	http.HandleFunc("/compileDocument", middleware(compileDocument))

	err := http.ListenAndServe(serverAddress, nil)
	log.Println(err)
}
