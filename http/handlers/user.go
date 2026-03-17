package handlers

import (
	"courseGolang/http/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	c.JSON(
		http.StatusOK, gin.H{
			"status":  "success",
			"data":    models.Users,
			"message": "Список пользователей получен",
		},
	)
}

func GetUser(c *gin.Context) {
	var userId = c.Param("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Неверный идентификатор пользователя",
		})
		return
	}
	for _, user := range models.Users {
		if user.ID == id {
			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"data":    user,
				"message": "Пользователь получен",
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"status":  "error",
		"message": "Пользователь не найден",
	})
}
