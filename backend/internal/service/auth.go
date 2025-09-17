package service

import (
	"context"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/config"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"time"
)

type Auth struct {
	cfg            config.Config
	userStore      domain.UserStore
	gameStateStore domain.GameStateStore
}

func NewAuth(cfg config.Config, userStore domain.UserStore, gameStateStore domain.GameStateStore) *Auth {
	return &Auth{cfg: cfg, userStore: userStore, gameStateStore: gameStateStore}
}

func (a *Auth) Login(ctx context.Context, username, password string) (string, error) {
	user, err := a.userStore.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if !domain.ValidatePassword(password, user.HashedPassword) {
		return "", domain.ErrUserNotFound
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id": fmt.Sprint(user.ID),
		"iat":     float64(time.Now().UTC().UnixNano()) / 1e9,
	})
	tokenString, err := token.SignedString(a.cfg.TokenSigningKeyBytes())
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

func (a *Auth) ValidateToken(ctx context.Context, tokenStr string) (*domain.User, bool) {
	token, err := jwt.Parse(
		tokenStr,
		func(token *jwt.Token) (interface{}, error) {
			return a.cfg.TokenSigningKeyBytes(), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return nil, false
	}
	if !token.Valid {
		return nil, false
	}
	if iat, err := token.Claims.GetIssuedAt(); err != nil || time.Since(iat.Time) > 16*time.Hour {
		return nil, false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, false
	}
	v, ok := claims["user_id"]
	if !ok {
		return nil, false
	}
	userIdStr, _ := v.(string)
	var userId int32
	_, err = fmt.Sscanf(userIdStr, "%d", &userId)
	if err != nil {
		return nil, false
	}
	user, err := a.userStore.Get(ctx, userId)
	if err != nil {
		slog.Error("failed to find valid user by id", "user_id", userId)
		return nil, false
	}
	return user, true
}

func (a *Auth) IsGamePaused(ctx context.Context) (bool, error) {
	return a.gameStateStore.GetIsPaused(ctx)
}
