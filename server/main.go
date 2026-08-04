package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Nickname string `json:"nickname,omitempty"`
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("升级连接失败: %v", err)
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		userID:   r.URL.Query().Get("userId"),
		nick:     r.URL.Query().Get("nick"),
		room:     r.URL.Query().Get("room"),
		lastPong: time.Now(),
	}

	if client.userID == "" {
		client.userID = conn.RemoteAddr().String()
	}
	if client.nick == "" {
		client.nick = "匿名用户"
	}
	if client.room == "" {
		client.room = "大厅"
	}

	go client.writePump()
	go client.readPump()

	hub.joinRoom(client, client.room)
}

func handleRegister(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "请求格式错误"})
			return
		}

		if req.Username == "" || req.Password == "" || req.Nickname == "" {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "所有字段不能为空"})
			return
		}

		if err := store.Register(req.Username, req.Password, req.Nickname); err != nil {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(AuthResponse{Success: true, Message: "注册成功", Nickname: req.Nickname})
	}
}

func handleLogin(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "请求格式错误"})
			return
		}

		user, err := store.Login(req.Username, req.Password)
		if err != nil {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(AuthResponse{Success: true, Message: "登录成功", Nickname: user.Nickname})
	}
}

func main() {
	store, err := NewStore("data")
	if err != nil {
		log.Fatalf("初始化存储失败: %v", err)
	}

	hub := NewHub(store)
	go hub.Run()

	http.HandleFunc("/api/register", handleRegister(store))
	http.HandleFunc("/api/login", handleLogin(store))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	http.Handle("/", http.FileServer(http.Dir("./static")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("聊天服务器运行在 :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
