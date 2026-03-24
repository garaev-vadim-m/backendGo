package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"go-users-api/db"
	"go-users-api/models"
	"go-users-api/utils"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
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

func Login(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var req LoginRequest

	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(response, "Invalid body", 400)
		return
	}

	row := db.DB.QueryRow("SELECT * FROM users WHERE email = ?", req.Email)

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

	if err != nil {
		http.Error(response, "Invalid credentials", 401)
		return
	}

	// 🔒 проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(response, "Invalid credentials", 401)
		return
	}

	// пока просто ответ
	user.Password = ""

	// генерируем JWT
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		http.Error(response, "Failed to generate token", 500)
		return
	}

	// возвращаем токен вместо user
	json.NewEncoder(response).Encode(map[string]string{
		"token": token,
	})
}

func Logout(response http.ResponseWriter, r *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")

	if token == "" {
		http.Error(response, "No token", 400)
		return
	}

	tokenBlacklist[token] = true

	response.Write([]byte(`{"message":"logged out"}`))
}
