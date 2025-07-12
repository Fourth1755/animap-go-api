package websocket

import (
	"log"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the
type ChatHub struct {
	// put registered clients into the room.
	clients map[string]map[*ChatClient]bool
	// Inbound messages from the clients.
	broadcast chan Message

	// Register requests from the clients.
	register chan *ChatClient

	// Unregister requests from clients.
	unregister chan *ChatClient
}

var H = &ChatHub{
	broadcast:  make(chan Message),
	register:   make(chan *ChatClient),
	unregister: make(chan *ChatClient),
	clients:    make(map[string]map[*ChatClient]bool),
}

func (h *ChatHub) Run() {
	for {
		select {
		case client := <-h.register:
			if h.clients[client.roomID] == nil {
				h.clients[client.roomID] = make(map[*ChatClient]bool)
			}
			h.clients[client.roomID][client] = true

			// Update user online status
			// h.updateUserOnlineStatus(client.userID, true)

			// Broadcast user joined event
			h.broadcastUserStatus(client.roomID, client.userID, "joined")

			log.Printf("Client registered to room %s", client.roomID)

		case client := <-h.unregister:
			if clients, ok := h.clients[client.roomID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)

					// Update user online status
					// h.updateUserOnlineStatus(client.userID, false)

					// Broadcast user left event
					h.broadcastUserStatus(client.roomID, client.userID, "left")

					log.Printf("Client unregistered from room %s", client.roomID)
				}
			}

		case message := <-h.broadcast:
			// Save message using hybrid store
			// ctx := context.Background()
			// if err := messageStore.SaveMessage(ctx, message); err != nil {
			// 	log.Printf("Failed to save message: %v", err)
			// }

			// Get user info for the message
			// message = getMessageWithUserInfo(message)

			// Broadcast to all clients in the room
			if clients, ok := h.clients[message.RoomID]; ok {
				for client := range clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(clients, client)
					}
				}
			}
		}
	}
}

// func (h *ChatHub) updateUserOnlineStatus(userID int, isOnline bool) {
// 	ctx := context.Background()

// 	// Update Redis
// 	userKey := fmt.Sprintf("user:%d", userID)
// 	if isOnline {
// 		h.redis.HSet(ctx, userKey, "online", "true", "last_seen", time.Now().Format(time.RFC3339))
// 	} else {
// 		h.redis.HSet(ctx, userKey, "online", "false", "last_seen", time.Now().Format(time.RFC3339))
// 	}

// 	// Update PostgreSQL
// 	query := `UPDATE users SET last_seen = $1 WHERE id = $2`
// 	db.Exec(query, time.Now(), userID)
// }

func (h *ChatHub) broadcastUserStatus(roomID string, userID string, status string) {
	// Broadcast user status change to room
	statusMessage := Message{
		RoomID:      roomID,
		UserID:      userID,
		Content:     status,
		MessageType: "user_status",
		CreatedAt:   time.Now(),
	}

	if clients, ok := h.clients[roomID]; ok {
		for client := range clients {
			select {
			case client.send <- statusMessage:
			default:
				close(client.send)
				delete(clients, client)
			}
		}
	}
}

// service
// func getMessageWithUserInfo(message Message) Message {
// 	// Try Redis first
// 	ctx := context.Background()
// 	userKey := fmt.Sprintf("user:%d", message.UserID)

// 	result, err := rdb.HMGet(ctx, userKey, "username", "profile_pic").Result()
// 	if err == nil && len(result) == 2 && result[0] != nil && result[1] != nil {
// 		message.Username = result[0].(string)
// 		message.ProfilePic = result[1].(string)
// 		return message
// 	}

// 	// Fallback to PostgreSQL
// 	query := `SELECT u.username, u.profile_pic FROM users u WHERE u.id = $1`
// 	err = db.QueryRow(query, message.UserID).Scan(&message.Username, &message.ProfilePic)
// 	if err != nil {
// 		log.Printf("Failed to get user info: %v", err)
// 	}

// 	// Cache in Redis for next time
// 	if message.Username != "" {
// 		rdb.HSet(ctx, userKey, "username", message.Username, "profile_pic", message.ProfilePic)
// 		rdb.Expire(ctx, userKey, time.Hour)
// 	}

// 	return message
// }
