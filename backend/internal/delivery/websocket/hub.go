package websocket

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/alligatorO15/taskMind/backend/internal/config"
	"github.com/alligatorO15/taskMind/backend/internal/domain/models"
	"github.com/alligatorO15/taskMind/backend/internal/infrastructure/logger"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Hub - центральный хаб для управления WebSocket-соединениями
// Отвечает за регистрацию клиентов, удаление отключившихся и рассылку уведомлений конкретным пользователям
type Hub struct {
	// clients - карта пользовательских соединений: userID -> набор подключений
	clients map[primitive.ObjectID]map[*Client]bool
	mu      sync.RWMutex

	// Каналы для регистрации и удаления клиентов
	register   chan *Client
	unregister chan *Client

	upgrader websocket.Upgrader
	cfg      config.WebSocketConfig
}

// Client - отдельное WebSocket-соединение конкретного пользователя
type Client struct {
	UserID primitive.ObjectID
	Conn   *websocket.Conn
	Send   chan []byte
}

// WSMessage - сообщение, отправляемое через WebSocket
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// NewHub - создает новый WebSocket-хаб c заданной конфигурацией
func NewHub(cfg config.WebSocketConfig) *Hub {
	return &Hub{
		clients:    make(map[primitive.ObjectID]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		cfg:        cfg,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// Разрешаем подключение с любого origin (в production нужно ограничить)
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// Run - запускает основной цикл обработки событий хаба
// Обрабатыввает регистрацию новых клиентов и удаление отключившихся
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
			h.mu.Unlock()
			logger.Logger.Infof("WebSocket: пользователь %s подключен", client.UserID.Hex())

		case client := <-h.unregister:
			h.mu.Lock()
			if conns, ok := h.clients[client.UserID]; ok {
				if _, exists := conns[client]; exists {
					delete(conns, client)
					close(client.Send)
					if len(conns) == 0 {
						delete(h.clients, client.UserID)
					}
				}
			}
			h.mu.Unlock()
			logger.Logger.Infof("WebSocket: пользователь %s отключен", client.UserID.Hex())
		}
	}
}

// HandleWebSocket - обрабатывает HTTP-запрос на обновление до WebSocket
// Создает клиента, регистрирует его и запускает горутины чтения/записи
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request, userID primitive.ObjectID) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Logger.Errorf("Ошибка обновлени до WebSocket: %v", err)
		return
	}

	client := &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	h.register <- client

	go h.writePump(client)
	go h.readPump(client)
}

// SendNotification — отправляет уведомление всем подключениям конкретного пользователя.
// Если пользователь не подключен — уведомление пропускается (сохранено в БД).
func (h *Hub) SendNotification(userID primitive.ObjectID, notification *models.Notification) {
	msg := WSMessage{
		Type:    "notification",
		Payload: notification,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Logger.Errorf("Ошибка сериализации WS-сообщения: %v", err)
		return
	}

	h.mu.RLock()
	clients, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok {
		return
	}

	for client := range clients {
		select {
		case client.Send <- data:
		default:
			// Буфер переполнен — отключаем клиента
			h.unregister <- client
		}
	}
}

// writePump - горутина записи сообщений в WebSocket-соединение клиента.
// Отправляет данные из канала Send и периодические ping-сообщения
func (h *Hub) writePump(client *Client) {
	ticker := time.NewTicker(h.cfg.PingInterval)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(h.cfg.WriteTimeout))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(h.cfg.WriteTimeout))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump - горутина чтения сообщений от WebSocket-клиента
// Обрабатывает pong-ответы и закрывает соединение при ошибках
func (h *Hub) readPump(client *Client) {
	defer func() {
		h.unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadDeadline(time.Now().Add(h.cfg.PongTimeout))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(h.cfg.PongTimeout))
		return nil
	})

	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				logger.Logger.Warnf("WebSocket ошибка чтения: %v", err)
			}
			break
		}
	}
}

// GetOnlineUsers - возвращает количество подключенных пользователей (для мониторинга)
func (h *Hub) GetOnlineUsers() int {
	h.mu.RLock()
	defer h.mu.Unlock()
	return len(h.clients)
}
