package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

type Edit struct {
	Position int    `json:"position"`
	Length   int    `json:"length"`
	Text     string `json:"text"`
	Type     string `json:"type"`
	ClientID string `json:"clientId"`
}

var (
	text string
	mu   sync.Mutex
	rdb  *redis.Client
)

func init() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
}

func applyEdit(text *string, edit Edit) error {
	if edit.Position < 0 || edit.Position > len(*text) {
		return fmt.Errorf("invalid position")
	}

	switch edit.Type {
	case "insert":
		*text = (*text)[:edit.Position] + edit.Text + (*text)[edit.Position:]
	case "delete":
		if edit.Position+edit.Length > len(*text) {
			return fmt.Errorf("invalid length")
		}
		*text = (*text)[:edit.Position] + (*text)[edit.Position+edit.Length:]
	default:
		return fmt.Errorf("unknown edit type")
	}
	return nil
}

func publishDelta(edit Edit) error {
	data, err := json.Marshal(edit)
	if err != nil {
		return err
	}

	broadcastToClients(data)
	return nil
}

var clients = make(map[*websocket.Conn]bool)
var clientsMu sync.Mutex

func addClient(conn *websocket.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	clients[conn] = true
}

func removeClient(conn *websocket.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	delete(clients, conn)
}

func broadcastToClients(message []byte) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	addClient(conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	removeClient(conn)
}

func handleEdit(w http.ResponseWriter, r *http.Request) {
	var edit Edit
	err := json.NewDecoder(r.Body).Decode(&edit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	err = applyEdit(&text, edit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = publishDelta(edit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(text))
}

func handleGetText(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Write([]byte(text))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/edit", handleEdit)
	http.HandleFunc("/text", handleGetText)
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":8080", nil)
}
