package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
	Avatar    string `json:"avatar"`
	MediaURL  string `json:"mediaUrl,omitempty"`
	MediaType string `json:"mediaType,omitempty"`
}

type Store struct {
	mu        sync.RWMutex
	dataDir   string
	users     map[string]*User
	messages  []StoredMessage
	following map[string]map[string]bool
}

func NewStore(dataDir string) (*Store, error) {
	os.MkdirAll(dataDir, 0755)
	s := &Store{
		dataDir:   dataDir,
		users:     make(map[string]*User),
		messages:  []StoredMessage{},
		following: make(map[string]map[string]bool),
	}
	if err := s.loadUsers(); err != nil {
		return nil, err
	}
	if err := s.loadMessages(); err != nil {
		return nil, err
	}
	if err := s.loadRelations(); err != nil {
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

func (s *Store) GetAvatarByNick(nick string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, u := range s.users {
		if u.Nickname == nick || u.Username == nick {
			return u.Avatar
		}
	}
	return ""
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

func (s *Store) SaveMessage(room, nick, content, avatar, mediaURL, mediaType string, createdAt int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages = append(s.messages, StoredMessage{
		Room:      room,
		Nick:      nick,
		Content:   content,
		Avatar:    avatar,
		MediaURL:  mediaURL,
		MediaType: mediaType,
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

func (s *Store) Follow(username, target string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[target]; !ok {
		return fmt.Errorf("用户不存在")
	}
	if username == target {
		return fmt.Errorf("不能关注自己")
	}
	if s.following[username] == nil {
		s.following[username] = make(map[string]bool)
	}
	s.following[username][target] = true
	return s.saveRelations()
}

func (s *Store) Unfollow(username, target string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.following[username] != nil {
		delete(s.following[username], target)
	}
	return s.saveRelations()
}

func (s *Store) IsFollowing(username, target string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.following[username] != nil && s.following[username][target]
}

func (s *Store) IsMutual(a, b string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.following[a] != nil && s.following[a][b] &&
		s.following[b] != nil && s.following[b][a]
}

func (s *Store) GetFollowing(username string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]string, 0, len(s.following[username]))
	for u := range s.following[username] {
		result = append(result, u)
	}
	return result
}

func (s *Store) GetFollowers(username string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []string
	for u, set := range s.following {
		if set[username] {
			result = append(result, u)
		}
	}
	return result
}

func (s *Store) SearchUsers(query string) []*User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*User
	for _, u := range s.users {
		if query == "" || strings.Contains(u.Username, query) || strings.Contains(u.Nickname, query) {
			result = append(result, u)
		}
	}
	return result
}

func (s *Store) loadRelations() error {
	data, err := os.ReadFile(s.dataDir + "/relations.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var raw map[string][]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for user, list := range raw {
		s.following[user] = make(map[string]bool)
		for _, u := range list {
			s.following[user][u] = true
		}
	}
	return nil
}

func (s *Store) saveRelations() error {
	raw := make(map[string][]string)
	for user, set := range s.following {
		list := make([]string, 0, len(set))
		for u := range set {
			list = append(list, u)
		}
		raw[user] = list
	}
	data, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.dataDir+"/relations.json", data, 0644)
}
