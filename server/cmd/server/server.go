package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/minhle003/Elsa-test/internal/firestore"
	"github.com/minhle003/Elsa-test/internal/handlers"
	"github.com/minhle003/Elsa-test/internal/socket"
	"github.com/minhle003/Elsa-test/pkg/logging"
	"net/http"
	"time"
)

func main() {
	logger := logging.LogrusLogger()

	firestoreClient, err := firestore.NewClient()
	if err != nil {
		logger.Critical("Failed to create Firestore client: %v", err)
		return
	}
	defer firestoreClient.Close()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowWebSockets:        true,
		AllowBrowserExtensions: true,
		AllowOrigins:           []string{"*"}, // allow all origins for demo
		AllowMethods:           []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:           []string{"Origin", "Content-Type", "Accept", "User-ID"},
		ExposeHeaders:          []string{"Content-Length"},
		AllowCredentials:       true,
		MaxAge:                 12 * time.Hour,
	}))

	go socket.RunHub()

	router.GET("/ws", socket.SetupWebSocketHandlers(firestoreClient, logger)) // WebSocket endpoint

	router.POST("/api/user", handlers.CreateUser(firestoreClient, logger))
	router.GET("/api/user/:username", handlers.GetUserByUserName(firestoreClient, logger))

	router.POST("/api/quiz", handlers.CreateQuiz(firestoreClient, logger))
	router.PUT("/api/quiz/:quizID", handlers.UpdateQuiz(firestoreClient, logger))
	router.GET("/api/quiz/quizzes", handlers.ListQuizzesByUser(firestoreClient, logger))

	router.POST("/api/session", handlers.CreateSession(firestoreClient, logger))
	router.GET("/api/session/:sessionId", handlers.GetSession(firestoreClient, logger))
	router.PATCH("/api/session/start", handlers.StartSession(firestoreClient, &socket.Hub, logger))
	router.PATCH("/api/session/join", handlers.JoinSession(firestoreClient, &socket.Hub, logger))
	router.PATCH("/api/session/change_question", handlers.ChangeQuestion(firestoreClient, &socket.Hub, logger))
	router.PATCH("/api/session/participant/update_score", handlers.UpdateScore(firestoreClient, &socket.Hub, logger))
	router.PATCH("/api/session/end", handlers.EndSession(firestoreClient, &socket.Hub, logger))

	logger.Info("Server is running on port 8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Critical("HTTP server error: %v", err)
	}

}
