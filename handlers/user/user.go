package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"go-users-api/db"
	"go-users-api/models"

	"github.com/go-chi/chi/v5"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var tokenBlacklist = make(map[string]bool)

func GetUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	rows, err := db.DB.Query("SELECT * FROM users")
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Login,
			&user.Email,
			&user.Age,
			&user.Country,
			&user.Password,
		)

		if err != nil {
			http.Error(response, err.Error(), 500)
			return
		}

		users = append(users, user)
	}

	json.NewEncoder(response).Encode(users)
}

func GetUserByID(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	paramId := chi.URLParam(request, "id")

	id, err := strconv.Atoi(paramId)
	if err != nil {
		http.Error(response, "Invalid ID", 400)
		return
	}

	row := db.DB.QueryRow("SELECT * FROM users WHERE id = ?", id)

	var user models.User
	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.Login,
		&user.Email,
		&user.Age,
		&user.Country,
		&user.Password,
	)

	if err == sql.ErrNoRows {
		http.Error(response, "User not found", 404)
		return
	}

	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	json.NewEncoder(response).Encode(user)
}
