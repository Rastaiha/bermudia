package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/config"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/Rastaiha/bermudia/internal/service"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strings"
)

type admin struct {
	cfg           config.Config
	adminService  *service.Admin
	playerService *service.Player
	gameState     domain.GameStateStore
	actives       func() map[string]int
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
		a.handleAdminError(w, err)
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

func (a *admin) GetTerritories(w http.ResponseWriter, r *http.Request) {
	result, err := a.adminService.GetTerritories(r.Context())
	if err != nil {
		a.handleAdminError(w, err)
		return
	}
	sendResult(w, result)
}

func (a *admin) SetTerritory(w http.ResponseWriter, r *http.Request) {
	var territory domain.Territory
	if err := json.NewDecoder(r.Body).Decode(&territory); err != nil {
		sendDecodeError(w)
		return
	}
	if err := a.adminService.SetTerritory(r.Context(), territory); err != nil {
		a.handleAdminError(w, err)
		return
	}
	sendResult(w, territory)
}

func (a *admin) GetBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookID")
	result, err := a.adminService.GetBook(r.Context(), bookID)
	if err != nil {
		a.handleAdminError(w, err)
		return
	}
	sendResult(w, result)
}

func (a *admin) GetIslandHeader(w http.ResponseWriter, r *http.Request) {
	islandID := chi.URLParam(r, "islandID")
	result, err := a.adminService.GetIslandHeader(r.Context(), islandID)
	if err != nil {
		a.handleAdminError(w, err)
		return
	}
	sendResult(w, result)
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
		a.handleAdminError(w, err)
		return
	}
	sendResult(w, result)
}

func (a *admin) GetPools(w http.ResponseWriter, r *http.Request) {
	result, err := a.adminService.GetPools(r.Context())
	if err != nil {
		a.handleAdminError(w, err)
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
		a.handleAdminError(w, err)
		return
	}
	sendResult(w, result)
}

func (a *admin) GetTerritoryIslandBindings(w http.ResponseWriter, r *http.Request) {
	territoryID := chi.URLParam(r, "territoryID")
	result, err := a.adminService.GetTerritoryIslandBindings(r.Context(), territoryID)
	if err != nil {
		a.handleAdminError(w, err)
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
		a.handleAdminError(w, err)
		return
	}
	sendResult(w, result)
}

func (a *admin) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user service.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendDecodeError(w)
		return
	}
	result, err := a.adminService.CreateUser(r.Context(), user)
	if err != nil {
		a.handleAdminError(w, err)
		return
	}
	sendResult(w, result)
}

func (a *admin) handleAdminError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		sendError(w, http.StatusRequestTimeout, "Request cancelled")
		return
	}
	var adminError service.AdminError
	if errors.As(err, &adminError) {
		sendError(w, http.StatusBadRequest, adminError.Error())
		return
	}
	var domainError domain.Error
	if errors.As(err, &domainError) {
		switch domainError.Reason() {
		case domain.ErrorReasonResourceNotFound:
			sendError(w, http.StatusNotFound, domainError.Error())
		case domain.ErrorReasonRuleViolation:
			sendError(w, http.StatusConflict, domainError.Error())
		default:
			sendError(w, http.StatusInternalServerError, domainError.Error())
		}
		return
	}
	switch {
	case errors.Is(err, domain.ErrTerritoryNotFound):
		sendError(w, http.StatusNotFound, err.Error())
		return
	case errors.Is(err, domain.ErrIslandNotFound):
		sendError(w, http.StatusNotFound, err.Error())
		return
	case errors.Is(err, domain.ErrBookNotFound):
		sendError(w, http.StatusNotFound, err.Error())
		return
	case errors.Is(err, domain.ErrPoolSettingsNotFound):
		sendError(w, http.StatusNotFound, err.Error())
		return
	}

	errText := "<nil>"
	if err != nil {
		errText = err.Error()
	}
	slog.Error("internal admin error", slog.String("error", errText))
	sendError(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error: %s", errText))
}

func (a *admin) GetUsers(w http.ResponseWriter, r *http.Request) {
	result, err := a.adminService.GetUsers(r.Context())
	if err != nil {
		a.handleAdminError(w, err)
		return
	}
	sendResult(w, result)
}

func (a *admin) GetGameState(w http.ResponseWriter, r *http.Request) {
	isPaused, err := a.gameState.GetIsPaused(r.Context())
	if err != nil {
		a.handleAdminError(w, err)
		return
	}
	sendResult(w, map[string]any{"isPaused": isPaused})
}

type setGameStateRequest struct {
	IsPaused bool `json:"isPaused"`
}

func (a *admin) SetGameState(w http.ResponseWriter, r *http.Request) {
	var req setGameStateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendDecodeError(w)
		return
	}
	if err := a.gameState.SetIsPaused(r.Context(), req.IsPaused); err != nil {
		a.handleAdminError(w, err)
		return
	}
	sendResult(w, map[string]any{"isPaused": req.IsPaused})
}

type broadcastRequest struct {
	Message string `json:"message"`
}

func (a *admin) Broadcast(w http.ResponseWriter, r *http.Request) {
	var req broadcastRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendDecodeError(w)
		return
	}
	if strings.TrimSpace(req.Message) == "" {
		sendError(w, http.StatusBadRequest, "message is required")
		return
	}
	count, err := a.playerService.BroadcastMessage(r.Context(), req.Message)
	if err != nil {
		a.handleAdminError(w, err)
		return
	}
	sendResult(w, map[string]any{"sentTo": count})
}

func (a *admin) GetConnections(w http.ResponseWriter, r *http.Request) {
	sendResult(w, a.actives())
}
