package services

import (
	"log"
	"sparkle-concept-sync/internal/models"
	"time"

	"github.com/gofiber/websocket/v2"
)

type WebSocketService struct {
	clients map[*websocket.Conn]bool
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		clients: make(map[*websocket.Conn]bool),
	}
}

// HandleWebSocket handles WebSocket connections
func (s *WebSocketService) HandleWebSocket(c *websocket.Conn) {
	defer func() {
		delete(s.clients, c)
		c.Close()
	}()

	// Add client to active connections
	s.clients[c] = true

	log.Printf("WebSocket client connected. Total clients: %d", len(s.clients))

	// Send welcome message
	welcome := models.WebSocketMessage{
		Type:      "connected",
		Data:      map[string]interface{}{"message": "WebSocket connected"},
		Timestamp: time.Now(),
	}
	c.WriteJSON(welcome)

	// Keep connection alive and handle incoming messages
	for {
		var msg map[string]interface{}
		if err := c.ReadJSON(&msg); err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		// Handle incoming messages (ping, subscribe, etc.)
		s.handleIncomingMessage(c, msg)
	}
}

// Broadcast sends a message to all connected clients
func (s *WebSocketService) Broadcast(message models.WebSocketMessage) {
	message.Timestamp = time.Now()

	for client := range s.clients {
		if err := client.WriteJSON(message); err != nil {
			log.Printf("WebSocket write error: %v", err)
			delete(s.clients, client)
			client.Close()
		}
	}
}

func (s *WebSocketService) handleIncomingMessage(c *websocket.Conn, msg map[string]interface{}) {
	msgType, ok := msg["type"].(string)
	if !ok {
		return
	}

	switch msgType {
	case "ping":
		response := models.WebSocketMessage{
			Type:      "pong",
			Data:      map[string]interface{}{"timestamp": time.Now()},
			Timestamp: time.Now(),
		}
		c.WriteJSON(response)

	case "subscribe":
		// Handle subscription to specific channels/devices
		// This would be implemented for more granular real-time updates

	default:
		log.Printf("Unknown WebSocket message type: %s", msgType)
	}
}
