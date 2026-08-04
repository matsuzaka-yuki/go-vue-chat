package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
}

type StoredMessage struct {
	Room      string `json:"room"`
	Nick      string `json:"nick"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"createdAt"`
}

type Store struct {
	mu       sync.RWMutex
	dataDir  string
	users    map[string]*User
	messages []StoredMessage
}

func NewStore(dataDir string) (*Store, error) {
	os.MkdirAll(dataDir, 0755)
	s := &Store{
		dataDir:  dataDir,
		users:    make(map[string]*User),
		messages: []StoredMessage{},
	}
	if err := s.loadUsers(); err != nil {
		return nil, err
	}
	if err := s.loadMessages(); err != nil {
		return nil, err
	}
	return s, nil
}

func hashPassword(pwd string) string {
	h := sha256.Sum256([]byte(pwd))
	return hex.EncodeToString(h[:])
}

func (s *Store) Register(username, password, nickname string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[username]; ok {
		return fmt.Errorf("用户名已存在")
	}
	s.users[username] = &User{
		Username: username,
		Password: hashPassword(password),
		Nickname: nickname,
	}
	return s.saveUsers()
}

func (s *Store) Login(username, password string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.users[username]
	if !ok {
		return nil, fmt.Errorf("用户名或密码错误")
	}
	if u.Password != hashPassword(password) {
		return nil, fmt.Errorf("用户名或密码错误")
	}
	return u, nil
}

func (s *Store) GetUser(username string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.users[username]
	if !ok {
		return nil, fmt.Errorf("用户不存在")
	}
	return u, nil
}

func (s *Store) UpdateProfile(username, email, nickname, bio string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.users[username]
	if !ok {
		return fmt.Errorf("用户不存在")
	}
	u.Email = email
	u.Nickname = nickname
	u.Bio = bio
	return s.saveUsers()
}

func (s *Store) UpdateAvatar(username, avatarPath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.users[username]
	if !ok {
		return fmt.Errorf("用户不存在")
	}
	u.Avatar = avatarPath
	return s.saveUsers()
}

func (s *Store) SaveMessage(room, nick, content string, createdAt int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages = append(s.messages, StoredMessage{
		Room:      room,
		Nick:      nick,
		Content:   content,
		CreatedAt: createdAt,
	})
	return s.saveMessages()
}

func (s *Store) LoadHistory(room string, limit int) ([]StoredMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []StoredMessage
	for i := len(s.messages) - 1; i >= 0; i-- {
		m := s.messages[i]
		if m.Room == room {
			result = append([]StoredMessage{m}, result...)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (s *Store) loadUsers() error {
	data, err := os.ReadFile(s.dataDir + "/users.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var users []*User
	if err := json.Unmarshal(data, &users); err != nil {
		return err
	}
	for _, u := range users {
		s.users[u.Username] = u
	}
	return nil
}

func (s *Store) saveUsers() error {
	users := make([]*User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.dataDir+"/users.json", data, 0644)
}

func (s *Store) loadMessages() error {
	data, err := os.ReadFile(s.dataDir + "/messages.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, &s.messages)
}

func (s *Store) saveMessages() error {
	data, err := json.MarshalIndent(s.messages, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.dataDir+"/messages.json", data, 0644)
}
