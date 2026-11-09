package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sankevichsima-lab/Currency-Tracker/internal/auth"
)

func AuthMiddleware(next http.HandlerFunc, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// 1. Сначала пробуем достать токен из Cookie (для браузера)
		cookie, err := r.Cookie("auth_token")
		if err == nil && cookie.Value != "" {
			tokenString = cookie.Value
		} else {
			// 2. Если куки нет, проверяем заголовок Authorization (для curl / мобильных клиентов)
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		// Если токена нет ни в куках, ни в заголовках
		if tokenString == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "missing or invalid token",
			})
			return
		}

		// Валидация токена
		claims, err := auth.ValidateToken(tokenString, secret)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid or expired token",
			})
			return
		}

		// Передаем управление дальше, если токен валиден
		_ = claims
		next(w, r)
	}
}
