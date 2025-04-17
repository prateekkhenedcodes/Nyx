package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/prateekkhenedcodes/Nyx/internal/auth"
	"github.com/prateekkhenedcodes/Nyx/sql/queries"
)

type Message struct {
	SenderID  string `json:"sender_id"`
	Content   string `json:"content"`
	ServerID  string `json:"server_id"`
	Timestamp string `json:"timestamp"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[string]map[*websocket.Conn]bool)

var mu sync.Mutex

var broadcast = make(chan Message)

func (cfg *apiConfig) JoinNyxServer(w http.ResponseWriter, r *http.Request) {

	serverId := r.URL.Query().Get("server_id")
	token := r.URL.Query().Get("token")
	if serverId == "" {
		respondWithError(w, 401, "need a proper query", fmt.Errorf("query for server id is empty"))
		return
	}
	if token == "" {
		respondWithError(w, 401, "Unauthorised access", fmt.Errorf("query for access token is empty"))
		return
	}

	_, err := auth.ValidateJWT(token, cfg.secretToken)
	if err != nil {
		respondWithError(w, 401, "Unauthorised", err)
		return
	}

	serverData, err := queries.GetNyxServerByServerId(cfg.db, serverId)
	if err != nil {
		respondWithError(w, 401, "the server if does not exists", err)
		return
	}

	expTime, err := time.Parse(time.RFC3339, serverData.ExpiresAt)
	if err != nil {
		respondWithError(w, 500, "could not parse the sting to time.time", err)
		return
	}

	if expTime.Before(time.Now()) {
		respondWithError(w, 401, "nyx_server has been expired", fmt.Errorf("server id has been expired"))
		return
	}

	mu.Lock()
	count := len(clients[serverId])
	mu.Unlock()

	if count >= serverData.MaxParticipants {
		respondWithError(w, 401, "Unauthorised", fmt.Errorf("the room is in its full capacity"))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		respondWithError(w, 500, "upgrade failed", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	if clients[serverId] == nil {
		clients[serverId] = make(map[*websocket.Conn]bool)
	}
	clients[serverId][conn] = true
	mu.Unlock()

	log.Printf("A client joined the server with server id %v", serverId)

	go func() {
		time.Sleep(time.Until(expTime))

		mu.Lock()
		if clients[serverId] != nil {
			for conn := range clients[serverId] {
				conn.WriteMessage(websocket.CloseMessage, []byte("nyx server has expired"))
				conn.Close()
			}
			delete(clients, serverId)
		}
		mu.Unlock()

		log.Printf("Server %v has expired. Closing all connections.", serverId)

	}()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			mu.Lock()
			delete(clients[serverId], conn)
			mu.Unlock()
			break
		}
		broadcast <- msg
	}
}

func HandleBroadcasts() {
	for {
		msg := <-broadcast
		mu.Lock()
		for conn := range clients[msg.ServerID] {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Println("Write error:", err)
				conn.Close()
				delete(clients[msg.ServerID], conn)
			}
		}
		mu.Unlock()
	}
}
