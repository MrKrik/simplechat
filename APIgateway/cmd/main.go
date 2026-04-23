package main

import (
	"log/slog"
	"os"
	"tcp-server/iternal/config"
	grpc "tcp-server/iternal/grpc/auth"
	auth "tcp-server/iternal/htpp-server/handlers"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := SetupLogger(cfg.Env)

	log.Info("Starting aplication", slog.Any("config", cfg))

	router := chi.NewRouter()

	authclient, err := grpc.New(cfg.AuthGRPC.Address, time.Duration(60))
	if err != nil {
		log.Error("Failed connect auth server")
	}

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Post("/login", auth.Login(authclient))
	router.Post("/register", auth.Register(authclient))

	// router.Use(mwLogerr.New(log))

}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
