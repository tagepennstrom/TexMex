package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"slices"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const frontendPort = "5173"
const filename = "document"

type Change struct {
	From   int    `json:"from"`   // Start index
	To     int    `json:"to"`     // Slut index
	Text   string `json:"text"`   // Tillagd text
	UserID int    `json:"userId"` // AnvändarID för CRDT
}

type EditDocMessage struct {
	Changes []Change `json:"changes"`
}

type Client struct {
	wscon          *websocket.Conn
	id             int
	messageChannel chan EditDocMessage
}

var document = `\documentclass{article}
\begin{document}
	abcd
\end{document}`
var connections []Client
var currId int = 0

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "12.34.56.78:90")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func getDocument(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(document))
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to write document: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
	}
}

func saveDocument(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorMessage := fmt.Sprintf("Error reading request body: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}
	r.Body.Close()

	log.Println("Saving document")
	document = string(body)
}

func broadcastMessage(message EditDocMessage, sender Client) {
	log.Printf("Broadcasting to %d clients\n", len(connections))
	for _, c := range connections {
		if c.id == sender.id {
			continue
		}
		select {
		case c.messageChannel <- message:
		default:
			log.Printf("Client %d's message channel is full. Most Likely Timed out. Closing connection\n", c.id)
			connections = removeConn(c)
			close(c.messageChannel)
			// Är det här rätt lösning? TODO: tänk över detta, varför can vår channel vara full?
		}
	}
}

func removeConn(connToDelete Client) []Client {
	log.Printf("Removing user with ID: %d\n", connToDelete.id)
	for i, c := range connections {
		if c == connToDelete {
			return slices.Delete(connections, i, i+1)
		}
	}
	return connections
}

func acceptConnection(w http.ResponseWriter, r *http.Request) Client {

	//ip, _ := getLocalIP()
	ip := "83.233.230.209"
	frontendHost := fmt.Sprintf("%s:%s", ip, "5173")
	opts := websocket.AcceptOptions{
		OriginPatterns: []string{frontendHost},
	}
	c, err := websocket.Accept(w, r, &opts)
	if err != nil {
		log.Printf("Failed to create websocket connection: %s", err)
	}
	currId++
	user := Client{
		wscon:          c,
		id:             currId,
		messageChannel: make(chan EditDocMessage, 10),
	} // Bra bufferstorlek??? Ny forskningsfråga!?!?

	go handleClientsMessages(user)

	return user

}

func handleClientsMessages(client Client) {
	defer client.wscon.CloseNow()

	for message := range client.messageChannel {
		log.Printf("Client %d received message", client.id)
		err := wsjson.Write(context.Background(), client.wscon, message)
		if err != nil {
			log.Printf("Failed to send message to client %d: %s", client.id, err)
			connections = removeConn(client)
			return
		}
	}
}

func updateDocument(changes []Change) {

	for _, change := range changes {

		document = document[:change.From] + change.Text + document[change.To:]

	}

}

func editDocWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	user := acceptConnection(w, r)
	log.Printf("User connected with ID: %d\n", user.id)
	connections = append(connections, user)

	initialMessage := struct {
		ID int `json:"id"`
	}{ID: user.id}

	ctx := context.Background()
	err := wsjson.Write(ctx, user.wscon, initialMessage)
	if err != nil {
		log.Printf("Failed to send connection ID to client: %s", err)
		user.wscon.CloseNow()
		return
	}

	go func() {
		var editDocMessage EditDocMessage

		for {
			err := wsjson.Read(ctx, user.wscon, &editDocMessage)

			if err != nil {
				log.Printf("Failed to read websocket message: %s", err)
				connections = removeConn(user)
				return
			}
			log.Printf("Changes made: %v\n", editDocMessage)

			updateDocument(editDocMessage.Changes)

			broadcastMessage(editDocMessage, user)

		}
	}()
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
	filenameLatex := fmt.Sprintf("%s.tex", filename)
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
	fmt.Fprintf(w, `{"pdfUrl": "/pdf"}`)
}

func servePdf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/pdf")
	filenamePdf := fmt.Sprintf("%s.pdf", filename)
	http.ServeFile(w, r, filenamePdf)

}

func middleware(handlerFunc http.HandlerFunc) http.HandlerFunc {
	// ip, _ := getLocalIP()
	ip := "83.233.230.209"
	frontendUrl := fmt.Sprintf("http://%s:%s", ip, frontendPort)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", frontendUrl)
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		log.Printf("%s %s\n", r.Method, r.RequestURI)
		if r.Method == "OPTIONS" {
			return
		}
		handlerFunc(w, r)
	}
}

func main() {
	const port = "8080"
	ip, _ := getLocalIP()
	serverAddress := fmt.Sprintf("%s:%s", ip, port)
	log.Printf("Server running on http://83.233.230.209/\n", serverAddress)

	http.HandleFunc("/document", middleware(getDocument))
	http.HandleFunc("/saveDocument", middleware(saveDocument))
	http.HandleFunc("/compileDocument", middleware(compileDocument))
	http.HandleFunc("/pdf", middleware(servePdf))

	http.HandleFunc("/editDocWebsocket", editDocWebsocketHandler)

	err := http.ListenAndServe(serverAddress, nil)
	log.Println(err)
}
