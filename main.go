package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"go-users-api/db"
	"go-users-api/handlers"

	"github.com/go-chi/cors"
)

func main() {
	db.InitDB()

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // потом сузим
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
	}))
	r.Route("/api", func(r chi.Router) {
		r.Get("/users", handlers.GetUsers)
		r.Get("/users/{id}", handlers.GetUserByID)
		r.Post("/login", handlers.Login)
	})

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
