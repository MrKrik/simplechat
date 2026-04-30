package app

import (
	grpcapp "auth/iternal/app/grpc"
	"auth/iternal/services/auth"
	"auth/iternal/services/validate"
	"auth/iternal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
	chatTokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	authService := auth.New(log, storage, storage, storage, tokenTTL, chatTokenTTL)

	validateService := validate.New(log, storage)

	grpcApp := grpcapp.New(log, authService, validateService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
