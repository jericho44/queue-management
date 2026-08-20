package websocket

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"queue-management-tenant/backend/pkg/logger"
	"queue-management-tenant/backend/pkg/redis"
)

type Client struct {
	Conn    *websocket.Conn
	Channel string
	Send    chan []byte
}

type WSHub struct {
	clients    map[string]map[*Client]bool // channel -> map of clients
	register   chan *Client
	unregister chan *Client
	broadcast  chan BroadcastMessage
	redis      *redis.RedisClient
	logger     *logger.Logger
	mu         sync.RWMutex
}

type BroadcastMessage struct {
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Data    any    `json:"data"`
}

func NewWSHub(redisClient *redis.RedisClient, log *logger.Logger) *WSHub {
	return &WSHub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan BroadcastMessage, 100),
		redis:      redisClient,
		logger:     log,
	}
}

func (h *WSHub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case client := <-h.register:
			h.mu.Lock()
			if _, exists := h.clients[client.Channel]; !exists {
				h.clients[client.Channel] = make(map[*Client]bool)
			}
			h.clients[client.Channel][client] = true
			h.mu.Unlock()
			h.logger.Info("WebSocket client connected to channel: %s", client.Channel)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.Channel]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.clients, client.Channel)
					}
				}
			}
			h.mu.Unlock()
			h.logger.Info("WebSocket client disconnected from channel: %s", client.Channel)

		case msg := <-h.broadcast:
			payload, err := json.Marshal(msg)
			if err != nil {
				continue
			}

			h.mu.RLock()
			if clients, ok := h.clients[msg.Channel]; ok {
				for client := range clients {
					select {
					case client.Send <- payload:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *WSHub) BroadcastToChannel(channel string, event string, data any) {
	msg := BroadcastMessage{
		Channel: channel,
		Event:   event,
		Data:    data,
	}
	h.broadcast <- msg

	// Also publish to Redis Pub/Sub for distributed instances
	if h.redis != nil {
		payload, err := json.Marshal(msg)
		if err == nil {
			_ = h.redis.Publish(context.Background(), channel, payload)
		}
	}
}

func (h *WSHub) HandleWS(c *websocket.Conn) {
	channel := c.Query("channel")
	if channel == "" {
		_ = c.Close()
		return
	}

	client := &Client{
		Conn:    c,
		Channel: channel,
		Send:    make(chan []byte, 256),
	}

	h.register <- client

	// Writer goroutine
	go func() {
		defer func() {
			_ = c.Close()
		}()
		for message := range client.Send {
			if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
				break
			}
		}
	}()

	// Reader loop (keep-alive)
	for {
		if _, _, err := c.ReadMessage(); err != nil {
			h.unregister <- client
			break
		}
	}
}
