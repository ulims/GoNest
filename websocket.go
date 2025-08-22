package gonest

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// WebSocketGateway represents a WebSocket gateway similar to NestJS
type WebSocketGateway interface {
	OnConnection(client *WebSocketClient) error
	OnDisconnection(client *WebSocketClient) error
	OnMessage(client *WebSocketClient, message *WebSocketMessage) error
	GetNamespace() string
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
	ID    string      `json:"id,omitempty"`
	Room  string      `json:"room,omitempty"`
}

// WebSocketClient represents a connected WebSocket client
type WebSocketClient struct {
	ID         string
	Connection *websocket.Conn
	Gateway    WebSocketGateway
	Rooms      map[string]bool
	Context    context.Context
	Logger     *logrus.Logger
	mutex      sync.RWMutex
}

// WebSocketServer manages WebSocket connections and gateways
type WebSocketServer struct {
	gateways map[string]WebSocketGateway
	clients  map[string]*WebSocketClient
	rooms    map[string]map[string]*WebSocketClient
	upgrader websocket.Upgrader
	logger   *logrus.Logger
	mutex    sync.RWMutex
}

// WebSocketConfig holds WebSocket configuration
type WebSocketConfig struct {
	ReadBufferSize  int
	WriteBufferSize int
	CheckOrigin     func(r *http.Request) bool
	Subprotocols    []string
}

// DefaultWebSocketConfig returns default WebSocket configuration
func DefaultWebSocketConfig() *WebSocketConfig {
	return &WebSocketConfig{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins in development
		},
		Subprotocols: []string{},
	}
}

// NewWebSocketServer creates a new WebSocket server
func NewWebSocketServer(config *WebSocketConfig, logger *logrus.Logger) *WebSocketServer {
	if config == nil {
		config = DefaultWebSocketConfig()
	}

	return &WebSocketServer{
		gateways: make(map[string]WebSocketGateway),
		clients:  make(map[string]*WebSocketClient),
		rooms:    make(map[string]map[string]*WebSocketClient),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  config.ReadBufferSize,
			WriteBufferSize: config.WriteBufferSize,
			CheckOrigin:     config.CheckOrigin,
			Subprotocols:    config.Subprotocols,
		},
		logger: logger,
	}
}

// RegisterGateway registers a WebSocket gateway
func (ws *WebSocketServer) RegisterGateway(gateway WebSocketGateway) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	namespace := gateway.GetNamespace()
	if namespace == "" {
		namespace = "/"
	}

	ws.gateways[namespace] = gateway
	ws.logger.Infof("Registered WebSocket gateway for namespace: %s", namespace)
}

// HandleConnection handles WebSocket connection upgrade
func (ws *WebSocketServer) HandleConnection(namespace string) echo.HandlerFunc {
	return func(c echo.Context) error {
		ws.mutex.RLock()
		gateway, exists := ws.gateways[namespace]
		ws.mutex.RUnlock()

		if !exists {
			return echo.NewHTTPError(http.StatusNotFound, "WebSocket gateway not found")
		}

		conn, err := ws.upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			ws.logger.WithError(err).Error("Failed to upgrade WebSocket connection")
			return err
		}

		client := &WebSocketClient{
			ID:         generateClientID(),
			Connection: conn,
			Gateway:    gateway,
			Rooms:      make(map[string]bool),
			Context:    c.Request().Context(),
			Logger:     ws.logger,
		}

		ws.mutex.Lock()
		ws.clients[client.ID] = client
		ws.mutex.Unlock()

		ws.logger.Infof("WebSocket client connected: %s", client.ID)

		// Call gateway connection handler
		if err := gateway.OnConnection(client); err != nil {
			ws.logger.WithError(err).Error("Gateway connection handler failed")
			client.Close()
			return err
		}

		// Start message handling
		go ws.handleClient(client)

		return nil
	}
}

// handleClient handles messages from a WebSocket client
func (ws *WebSocketServer) handleClient(client *WebSocketClient) {
	defer func() {
		ws.removeClient(client)
		client.Gateway.OnDisconnection(client)
		client.Close()
	}()

	client.Connection.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Connection.SetPongHandler(func(string) error {
		client.Connection.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var message WebSocketMessage
		err := client.Connection.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				ws.logger.WithError(err).Error("WebSocket read error")
			}
			break
		}

		// Handle the message
		if err := client.Gateway.OnMessage(client, &message); err != nil {
			ws.logger.WithError(err).Error("Gateway message handler failed")
			client.EmitError("message_error", err.Error())
		}
	}
}

// removeClient removes a client from the server
func (ws *WebSocketServer) removeClient(client *WebSocketClient) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	delete(ws.clients, client.ID)

	// Remove from all rooms
	for room := range client.Rooms {
		if roomClients, exists := ws.rooms[room]; exists {
			delete(roomClients, client.ID)
			if len(roomClients) == 0 {
				delete(ws.rooms, room)
			}
		}
	}

	ws.logger.Infof("WebSocket client disconnected: %s", client.ID)
}

// BroadcastToRoom broadcasts a message to all clients in a room
func (ws *WebSocketServer) BroadcastToRoom(room string, message *WebSocketMessage) {
	ws.mutex.RLock()
	roomClients, exists := ws.rooms[room]
	ws.mutex.RUnlock()

	if !exists {
		return
	}

	for _, client := range roomClients {
		if err := client.Emit(message.Event, message.Data); err != nil {
			ws.logger.WithError(err).Errorf("Failed to send message to client %s", client.ID)
		}
	}
}

// BroadcastToNamespace broadcasts a message to all clients in a namespace
func (ws *WebSocketServer) BroadcastToNamespace(namespace string, message *WebSocketMessage) {
	ws.mutex.RLock()
	gateway, exists := ws.gateways[namespace]
	ws.mutex.RUnlock()

	if !exists {
		return
	}

	for _, client := range ws.clients {
		if client.Gateway == gateway {
			if err := client.Emit(message.Event, message.Data); err != nil {
				ws.logger.WithError(err).Errorf("Failed to send message to client %s", client.ID)
			}
		}
	}
}

// WebSocketClient methods

// Emit sends a message to the client
func (wc *WebSocketClient) Emit(event string, data interface{}) error {
	wc.mutex.Lock()
	defer wc.mutex.Unlock()

	message := WebSocketMessage{
		Event: event,
		Data:  data,
	}

	return wc.Connection.WriteJSON(message)
}

// EmitError sends an error message to the client
func (wc *WebSocketClient) EmitError(event string, errorMsg string) error {
	return wc.Emit(event, map[string]string{"error": errorMsg})
}

// JoinRoom adds the client to a room
func (wc *WebSocketClient) JoinRoom(room string) {
	wc.mutex.Lock()
	wc.Rooms[room] = true
	wc.mutex.Unlock()

	// Add to server rooms
	// Note: This would need access to the server instance
	wc.Logger.Infof("Client %s joined room: %s", wc.ID, room)
}

// LeaveRoom removes the client from a room
func (wc *WebSocketClient) LeaveRoom(room string) {
	wc.mutex.Lock()
	delete(wc.Rooms, room)
	wc.mutex.Unlock()

	wc.Logger.Infof("Client %s left room: %s", wc.ID, room)
}

// Close closes the WebSocket connection
func (wc *WebSocketClient) Close() error {
	return wc.Connection.Close()
}

// GetIP returns the client's IP address
func (wc *WebSocketClient) GetIP() string {
	return wc.Connection.RemoteAddr().String()
}

// Utility functions

// generateClientID generates a unique client ID
func generateClientID() string {
	return fmt.Sprintf("client_%d", time.Now().UnixNano())
}

// WebSocket decorators and utilities

// WebSocketNamespace decorator (simulated with struct embedding)
type WebSocketNamespace struct {
	Namespace string
}

// WebSocketEvent decorator (simulated with method naming convention)
type WebSocketEvent struct {
	Event string
}

// BaseWebSocketGateway provides a base implementation for WebSocket gateways
type BaseWebSocketGateway struct {
	namespace string
	logger    *logrus.Logger
}

// NewBaseWebSocketGateway creates a new base WebSocket gateway
func NewBaseWebSocketGateway(namespace string, logger *logrus.Logger) *BaseWebSocketGateway {
	return &BaseWebSocketGateway{
		namespace: namespace,
		logger:    logger,
	}
}

// GetNamespace returns the gateway namespace
func (bg *BaseWebSocketGateway) GetNamespace() string {
	return bg.namespace
}

// OnConnection default implementation
func (bg *BaseWebSocketGateway) OnConnection(client *WebSocketClient) error {
	bg.logger.Infof("Client connected to namespace %s: %s", bg.namespace, client.ID)
	return nil
}

// OnDisconnection default implementation
func (bg *BaseWebSocketGateway) OnDisconnection(client *WebSocketClient) error {
	bg.logger.Infof("Client disconnected from namespace %s: %s", bg.namespace, client.ID)
	return nil
}

// OnMessage default implementation
func (bg *BaseWebSocketGateway) OnMessage(client *WebSocketClient, message *WebSocketMessage) error {
	bg.logger.Infof("Received message on namespace %s from client %s: %s", bg.namespace, client.ID, message.Event)
	return nil
}

// WebSocket middleware for authentication
type WebSocketAuthMiddleware struct {
	authenticator func(r *http.Request) (interface{}, error)
}

// NewWebSocketAuthMiddleware creates a new WebSocket auth middleware
func NewWebSocketAuthMiddleware(authenticator func(r *http.Request) (interface{}, error)) *WebSocketAuthMiddleware {
	return &WebSocketAuthMiddleware{
		authenticator: authenticator,
	}
}

// Middleware returns the middleware function
func (wsam *WebSocketAuthMiddleware) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if wsam.authenticator != nil {
				user, err := wsam.authenticator(c.Request())
				if err != nil {
					return echo.NewHTTPError(http.StatusUnauthorized, "WebSocket authentication failed")
				}
				c.Set("user", user)
			}
			return next(c)
		}
	}
}
