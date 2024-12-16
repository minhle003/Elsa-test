package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/minhle003/Elsa-test/internal/services/session_service"
	"github.com/minhle003/Elsa-test/internal/socket"
	"github.com/minhle003/Elsa-test/pkg/logging"
	"net/http"

	"cloud.google.com/go/firestore"
)

func CreateSession(client *firestore.Client, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("User-ID")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id"})
			return
		}
		var request struct {
			QuizId string `json:"quizId" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sessionService := session_service.NewSessionService(client, logger, c)
		newSession, err := sessionService.CreateSession(request.QuizId, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, newSession)
	}
}

func GetSession(client *firestore.Client, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("sessionId")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing session id"})
			return
		}
		userID := c.GetHeader("User-ID")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id"})
			return
		}

		sessionService := session_service.NewSessionService(client, logger, c)
		session, err := sessionService.GetSession(sessionID, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, session)
	}
}

func JoinSession(client *firestore.Client, hub *socket.WebSocketHub, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			SessionId string `json:"sessionId" binding:"required"`
			Name      string `json:"name" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sessionService := session_service.NewSessionService(client, logger, c)
		session, participantId, err := sessionService.JoinSession(request.Name, request.SessionId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hub.Broadcast <- session
		c.JSON(http.StatusOK, gin.H{"session": session, "participantId": participantId})
	}
}

func StartSession(client *firestore.Client, hub *socket.WebSocketHub, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("User-ID")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id"})
			return
		}
		var request struct {
			SessionId string `json:"sessionId" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sessionService := session_service.NewSessionService(client, logger, c)
		session, err := sessionService.StartSession(request.SessionId, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hub.Broadcast <- session
		c.JSON(http.StatusOK, session)
	}
}

func ChangeQuestion(client *firestore.Client, hub *socket.WebSocketHub, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("User-ID")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id"})
			return
		}
		var request struct {
			SessionId     string `json:"sessionId" binding:"required"`
			QuestionIndex int    `json:"questionIndex" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sessionService := session_service.NewSessionService(client, logger, c)
		session, err := sessionService.ChangeQuestion(request.SessionId, request.QuestionIndex, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hub.Broadcast <- session
		c.JSON(http.StatusOK, session)
	}
}

func UpdateScore(client *firestore.Client, hub *socket.WebSocketHub, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			SessionId     string `json:"sessionId" binding:"required"`
			ParticipantId string `json:"participantId" binding:"required"`
			Score         int    `json:"score" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sessionService := session_service.NewSessionService(client, logger, c)
		session, err := sessionService.UpdateScore(request.SessionId, request.ParticipantId, request.Score)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hub.Broadcast <- session
		c.JSON(http.StatusOK, session)
	}
}

func EndSession(client *firestore.Client, hub *socket.WebSocketHub, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("User-ID")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id"})
			return
		}
		var request struct {
			SessionId string `json:"sessionId" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sessionService := session_service.NewSessionService(client, logger, c)
		session, err := sessionService.EndSession(request.SessionId, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hub.Broadcast <- session
		c.JSON(http.StatusOK, session)
	}
}
