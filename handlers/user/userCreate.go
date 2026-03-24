package handlers

import (
	"encoding/json"
	"go-users-api/db"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      int    `json:"age"`
	Country  string `json:"country"`
}

// func CreateUser(w http.ResponseWriter, r *http.Request)
// w — сюда ты пишешь ответ (response)
// r — входящий запрос (request)
func CreateUser(w http.ResponseWriter, r *http.Request) {
	//Говорим клиенту: “я верну JSON”
	w.Header().Set("Content-Type", "application/json")
	//Создаём структуру под входящие данные (struct)
	var req CreateUserRequest
	//Читаем body запроса
	// Пример:
	// 	{
	//   "email": "test@test.ru",
	//   "password": "1234"
	// }
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid body", 400)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password required", 400)
		return
	}
	//Проверка на существование пользователя
	var exists int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", req.Email).Scan(&exists)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if exists > 0 {
		http.Error(w, "User already exists", 400)
		return
	}
	//Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", 500)
		return
	}
	//Вставка в БД
	result, err := db.DB.Exec(`
		INSERT INTO users (name, login, email, age, country, password)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		req.Name,
		req.Login,
		req.Email,
		req.Age,
		req.Country,
		string(hashedPassword),
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	id, _ := result.LastInsertId()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    id,
		"email": req.Email,
	})
}
