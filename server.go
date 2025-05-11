package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"slices"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"websocket-server/crdt"
)

const frontendPort = "5173"

type Change struct {
	FromA  int    `json:"fromA"`  // Start index
	ToA    int    `json:"toA"`    // Slut index
	FromB  int    `json:"fromB"`  // Start index
	ToB    int    `json:"toB"`    // Slut index
	Text   string `json:"text"`   // Tillagd text
	UserID int    `json:"userId"` // AnvändarID för CRDT
}

type EditDocMessage struct {
	document     string
	cursorIndex  int
	jsonCChanges string
}

type Envelope struct {
	Type       string         `json:"type"`                 // "operation", "stateRequest", "stateResponse"
	EditDocMsg EditDocMessage `json:"editDocMsg,omitempty"` // changes
	ByteState  []byte         `json:"byteState,omitempty"`  // for stateResponse
	UserID     int            `json:"userId,omitempty"`     // optional sender ID
}

type Client struct {
	wscon *websocket.Conn
	id    int
}

// *
// Globala variabler
// *

var connections []Client
var currId int = 0
var globalDocument crdt.Document

// *
// Funktioner
// *

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "12.34.56.78:90")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func broadcastMessage(ctx context.Context, message EditDocMessage, sender Client) {

	log.Printf("Broadcasting to %d clients\n", len(connections))

	for _, c := range connections {
		if c.id == sender.id {
			continue
		}
		resp := Envelope{
			Type:       "operation",
			EditDocMsg: message,
		}

		err := wsjson.Write(ctx, c.wscon, resp)

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
		ID   int    `json:"id"`
		Type string `json:"type"`
	}{ID: user.id, Type: "user_connected"}

	// Send the ID to the client (use wsjson.Write to send a JSON message)
	ctx := context.Background()
	err := wsjson.Write(ctx, user.wscon, initialMessage)
	if err != nil {
		log.Printf("Failed to send connection ID to client: %s", err)
		user.wscon.CloseNow()
		return
	}
	defer user.wscon.CloseNow()

	var env Envelope

	for {
		err := wsjson.Read(ctx, user.wscon, &env)

		if err != nil {
			log.Printf("Failed to read websocket message: %s", err)
			connections = removeConn(user)
			return
		}

		switch env.Type {

		case "operation":

			println("operation case (server.go)")

			// uppdatera globala versionen
			println(
				"a)", env.EditDocMsg.document,
				"b)", env.EditDocMsg.cursorIndex,
				"c)", env.EditDocMsg.jsonCChanges,
			)
			fmt.Println("d)", env.EditDocMsg)
			fmt.Println("e)", env)

			jsonCChange := env.EditDocMsg.jsonCChanges

			globalDocument.HandleCChange(jsonCChange)

			println("-------Global Update--------")
			globalDocument.PrintDocument(false)

			broadcastMessage(ctx, env.EditDocMsg, user)
			break

		case "stateRequest":
			println("Global CRDT state is:", globalDocument.ToString())

			data, err := globalDocument.Snapshot()

			if err != nil {
				log.Printf("snapshot error: %v", err)
				break
			}

			resp := Envelope{
				Type:      "stateResponse",
				ByteState: data,
			}

			wsjson.Write(ctx, user.wscon, resp)
			println("Request sent")
			break

		default:
			println("Error. No case for switch statment. (server.go)")
			break
		}

	}
}

func middleware(handlerFunc http.HandlerFunc) http.HandlerFunc {
	ip, _ := getLocalIP()
	frontendUrl := fmt.Sprintf("http://%s:%s", ip, frontendPort)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", frontendUrl)
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		log.Printf("%s %s\n", r.Method, r.RequestURI)
		if r.Method == "OPTIONS" {
			return
		}
		handlerFunc(w, r)
	}
}

func main() {

	// todo: det nedan är tillfälligt för att testa crdt synkning
	// ***
	filler := "ABC"
	globalDocument = crdt.DocumentFromStr(filler)

	// ***

	const port = "8080"
	ip, _ := getLocalIP()
	serverAddress := fmt.Sprintf("%s:%s", ip, port)
	log.Printf("Server running on http://%s/\n", serverAddress)

	mux := http.NewServeMux()

	mux.HandleFunc("/projects", middleware(getProjects))
	mux.HandleFunc("/projects/{name}", middleware(projectHandler))
	mux.HandleFunc("/projects/{projectName}/documents/{documentName}", middleware(projectDocumentHandler))
	mux.HandleFunc("/projects/{name}/pdf", middleware(getProjectPdf))

	mux.HandleFunc("/editDocWebsocket", editDocWebsocketHandler)

	err := http.ListenAndServe(serverAddress, mux)
	log.Println(err)
}
