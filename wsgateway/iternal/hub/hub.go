package hub

import (
	"log"
	"sync"
)

type Hub struct {
	Register   chan *Client
	unregister chan *Client

	// userID -> map[connID]*Client
	clients map[string]map[string]*Client

	broadcast chan MessageToBroadcast
	// conversation_id -> list of userIDs
	Rooms map[string][]string

	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan MessageToBroadcast),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[string]*Client),
		Rooms:      make(map[string][]string),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()

			if _, exists := h.clients[client.UserID]; !exists {
				h.clients[client.UserID] = make(map[string]*Client)
			}
			h.clients[client.UserID][client.ConnID] = client

			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()

			if conns, exists := h.clients[client.UserID]; exists {
				delete(conns, client.ConnID)
				if len(conns) == 0 {
					delete(h.clients, client.UserID)
				}
			}

			close(client.Send)

			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.deliverMessage(msg)
		}
	}
}

func (h *Hub) deliverMessage(msg MessageToBroadcast) {
	h.mu.RLock()

	userIDs, ok := h.Rooms[msg.RoomID]

	h.mu.RUnlock()
	log.Println(msg.Text)
	if !ok {
		return // Nil room
	}
	for _, uID := range userIDs {

		if msg.From == uID {
			continue
		}

		if devices, online := h.clients[uID]; online {
			for _, client := range devices {
				select {
				case client.Send <- []byte(msg.Text):

				default:
				}
			}
		}
	}
}
