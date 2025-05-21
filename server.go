package main

import (
	"context"
	"encoding/json"
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
	CursorIndex  int    `json:"cursorIndex"`
	ByteCChanges []byte `json:"byteCChanges"`
}

type Envelope struct {
	Type       string         `json:"type"`                 // "operation", "stateRequest", "stateResponse"
	EditDocMsg EditDocMessage `json:"editDocMsg,omitempty"` // changes
	ByteState  []byte         `json:"byteState,omitempty"`  // for stateResponse
}

type Client struct {
	wscon        *websocket.Conn
	id           int
	projectName  string
	documentName string
}

type ProjectData struct {
	projectName  string
	documentName string
	docu         *crdt.Document
}

// *
// Globala variabler
// *

var connections []Client
var currIDCounter int = 0

var globalProjects map[string]ProjectData

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

	brdcstCount := 0

	for _, c := range connections {
		if c.id == sender.id {
			continue
		}
		if c.projectName != sender.projectName {
			continue
		}
		resp := Envelope{
			Type:       "operation",
			EditDocMsg: message,
		}

		err := wsjson.Write(ctx, c.wscon, resp)

		if err != nil {
			log.Printf("Failed to write websocket message: %s", err)
		} else {
			brdcstCount++
		}
	}
	log.Printf("Broadcasted to %d clients connected to project: '%s'\n", brdcstCount, sender.projectName)

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

	c.SetReadLimit(10 << 20) // 10mib ( 10 * 2^20 )

	project := r.URL.Query().Get("projectName")
	doc := r.URL.Query().Get("documentName")

	currIDCounter++
	user := Client{
		wscon:        c,
		id:           currIDCounter,
		projectName:  project,
		documentName: doc,
	}

	_, exists := globalProjects[project]

	if !exists {
		new := crdt.NewDocument()
		new.Active = false

		newEntry := ProjectData{
			projectName:  project,
			documentName: doc,
			docu:         &new,
		}

		// mappa
		globalProjects[project] = newEntry
	}

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
		_, msg, err := user.wscon.Read(ctx)

		if err != nil {
			log.Printf("Failed to read websocket message: %s", err)
			connections = removeConn(user)
			return
		}

		if err := json.Unmarshal(msg, &env); err != nil {
			log.Printf("unmarshal error (beginning): %v", err)
			return
		}

		switch env.Type {

		case "operation":
			// uppdatera globala CRDTn
			ByteCChanges := env.EditDocMsg.ByteCChanges

			selectedDoc := globalProjects[user.projectName].docu
			selectedDoc.HandleCChange(string(ByteCChanges))

			saveProjectDocumentServerSide(user.projectName, user.documentName)

			broadcastMessage(ctx, env.EditDocMsg, user)

			break

		case "stateRequest":

			selectedDoc := globalProjects[user.projectName].docu

			s1 := *selectedDoc

			data, err := s1.Snapshot()

			if err != nil {
				log.Printf("snapshot error: %v", err)
				break
			}

			resp := Envelope{
				Type:      "stateResponse",
				ByteState: data,
			}

			wsjson.Write(ctx, user.wscon, resp)
			println("Response to request sent")
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
	const port = "8080"
	ip, _ := getLocalIP()
	serverAddress := fmt.Sprintf("%s:%s", ip, port)
	log.Printf("Server running on http://%s/\n", serverAddress)

	globalProjects = make(map[string]ProjectData) // init projects mapping

	mux := http.NewServeMux()

	mux.HandleFunc("/projects", middleware(getProjects))
	mux.HandleFunc("/projects/{name}", middleware(projectHandler))
	mux.HandleFunc("/projects/{projectName}/documents/{documentName}", middleware(projectDocumentHandler))
	mux.HandleFunc("/projects/{name}/pdf", middleware(getProjectPdf))
	mux.HandleFunc("/projects/uploadFile", middleware(uploadFileToProject))
	mux.HandleFunc("/projects/getfiles", middleware(getFilesFromProject))
	mux.HandleFunc("/projects/delFile", middleware(deleteFileFromProject))

	mux.HandleFunc("/editDocWebsocket", editDocWebsocketHandler)

	err := http.ListenAndServe(serverAddress, mux)
	log.Println(err)
}
