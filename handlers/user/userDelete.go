package handlers

import (
	"encoding/json"
	"go-users-api/db"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// 👉 Говорим клиенту, что ответ будет в формате JSON
	w.Header().Set("Content-Type", "application/json")

	// 👉 Достаём параметр id из URL
	// Например: /users/5 → "5"
	paramId := chi.URLParam(r, "id")

	// 👉 Преобразуем строку "5" → число 5
	id, err := strconv.Atoi(paramId)
	if err != nil {
		// 👉 Если не получилось (например "abc"), возвращаем ошибку 400
		http.Error(w, "Invalid ID", 400)
		return
	}

	// 👉 Проверяем, существует ли пользователь
	// QueryRow вернёт одну строку
	var exists int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", id).Scan(&exists)
	if err != nil {
		// 👉 Ошибка при работе с БД
		http.Error(w, err.Error(), 500)
		return
	}

	if exists == 0 {
		// 👉 Если пользователя нет — 404
		http.Error(w, "User not found", 404)
		return
	}

	// 👉 Удаляем пользователя из БД
	result, err := db.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		// 👉 Если SQL сломался
		http.Error(w, err.Error(), 500)
		return
	}

	// 👉 Получаем количество удалённых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 👉 Дополнительная защита (на всякий случай)
	if rowsAffected == 0 {
		http.Error(w, "User not found", 404)
		return
	}

	// 👉 Отправляем успешный ответ
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deleted",
	})
}
