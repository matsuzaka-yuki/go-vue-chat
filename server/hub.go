package main

import (
	"log"
	"sync"
	"time"
)

const (
	staleCheckInterval = 30 * time.Second
	staleTimeout       = 90 * time.Second
	historyLimit       = 50
)

type Hub struct {
	mu         sync.RWMutex
	rooms      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
	store      *Store
}

func NewHub(store *Store) *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message, 256),
		store:      store,
	}
}

func (h *Hub) Run() {
	staleTicker := time.NewTicker(staleCheckInterval)
	defer staleTicker.Stop()

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.rooms[client.room]; !ok {
				h.rooms[client.room] = make(map[*Client]bool)
			}
			h.rooms[client.room][client] = true
			h.mu.Unlock()

			h.sendHistory(client)
			h.notifyRoomUsers(client.room)
			log.Printf("%s 加入了房间 %s", client.nick, client.room)

		case client := <-h.unregister:
			h.removeClient(client)

		case msg := <-h.broadcast:
			if err := h.store.SaveMessage(msg.Room, msg.Nick, msg.Content, msg.Time); err != nil {
				log.Printf("保存消息失败: %v", err)
			}

			h.mu.RLock()
			clients := h.rooms[msg.Room]
			h.mu.RUnlock()

			data := msg.encode()
			for client := range clients {
				select {
				case client.send <- data:
				default:
					h.removeClient(client)
				}
			}

		case <-staleTicker.C:
			h.cleanStaleClients()
		}
	}
}

func (h *Hub) removeClient(client *Client) {
	h.mu.Lock()
	if client.disconnected {
		h.mu.Unlock()
		return
	}
	client.disconnected = true

	if clients, ok := h.rooms[client.room]; ok {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			close(client.send)
			if len(clients) == 0 {
				delete(h.rooms, client.room)
			}
		}
	}
	h.mu.Unlock()

	if client.room != "" {
		h.notifyRoomUsers(client.room)
		log.Printf("%s 离开了房间 %s", client.nick, client.room)
	}
}

func (h *Hub) cleanStaleClients() {
	h.mu.RLock()
	var stale []*Client
	now := time.Now()
	for _, clients := range h.rooms {
		for c := range clients {
			if now.Sub(c.lastPong) > staleTimeout {
				stale = append(stale, c)
			}
		}
	}
	h.mu.RUnlock()

	for _, c := range stale {
		log.Printf("检测到超时客户端 %s，强制断开", c.nick)
		c.conn.Close()
	}
}

func (h *Hub) sendHistory(client *Client) {
	stored, err := h.store.LoadHistory(client.room, historyLimit)
	if err != nil {
		log.Printf("加载历史消息失败: %v", err)
		return
	}

	for _, sm := range stored {
		m := Message{
			Type:    "message",
			Nick:    sm.Nick,
			Content: sm.Content,
			Time:    sm.CreatedAt,
		}
		select {
		case client.send <- m.encode():
		default:
			return
		}
	}
}

func (h *Hub) joinRoom(c *Client, room string) {
	h.register <- c
}

func (h *Hub) leaveRoom(c *Client) {
	h.unregister <- c
}

func (h *Hub) listRooms() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	rooms := make([]string, 0, len(h.rooms))
	for room := range h.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}

func (h *Hub) notifyRoomUsers(room string) {
	h.mu.RLock()
	clients := h.rooms[room]
	h.mu.RUnlock()

	users := make([]string, 0, len(clients))
	for c := range clients {
		users = append(users, c.nick)
	}

	msg := Message{
		Type:  "users",
		Room:  room,
		Users: users,
	}
	data := msg.encode()
	for c := range clients {
		select {
		case c.send <- data:
		default:
		}
	}
}
