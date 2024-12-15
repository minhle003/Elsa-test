package socket

import (
	"cloud.google.com/go/firestore"
	"context"
	socketio "github.com/googollee/go-socket.io"
	"github.com/minhle003/Elsa-test/internal/services/session_service"
	"github.com/minhle003/Elsa-test/pkg/logging"
)

func SetupSocketHandlers(server *socketio.Server, firestoreClient *firestore.Client, logger logging.Logger) {
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		logger.Info("Client connected: ", s.ID())
		return nil
	})

	server.OnEvent("/", "get_session", func(s socketio.Conn, msg map[string]string) {
		sessionID := msg["sessionId"]
		userID := msg["userId"]

		sessionService := session_service.NewSessionService(firestoreClient, logger, context.Background())
		session, err := sessionService.GetSession(sessionID, userID)
		if err != nil {
			s.Emit("error", err.Error())
			return
		}

		s.Emit("session_data", session)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		logger.Info("Client disconnected: ", s.ID(), " Reason: ", reason)
	})
}
