package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, add origin validation
	},
}

type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
	Type     string `json:"type"` // "message", "join", "leave"
	CourseID string `json:"course_id"`
}

type Client struct {
	Hub       *Hub
	Conn      *websocket.Conn
	Send      chan Message
	CourseID  string
	Username  string
	closeOnce sync.Once
}

type Hub struct {
	Rooms      map[string]map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.Rooms[client.CourseID] == nil {
				h.Rooms[client.CourseID] = make(map[*Client]bool)
			}
			h.Rooms[client.CourseID][client] = true
			h.mu.Unlock()
			log.Printf("Client joined room: %s", client.CourseID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.Rooms[client.CourseID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					client.closeOnce.Do(func() {
						close(client.Send)
					})
					if len(clients) == 0 {
						delete(h.Rooms, client.CourseID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client left room: %s", client.CourseID)

		case message := <-h.Broadcast:
			h.mu.Lock()
			clients := h.Rooms[message.CourseID]
			for client := range clients {
				select {
				case client.Send <- message:
				default:
					client.closeOnce.Do(func() {
						close(client.Send)
					})
					delete(clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
		c.Hub.Broadcast <- Message{
			Username: c.Username,
			Text:     "Kullanıcı ayrıldı",
			Type:     "leave",
			CourseID: c.CourseID,
		}
	}()

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		msg.CourseID = c.CourseID
		msg.Username = c.Username
		msg.Type = "message"
		c.Hub.Broadcast <- msg
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteJSON(message)
		}
	}
}

func (h *Hub) HandleWebSocket(c *gin.Context) {
	courseID := c.Param("courseId")
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)
	
	username := fmt.Sprintf("User_%d", userID)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		Hub:      h,
		Conn:     conn,
		Send:     make(chan Message, 256),
		CourseID: courseID,
		Username: username,
	}

	h.Register <- client

	// Broadcast join message
	h.Broadcast <- Message{
		Username: username,
		Text:     "Kullanıcı katıldı",
		Type:     "join",
		CourseID: courseID,
	}

	go client.WritePump()
	go client.ReadPump()
}
