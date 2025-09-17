package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Rastaiha/bermudia/internal/domain"
	"net/http"
	"strings"
)

const (
	userContextKey = "user"
)

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if tokenStr == "" {
			sendError(w, http.StatusUnauthorized, "Missing auth token")
			return
		}
		user, ok := h.authService.ValidateToken(r.Context(), tokenStr)
		if !ok {
			sendError(w, http.StatusUnauthorized, "Invalid auth token")
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), userContextKey, user))

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) pauseCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isPaused, err := h.authService.IsGamePaused(r.Context())
		if err != nil {
			handleError(w, err)
			return
		}
		if isPaused {
			sendError(w, http.StatusLocked, "بازی در حال حاضر متوقف شده است.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendDecodeError(w)
		return
	}
	token, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			sendError(w, http.StatusNotFound, "نام کاربری یا کلمه عبور اشتباه است")
			return
		}
		handleError(w, err)
		return
	}
	sendResult(w, map[string]any{
		"token": token,
	})
}

func getUser(ctx context.Context) (*domain.User, error) {
	v := ctx.Value(userContextKey)
	user, ok := v.(*domain.User)
	if !ok {
		return nil, errors.New("expected user was not found in context")
	}
	return user, nil
}
