package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Client struct {
	hub          *Hub
	conn         *websocket.Conn
	send         chan []byte
	room         string
	nick         string
	username     string
	userID       string
	lastPong     time.Time
	disconnected bool
}

type Message struct {
	Type      string   `json:"type"`
	UserID    string   `json:"userId,omitempty"`
	Nick      string   `json:"nick,omitempty"`
	Content   string   `json:"content,omitempty"`
	Room      string   `json:"room,omitempty"`
	To        string   `json:"to,omitempty"`
	Rooms     []string `json:"rooms,omitempty"`
	Users     []string `json:"users,omitempty"`
	Time      int64    `json:"time,omitempty"`
	Avatar    string   `json:"avatar,omitempty"`
	MediaURL  string   `json:"mediaUrl,omitempty"`
	MediaType string   `json:"mediaType,omitempty"`
	FileSize  int64    `json:"fileSize,omitempty"`
}

func (m Message) encode() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.lastPong = time.Now()
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msgBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("读取错误: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("无效消息: %v", err)
			continue
		}

		msg.UserID = c.username
		msg.Nick = c.nick
		msg.Time = time.Now().UnixMilli()

		switch msg.Type {
		case "join":
			c.room = msg.Room
			c.nick = msg.Nick
			c.hub.joinRoom(c, msg.Room)
		case "leave":
			c.hub.leaveRoom(c)
		case "message":
			msg.Type = "message"
			c.hub.broadcast <- msg
		case "private":
			msg.Type = "private"
			c.hub.broadcast <- msg
		case "rooms":
			msg.Type = "rooms"
			msg.Rooms = c.hub.listRooms()
			c.send <- c.encode(msg)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) encode(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}
