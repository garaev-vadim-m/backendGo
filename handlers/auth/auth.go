package handlers

import (
	"encoding/json"
	"go-users-api/db"
	"go-users-api/middleware"
	"go-users-api/models"
	"go-users-api/utils"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(response, "No token", 400)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	middleware.AddToBlacklist(token)

	response.Write([]byte(`{"message":"logged out"}`))
}
