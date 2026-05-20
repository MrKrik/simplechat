package app

import (
	grpcapp "auth/iternal/app/grpc"
	"auth/iternal/config"
	"auth/iternal/services/auth"
	"auth/iternal/services/validate"
	"auth/iternal/storage/redis"
	"auth/iternal/storage/sqlite"
	"context"
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
	redisConfig config.RedisConfig,
) *App {
	storage, err := sqlite.New(storagePath)

	if err != nil {
		panic(err)
	}

	tokenStorage, err := redis.NewTokenStore(context.Background(), redisConfig)

	if err != nil {
		panic(err)
	}
	authService := auth.New(log, storage, storage, storage, tokenTTL, chatTokenTTL, tokenStorage)

	validateService := validate.New(log, storage, tokenStorage)

	grpcApp := grpcapp.New(log, authService, validateService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
