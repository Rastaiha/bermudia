package handler

import (
	"encoding/json"
	"errors"
	"github.com/Rastaiha/bermudia/internal/config"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/Rastaiha/bermudia/internal/service"
	"net/http"
	"strings"
)

type admin struct {
	cfg          config.Config
	adminService *service.Admin
}

type adminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *admin) Login(w http.ResponseWriter, r *http.Request) {
	var req adminLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendDecodeError(w)
		return
	}
	token, err := a.adminService.Login(r.Context(), req.Username, req.Password)
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

func (a *admin) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if tokenStr == "" {
			sendError(w, http.StatusUnauthorized, "Missing auth token")
			return
		}
		ok := a.adminService.ValidateToken(r.Context(), tokenStr)
		if !ok {
			sendError(w, http.StatusUnauthorized, "Invalid auth token")
			return
		}

		next.ServeHTTP(w, r)
	})
}
