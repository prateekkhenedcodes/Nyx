package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
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
	if serverId == "" {
		respondWithError(w, 401, "need a proper query ", fmt.Errorf("query is empty"))
		return
	}

	_, err := queries.GetNyxServerByServerId(cfg.db, serverId)
	if err != nil {
		respondWithError(w, 401, "the server if does not exists", err)
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
