package trainerchat

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Izinkan semua origin (sesuaikan untuk production)
	CheckOrigin: func(r *http.Request) bool { return true },
}


// WSHandler sekarang butuh service untuk persist pesan ke DB
func WSHandler(hub *Hub, service Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := parseUintParam(c.Param("session_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid session id"})
			return
		}

		role := c.Query("role")
		if role != "user" && role != "trainer" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "role must be 'user' or 'trainer'"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("ws upgrade error: %v", err)
			return
		}

		client := &Client{
			SessionID: sessionID,
			Role:      role,
			Send:      make(chan []byte, 64),
		}
		hub.Register(client)

		go writePump(client, conn)
		readPump(client, conn, hub, service) // <-- service diteruskan
	}
}

func readPump(client *Client, conn *websocket.Conn, hub *Hub, service Service) {
	defer func() {
		hub.Unregister(client)
		conn.Close()
	}()

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
    _, raw, err := conn.ReadMessage()
    if err != nil {
        if websocket.IsUnexpectedCloseError(err,
            websocket.CloseGoingAway,
            websocket.CloseAbnormalClosure,
        ) {
            log.Printf("ws read error session=%d role=%s: %v", client.SessionID, client.Role, err)
        }
        break
    }

    var msg WSMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}
		// Fallback: kalau Content kosong, coba Message
		if msg.Content == "" {
			msg.Content = msg.Message
		}
		msg.Sender = client.Role

    if msg.Type != "message" {
        // typing/stop_typing/ping tetap broadcast langsung tanpa disimpan
        out, err := json.Marshal(msg)
        if err == nil {
            hub.Broadcast(client, out)
        }
        continue
    }

    // --- persist ke DB lewat service yang sudah ada ---
    saved, err := service.SendMessage(client.SessionID, SendMessageRequest{
        Sender:  client.Role,
        Message: msg.Content,
    })
    if err != nil {
        log.Printf("failed to persist ws message session=%d role=%s: %v", client.SessionID, client.Role, err)
        continue
    }

    out, err := json.Marshal(WSMessage{
        Type:      "message",
        Sender:    saved.Sender,
        Content:   saved.Message, // <-- ganti dari Message ke Content
        Timestamp: saved.CreatedAt.Format("15:04"),
    })
    if err != nil {
        continue
    }

    hub.Broadcast(client, out)
	log.Printf("ws broadcast session=%d sender=%s content=%s clients_in_session=%d", 
    client.SessionID, client.Role, msg.Content, hub.ClientCount(client.SessionID))
}
}

// writePump menulis pesan dari channel Send ke WebSocket
func writePump(client *Client, conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// sessionIDStr helper untuk logging
func sessionIDStr(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}