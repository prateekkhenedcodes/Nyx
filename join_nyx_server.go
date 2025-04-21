package main

import (
	"encoding/json"
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
	Pseudonym string          `json:"pseudonym"`
	Content   string          `json:"content"`
	ServerID  string          `json:"server_id"`
	Timestamp string          `json:"timestamp"`
	Conn      *websocket.Conn `json:"-"` // omit from JSON
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type clientInfo struct {
	websocketConn *websocket.Conn
	accessToken   string
}

var clients = make(map[string][]clientInfo)

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

	mu.Lock()
	if clients[serverId] == nil {
		clients[serverId] = make([]clientInfo, 0)
	}
	clients[serverId] = append(clients[serverId], clientInfo{
		websocketConn: conn,
		accessToken:   token,
	})
	mu.Unlock()

	log.Printf("A client joined the server with server id %v", serverId)

	go func() {
		time.Sleep(time.Until(expTime))

		mu.Lock()
		if clients[serverId] != nil {
			for _, clientInfo := range clients[serverId] {
				err := clientInfo.websocketConn.WriteMessage(websocket.CloseMessage, []byte("nyx server has expired"))
				if err != nil {
					log.Print("could not close the nyx server chat room")
				}
				err = clientInfo.websocketConn.Close()
				if err != nil {
					log.Print("could not close the websocket connection")
				}
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
			for i, clientInfo := range clients[serverId] {
				if clientInfo.websocketConn == conn {
					mu.Lock()
					clients[serverId] = append(clients[serverId][:i], clients[serverId][i+1:]...)
					mu.Unlock()
					break
				}
			}
		}
		broadcast <- Message{
			Pseudonym: msg.Pseudonym,
			Content:   msg.Content,
			ServerID:  msg.ServerID,
			Timestamp: msg.Timestamp,
			Conn:      conn,
		}
	}
}

func HandleBroadcasts() {
	for {
		msg := <-broadcast
		mu.Lock()
		for i, clientInfo := range clients[msg.ServerID] {
			if clientInfo.websocketConn == msg.Conn {
				continue
			}
			err := clientInfo.websocketConn.WriteJSON(msg)
			if err != nil {
				log.Println("Write error:", err)
				err = clientInfo.websocketConn.Close()
				if err != nil {
					log.Print("could not close the websocket connection")
				}
				clients[msg.ServerID] = append(clients[msg.ServerID][:i], clients[msg.ServerID][i+1:]...)
			}
		}
		mu.Unlock()
	}
}

func (cfg *apiConfig) DisconnectNyxServer(w http.ResponseWriter, r *http.Request) {

	type parameter struct {
		ServerID string `json:"server_id"`
	}

	header := r.Header
	token, err := auth.GetBearerToken(header)
	if err != nil {
		respondWithError(w, 401, "Unauthorised", err)
		return
	}

	_, err = auth.ValidateJWT(token, cfg.secretToken)
	if err != nil {
		respondWithError(w, 401, "Unauthorised", err)
		return
	}

	params := parameter{}

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "could not decode the parameters", err)
		return
	}
	mu.Lock()
	for i, clientInfo := range clients[params.ServerID] {
		if clientInfo.accessToken == token {
			clients[params.ServerID] = append(clients[params.ServerID][:i], clients[params.ServerID][i+1:]...)
			clientInfo.websocketConn.Close()
			mu.Unlock()
			respondWithJSON(w, 200, nil)
			return
		}
	}
	mu.Unlock()
	respondWithError(w, 200, "you are not connected to this server", fmt.Errorf("you are not in this server, check the serverID once"))
}
