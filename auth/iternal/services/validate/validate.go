package validate

import (
	"auth/iternal/domain/models"
	"auth/iternal/lib/jwt"
	"context"
	"fmt"
	"log/slog"
)

type Validate struct {
	log           *slog.Logger
	appProvider   AppProvider
	tokenProvider TokenProvider
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

type TokenProvider interface {
	TokenExists(ctx context.Context, key string) (bool, error)
}

func New(log *slog.Logger,
	appProvider AppProvider,
	tokenProvider TokenProvider,
) *Validate {
	return &Validate{
		log:           log,
		appProvider:   appProvider,
		tokenProvider: tokenProvider,
	}
}

func (v *Validate) ValidateToken(ctx context.Context, token string) (ok bool, err error) {

	const op = "validate.ValidateToken"

	log := v.log.With(
		slog.String("op", op),
		slog.String("chatToken", token),
	)

	app, err := v.appProvider.App(ctx, 1)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	secret := app.Secret

	log.Info("validate chat token")

	ok, err = jwt.ValidateToken(token, secret)
	if !ok {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	ok, err = v.tokenProvider.TokenExists(ctx, token)
	if !ok {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("chat token valid")

	return true, nil
}
