package handler

import (
	"encoding/json"
	"errors"
	"github.com/Rastaiha/bermudia/internal/config"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/Rastaiha/bermudia/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
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

func (a *admin) SetTerritory(w http.ResponseWriter, r *http.Request) {
	var territory domain.Territory
	if err := json.NewDecoder(r.Body).Decode(&territory); err != nil {
		sendDecodeError(w)
		return
	}
	if err := a.adminService.SetTerritory(r.Context(), territory); err != nil {
		handleError(w, err)
		return
	}
	sendResult(w, territory)
}

func (a *admin) SetBookAndBindToIsland(w http.ResponseWriter, r *http.Request) {
	islandID := chi.URLParam(r, "islandID")
	var input service.BookInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendDecodeError(w)
		return
	}
	result, err := a.adminService.SetBookAndBindToIsland(r.Context(), islandID, input)
	if err != nil {
		handleError(w, err)
		return
	}
	sendResult(w, result)
}

func (a *admin) SetBookAndBindToPool(w http.ResponseWriter, r *http.Request) {
	poolID := chi.URLParam(r, "poolID")
	var input service.BookInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendDecodeError(w)
		return
	}
	result, err := a.adminService.SetBookAndBindToPool(r.Context(), poolID, input)
	if err != nil {
		handleError(w, err)
		return
	}
	sendResult(w, result)
}

func (a *admin) GetTerritoryIslandBindings(w http.ResponseWriter, r *http.Request) {
	territoryID := chi.URLParam(r, "territoryID")
	result, err := a.adminService.GetTerritoryIslandBindings(r.Context(), territoryID)
	if err != nil {
		handleError(w, err)
		return
	}
	sendResult(w, result)
}

func (a *admin) SetTerritoryIslandBindings(w http.ResponseWriter, r *http.Request) {
	var bindings service.TerritoryIslandBindings
	if err := json.NewDecoder(r.Body).Decode(&bindings); err != nil {
		sendDecodeError(w)
		return
	}
	result, err := a.adminService.SetTerritoryIslandBindings(r.Context(), bindings)
	if err != nil {
		handleError(w, err)
		return
	}
	sendResult(w, result)
}

func (a *admin) CreateUser(w http.ResponseWriter, r *http.Request) {
	indexStr := r.URL.Query().Get("index")
	index := 0
	if indexStr != "" {
		var err error
		index, err = strconv.Atoi(indexStr)
		if err != nil {
			sendError(w, http.StatusBadRequest, "invalid index query parameter")
			return
		}
	}
	var user service.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendDecodeError(w)
		return
	}
	result, err := a.adminService.CreateUser(r.Context(), index, user)
	if err != nil {
		handleError(w, err)
		return
	}
	sendResult(w, result)
}
