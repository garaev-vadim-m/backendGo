package main

import (
	"courseGolang/http/config"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfig()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API работает!",
			"version": "1.0.0",
		})
	})

	api := router.Group("/api")
	{
		api.GET("/users", handlers.GetAllUsers)
		api.GET("/users/:id", handlers.GetUserByID)
	}

	log.Printf("Сервер запущен на порту %s", cfg.Port)
	if err := router.Run(cfg.Port); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
