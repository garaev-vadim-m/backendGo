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

var tokenBlacklist = make(map[string]bool)

func GetUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	rows, err := db.DB.Query("SELECT u.id, u.name, u.login, u.email, u.age, u.country, u.password, r.id as role_id, r.name as role_name FROM users u JOIN roles r ON u.role_id = r.id;")
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
			&user.Password, // обязательно!
			&user.Role.ID,
			&user.Role.Name,
		)

		if err != nil {
			http.Error(response, err.Error(), 500)
			return
		}

		users = append(users, user)
	}

	json.NewEncoder(response).Encode(users)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	paramId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	row := db.DB.QueryRow(`
		SELECT u.id, u.name, u.login, u.email, u.age, u.country, u.password, r.id, r.name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.id = ?
	`, id)

	var user models.User
	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.Login,
		&user.Email,
		&user.Age,
		&user.Country,
		&user.Password,
		&user.Role.ID,
		&user.Role.Name,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", 404)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(user)
}
