package controller

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type ctxKey string

const (
	authorizationHeader        = "Authorization"
	userCtx             ctxKey = "userId"
)

func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)
		if header == "" {
			http.Error(w, "empty auth header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		if len(headerParts[1]) == 0 {
			http.Error(w, "token is empty", http.StatusUnauthorized)
			return
		}

		// Валидация и получение userId по токену
		userId, err := h.service.Authorization.ParseToken(headerParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Добавление userId в контекст запроса
		ctx := context.WithValue(r.Context(), userCtx, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Получение userId из контекста
func getUserId(r *http.Request) (int, error) {
	id, ok := r.Context().Value(userCtx).(int)
	if !ok {
		return 0, errors.New("user id not found or invalid type")
	}

	return id, nil
}
