package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type user struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token   string `json:"token,omitempty"`
	Message string `json:"message"`
}

type Server struct {
	mu     sync.RWMutex
	users  map[string]string
	tokens map[string]string
}

func NewServer() *Server {
	s := &Server{
		users:  make(map[string]string),
		tokens: make(map[string]string),
	}

	// default admin user
	s.users["admin"] = "admin"
	s.seedUsers(5)
	return s
}

func (s *Server) seedUsers(count int) {
	generated := 0
	for generated < count {
		username := randomUsername()
		if _, exists := s.users[username]; exists {
			continue
		}
		s.users[username] = randomPassword()
		generated++
	}
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "missing credentials", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[req.Username]; exists {
		http.Error(w, "user exists", http.StatusConflict)
		return
	}

	s.users[req.Username] = req.Password
	resp := authResponse{
		Message: "user registered",
	}
	writeJSON(w, resp)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	pass, ok := s.users[req.Username]
	if !ok || pass != req.Password {
		s.mu.RUnlock()
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	s.mu.RUnlock()

	token := s.issueToken(req.Username)
	resp := authResponse{
		Token:   token,
		Message: "login success",
	}
	writeJSON(w, resp)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, token, ok := s.validateToken(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	s.mu.Lock()
	delete(s.tokens, token)
	s.mu.Unlock()

	resp := authResponse{
		Message: username + " logged out",
	}
	writeJSON(w, resp)
}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, _, ok := s.validateToken(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]user, 0, len(s.users))
	for name := range s.users {
		users = append(users, user{Username: name})
	}

	resp := struct {
		Users []user `json:"users"`
		Me    string `json:"me"`
	}{
		Users: users,
		Me:    username,
	}

	writeJSON(w, resp)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) issueToken(username string) string {
	token := generateToken()
	s.mu.Lock()
	s.tokens[token] = username
	s.mu.Unlock()
	return token
}

func (s *Server) validateToken(r *http.Request) (string, string, bool) {
	token, ok := tokenFromHeader(r)
	if !ok {
		return "", "", false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	username, exists := s.tokens[token]
	return username, token, exists
}

func tokenFromHeader(r *http.Request) (string, bool) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", false
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}

	token := strings.TrimPrefix(header, prefix)
	return token, true
}

func generateToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func randomUsername() string {
	adjectives := []string{"blue", "fast", "bright", "silent", "swift", "lucky", "sunny", "calm"}
	nouns := []string{"whale", "fox", "eagle", "panda", "koala", "tiger", "otter", "lynx"}
	return fmt.Sprintf("%s-%s-%02d",
		adjectives[rand.Intn(len(adjectives))],
		nouns[rand.Intn(len(nouns))],
		rand.Intn(100),
	)
}

func randomPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func writeJSON(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	server := NewServer()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/register", server.handleRegister)
	mux.HandleFunc("/api/login", server.handleLogin)
	mux.HandleFunc("/api/logout", server.handleLogout)
	mux.HandleFunc("/api/users", server.handleUsers)
	mux.HandleFunc("/healthz", server.handleHealth)

	log.Println("backend listening on :8080")
	if err := http.ListenAndServe(":8080", withCORS(mux)); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
