package main

import (
	"courseGolang/http/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Создаем роутер с дефолтными настройками (включает Logger и Recovery middleware)
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API работает!",
			"version": "1.0.0",
		})
	})
	// Группируем маршруты для API
	api := router.Group("/api")
	{
		api.GET("/users", handlers.GetUsers)
		api.GET("/users/:id", handlers.GetUser)
	}

	// Запускаем сервер на порту 8080
	log.Println("Сервер запущен на порту 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
