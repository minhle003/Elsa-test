package handlers

import (
	"github.com/minhle003/Elsa-test/internal/services/quiz_service"
	"github.com/minhle003/Elsa-test/pkg/logging"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func CreateQuiz(client *firestore.Client, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var quiz quiz_service.Quiz
		if err := c.ShouldBindJSON(&quiz); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"invalid request": err.Error()})
			return
		}

		quizService := quiz_service.NewQuizService(client, logger, c)
		newQuiz, err := quizService.CreateQuiz(quiz)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error create new quiz": err.Error()})
			return
		}
		c.JSON(http.StatusOK, quiz_service.Quiz{ID: newQuiz.ID})
	}
}

func UpdateQuiz(client *firestore.Client, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("User-ID")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id"})
			return
		}
		quizID := c.Param("quizID")
		if quizID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"invalid request": "Quiz ID is required"})
			return
		}

		var quiz quiz_service.Quiz
		if err := c.ShouldBindJSON(&quiz); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Ensure the quiz ID from the URL matches the quiz ID in the payload
		if quiz.ID != quizID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Quiz ID in URL and payload do not match"})
			return
		}

		quizService := quiz_service.NewQuizService(client, logger, c)
		_, err := quizService.UpdateQuiz(quiz, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}

func ListQuizzesByUser(client *firestore.Client, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("User-ID")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id"})
			return
		}

		quizService := quiz_service.NewQuizService(client, logger, c)

		quizzes, err := quizService.ListQuizzesByUser(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, quizzes)
	}
}
