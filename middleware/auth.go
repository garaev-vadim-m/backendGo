package middleware

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("secret_key")

// Blacklist
var (
	tokenBlacklist = make(map[string]bool)
	blacklistMu    sync.RWMutex
)

func AddToBlacklist(token string) {
	blacklistMu.Lock()
	defer blacklistMu.Unlock()
	tokenBlacklist[token] = true
}

func IsBlacklisted(token string) bool {
	blacklistMu.RLock()
	defer blacklistMu.RUnlock()
	return tokenBlacklist[token]
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", 401)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format", 401)
			return
		}

		tokenStr := parts[1]

		// ✅ Проверка blacklist
		if IsBlacklisted(tokenStr) {
			http.Error(w, "Token revoked", 401)
			return
		}

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			// ✅ Проверка алгоритма — защита от alg:none атаки
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", 401)
			return
		}

		// ✅ Передаём userID в контекст, чтобы использовать в хендлерах
		if userID, ok := (*claims)["user_id"]; ok {
			role, _ := (*claims)["role"].(string)
			ctx := context.WithValue(r.Context(), "userID", userID)
			ctx = context.WithValue(ctx, "role", role)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
