package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Протокол обмена данными
type Command struct {
	Type     string  `json:"type"`      // "mouse_move", "lock_toggle", "publish_task"
	TargetID string  `json:"target_id"` // "student1" или "all"
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Payload  string  `json:"payload"`
}

type Client struct {
	ID   string
	Role string
	Conn *websocket.Conn
	Send chan []byte
}

type Hub struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Hub) handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	id := r.URL.Query().Get("id")
	role := r.URL.Query().Get("role")
	client := &Client{ID: id, Role: role, Conn: conn, Send: make(chan []byte, 256)}

	h.mu.Lock()
	h.clients[id] = client
	h.mu.Unlock()

	log.Printf("Подключен: %s [%s]", id, role)

	// Горутина записи
	go func() {
		defer conn.Close()
		for message := range client.Send {
			client.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}()

	// Цикл чтения и маршрутизации
	defer func() {
		h.mu.Lock()
		delete(h.clients, id)
		h.mu.Unlock()
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var cmd Command
		if err := json.Unmarshal(msg, &cmd); err != nil {
			continue
		}

		// ЛОГИКА МАРШРУТИЗАЦИИ
		h.mu.RLock()
		if cmd.TargetID == "all" {
			// Рассылка всем ученикам (например, новое задание)
			for _, target := range h.clients {
				if target.Role == "student" {
					target.Send <- msg
				}
			}
		} else if target, ok := h.clients[cmd.TargetID]; ok {
			// Адресный перехват управления
			target.Send <- msg
		}
		h.mu.RUnlock()
	}
}

func main() {
	hub := &Hub{clients: make(map[string]*Client)}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/ws", hub.handleConnections)

	fmt.Println("🚀 Сервер 'МЕЛ' запущен: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
