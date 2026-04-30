package validate

import (
	"auth/iternal/domain/models"
	"auth/iternal/lib/jwt"
	"context"
	"fmt"
	"log/slog"
)

type Validate struct {
	log         *slog.Logger
	appProvider AppProvider
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

func New(log *slog.Logger,
	appProvider AppProvider,
) *Validate {
	return &Validate{
		log:         log,
		appProvider: appProvider,
	}
}

func (v *Validate) ValidateToken(ctx context.Context, token string) (ok bool, err error) {

	const op = "validate.Validate token"

	log := v.log.With(
		slog.String("op", op),
		slog.String("authToken", token),
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

	log.Info("chat token valid")

	return true, nil
}
