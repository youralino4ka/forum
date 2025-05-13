package message

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // В production замените на проверку origin
	},
}

type WSClient struct {
	conn     *websocket.Conn
	send     chan []byte
	userID   int
	username string
}

type WSHub struct {
	clients    map[*WSClient]bool
	broadcast  chan []byte
	register   chan *WSClient
	unregister chan *WSClient
	mu         sync.Mutex
	service    *Service
	log        zerolog.Logger
}

func NewWSHub(service *Service, log zerolog.Logger) *WSHub {
	return &WSHub{
		broadcast:  make(chan []byte),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		clients:    make(map[*WSClient]bool),
		service:    service,
		log:        log,
	}
}

func (h *WSHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

			// Отправляем историю сообщений новому клиенту
			go func() {
				messages, err := h.service.GetRecentMessages(context.Background(), 50)
				if err != nil {
					h.log.Error().Err(err).Msg("Failed to get recent messages")
					return
				}

				for _, msg := range messages {
					messageData := map[string]interface{}{
						"type":    "message",
						"id":      msg.ID,
						"user_id": msg.UserID,
						"content": msg.Content,
						"time":    msg.CreatedAt.Format(time.RFC3339),
					}

					data, err := json.Marshal(messageData)
					if err != nil {
						h.log.Error().Err(err).Msg("Failed to marshal message")
						continue
					}

					client.send <- data
				}
			}()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *WSHub) HandleWebSocket(w http.ResponseWriter, r *http.Request, userID int, username string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to upgrade connection to WebSocket")
		return
	}

	client := &WSClient{
		conn:     conn,
		send:     make(chan []byte, 256),
		userID:   userID,
		username: username,
	}

	h.register <- client

	go client.writePump()
	go client.readPump(h)
}

func (c *WSClient) readPump(h *WSHub) {
	defer func() {
		h.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.log.Error().Err(err).Msg("WebSocket read error")
			}
			break
		}

		// Обработка нового сообщения
		var msgData struct {
			Content string `json:"content"`
		}
		if err := json.Unmarshal(message, &msgData); err != nil {
			h.log.Error().Err(err).Msg("Failed to unmarshal message")
			continue
		}

		// Сохраняем сообщение в БД
		msg, err := h.service.PostMessage(context.Background(), c.userID, msgData.Content)
		if err != nil {
			h.log.Error().Err(err).Msg("Failed to post message")
			continue
		}

		// Формируем сообщение для рассылки
		response := map[string]interface{}{
			"type":    "message",
			"id":      msg.ID,
			"user_id": msg.UserID,
			"content": msg.Content,
			"time":    msg.CreatedAt.Format(time.RFC3339),
		}

		data, err := json.Marshal(response)
		if err != nil {
			h.log.Error().Err(err).Msg("Failed to marshal response")
			continue
		}

		h.broadcast <- data
	}
}

func (c *WSClient) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
