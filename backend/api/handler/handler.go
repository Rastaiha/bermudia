package handler

import (
	"context"
	"errors"
	"github.com/Rastaiha/bermudia/api/hub"
	"github.com/Rastaiha/bermudia/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"log"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"time"
)

type Handler struct {
	server           *http.Server
	wsUpgrader       websocket.Upgrader
	authService      *service.Auth
	territoryService *service.Territory
	islandService    *service.Island
	playerService    *service.Player
	playerHub        *hub.Hub
	tradeHub         *hub.Hub
	inboxHub         *hub.Hub
}

func New(authService *service.Auth, territoryService *service.Territory, islandService *service.Island, playerService *service.Player) *Handler {
	return &Handler{
		authService:      authService,
		territoryService: territoryService,
		islandService:    islandService,
		playerService:    playerService,
		playerHub:        hub.NewHub(),
		tradeHub:         hub.NewHub(),
		inboxHub:         hub.NewHub(),
	}
}

func (h *Handler) Start() {
	r := chi.NewRouter()

	r.Use(logger())
	r.Use(corsMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/territories/{territoryID}", h.GetTerritory) // TODO: make it authenticated
		r.Post("/login", h.Login)

		// ws endpoints
		r.HandleFunc("/events", h.StreamPlayerEvents)
		r.HandleFunc("/trade/events", h.StreamTradeEvents)
		r.HandleFunc("/inbox/events", h.StreamInboxEvents)

		// Authenticated endpoints
		r.Group(func(r chi.Router) {
			r.Use(h.authMiddleware)
			r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
				user, err := getUser(r.Context())
				if err != nil {
					handleError(w, err)
					return
				}
				sendResult(w, user)
			})
			r.Get("/islands/{islandID}", h.GetIsland)
			r.Post("/answer/{inputID}", h.SubmitAnswer)
			r.Get("/player", h.GetPlayer)
			r.Post("/travel_check", h.TravelCheck)
			r.Post("/travel", h.Travel)
			r.Post("/refuel_check", h.RefuelCheck)
			r.Post("/refuel", h.Refuel)
			r.Post("/anchor_check", h.AnchorCheck)
			r.Post("/anchor", h.Anchor)
			r.Post("/migrate_check", h.MigrateCheck)
			r.Post("/migrate", h.Migrate)
			r.Post("/unlock_treasure_check", h.UnlockTreasureCheck)
			r.Post("/unlock_treasure", h.UnlockTreasure)
			r.Route("/trade", func(r chi.Router) {
				r.Post("/make_offer_check", h.MakeOfferCheck)
				r.Post("/make_offer", h.MakeOffer)
				r.Post("/accept_offer", h.AcceptOffer)
				r.Post("/delete_offer", h.DeleteOffer)
				r.Get("/offers", h.GetTradeOffers)
			})
			r.Get("/inbox/messages", h.GetInboxMessages)
		})
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	h.playerService.OnPlayerUpdate(h.HandlePlayerUpdateEvent)
	h.playerService.OnTradeEventBroadcast(h.HandleTradeEventBroadcast)
	h.playerService.OnInboxEvent(h.HandleInboxEvent)

	slog.Info("Server starting")
	h.server = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	h.wsUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	go func() {
		err := h.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Error starting server:", err)
		}
	}()
}

func (h *Handler) Stop() {
	if err := h.server.Shutdown(context.Background()); err != nil {
		slog.Error("Error stopping server:", err)
	}
}

// Simple CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

var tokenPattern = regexp.MustCompile(`token?=(\S+) `)

func logger() func(http.Handler) http.Handler {
	return middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger: slog.NewLogLogger(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if a.Key != slog.MessageKey {
					return a
				}
				msg := []byte(a.Value.String())
				submatch := tokenPattern.FindSubmatchIndex(msg)
				if len(submatch) < 4 {
					return a
				}
				tokenRemoved := append(make([]byte, 0, len(msg)), msg[:submatch[2]]...)
				tokenRemoved = append(tokenRemoved, []byte("***")...)
				tokenRemoved = append(tokenRemoved, msg[submatch[3]:]...)
				return slog.String(a.Key, string(tokenRemoved))
			},
		}), slog.LevelInfo),
	})
}
