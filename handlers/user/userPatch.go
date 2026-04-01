package handlers

import (
	"encoding/json"
	"go-users-api/db"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UpdateUserRequest struct {
	Name    *string `json:"name"`
	Login   *string `json:"login"`
	Email   *string `json:"email"`
	Age     *int    `json:"age"`
	Country *string `json:"country"`
}

func UpdateUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	paramId := chi.URLParam(request, "id")

	id, err := strconv.Atoi(paramId)
	if err != nil {
		http.Error(response, "Invalid user ID", http.StatusBadRequest)
		return
	}
	var req UpdateUserRequest

	err = json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(response, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := "UPDATE users SET "
	args := []interface{}{}

	// если пришло имя
	if req.Name != nil {
		query += "name = ?,"
		args = append(args, *req.Name)
	}

	// если пришёл login
	if req.Login != nil {
		query += "login = ?,"
		args = append(args, *req.Login)
	}

	// если пришёл email
	if req.Email != nil {
		query += "email = ?,"
		args = append(args, *req.Email)
	}

	// если пришёл возраст
	if req.Age != nil {
		query += "age = ?,"
		args = append(args, *req.Age)
	}

	// если пришла страна
	if req.Country != nil {
		query += "country = ?,"
		args = append(args, *req.Country)
	}
	// убираем последнюю запятую
	query = query[:len(query)-1]

	// добавляем WHERE
	query += " WHERE id = ?"

	// добавляем id в аргументы
	args = append(args, id)

	result, err := db.DB.Exec(query, args...)

	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	// проверяем сколько строк обновилось
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	// если 0 — пользователь не найден
	if rowsAffected == 0 {
		http.Error(response, "User not found", 404)
		return
	}

	// отправляем ответ
	json.NewEncoder(response).Encode(map[string]string{
		"message": "User updated",
	})
}
