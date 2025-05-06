package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"slices"
	"websocket-server/crdt"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const frontendPort = "5173"

type EditInstruction struct {
	operation string
	insertion string
	uID       int
}

type Client struct {
	wscon *websocket.Conn
	id    int
}

var document = crdt.DocumentFromStr(`\documentclass{article}
\begin{document}
abcd
\end{document}`)
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
	_, err := w.Write([]byte(document.ToString()))
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to write document: %s", err)
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
	}
}

func broadcastMessage(ctx context.Context, message EditInstruction, sender Client) {

	log.Printf("Broadcasting to %d clients\n", len(connections))

	for _, c := range connections {
		if c.id == sender.id {
			continue
		}
		err := wsjson.Write(ctx, c.wscon, message)

		if err != nil {
			log.Printf("Failed to write websocket message: %s", err)
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

	ip, _ := getLocalIP()
	frontendHost := fmt.Sprintf("%s:%s", ip, "5173")
	opts := websocket.AcceptOptions{
		OriginPatterns: []string{frontendHost},
	}
	c, err := websocket.Accept(w, r, &opts)
	if err != nil {
		log.Printf("Failed to create websocket connection: %s", err)
	}
	currId++
	user := Client{wscon: c, id: currId}
	return user

}

func editDocWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	user := acceptConnection(w, r)
	log.Printf("User connected with ID: %d\n", user.id)
	connections = append(connections, user)

	initialMessage := struct {
		ID int `json:"id"`
	}{ID: user.id}

	// Send the ID to the client (use wsjson.Write to send a JSON message)
	ctx := context.Background()
	err := wsjson.Write(ctx, user.wscon, initialMessage)
	if err != nil {
		log.Printf("Failed to send connection ID to client: %s", err)
		user.wscon.CloseNow()
		return
	}
	defer user.wscon.CloseNow()

	var newChange EditInstruction

	for {
		err := wsjson.Read(ctx, user.wscon, &newChange)

		if err != nil {
			log.Printf("Failed to read websocket message: %s", err)
			connections = removeConn(user)
			return
		}

		log.Printf("Operation made: %s\n", newChange.operation)
		broadcastMessage(ctx, newChange, user)

	}
}

/* func compileDocument(w http.ResponseWriter, r *http.Request) {
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

} */

func middleware(handlerFunc http.HandlerFunc) http.HandlerFunc {
	ip, _ := getLocalIP()
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
	log.Printf("Server running on http://%s/\n", serverAddress)

	http.HandleFunc("/document", middleware(getDocument))
	//http.HandleFunc("/compileDocument", middleware(compileDocument))
	//http.HandleFunc("/pdf", middleware(servePdf))

	http.HandleFunc("/editDocWebsocket", editDocWebsocketHandler)

	err := http.ListenAndServe(serverAddress, nil)
	log.Println(err)
}
