package redis

import (
	"auth/iternal/config"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenStore struct {
	Client *redis.Client
}

func NewTokenStore(ctx context.Context, cfg config.RedisConfig) (*TokenStore, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.User,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Timeout) * time.Second,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		fmt.Printf("failed to connect to redis server: %s\n", err.Error())
		return nil, err
	}

	return &TokenStore{Client: db}, nil
}

func (s *TokenStore) Stop() error {
	return s.Client.Close()
}

func (s *TokenStore) SaveToken(ctx context.Context, key string, value string, ttl time.Duration) error {
	return s.Client.Set(ctx, key, value, ttl).Err()
}

func (s *TokenStore) TokenExists(ctx context.Context, key string) (bool, error) {
	exists, err := s.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (s *TokenStore) DeleteToken(ctx context.Context, key string) error {
	return s.Client.Del(ctx, key).Err()
}
