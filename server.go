package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"slices"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const frontendHost = "localhost:8000"
const serverAddress = "localhost:8080"

type EditDocMessage struct {
	Document string `json:"document"`
}

var document = ""
var connections []*websocket.Conn 

func getDocument(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(document))
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to write document: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
	}
}

func broadcastMessage(ctx context.Context, message EditDocMessage) {
	log.Printf("Broadcasting to %d clients\n", len(connections))
	log.Printf("Broadcasting message: %v\n", message)
	for _, c := range connections {
		err := wsjson.Write(ctx, c, message)
		if err != nil {
			log.Printf("Failed to write websocket message: %s", err)
		}
	}
}

func removeConn(connToDelete *websocket.Conn) []*websocket.Conn {
	for i, c := range connections {
		if c == connToDelete {
			return slices.Delete(connections, i, i + 1)
		}
	}
	return connections;
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
	defer c.CloseNow()
	connections = append(connections, c)

	ctx := context.Background()
	var editDocMessage EditDocMessage
	for {
		err := wsjson.Read(ctx, c, &editDocMessage)
		if err != nil {
			log.Printf("Failed to read websocket message: %s", err)
			connections = removeConn(c)
			return
		}

		log.Printf("Received: %v\n", editDocMessage)
		document = editDocMessage.Document
		broadcastMessage(ctx, editDocMessage)
	}
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
