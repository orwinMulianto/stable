package trainerchat

import (
	"sync"
)

// Client merepresentasikan satu koneksi WebSocket (user atau trainer)
type Client struct {
	SessionID uint
	Role      string // "user" atau "trainer"
	Send      chan []byte
}

// Hub mengelola semua client yang aktif, dikelompokkan per session
type Hub struct {
	mu       sync.RWMutex
	sessions map[uint]map[*Client]bool // sessionID → set of clients
}

func NewHub() *Hub {
	return &Hub{
		sessions: make(map[uint]map[*Client]bool),
	}
}

// Register mendaftarkan client baru ke session
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.sessions[client.SessionID] == nil {
		h.sessions[client.SessionID] = make(map[*Client]bool)
	}
	h.sessions[client.SessionID][client] = true
}

// Unregister menghapus client dari session dan menutup channel-nya
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if room, ok := h.sessions[client.SessionID]; ok {
		if _, ok := room[client]; ok {
			delete(room, client)
			close(client.Send)
		}
		if len(room) == 0 {
			delete(h.sessions, client.SessionID)
		}
	}
}

// Broadcast mengirim pesan ke semua client dalam session, kecuali pengirim
func (h *Hub) Broadcast(sender *Client, message []byte) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    for client := range h.sessions[sender.SessionID] {
        select {
        case client.Send <- message:
        default:
        }
    }
}

func (h *Hub) ClientCount(sessionID uint) int {
    h.mu.RLock()
    defer h.mu.RUnlock()
    return len(h.sessions[sessionID])
}

var globalHub = NewHub()

func GetGlobalHub() *Hub {
    return globalHub
}