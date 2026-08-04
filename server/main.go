package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Username string `json:"username,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Email    string `json:"email,omitempty"`
	Bio      string `json:"bio,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

type ProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Bio      string `json:"bio"`
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

		json.NewEncoder(w).Encode(AuthResponse{
			Success:  true,
			Message:  "登录成功",
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Bio:      user.Bio,
			Avatar:   user.Avatar,
		})
	}
}

func handleGetProfile(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		username := r.URL.Query().Get("username")
		if username == "" {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "缺少用户名"})
			return
		}

		user, err := store.GetUser(username)
		if err != nil {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(AuthResponse{
			Success:  true,
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Bio:      user.Bio,
			Avatar:   user.Avatar,
		})
	}
}

func handleUpdateProfile(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req ProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "请求格式错误"})
			return
		}

		if req.Username == "" {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "缺少用户名"})
			return
		}

		if err := store.UpdateProfile(req.Username, req.Email, req.Nickname, req.Bio); err != nil {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(AuthResponse{Success: true, Message: "更新成功", Nickname: req.Nickname, Email: req.Email, Bio: req.Bio})
	}
}

func handleUploadAvatar(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		username := r.FormValue("username")
		if username == "" {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "缺少用户名"})
			return
		}

		r.ParseMultipartForm(2 << 20)
		file, handler, err := r.FormFile("avatar")
		if err != nil {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "读取文件失败"})
			return
		}
		defer file.Close()

		ext := filepath.Ext(handler.Filename)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "仅支持 jpg/png/gif 格式"})
			return
		}

		avatarDir := "data/avatars"
		os.MkdirAll(avatarDir, 0755)

		hash := sha256.Sum256([]byte(username + time.Now().String()))
		filename := fmt.Sprintf("%s%s", hex.EncodeToString(hash[:8]), ext)
		dst, err := os.Create(filepath.Join(avatarDir, filename))
		if err != nil {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "保存文件失败"})
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: "写入文件失败"})
			return
		}

		avatarPath := "/avatars/" + filename
		if err := store.UpdateAvatar(username, avatarPath); err != nil {
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(AuthResponse{Success: true, Message: "头像上传成功", Avatar: avatarPath})
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
	http.HandleFunc("/api/profile", handleGetProfile(store))
	http.HandleFunc("/api/profile/update", handleUpdateProfile(store))
	http.HandleFunc("/api/avatar/upload", handleUploadAvatar(store))

	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("data/avatars"))))

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
