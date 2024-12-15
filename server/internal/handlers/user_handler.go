package handlers

import (
	"github.com/minhle003/Elsa-test/internal/services/user_service"
	"github.com/minhle003/Elsa-test/pkg/logging"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func CreateUser(client *firestore.Client, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user user_service.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userService := user_service.NewUserService(client, logger, c)
		newUser, err := userService.CreateUser(user.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, newUser)
	}
}

func GetUserByUserName(client *firestore.Client, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
			return
		}
		userService := user_service.NewUserService(client, logger, c)
		user, err := userService.GetUserByUserName(username)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
