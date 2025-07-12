package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// var (
//
//	newline = []byte{'\n'}
//	space   = []byte{' '}
//
// )
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (c *ChatClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var message Message
		err := c.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		message.RoomID = c.roomID
		message.UserID = c.userID
		message.CreatedAt = time.Now()

		c.hub.broadcast <- message
	}
}

func (c *ChatClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			messageJSON, _ := json.Marshal(message)
			w.Write(messageJSON)

			if err := w.Close(); err != nil {
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

// func ServeWs(w http.ResponseWriter, r *http.Request) {
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	//Get room's id from client...
// 	queryValues := r.URL.Query()
// 	roomId := queryValues.Get("roomId")
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	c := &connection{send: make(chan []byte, 256), ws: ws}
// 	s := subscription{c, roomId}
// 	H.register <- s
// 	go s.writePump()
// 	go s.readPump()
// }

type Message struct {
	ID          string    `json:"id" db:"id"`
	RoomID      string    `json:"room_id" db:"room_id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Content     string    `json:"content" db:"content"`
	MessageType string    `json:"message_type" db:"message_type"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`

	// Join fields for display
	Username   string `json:"username,omitempty"`
	ProfilePic string `json:"profile_pic,omitempty"`

	// Metadata
	EditedAt  *time.Time     `json:"edited_at,omitempty"`
	ReplyToID *string        `json:"reply_to_id,omitempty"`
	Reactions map[string]int `json:"reactions,omitempty"`
}

type ChatClient struct {
	conn   *websocket.Conn
	send   chan Message
	roomID string
	userID string
	hub    *ChatHub
}

// WebSocket handler
func HandleWebSocket(c *gin.Context) {
	roomID := c.Param("roomId")

	userID := c.Query("userId")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &ChatClient{
		conn:   conn,
		send:   make(chan Message, 256),
		roomID: roomID,
		userID: userID,
		hub:    H,
	}

	H.register <- client

	go client.writePump()
	go client.readPump()
}
