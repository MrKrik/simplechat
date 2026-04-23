package suite

import (
	"auth/iternal/config"
	"context"
	"net"
	"strconv"
	"testing"

	auth1 "github.com/MrKrik/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient auth1.AuthClient
}

const grpcHost = "localhost"

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath("../config/local.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials())) //unsecure connection
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}
	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: auth1.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
