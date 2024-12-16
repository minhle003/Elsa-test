package socket

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/minhle003/Elsa-test/internal/services/session_service"
	"github.com/minhle003/Elsa-test/pkg/logging"
	"log"
	"net/http"
	"sync"
)

type WebSocketHub struct {
	Clients    map[*websocket.Conn]bool
	Broadcast  chan session_service.Session
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	Mutex      sync.RWMutex
}

var Hub = WebSocketHub{
	Broadcast:  make(chan session_service.Session),
	Register:   make(chan *websocket.Conn),
	Unregister: make(chan *websocket.Conn),
	Clients:    make(map[*websocket.Conn]bool),
}

type ClientMessage struct {
	SessionId string `json:"sessionId"`
	UserId    string `json:"userId"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { // Allow all origins for development (adjust for production)
		return true
	},
}

func SetupWebSocketHandlers(firestoreClient *firestore.Client, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket:", err)
			return
		}
		defer conn.Close()

		Hub.Register <- conn

		defer func() {
			Hub.Unregister <- conn
			conn.Close()
		}()

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}

			if messageType == websocket.TextMessage {
				var clientMsg ClientMessage
				if err := json.Unmarshal(message, &clientMsg); err != nil {
					log.Println("Error unmarshalling client message:", err)
					continue
				}

				sessionService := session_service.NewSessionService(firestoreClient, logger, context.Background())
				session, err := sessionService.GetSession(clientMsg.SessionId, clientMsg.UserId)

				response, err := json.Marshal(session)
				if err != nil {
					log.Println("Error marshalling server message:", err)
					continue
				}

				if err := conn.WriteMessage(websocket.TextMessage, response); err != nil {
					log.Println("Error writing message:", err)
					break
				}
			}
		}
	}
}

func RunHub() {
	for {
		select {
		case client := <-Hub.Register:
			Hub.Mutex.Lock()
			Hub.Clients[client] = true
			Hub.Mutex.Unlock()

		case client := <-Hub.Unregister:
			Hub.Mutex.Lock()
			if _, ok := Hub.Clients[client]; ok {
				delete(Hub.Clients, client)
			}
			Hub.Mutex.Unlock()

		case message := <-Hub.Broadcast:
			Hub.Mutex.RLock()
			for client := range Hub.Clients {
				err := client.WriteJSON(message)
				if err != nil {
					Hub.Unregister <- client
				}
			}
			Hub.Mutex.RUnlock()
		}
	}
}
