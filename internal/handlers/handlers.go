package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[WebSocketConnection]string)
var wsChan = make(chan WsPayload)

type WebSocketConnection struct {
	*websocket.Conn
}

// Ws response defines the response sent back from websocket
type WsJsonResponse struct {
	Action            string   `json:"action"`
	Message           string   `json:"message"`
	MessageType       string   `json:"message_type"`
	ConnectedUserList []string `json:"users_list"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

// WsEndpoint upgrades connection to websocket
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WsEndpoint: Upgrade error: ", err)
	}

	log.Println("Client connected to endpoint", ws.RemoteAddr().String())

	var response WsJsonResponse

	response.Action = "welcome"
	response.Message = "Connected to server"
	err = ws.WriteJSON(response)
	if err != nil {
		log.Println("WsEndpoint: Error trying to WriteJson", err)
	}

	var conn = WebSocketConnection{Conn: ws}
	clients[conn] = ""

	go ListenWs(&conn)
}

// Input
func ListenWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovering from ListenWs", r)
		}
	}()

	for {
		var payload WsPayload
		err := conn.Conn.ReadJSON(&payload)
		if err != nil {
			log.Println("Not sent to channel, error trying to read payload json:", err)
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

// Output
func ListenPayloadChannel() {
	var response WsJsonResponse
	for {
		e := <-wsChan

		log.Println("Listen payload: ", e)
		switch e.Action {
		case "username":
			clients[e.Conn] = e.Username

			response.Action = "userList"
			response.ConnectedUserList = GetUserNameList()
			broadcast(response)
		case "left":
			delete(clients, e.Conn)

			response.Action = "userList"
			response.ConnectedUserList = GetUserNameList()
			broadcast(response)
		}
	}
}

func broadcast(response WsJsonResponse) {
	for client := range clients {
		// fix: it is sending messages to myself lol
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Error trying to broadcast to client: ", err)
			closeErr := client.Close()
			if closeErr != nil {
				log.Println("Error trying to Close client", closeErr)
			}
			delete(clients, client)
		}
	}
}

func GetUserNameList() []string {
	var userList []string
	for _, client := range clients {
		userList = append(userList, client)
	}
	return userList
}
